package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	Firstname string
	Lastname  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(firstname, lastname, email, password string) (*User, error) {

	harshedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  string(harshedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
