package main

import (
	"log"
	"os"

	"github.com/amartyaa/go-url-shortner/router"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Routes(app *fiber.App) {
	app.Get("/shorten", router.Resolve)
	app.Post("/shorten", router.Shorten)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	Routes(app)
	app.Listen(os.Getenv("APP_PORT"))
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
