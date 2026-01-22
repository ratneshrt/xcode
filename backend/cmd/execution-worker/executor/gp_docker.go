package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ratneshrt/xcode/cmd/execution-worker/utils"
)

type GoDockerExecutor struct {
	TimeLimit time.Duration
}

func (g *GoDockerExecutor) Execute(
	code string,
	testInputs []string,
	expectedOutputs []string,
) ExecutionResult {

	tmpDir, err := os.MkdirTemp("", "xcode-docker-go-*")

	if err != nil {
		return ExecutionResult{Error: err.Error()}
	}
	defer os.RemoveAll(tmpDir)

	sourcePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(sourcePath, []byte(code), 0644); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	for i := 0; i < len(testInputs); i++ {
		actual, err := g.runInDocker(tmpDir, testInputs[i])
		if err != nil {
			return ExecutionResult{
				Error:    err.Error(),
				FailedAt: i + 1,
			}
		}

		if utils.Normalize(actual) != utils.Normalize(expectedOutputs[i]) {
			return ExecutionResult{
				Passed:   false,
				FailedAt: i + 1,
				Output:   actual,
			}
		}
	}

	return ExecutionResult{Passed: true}
}

func (g *GoDockerExecutor) runInDocker(dir string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.TimeLimit)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1",
		"-v", fmt.Sprintf("%s:/app:ro", dir),
		"-w", "/app",
		"-i",
		"golang:1.22-alpine",
		"sh", "-c",
		"go run main.go",
	)

	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit exceeded")
	}

	if err != nil {
		err_s := stderr.String()
		return "", fmt.Errorf("runtime error: %w", &err_s)
	}

	return stdout.String(), nil
}
