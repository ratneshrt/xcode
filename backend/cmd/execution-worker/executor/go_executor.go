package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ratneshrt/xcode/cmd/execution-worker/utils"
)

type GoExecutor struct {
	TimeLimit time.Duration
}

func (g *GoExecutor) Execute(
	code string,
	testInputs []string,
	expectedOutputs []string,
) ExecutionResult {
	tmpDir, err := os.MkdirTemp("", "xcode-go-*")
	if err != nil {
		return ExecutionResult{Error: err.Error()}
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(sourceFile, []byte(code), 0644); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	binaryPath := filepath.Join(tmpDir, "app")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, sourceFile)
	var buildErr bytes.Buffer
	buildCmd.Stderr = &buildErr

	if err := buildCmd.Run(); err != nil {
		return ExecutionResult{
			Error: "compile error: " + buildErr.String(),
		}
	}

	for i, input := range testInputs {
		cmd := exec.Command(binaryPath)

		cmd.Stdin = bytes.NewBufferString(input)

		var out bytes.Buffer
		cmd.Stdout = &out

		if err := runwithTimeout(cmd, g.TimeLimit); err != nil {
			return ExecutionResult{
				Error:    "time limti exceed",
				FailedAt: i + 1,
			}
		}

		actual := utils.Normalize(out.String())
		expected := utils.Normalize(expectedOutputs[i])

		if actual != expected {
			return ExecutionResult{
				Passed:   false,
				FailedAt: i + 1,
				Output:   actual,
			}
		}
	}

	return ExecutionResult{Passed: true}
}

func runwithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return fmt.Errorf("timeout")
	case err := <-done:
		return err
	}
}
