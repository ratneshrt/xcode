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
	tmpDir, err := os.MkdirTemp("", "xcode-go-*")
	if err != nil {
		return ExecutionResult{Error: err.Error()}
	}
	defer os.RemoveAll(tmpDir)

	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	container := fmt.Sprintf("xcode-g-%d", time.Now().UnixNano())

	if err := run("docker", "create", "--name", container, "xcode-go-runner"); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	defer run("docker", "rm", "-f", container)

	if err := run("docker", "start", container); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	if err := run("docker", "cp", mainFile, container+":/app/main.go"); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	if err := run(
		"docker", "exec", container,
		"sh", "-c", "GO111MODULE=off go build -o app main.go",
	); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	for i := range testInputs {
		out, err := runwithInput(
			g.TimeLimit,
			testInputs[i],
			"docker", "exec", "-i", container, "./app",
		)
		if err != nil {
			return ExecutionResult{
				Error:    err.Error(),
				FailedAt: i + 1,
			}
		}

		if utils.Normalize(out) != utils.Normalize(expectedOutputs[i]) {
			return ExecutionResult{
				Passed:   false,
				FailedAt: i + 1,
				Output:   out,
			}
		}
	}

	return ExecutionResult{Passed: true}
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	var errBuf bytes.Buffer
	c.Stderr = &errBuf
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s", errBuf.String())
	}

	return nil
}

func runwithInput(timeout time.Duration, input string, cmd string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdin = strings.NewReader(input)

	var out, errBuf bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errBuf

	err := c.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit exceeded")
	}

	if err != nil {
		return "", fmt.Errorf("runtime error: %s", errBuf.String())
	}

	return out.String(), nil
}
