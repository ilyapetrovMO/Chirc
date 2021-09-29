package commands

import (
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
			&ErrNicknameInUse{Nickname: "Wiz"},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			srv, clnt := testConn(t)
			defer srv.Close()
			defer clnt.Close()

			test.state.Conn = srv
			m := users.NewMap()
			for _, u := range test.users {
				m.InsertNew(u)
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
