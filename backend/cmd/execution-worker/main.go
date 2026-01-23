package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ratneshrt/xcode/cmd/execution-worker/executor"
	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/models"
	"github.com/ratneshrt/xcode/queue"
)

func main() {
	log.Println("Execution worker started...")

	database.ConnectAuthDB()
	database.ConnectProblemDB()
	queue.ConnectRedis()

	log.Println("worker connected to dbs and redis")

	runWorker()
}

func runWorker() {
	for {
		log.Println("Waiting for submission job...")

		res, err := queue.RDB.BRPop(
			context.Background(),
			0,
			queue.SubmissionQueue,
		).Result()

		if err != nil {
			log.Println("Redis BRPOP error:", err)
			time.Sleep(time.Second)
			continue
		}

		if len(res) != 2 {
			log.Println("invalid queue payload:", res)
			continue
		}

		log.Println("====runWorker executed & now handleJob will be executed====")

		payload := res[1]
		handleJob(payload)
	}
}

/* func handleJob(payload string) {
	log.Println("job received: ", payload)

	var job queue.SubmissionJob
	if err := json.Unmarshal([]byte(payload), &job); err != nil {
		log.Println("failed to parse job: ", err)
		return
	}

	log.Printf("processing submission_id=%d\n", job.SubmissionID)

	var submission models.Submission
	if err := database.AuthDB.First(&submission, job.SubmissionID).Error; err != nil {
		log.Println("submission not found:", err)
		return
	}

	submission.Status = "running"
	database.AuthDB.Save(&submission)

	log.Printf("submission loaded: id=%d user=%d problem=%d language=%s status=%s\n", submission.ID, submission.UserID, submission.ProblemID, submission.Language, submission.Status)

	var testCases []models.ProblemTestCase
	if err := database.ProblemDB.Where("problem_id = ?", submission.ProblemID).Find(&testCases).Error; err != nil {
		log.Println("failed to fetch test cases:", err)
		submission.Status = "runtime_error"
		submission.Result = "failed to fetch test cases"
		database.AuthDB.Save(&submission)
		return
	}

	if len(testCases) == 0 {
		log.Println("no test cases found")
		submission.Status = "runtime_error"
		submission.Result = "no result cases found"
		database.AuthDB.Save(&submission)
		return
	}

	inputs := make([]string, 0, len(testCases))
	expected := make([]string, 0, len(testCases))

	log.Printf("found %d test cases\n", len(testCases))

	for i, tc := range testCases {
		inputs = append(inputs, tc.Input)
		expected = append(expected, tc.ExpectedOutput)
		log.Printf(
			"testCase #%d hidden=%t question=%q expected=%q\n",
			i+1,
			tc.IsHidden,
			tc.Input,
			tc.ExpectedOutput,
		)
	}

	exec := &executor.GoExecutor{
		TimeLimit: 2 * time.Second,
	}

	result := exec.Execute(submission.Code, inputs, expected)

	if result.Error != "" {
		submission.Status = "runtime_error"
		submission.Result = result.Error
	} else if !result.Passed {
		submission.Status = "wrong_answer"
		submission.Result = fmt.Sprintf(
			"wrong answer at  test case %d",
			result.FailedAt,
		)
	} else {
		submission.Status = "accepted"
		submission.Result = "all test cases passed"
	}

	if err := database.AuthDB.Save(&submission).Error; err != nil {
		log.Println("failed to update submission result: ", err)
		return
	}

	log.Printf("submission %d finished with status=%s\n", submission.ID, submission.Status)
} */

func handleJob(payload string) {
	var job queue.SubmissionJob
	if err := json.Unmarshal([]byte(payload), &job); err != nil {
		log.Println("invalid job payload")
		return
	}

	var submission models.Submission
	if err := database.AuthDB.First(&submission, job.SubmissionID).Error; err != nil {
		log.Println("submission not found")
		return
	}

	submission.Status = "running"
	database.AuthDB.Save(&submission)

	var testCases []models.ProblemTestCase
	database.ProblemDB.Where("problem_id = ?", submission.ProblemID).Order("id ASC").Find(&testCases)

	inputs := make([]string, 0, len(testCases))
	expected := make([]string, 0, len(testCases))

	for _, tc := range testCases {
		inputs = append(inputs, tc.Input)
		expected = append(expected, tc.ExpectedOutput)
	}

	exec := &executor.GoDockerExecutor{
		TimeLimit: 2 * time.Second,
	}

	res := exec.Execute(submission.Code, inputs, expected)

	if res.Error != "" {
		submission.Status = "runtime_error"
		submission.Result = res.Error
	} else if !res.Passed {
		submission.Status = "wrong answer"
		submission.Result = fmt.Sprintf("failed at test case %d", res.FailedAt)
	} else {
		submission.Status = "accepted"
		submission.Result = "all test cases passed"
	}

	database.AuthDB.Save(&submission)

	log.Printf("submission %d finished: %s and the error is: %s\n", submission.ID, submission.Status, submission.Result)
}
