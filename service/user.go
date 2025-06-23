package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wycliff-ochieng/database"

	"github.com/wycliff-ochieng/internal"
	model "github.com/wycliff-ochieng/models"
)

type UserService struct {
	db database.DBInterface
}

var (
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidPassword = errors.New("incorrect password")
	ErrUserNotFound    = errors.New("user not found")
)

func NewUserService(db database.DBInterface) UserService {
	return UserService{db: db}
}

func (u UserService) Register(id int64, firstname string, lastname string, email string, password string) (*model.UserResponse, error) {
	//check if email exists
	var exists bool

	err := u.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%v email already exists", email)
	}

	//creating user
	user, err := model.NewUser(id, firstname, lastname, email, password)
	if err != nil {
		return nil, err
	}

	//insert into db
	query := `INSERT INTO users(firstname,lastname,email,password,createdat,updatedat) values($1,$2,$3,$4,$5,$6)`

	_, err = u.db.Exec(query, user.Firstname, user.Lastname, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &model.UserResponse{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *UserService) Login(email string, password string) (*internal.TokenPair, *model.UserResponse, error) {
	var user model.User

	query := `SELECT firstname,lastname,email,password,createdat,updatedat FROM users where email = $1`

	err := u.db.QueryRow(query, email).Scan(
		//&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch user: %v", err)
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, nil, ErrInvalidPassword
	}
	//return nil, nil, nil

	//var jwtSecret := []byte("mydogsnameisrufus")
	//var refreshSecret := []byte("myotherdogiscalledbuckeye")

	//generate tokens
	tokens, err := internal.GenerateTokenPair(
		user.ID,
		user.Email,
		string(internal.JwtSecret),
		string(internal.RefreshSecret),
		time.Hour*24,
		time.Hour*24*7,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating tokens: %v", err)
	}

	return tokens, &model.UserResponse{
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		CreatedAt: user.CreatedAt,
	}, nil

}
