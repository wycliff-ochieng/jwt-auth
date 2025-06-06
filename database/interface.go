package database

import (
	"database/sql"

	model "github.com/wycliff-ochieng/models"
)

type DBInterface interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}

type Storage interface {
	register(email, fistaname, lastname, password string) (*model.UserResponse, error)
}
