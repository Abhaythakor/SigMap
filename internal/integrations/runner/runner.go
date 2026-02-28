package runner

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
)

// Runner handles execution of external CLI tools.
type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

// Execute runs a command and returns its output.
func (r *Runner) Execute(ctx context.Context, name string, args ...string) ([]byte, error) {
	log.Printf("Runner: Executing %s with args %v", name, args)
	
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("command failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
