package model

import (
	"time"

	//"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Firstname string
	Lastname  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id int64, firstname, lastname, email, password string) (*User, error) {

	harshedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        id,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  string(harshedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type UserResponse struct {
	Firstname string
	Lastname  string
	Email     string
	CreatedAt time.Time
}
