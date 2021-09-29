package commands

import (
	"errors"
	"fmt"
)

var (
	ErrInParams        = errors.New("error in parameter list")
	ErrNoNicknameGiven = errors.New("431 :No nickname given")
)

type ErrUnknownCommand struct {
	Command string
}

func (e *ErrUnknownCommand) Error() string {
	return fmt.Sprintf("421 %s :unknown command", e.Command)
}

type ErrNicknameInUse struct {
	Nickname string
}

func (e *ErrNicknameInUse) Error() string {
	return fmt.Sprintf("433 %s :Nickname already in use", e.Nickname)
}
