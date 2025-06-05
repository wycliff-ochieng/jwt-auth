package main

import "github.com/wycliff-ochieng/cmd/api"

func main() {

	server := api.NewAPIServer(":8000")
	server.Run()
}
