package executor

type ExecutionResult struct {
	Passed   bool
	Error    string
	FailedAt int
	Output   string
}

type CodeExecutor interface {
	Execute(
		code string,
		testInputs []string,
		expectedOutpuut []string,
	) ExecutionResult
}
