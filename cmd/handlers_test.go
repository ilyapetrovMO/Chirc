package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/ilyapetrovMO/Chirc/internal/options"
	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func TestHandleConnection(t *testing.T) {
	client, server := testConn(t)
	defer client.Close()
	defer server.Close()

	client.SetReadDeadline(time.Now().Add(time.Second * 1))

	app := &application{
		Users:  users.NewMap(log.Default()),
		Logger: log.Default(),
		Options: &options.Options{
			Port: 6697,
		},
		LocalAddr: server.LocalAddr(),
	}

	// t.Run("PING PONG sequence", func(t *testing.T) {
	// 	client.Write([]byte("PING\r\n"))
	// 	app.handleConnection(server)

	// 	s := bufio.NewScanner(client)
	// 	s.Scan()
	// 	got := s.Text()
	// 	want := fmt.Sprintf("PONG %s", app.LocalAddr.String())

	// 	if got != want {
	// 		t.Errorf("got %s want %s", got, want)
	// 	}
	// })

	t.Run("NICK USER sequence for new user", func(t *testing.T) {
		nick := "gosha"
		usrname := "bigboy"

		w := bufio.NewWriter(client)
		w.WriteString(fmt.Sprintf("NICK %s\n\r", nick))
		w.WriteString(fmt.Sprintf("USER %s * * :Full Name\n\r", usrname))

		if err := w.Flush(); err != nil {
			t.Fatalf("%s", err)
		}

		go app.handleConnection(server)

		s := bufio.NewScanner(client)
		s.Scan()
		got := s.Text()
		want := fmt.Sprintf(":%s 001 %s :Welcome to the Internet Relay Network %s!%s@%s",
			app.LocalAddr.String(), nick, nick, usrname, client.LocalAddr().String())

		if got != want {
			t.Errorf("\ngot\n%s want\n%s", got, want)
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
