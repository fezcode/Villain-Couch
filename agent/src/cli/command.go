package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// CommandRunner manages running a command in the background.
type CommandRunner struct {
	mu  sync.Mutex
	cmd *exec.Cmd

	// done is closed when the command finishes. The value sent is the result
	// of cmd.Wait().
	done chan error
}

// New creates a new CommandRunner for the given command and arguments.
func New(name string, arg ...string) *CommandRunner {
	return &CommandRunner{
		cmd: exec.Command(name, arg...),
	}
}

func NewCommandRunnerForVLC(args VLCRunnerArguments) *CommandRunner {
	return &CommandRunner{
		cmd: PrepareVLCCommand(args),
	}
}

// runs the command in the background. It is safe to call this function multiple times.
func (c *CommandRunner) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If done is not nil, the command has already been started.
	if c.done != nil {
		return fmt.Errorf("command has already been started")
	}

	// Pipe the command's stdout and stderr to the parent process to see its output.
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Initialize the done channel and start a goroutine to wait for the command to finish.
	c.done = make(chan error, 1)
	go func() {
		c.done <- c.cmd.Wait()
		close(c.done)
	}()

	return nil
}

// sends an interrupt signal to the running command.
// It returns an error if the command is not running or if the signal fails.
func (c *CommandRunner) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If the process is nil, the command hasn't been started or has already finished.
	if c.cmd.Process == nil {
		return fmt.Errorf("command is not running")
	}

	// Send the interrupt signal (like Ctrl+C).
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		// Use Process.Kill() for a more reliable termination on Windows.
		if err := c.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}

	}

	return nil
}

// returns a channel that is closed when the command finishes.
// The value sent on the channel is the error result from cmd.Wait().
func (c *CommandRunner) Done() <-chan error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.done
}

func (c *CommandRunner) PID() int {
	return c.cmd.Process.Pid
}
