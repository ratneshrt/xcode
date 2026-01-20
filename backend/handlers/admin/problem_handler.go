package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/handlers/admin/dto"
	"github.com/ratneshrt/xcode/models"
)

func CreateProblem(c *gin.Context) {
	var req dto.CreateProblemRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := database.ProblemDB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to start transaction",
		})
		return
	}

	problem := models.Problem{
		Title:        req.Title,
		Slug:         req.Slug,
		Difficulty:   req.Difficulty,
		Description:  req.Description,
		Constraints:  req.Constraints,
		InputFormat:  req.InputFormat,
		OutputFormat: req.OutputFormat,
		Status:       "draft",
	}

	if err := tx.Create(&problem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, ex := range req.Examples {
		example := models.ProblemExample{
			ProblemID:   problem.ID,
			Input:       ex.Input,
			Output:      ex.Output,
			Explanation: ex.Explanation,
		}

		if err := tx.Create(&example).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	for _, tc := range req.TestCases {
		testCase := models.ProblemTestCase{
			ProblemID:      problem.ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
		}

		if err := tx.Create(&testCase).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"problem_id": problem.ID,
		"status":     "draft",
	})
}

func PublishProblem(c *gin.Context) {
	id := c.Param("id")

	var problem models.Problem

	if err := database.ProblemDB.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "problem not found",
		})
		return
	}

	if problem.Status == "published" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "problem is already published",
		})
		return
	}

	now := time.Now()
	problem.Status = "published"
	problem.PublishedAt = &now

	if err := database.ProblemDB.Save(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "problem published",
		"problem_id":   problem.ID,
		"published_at": now,
	})
}
