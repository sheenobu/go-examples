package users

import "github.com/sheenobu/go-examples/rest"

// Storage defines the interface for storing and retrieving users
type Storage interface {
	List() ([]rest.User, error)
	Get(username string) (*rest.User, error)
	Save(user *rest.User) error
	Delete(user *rest.User) error
}

// NewStorage creates a new storage object
func NewStorage() Storage {
	return &memoryStorage{
		users: make(map[string]rest.User),
	}
}
