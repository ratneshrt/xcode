package main

import "fmt"

func ValidateProblem(p *ProblemYAML) error {
	if p.Title == "" {
		return fmt.Errorf("title is required")
	}

	if p.Slug == "" {
		return fmt.Errorf("slug is required")
	}

	if p.Difficulty != "easy" && p.Difficulty != "medium" && p.Difficulty != "hard" {
		return fmt.Errorf("invalid difficulty: %s", p.Difficulty)
	}

	if len(p.TestCases) == 0 {
		return fmt.Errorf("at least one test case is required")
	}

	return nil
}
