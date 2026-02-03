package tmux

type TmuxManager interface {
	// CreateSessionAndSwitch checks if the provided session already exists on
	// the default server. If the session does not exist, it creates the session
	// and sets the default path to the provided path value and attaches the
	// client to the session.
	CreateSessionAndSwitch(name, path string) error
}

type tmuxManager struct {
	commands TmuxCommands
}

// NewTmuxManager creates a new instance of tmux commands manager.
func NewTmuxManager(commands TmuxCommands) *tmuxManager {
	return &tmuxManager{commands}
}

func (m *tmuxManager) CreateSessionAndSwitch(name, path string) error {
	err := m.commands.CreateSession(name, path)
	if err != nil {
		return err
	}

	if m.commands.IsInTmuxSession() {
		return m.commands.SwitchClient(name)
	}
	return m.commands.AttachSession(name)
}
