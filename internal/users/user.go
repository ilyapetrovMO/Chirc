package users

// TODO: replace map with persistant storage
import (
	"errors"
	"net"
	"sync"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("specified user does not exist")
	ErrUserInvalid  = errors.New("user struct malformed")
)

type User struct {
	Pass     string
	Username string
	Nickname string
	FullName string
	Conn     net.Conn
}

type UserState struct {
	Pass     string
	User     *User
	LoggedIn bool
	Conn     net.Conn
}

type Map struct {
	sync.RWMutex
	m map[string]*User
}

func NewMap() *Map {
	m := &Map{
		m: make(map[string]*User),
	}
	return m
}

func (m *Map) InsertNew(user *User) error {
	if m.NickExists(user.Nickname) {
		return ErrUserExists
	}

	if !user.isValid() {
		return ErrUserInvalid
	}

	m.Lock()
	defer m.Unlock()

	m.m[user.Nickname] = user

	return nil
}

func (m *Map) Update(user *User) error {
	if !m.NickExists(user.Nickname) {
		return ErrUserNotFound
	}

	if !user.isValid() {
		return ErrUserInvalid
	}

	m.Lock()
	defer m.Unlock()

	m.m[user.Nickname] = user
	return nil
}

func (m *Map) UsernameExists(username string) bool {
	m.RLock()
	defer m.RUnlock()

	for _, user := range m.m {
		if user.Username == username {
			return true
		}
	}
	return false
}

// CAUTION: Change nick does not check if chosen nick is already in use
func (m *Map) ChangeNick(nickname string) error {
	m.Lock()
	defer m.Unlock()

	usr := m.m[nickname]
	usr.Nickname = nickname
	m.m[nickname] = usr
	return nil
}

func (m *Map) NickExists(nickname string) bool {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.m[nickname]; ok {
		return true
	}
	return false
}

func (u *User) isValid() bool {
	if u.Nickname == "" || u.Username == "" {
		return false
	}
	return true
}
