package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/ratneshrt/xcode/database"
)

func main() {
	log.Println("problem worker started")
	database.ConnectProblemDB()

	root := "problems"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || info.Name() != "problem.yaml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var p ProblemYAML
		if err := yaml.Unmarshal(data, &p); err != nil {
			return err
		}

		if err := ValidateProblem(&p); err != nil {
			return err
		}

		if err := LoadProblem(&p); err != nil {
			return err
		}

		log.Println("loaded", p.Slug)
		return nil
	})

	if err != nil {
		log.Fatal("worker failed:", err)
	}
}
