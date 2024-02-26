package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting the server...")
	server.SetupDatabase()
	fmt.Println("Database connected.")
	server.StartServer()
	fmt.Println("Server started.")
}
