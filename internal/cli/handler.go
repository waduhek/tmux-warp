package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/waduhek/tmux-warp/internal/tmux"
	"github.com/waduhek/tmux-warp/internal/warp"
)

var ErrHomeEnvNotSet = errors.New("value of environment variable HOME not set")

type CLIHandler struct {
	manager tmux.TmuxManager
	warp    warp.Warper
}

// NewCLIHandler creates a new instance of CLIHandler.
func NewCLIHandler(
	manager tmux.TmuxManager,
	warp warp.Warper,
) *CLIHandler {
	return &CLIHandler{manager, warp}
}

// StartSession checks if there is a tmux session with the provided name. If the
// session does not exist, creates a new session by reading the warp config file
// and finding the target path and then switches to the session.
func (h *CLIHandler) StartSession(name string) error {
	warpConfigPath, err := h.getWarpConfigPath()
	if err != nil {
		return err
	}

	ctx := context.Background()
	path, err := h.warp.GetWarpPathByName(ctx, warpConfigPath, name)
	if err != nil {
		return err
	}

	err = h.manager.CreateSessionAndSwitch(name, path)
	if err != nil {
		return err
	}

	return nil
}

func (h *CLIHandler) getWarpConfigPath() (string, error) {
	homeEnv := os.Getenv("HOME")
	if homeEnv == "" {
		return "", ErrHomeEnvNotSet
	}

	warpPath := filepath.Join(homeEnv, ".warprc")
	return warpPath, nil
}
