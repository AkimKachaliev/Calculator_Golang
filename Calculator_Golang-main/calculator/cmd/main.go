package main

import (
	"fmt"

	"github.com/AkimKachaliev/Calculator_Golang/Calculator_Golang-main/calculator/server"
)

func main() {
	fmt.Println("Starting the server...")
	server.SetupDatabase()
	fmt.Println("Database connected.")
	server.StartServer()
	fmt.Println("Server started.")
}
