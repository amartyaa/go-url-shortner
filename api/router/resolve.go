package router

import (
	"fmt"

	"github.com/amartyaa/go-url-shortner/db"
	"github.com/amartyaa/go-url-shortner/helpers"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func Resolve(c *fiber.Ctx) error {
	fmt.Println("Resolve")
	url := c.Request().URI().QueryString()
	param := helpers.ValidParams(string(url))
	if !param.Check {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	client := db.Connect(0)
	defer client.Close()
	val, err := client.Get(db.Dbctx, param.URL).Result()
	if err == redis.Nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Shortened URL not found",
		})
	} else if err != nil {
		return c.Status(502).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	//Log the request
	client2 := db.Connect(1)
	defer client2.Close()
	client2.Incr(db.Dbctx, "total_requests")

	return c.Redirect(val, 301) //301 is the status code for permanent redirect
}
