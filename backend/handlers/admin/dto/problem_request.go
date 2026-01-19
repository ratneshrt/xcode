package dto

type CreateProblemRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Difficulty  string `json:"diff" binding:"required"`
	Description string `json:"desc" binding:"required"`

	Constraints  string `json:"constraints"`
	InputFormat  string `json:"input_format" binding:"required"`
	OutputFormat string `json:"output_format" binding:"required"`

	Examples  []ProblemExampleRequest  `json:"examples" binding:"required"`
	TestCases []ProblemTestCaseRequest `json:"test_cases" binding:"required"`
}

type ProblemExampleRequest struct {
	Input       string `json:"input" binding:"required"`
	Output      string `json:"output" binding:"required"`
	Explanation string `json:"explanation"`
}

type ProblemTestCaseRequest struct {
	Input          string `json:"input" binding:"required"`
	ExpectedOutput string `json:"expected_output" binding:"required"`
	IsHidden       bool   `json:"is_hidden"`
}
