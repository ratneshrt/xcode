package main

import (
	"time"

	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/models"
)

func LoadProblem(p *ProblemYAML) error {
	tx := database.ProblemDB.Begin()

	var problem models.Problem
	err := tx.Where("slug = ?", p.Slug).First(&problem).Error

	if err != nil {
		problem = models.Problem{
			Title:        p.Title,
			Slug:         p.Slug,
			Difficulty:   p.Difficulty,
			Description:  p.Description,
			Constraints:  p.Constraints,
			InputFormat:  p.InputFormat,
			OutputFormat: p.OutputFormat,
			Status:       p.Status,
		}

		if p.Status == "published" {
			now := time.Now()
			problem.PublishedAt = &now
		}

		if err := tx.Create(&problem).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		problem.Title = p.Title
		problem.Description = p.Description
		problem.Constraints = p.Constraints
		tx.Save(&problem)
	}

	tx.Where("problem_id = ?", problem.ID).Delete(&models.ProblemExample{})
	tx.Where("problem_id = ?", problem.ID).Delete(&models.ProblemTestCase{})

	for _, ex := range p.Examples {
		tx.Create(&models.ProblemExample{
			ProblemID:   problem.ID,
			Input:       ex.Input,
			Output:      ex.Output,
			Explanation: ex.Explanation,
		})
	}

	for _, tc := range p.TestCases {
		tx.Create(&models.ProblemTestCase{
			ProblemID:      problem.ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.Hidden,
		})
	}

	return tx.Commit().Error
}
