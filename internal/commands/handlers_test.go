package commands

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func TestHandleNICK(t *testing.T) {
	tt := []struct {
		name      string
		command   *Command
		users     []*users.User
		state     *users.UserState
		want      *users.UserState
		wantError error
	}{
		{
			"set nick when not registered and nick avalible",
			&Command{
				Command:    "NICK",
				Parameters: []string{"Wiz"},
			},
			[]*users.User{},
			&users.UserState{User: &users.User{}},
			&users.UserState{User: &users.User{Nickname: "Wiz"}},
			nil,
		}, {
			"set nick when not registered and nick not available",
			&Command{
				Command:    "NICK",
				Parameters: []string{"Wiz"},
			},
			[]*users.User{{Username: "1", Nickname: "Wiz"}},
			&users.UserState{User: &users.User{}},
			&users.UserState{User: &users.User{}},
			nil,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			srv, clnt := testConn(t)
			defer srv.Close()
			defer clnt.Close()

			test.state.Conn = srv
			m := users.NewMap(log.Default())
			for _, u := range test.users {
				m.ReserveNick(u.Nickname)
				m.RegisterUser(u)
			}

			err := test.command.handleNICK(test.state, m)
			if err != nil {
				if test.wantError != nil {
					if test.wantError.Error() != err.Error() {
						t.Errorf("got %s want %s", err, test.wantError)
						return
					}
					return
				}

				t.Fatalf("unexpected error %s", err)
			}

			test.want.Conn = srv

			if !reflect.DeepEqual(test.want, test.state) {
				t.Errorf("got %v want %v", test.state, test.want)
			}
		})
	}
}

func TestHandleUSER(t *testing.T) {
	tt := []struct {
		name      string
		command   *Command
		users     []*users.User
		state     *users.UserState
		wantState *users.UserState
		wantError error
	}{
		{
			"set username when not registered",
			&Command{
				Command:    "USER",
				Parameters: []string{"josh", "*", "*"},
				Trailing:   "Josh Trailing",
			},
			[]*users.User{{Nickname: "Wiz"}},
			&users.UserState{User: &users.User{Nickname: "Wiz"}},
			&users.UserState{
				User: &users.User{
					Nickname: "Wiz",
					Username: "josh",
					FullName: "Josh Trailing"},
				IsRegistered: true,
			},
			nil,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			srv, clnt := testConn(t)
			defer srv.Close()
			defer clnt.Close()

			test.state.Conn = srv
			m := users.NewMap(log.Default())
			for _, u := range test.users {
				m.ReserveNick(u.Nickname)
			}

			err := test.command.handleUSER(test.state, m)
			if err != nil {
				if test.wantError != nil {
					if test.wantError.Error() != err.Error() {
						t.Errorf("got %s want %s", err, test.wantError)
						return
					}
					return
				}

				t.Fatalf("unexpected error %s", err)
			}

			test.wantState.Conn = srv

			if !reflect.DeepEqual(test.wantState, test.state) {
				t.Errorf("got %v want %v", test.state, test.wantState)
			}
		})
	}

	t.Run("test RPL_WELCOME recieved at the end of successful USER command", func(t *testing.T) {
		srv, clnt := testConn(t)
		defer srv.Close()
		defer clnt.Close()

		state := &users.UserState{
			User: &users.User{
				Nickname: "Wiz",
				Username: "josh",
				Conn:     srv,
			},
			Conn:         srv,
			IsRegistered: false,
		}

		cmd := &Command{
			Command:    "USER",
			Parameters: []string{"josh", "*", "*"},
			Trailing:   "Josh Trailing",
		}

		m := users.NewMap(log.Default())
		m.ReserveNick("Wiz")
		cmd.handleUSER(state, m)

		r := bufio.NewScanner(clnt)
		r.Scan()
		got := r.Text()

		want := fmt.Sprintf(":%s 001 %s :Welcome to the Internet Relay Network %s!%s@%s",
			srv.LocalAddr().String(), state.User.Nickname, state.User.Nickname, state.User.Username, clnt.LocalAddr().String())

		if got != want {
			t.Errorf("\ngot \n%s want \n%s", got, want)
		}
	})
}

func testConn(t *testing.T) (net.Conn, net.Conn) {
	t.Helper()

	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	var server net.Conn
	go func() {
		defer ln.Close()
		server, _ = ln.Accept()
	}()

	client, _ := net.Dial("tcp", ln.Addr().String())
	time.Sleep(time.Millisecond * 20)
	return client, server
}
