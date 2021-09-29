package commands

import (
	"net"
	"testing"

	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func TestHandleNICK(t *testing.T) {
	t.Run("set nick when not registered and nick available", func(t *testing.T) {
		cmd := &Command{
			Command:    "NICK",
			Parameters: []string{"Wiz"},
		}

		m := users.NewMap()

		srv, clnt := testConn(t)
		defer srv.Close()
		defer clnt.Close()
		state := &users.UserState{
			Conn: srv,
		}
		err := cmd.handleNICK(state, m)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
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
	return server, client
}
