package tmux_test

import (
	"errors"
	"testing"

	"github.com/waduhek/tmux-warp/internal/tmux"
)

type testCase struct {
	name     string
	commands tmux.TmuxCommands
}

type inTmuxSessionCommands struct {
}

func (c *inTmuxSessionCommands) IsInTmuxSession() bool {
	return true
}

func (c *inTmuxSessionCommands) CreateSession(name, path string) error {
	return nil
}

func (c *inTmuxSessionCommands) AttachSession(name string) error {
	return errors.New("attach session should not be called")
}

func (c *inTmuxSessionCommands) SwitchClient(name string) error {
	return nil
}

// ---

type outsideTmuxSessionCommands struct {
}

func (c *outsideTmuxSessionCommands) IsInTmuxSession() bool {
	return false
}

func (c *outsideTmuxSessionCommands) CreateSession(name, path string) error {
	return nil
}

func (c *outsideTmuxSessionCommands) AttachSession(name string) error {
	return nil
}

func (c *outsideTmuxSessionCommands) SwitchClient(name string) error {
	return errors.New("switch client should not be called")
}

// ---

var testCases = []testCase{
	{
		name:     "inside_tmux_session",
		commands: &inTmuxSessionCommands{},
	},
	{
		name:     "outside_tmux_session",
		commands: &outsideTmuxSessionCommands{},
	},
}

func TestCreateAndSwitch(t *testing.T) {
	sessionName := "test"
	sessionPath := "/path/to/test"

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			manager := tmux.NewTmuxManager(test.commands)

			err := manager.CreateSessionAndSwitch(sessionName, sessionPath)
			if err != nil {
				t.Errorf("error while creating session and switching: %s", err)
			}
		})
	}
}
