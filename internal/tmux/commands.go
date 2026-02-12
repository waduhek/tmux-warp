package tmux

import (
	"os"
	"os/exec"
)

type TmuxCommands interface {
	// IsInTmuxSession checks if the command has been run inside a tmux session.
	IsInTmuxSession() bool

	// HasSession checks if the provided session name exists.
	HasSession(name string) bool

	// CreateSession creates a new tmux session with the provided name and
	// initialises the session path to the provided value.
	CreateSession(name, path string) error

	// AttachSession attaches the caller to the provided tmux session. Ensure
	// that this command is called outside of a tmux session.
	AttachSession(name string) error

	// SwitchClient switches the client to the provided tmux session. Ensure
	// that this command is called inside a tmux session.
	SwitchClient(name string) error
}

type tmuxCommands struct {
}

func NewTmuxCommands() *tmuxCommands {
	return &tmuxCommands{}
}

func (t *tmuxCommands) CreateSession(name, path string) error {
	err := exec.Command(
		"tmux", "new-session", "-d",
		"-s", name,
		"-c", path,
	).Run()
	if err != nil {
		return err
	}

	return nil
}

func (t *tmuxCommands) IsInTmuxSession() bool {
	_, ok := os.LookupEnv("TMUX")
	return ok
}

func (t *tmuxCommands) SwitchClient(name string) error {
	err := exec.Command(
		"tmux", "switch-client",
		"-t", name,
	).Run()
	if err != nil {
		return err
	}

	return nil
}

func (t *tmuxCommands) AttachSession(name string) error {
	cmd := exec.Command(
		"tmux", "attach-session",
		"-t", name,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	return err
}

func (t *tmuxCommands) HasSession(name string) bool {
	err := exec.Command("tmux", "has-session", "-t", name).Run()
	return err == nil
}
