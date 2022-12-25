package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amartyaa/go-url-shortner/dns"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Routes(app *fiber.App) {
	app.Get("/shorten", dns.Resolve)
	app.Post("/shorten", dns.Shorten)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Hello, playground")
	app := fiber.New()
	Routes(app)
	app.Listen(os.Getenv("APP_PORT"))
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
