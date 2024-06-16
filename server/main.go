package main

import (
	"fmt"
	"mutebotx/server"
	"net/http"
	"os"
)

func main() {
	server := server.CreateServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started on port " + port)
	http.ListenAndServe(":"+port, server)
}
