package main

import (
	"fmt"
	"log"

	database "example.com/connectivity"
	routes "example.com/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
)

func main() {

	DbUser := database.GetConnectionString()

	// Connect to database
	if err := database.GetConnection(DbUser); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Use(helmet.New())

	fmt.Println("Server started")

	// API Routes
	routes.Routes(app)

	if err := app.Listen(":4401"); err != nil {
		log.Fatalln(err)
	}
}
