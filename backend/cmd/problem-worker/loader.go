package main

import (
	"errors"
	"time"

	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/models"
	"gorm.io/gorm"
)

func LoadProblem(p *ProblemYAML) error {
	tx := database.ProblemDB.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	var problem models.Problem
	err := tx.Where("slug = ?", p.Slug).First(&problem).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
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
	} else if err != nil {
		tx.Rollback()
		return err
	} else {
		problem.Title = p.Title
		problem.Description = p.Description
		problem.Constraints = p.Constraints
		problem.Status = p.Status
		problem.Difficulty = p.Difficulty
		problem.InputFormat = p.InputFormat
		problem.OutputFormat = p.OutputFormat

		if p.Status == "published" && problem.PublishedAt == nil {
			now := time.Now()
			problem.PublishedAt = &now
		}

		if p.Status == "draft" {
			problem.PublishedAt = nil
		}

		if err := tx.Save(&problem).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("problem_id = ?", problem.ID).Delete(&models.ProblemExample{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("problem_id = ?", problem.ID).Delete(&models.ProblemTestCase{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, ex := range p.Examples {
		e := models.ProblemExample{
			ProblemID:   problem.ID,
			Input:       ex.Input,
			Output:      ex.Output,
			Explanation: ex.Explanation,
		}
		if err := tx.Create(&e).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, tc := range p.TestCases {
		t := models.ProblemTestCase{
			ProblemID:      problem.ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.Hidden,
		}
		if err := tx.Create(&t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
