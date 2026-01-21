package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/handlers/public/dto"
	"github.com/ratneshrt/xcode/models"
)

func CreateSubmission(c *gin.Context) {
	var req dto.CreateSubmissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDAny, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDAny.(uint)

	submission := models.Submission{
		UserID:    userID,
		ProblemID: req.ProblemID,
		Language:  req.Language,
		Code:      req.Code,
		Status:    "pending",
	}

	if err := database.AuthDB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create submission",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submission_id": submission.ID,
		"status":        submission.Status,
	})
}
