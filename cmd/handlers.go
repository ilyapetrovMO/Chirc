package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/ilyapetrovMO/Chirc/internal/commands"
	"github.com/ilyapetrovMO/Chirc/internal/options"
	"github.com/ilyapetrovMO/Chirc/internal/users"
)

type application struct {
	Logger    *log.Logger
	Options   *options.Options
	Users     *users.Map
	LocalAddr net.Addr
}

func (a *application) StartAndListen() error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(a.Options.Port))
	if err != nil {
		a.Logger.Printf("ERROR: could not start listening,\n%s", err.Error())
		return err
	}

	a.Logger.Printf("server listening on port %d", a.Options.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			a.Logger.Printf("ERROR: could not accept connection,\n%s", err)
		}

		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 5)

		a.Logger.Printf("connected to %s", conn.RemoteAddr().String())

		go a.handleConnection(conn)
	}
}

func (a *application) handleConnection(conn net.Conn) {
	defer conn.Close()

	state := &users.UserState{
		User: &users.User{},
		Conn: conn,
	}

	sc := bufio.NewScanner(conn)
	for {
		sc.Scan()
		str := sc.Text()

		err := a.handleCmd(state, str)
		if err != nil {
			a.Logger.Printf("encountered error while handling connection: %s", err)
			break
		}
	}

	if state.User.Nickname != "" {
		a.Users.Delete(state.User.Nickname)
	}
}

func (a *application) handleCmd(state *users.UserState, str string) error {
	cmd, err := commands.NewCommand(str)
	if err != nil {
		return err
	}

	err = cmd.Handle(state, a.Users)
	return err
}
