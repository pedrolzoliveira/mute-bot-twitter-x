package main

import (
	"fmt"
	"mutebotx/server"
	"net/http"
)

func main() {
	server := server.CreateServer()
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", server)
}
