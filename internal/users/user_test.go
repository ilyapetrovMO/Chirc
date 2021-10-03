package users

import (
	"log"
	"testing"
)

func TestUserMap(t *testing.T) {
	t.Run("reserve nick", func(t *testing.T) {
		m := NewMap(log.Default())
		err := m.ReserveNick("Wiz")

		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
	})

	t.Run("try to reserve existing nick", func(t *testing.T) {
		m := NewMap(log.Default())
		m.ReserveNick("Wiz")
		err := m.ReserveNick("Wiz")

		if err != ErrNickExists {
			t.Fatalf("expected error")
		}
	})

	// t.Run("try to insert invalid user", func(t *testing.T) {
	// 	m := NewMap(log.Default())
	// 	usr := &User{}
	// 	err := m.InsertNew(usr)
	// 	want := ErrUserInvalid

	// 	if err == nil {
	// 		t.Fatalf("expected error, got nil")
	// 	}

	// 	if err != want {
	// 		t.Errorf("got %s want %s", err, want)
	// 	}
	// })
}
