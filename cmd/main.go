package main

import (
	"fmt"

	"github.com/wycliff-ochieng/cmd/api"
)

func main() {
	fmt.Println("Starting server now")

	server := api.NewAPIServer(":8000")
	server.Run()
}
