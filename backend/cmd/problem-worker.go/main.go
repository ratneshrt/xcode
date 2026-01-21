package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/ratneshrt/xcode/database"
)

func main() {
	database.ConnectProblemDB()

	files, err := filepath.Glob("problems/*/problem.yaml") //specific the file path
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		var problem ProblemYAML
		if err := yaml.Unmarshal(data, &problem); err != nil {
			log.Fatal(err)
		}

		if err := ValidateProblem(&problem); err != nil {
			log.Fatalf("validation failed for %s: %v", file, err)
		}

		if err := LoadProblem(&problem); err != nil {
			log.Fatalf("Load failed for %s: %v", file, err)
		}

		log.Println("loaded:", problem.Slug)
	}
}
