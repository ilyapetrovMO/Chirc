package commands

import (
	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func (c *Command) handleNICK(state *users.UserState, m *users.Map) error {
	if len(c.Parameters) != 1 {
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

	m.ChangeNick(state.User.Username, nick)

	return nil
}
