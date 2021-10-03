package commands

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"github.com/ilyapetrovMO/Chirc/internal/users"
)

var (
	ErrMalformedCommandString = errors.New("unexpected error")
	ErrCommandStringEmpty     = errors.New("could not parse command, got empty string")
)

func replyErrNicknameInUse(nick string, conn net.Conn) {
	msg := nick + " :Nickname is already in use"
	replyWithError(433, "", msg, conn)
}

func replyErrNoNicknameGiven(state *users.UserState) {
	msg := ":No nickname given"
	replyWithError(431, "", msg, state.Conn)
}

func replyErrAlreadyRegistered(state *users.UserState) {
	msg := ":Unauthorized command (already registered)"
	replyWithError(462, state.User.Nickname, msg, state.Conn)
}

func (c *Command) replyErrNeedMoreParams(state *users.UserState) {
	msg := fmt.Sprintf("%s :Not enough parameters", c.Command)
	replyWithError(461, state.User.Nickname, msg, state.Conn)
}

func (c *Command) replyErrUnknownCommand(state *users.UserState) {
	msg := fmt.Sprintf("%s :Unknown command", c.Command)
	replyWithError(421, state.User.Nickname, msg, state.Conn)
}

func replyWithError(command int, nick, msg string, conn net.Conn) {
	if nick == "" {
		nick = "*"
	}

	w := bufio.NewWriter(conn)
	str := fmt.Sprintf(":%s %d %s %s\r\n", conn.LocalAddr().String(), command, nick, msg)
	w.WriteString(str)
	w.Flush()
}
