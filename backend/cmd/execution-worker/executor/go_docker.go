package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ratneshrt/xcode/cmd/execution-worker/utils"
)

type GoDockerExecutor struct {
	TimeLimit time.Duration
}

var pullOnce sync.Once

func (g *GoDockerExecutor) Execute(
	code string,
	testInputs []string,
	expectedOutputs []string,
) ExecutionResult {

	pullOnce.Do(func() {
		pullImage("golang:1.22-alpine")
		pullImage("alpine")
	})

	volumeName := fmt.Sprintf("xcode-go-%d", time.Now().UnixNano())
	createVolume(volumeName)
	defer removeVolume(volumeName)

	if err := buildInContainer(volumeName, code, g.TimeLimit); err != nil {
		return ExecutionResult{Error: err.Error()}
	}

	for i := 0; i < len(testInputs); i++ {
		actual, err := runBinary(volumeName, testInputs[i], g.TimeLimit)
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

func runBinary(volume, input string, limit time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1",
		"-v", volume+":/app",
		"-w", "/app",
		"alpine",
		"./app",
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
		return "", fmt.Errorf("runtime error: %s", err_s)
	}

	return stdout.String(), nil
}

/* func (g *GoDockerExecutor) buildBinary(dir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.TimeLimit)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm",
		"--network", "none",
		"-v", fmt.Sprintf("%s:/app", dir),
		"-w", "/app",
		"golang:1.22-alpine",
		"go", "build", "-o", "app", "main.go",
	)

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compile error: %s", errBuf.String())
	}

	return nil
}*/

func buildInContainer(volume, code string, limit time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "--rm",
		"--network", "none",
		"-v", volume+":/app",
		"-w", "/app",
		"-e", "CODE="+code,
		"golang:1.22-alpine",
		"sh", "-c",
		`printf "%s" "$CODE" > main.go && GO111MODULE=off go build -o app main.go`,
	)

	cmd.Stdin = strings.NewReader(code)

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compile error: %s", errBuf.String())
	}

	return nil
}

func pullImage(image string) {
	cmd := exec.Command("docker", "pull", image)
	_ = cmd.Run()
}

func createVolume(name string) {
	_ = exec.Command("docker", "volume", "create", name).Run()
}

func removeVolume(name string) {
	_ = exec.Command("docker", "volume", "rm", "-f", name).Run()
}
