package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/wycliff-ochieng/database"
	"github.com/wycliff-ochieng/handlers"
	"github.com/wycliff-ochieng/service"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() {

	l := log.New(os.Stdout, "AUTHETICATION IN PROGRESS", log.LstdFlags)

	db, err := database.NewPostgres()
	if err != nil {
		fmt.Printf("failed to connect to db %v", err)
	}

	if err := db.Init(); err != nil {
		panic(err)
	}

	userService := service.NewUserService(db)

	uh := handlers.NewAuthHandle(l, &userService)

	router := mux.NewRouter()

	postRouter := router.Methods("POST").Subrouter()
	postRouter.HandleFunc("/register", uh.RegisterUser)
	http.ListenAndServe(s.addr, router)
}
