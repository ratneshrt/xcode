package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/models"
)

type ProblemListResponse struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Difficulty string `json:"difficulty"`
}

type ProblemExampleResponse struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation"`
}

type ProbelmDetailResponse struct {
	Title        string                   `json:"title"`
	Difficulty   string                   `json:"difficulty"`
	Description  string                   `json:"description"`
	Constraints  string                   `json:"constraints"`
	InputFormat  string                   `json:"input_format"`
	OutputFormat string                   `json:"output_format"`
	Examples     []ProblemExampleResponse `json:"examples"`
}

func ListProblems(c *gin.Context) {
	var problems []ProblemListResponse

	err := database.ProblemDB.Model(&models.Problem{}).Select("title, slug, difficulty").Where("status = ?", "published").Order("published_at DESC").Scan(&problems).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch problems",
		})
		return
	}

	c.JSON(http.StatusOK, problems)

}

func GetProblemBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var problem models.Problem

	if err := database.ProblemDB.Where("slug = ? AND status = ?", slug, "published").First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "problem not found",
		})
		return
	}

	var examples []ProblemExampleResponse

	database.ProblemDB.Model(&models.ProblemExample{}).Select("input, output, explanation").Where("problem_id = ?", problem.ID).Scan(&examples)

	resp := ProbelmDetailResponse{
		Title:        problem.Title,
		Difficulty:   problem.Difficulty,
		Description:  problem.Description,
		Constraints:  problem.Constraints,
		InputFormat:  problem.InputFormat,
		OutputFormat: problem.OutputFormat,
		Examples:     examples,
	}

	c.JSON(http.StatusOK, resp)
}
