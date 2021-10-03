package commands

import (
	"strings"
	"unicode"

	"github.com/ilyapetrovMO/Chirc/internal/users"
)

type Command struct {
	Prefix         string
	Command        string
	Parameters     []string
	Trailing       string
	unknownCommand bool
}

func NewCommand(cmd string) (*Command, error) {
	if cmd == "" {
		return nil, ErrCommandStringEmpty
	}

	c := &Command{
		Parameters: []string{},
	}

	fields, err := fields(cmd)
	if err != nil {
		return nil, err
	}

	idx := 0
	if fields[0][0] == ':' {
		c.Prefix = strings.TrimPrefix(fields[0], ":")
		idx++
	}

	if isValidCommand(fields[idx]) {
		c.Command = fields[idx]
		idx++
	} else {
		c.unknownCommand = true
	}

	for ; idx < len(fields); idx++ {
		if fields[idx][0] == ':' {
			c.Trailing = strings.TrimPrefix(fields[idx], ":")
			break
		}

		c.Parameters = append(c.Parameters, fields[idx])
	}

	return c, nil
}

func (c *Command) Handle(state *users.UserState, users *users.Map) error {
	if c.unknownCommand {
		c.replyErrUnknownCommand(state)
		return nil
	}
	var err error
	switch c.Command {
	// PASS is non-functional
	case "PASS":
		err = c.handlePASS(state, users)
	case "NICK":
		err = c.handleNICK(state, users)
	case "USER":
		err = c.handleUSER(state, users)
	}
	return err
}

// Works like strings.Fields, except if it finds a ':' rune,
// everything after it will be returned as a single string
func fields(cmd string) ([]string, error) {
	res := []string{""}
	cmd = strings.TrimSpace(cmd)

	idx := 0
	midWord := false
	for i, r := range cmd {
		if unicode.IsSpace(r) {
			idx++
			res = append(res, "")
			midWord = false
			continue
		}
		if r == ':' && midWord {
			return nil, ErrMalformedCommandString
		}
		if r == ':' && idx != 0 {
			if i+1 >= len(cmd) {
				return nil, ErrMalformedCommandString
			}
			res[idx] = cmd[i:]
			break
		}

		midWord = true
		res[idx] += string(r)
	}

	return res, nil
}

func isValidCommand(cmd string) bool {
	switch cmd {
	case "PASS":
	case "NICK":
	case "USER":
	case "PING":
	case "PONG":
	default:
		return false
	}
	return true
}
