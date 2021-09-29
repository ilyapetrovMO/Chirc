package main

import (
	"bufio"
	"fmt"

	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func (a *application) ReplyWithError(state *users.UserState, cmd string) error {
	w := bufio.NewWriter(state.Conn)

	w.WriteString(fmt.Sprintf(":%s %s\r\n", state.Conn.LocalAddr().String(), cmd))
	err := w.Flush()
	if err != nil {
		return err
	}

	return nil
}
