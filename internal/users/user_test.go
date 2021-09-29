package users

import (
	"testing"
)

func TestUserMap(t *testing.T) {
	t.Run("insert new user", func(t *testing.T) {
		m := NewMap()
		usr := User{
			Nickname: "Nick",
			Username: "User",
		}
		err := m.InsertNew(usr)

		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
	})

	t.Run("try to insert invalid user", func(t *testing.T) {
		m := NewMap()
		usr := User{}
		err := m.InsertNew(usr)
		want := ErrUserInvalid

		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if err != want {
			t.Errorf("got %s want %s", err, want)
		}
	})
}
