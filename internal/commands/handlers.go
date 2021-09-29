package commands

import (
	"github.com/ilyapetrovMO/Chirc/internal/users"
)

// TODO: reusable Command validation logic
func (c *Command) handlePASS(state *users.UserState, m *users.Map) error {
	if len(c.Parameters) < 1 {
		return &ErrNeedMoreParams{c.Command}
	}

	if state.LoggedIn {
		return ErrAlreadyRegistered
	}

	state.User.Pass = c.Parameters[0]
	return nil
}

func (c *Command) handleNICK(state *users.UserState, m *users.Map) error {
	if len(c.Parameters) < 1 {
		return ErrNoNicknameGiven
	}

	nick := c.Parameters[0]

	if !state.LoggedIn {
		if m.NickExists(nick) {
			return &ErrNicknameInUse{nick}
		}

		state.User.Nickname = nick
		return nil
	}

	err := m.ChangeNick(nick)

	return err
}

// TODO: in progress
func (c *Command) handleUSER(state *users.UserState, m *users.Map) error {
	if len(c.Parameters) < 3 {
		return &ErrNeedMoreParams{c.Command}
	}

	if c.Trailing == "" {
		return &ErrNeedMoreParams{c.Command}
	}

	return nil
}
