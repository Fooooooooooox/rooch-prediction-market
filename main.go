package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rooch-prediction-market/backend/config"
	"github.com/rooch-prediction-market/backend/server"
)

func getPort() string {
	c := config.New("", "")
	port := c.Port

	if port == "" {
		port = "50051"
	}

	return port
}

func main() {
	m := server.CreateServer()

	server.Initialize(":"+getPort(), m) // Assume `yourHandler()` sets up your routes
	fmt.Println("Server is running at http://localhost:" + getPort())
	server.StartServer()

	// Handling graceful shutdown in a goroutine
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	server.StopServer()
}
