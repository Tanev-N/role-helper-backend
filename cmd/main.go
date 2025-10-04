package main

import (
	"role-helper/internal/delivery/http"
)

func main() {
	server := httpserver.NewHTTPServer()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
