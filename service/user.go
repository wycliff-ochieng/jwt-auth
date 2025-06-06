package service

import (
	"errors"
	"fmt"

	"github.com/wycliff-ochieng/database"
	model "github.com/wycliff-ochieng/models"
)

type UserService struct {
	db database.DBInterface
}

var (
	ErrEmailExists = errors.New("email already exists")
)

func NewUserService(db database.DBInterface) UserService {
	return UserService{db: db}
}

func (u UserService) Register(firstname string, lastname string, email string, password string) (*model.UserResponse, error) {
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
	user, err := model.NewUser(firstname, lastname, email, password)
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

func (u *UserService) Login(email string, password string) (*model.LoginResqponse, error) {
	u.db.QueryRow("SELECT firstname ")
}
