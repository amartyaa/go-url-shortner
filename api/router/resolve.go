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
	param_chan := make(chan helpers.Params)
	go func() {
		param_chan <- helpers.ValidParams(string(url))
		close(param_chan)
	}()
	param := <-param_chan
	if !param.Check {
		fmt.Println(param.URL)
		fmt.Println("Invalid request")
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	fmt.Println(param.URL)
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
	fmt.Println(val)
	//Log the request
	client2 := db.Connect(1)
	defer client2.Close()
	client2.Incr(db.Dbctx, "total_requests")
	resp := response{
		Url:             val,
		CustomShort:     param.URL,
		Expiry:          client.TTL(db.Dbctx, param.URL).Val(),
		XRateRemaining:  0,
		XRateLimitReset: 0,
	}
	xrr, _ := client2.Decr(db.Dbctx, c.IP()).Result()
	resp.XRateRemaining = int(xrr)
	resp.XRateLimitReset = client2.TTL(db.Dbctx, c.IP()).Val()
	return c.Status(200).JSON(resp)
}
