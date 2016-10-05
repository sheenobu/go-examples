package users

import (
	"fmt"
	"sync"

	"github.com/sheenobu/go-examples/rest"
)

// memory based storage for the user.

type memoryStorage struct {
	users map[string]rest.User
	lock  sync.RWMutex
}

func (ms *memoryStorage) List() (ux []rest.User, err error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	for _, val := range ms.users {
		ux = append(ux, val)
	}

	return
}

func (ms *memoryStorage) Get(username string) (*rest.User, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()

	u, ok := ms.users[username]
	if !ok {
		return nil, &notFound{username: username}
	}

	return &u, nil
}

func (ms *memoryStorage) Save(user *rest.User) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.users[user.Username] = *user

	return nil
}

func (ms *memoryStorage) Delete(user *rest.User) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	_, ok := ms.users[user.Username]
	if !ok {
		return &notFound{username: user.Username}
	}

	delete(ms.users, user.Username)

	return nil
}

type notFound struct {
	username string
}

func (nf *notFound) Error() string {
	return fmt.Sprintf("User '%v' not found", nf.username)
}

func (nf *notFound) Code() int {
	return 404
}
