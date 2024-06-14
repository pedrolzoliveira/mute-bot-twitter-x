package main

import (
	"mutebotx/server"
	"net/http"
)

func main() {
	server := server.CreateServer()
	http.ListenAndServe(":8080", server)
}
