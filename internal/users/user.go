package users

import (
	"errors"
	"log"
	"net"
	"sync"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrNickExists   = errors.New("nick already exists")
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
	Pass         string
	User         *User
	IsRegistered bool
	Conn         net.Conn
}

type Map struct {
	log *log.Logger
	sync.RWMutex
	m map[string]*User
}

func NewMap(logger *log.Logger) *Map {
	m := &Map{
		log: logger,
		m:   make(map[string]*User),
	}
	return m
}

func (m *Map) ReserveNick(nick string) error {
	if m.nickExists(nick) {
		return ErrNickExists
	}

	err := m.insertNew(&User{Nickname: nick})
	return err
}

func (m *Map) RegisterUser(user *User) error {
	if !m.nickExists(user.Nickname) {
		return ErrUserNotFound
	}

	err := m.update(user)
	return err
}

func (m *Map) insertNew(user *User) error {
	if m.nickExists(user.Nickname) {
		return ErrUserExists
	}

	m.Lock()
	defer m.Unlock()

	m.m[user.Nickname] = user
	m.log.Printf("reserved nick %s", user.Nickname)

	return nil
}

func (m *Map) update(user *User) error {
	if !m.nickExists(user.Nickname) {
		return ErrUserNotFound
	}

	if !user.isValid() {
		return ErrUserInvalid
	}

	m.Lock()
	defer m.Unlock()

	m.m[user.Nickname] = user
	m.log.Printf("completed registration for %s", user.Nickname)
	return nil
}

func (m *Map) Delete(nickname string) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, nickname)
}

func (m *Map) usernameExists(username string) bool {
	m.RLock()
	defer m.RUnlock()

	for _, user := range m.m {
		if user.Username == username {
			return true
		}
	}
	return false
}

func (m *Map) nickExists(nickname string) bool {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.m[nickname]; ok {
		return true
	}
	return false
}

// TODO: Change nick does not check if chosen nick is already in use
func (m *Map) ChangeNick(nickname string) error {
	m.Lock()
	defer m.Unlock()

	usr := m.m[nickname]
	m.Delete(nickname)

	usr.Nickname = nickname
	m.m[nickname] = usr
	return nil
}

func (u *User) isValid() bool {
	if u.Nickname == "" || u.Username == "" {
		return false
	}
	return true
}
