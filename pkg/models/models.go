package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials is used if a user tries to login with an incorrect email or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Gist struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

type GistModelInterface interface {
	Insert(string, string, string) (int, error)
	Get(int) (*Gist, error)
	Latest() ([]*Gist, error)
}

type UserModelInterface interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Get(int) (*User, error)
}
