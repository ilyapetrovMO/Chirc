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
	User     User
	LoggedIn bool
	Conn     net.Conn
}

type Map struct {
	sync.RWMutex
	m map[string]User
}

func NewMap() *Map {
	m := &Map{
		m: make(map[string]User),
	}
	return m
}

func (m *Map) InsertNew(user User) error {
	if m.UsernameExists(user.Username) {
		return ErrUserExists
	}

	if !user.isValid() {
		return ErrUserInvalid
	}

	m.Lock()
	defer m.Unlock()

	m.m[user.Username] = user

	return nil
}

func (m *Map) Update(user User) error {
	if !m.UsernameExists(user.Username) {
		return ErrUserNotFound
	}

	if !user.isValid() {
		return ErrUserInvalid
	}

	m.Lock()
	defer m.Unlock()

	m.m[user.Username] = user
	return nil
}

func (m *Map) UsernameExists(username string) bool {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.m[username]; ok {
		return true
	}
	return false
}

// CAUTION: Change nick does not check if chosen nick is already in use
func (m *Map) ChangeNick(username, nickname string) error {
	m.Lock()
	defer m.Unlock()

	usr := m.m[username]
	usr.Nickname = nickname
	m.m[username] = usr
	return nil
}

func (m *Map) NickExists(nick string) bool {
	m.RLock()
	defer m.RUnlock()

	for _, u := range m.m {
		if u.Nickname == nick {
			return true
		}
	}
	return false
}

func (u *User) isValid() bool {
	if u.Nickname == "" || u.Username == "" {
		return false
	}
	return true
}
