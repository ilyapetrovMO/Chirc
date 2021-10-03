package commands

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/ilyapetrovMO/Chirc/internal/users"
)

// TODO: reusable Command validation logic
func (c *Command) handlePASS(state *users.UserState, m *users.Map) error {
	if len(c.Parameters) < 1 {
		c.replyErrNeedMoreParams(state)
		return nil
	}

	if state.IsRegistered {
		replyErrAlreadyRegistered(state)
		return nil
	}

	state.User.Pass = c.Parameters[0]
	return nil
}

func (c *Command) handleNICK(state *users.UserState, m *users.Map) error {
	if len(c.Parameters) < 1 {
		replyErrNoNicknameGiven(state)
		return nil
	}

	nick := c.Parameters[0]

	if !state.IsRegistered {
		err := m.ReserveNick(nick)
		if err == users.ErrNickExists {
			replyErrNicknameInUse(nick, state.Conn)
			return nil
		} else if err != nil {
			return err
		}

		state.User.Nickname = nick
		return nil
	}

	err := m.ChangeNick(nick)
	return err
}

func (c *Command) handleUSER(state *users.UserState, m *users.Map) error {
	if state.IsRegistered {
		replyErrAlreadyRegistered(state)
		return nil
	}

	if state.User.Nickname == "" {
		return errors.New("USER command used before NICK")
	}

	if len(c.Parameters) < 3 || c.Trailing == "" {
		c.replyErrNeedMoreParams(state)
		return nil
	}

	state.User.Username = c.Parameters[0]
	state.User.FullName = c.Trailing
	err := m.RegisterUser(state.User)
	if err != nil {
		return err
	}

	state.IsRegistered = true
	err = sendRplWelcome(state)

	return err
}

func sendRplWelcome(state *users.UserState) error {
	wlcm := fmt.Sprintf(":%s 001 %s :Welcome to the Internet Relay Network %s!%s@%s\r\n",
		state.Conn.LocalAddr().String(), state.User.Nickname, state.User.Nickname, state.User.Username, state.Conn.RemoteAddr().String())

	w := bufio.NewWriter(state.Conn)
	w.WriteString(wlcm)
	err := w.Flush()
	return err
}
