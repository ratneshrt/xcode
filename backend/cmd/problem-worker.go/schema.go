package main

type ProblemYAML struct {
	Title      string `yaml:"title"`
	Slug       string `yaml:"slug"`
	Difficulty string `yaml:"difficulty"`
	Status     string `yaml:"status"`

	Description string `yaml:"description"`
	Constraints string `yaml:"constraints"`

	InputFormat  string `yaml:"input_format"`
	OutputFormat string `yaml:"output_format"`

	Examples  []ExamplesYAML `yaml:"examples"`
	TestCases []TestCaseYAML `yaml:"test_cases"`
}

type ExamplesYAML struct {
	Input       string `yaml:"input"`
	Output      string `yaml:"output"`
	Explanation string `yaml:"explanation"`
}

type TestCaseYAML struct {
	Input          string `yaml:"input"`
	ExpectedOutput string `yaml:"expected_output"`
	Hidden         bool   `yaml:"is_hidden"`
}
