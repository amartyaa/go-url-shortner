package router

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/amartyaa/go-url-shortner/db"
	"github.com/amartyaa/go-url-shortner/helpers"
	"github.com/gofiber/fiber/v2"

	"github.com/asaskevich/govalidator"
	// "github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
)

type request struct {
	Url         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	Url             string        `json:"url"`
	CustomShort     string        `json:"customShort"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"x-rate-remaining"`
	XRateLimitReset time.Duration `json:"x-rate-limit-reset"`
}

func Shorten(c *fiber.Ctx) error {
	fmt.Println("Shorten")
	var req request
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(req)
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid Request bole na",
		})
	}
	//IMPLEMENT RATE LIMITING

	r2 := db.Connect(1)
	defer r2.Close()
	val, err := r2.Get(db.Dbctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(db.Dbctx, c.IP(), os.Getenv("QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			ttl := r2.TTL(db.Dbctx, c.IP()).Val()
			return c.Status(429).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"retry-after": ttl,
			})
		}
	}

	//IMPLEMENT URL VALIDATION
	if !govalidator.IsURL(req.Url) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid URL goValidator",
		})
	}
	if !helpers.LoopDomain(req.Url) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid URL LoopDomain",
		})
	}
	//IMPLEMENT HTTPS/SSLv3/SSLv2 CHECK

	req.Url = helpers.EnforeceHTTPS(req.Url)

	//IMPLEMENT CUSTOM SHORT URL

	var id string
	if req.CustomShort != "" {
		id = req.CustomShort
	} else {
		id = helpers.GenerateID()
	}

	r1 := db.Connect(0)
	defer r1.Close()

	val, err = r1.Get(db.Dbctx, id).Result()
	if err == redis.Nil {
		_ = r1.Set(db.Dbctx, id, req.Url, req.Expiry*time.Second).Err()
	} else if err != nil {
		return c.Status(502).JSON(fiber.Map{
			"error": "Internal server error",
		})
	} else if val != "" {
		return c.Status(409).JSON(fiber.Map{
			"error": "Shortened URL already exists",
			"url":   val,
		})
	}

	err = r1.Set(db.Dbctx, id, req.Url, req.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(502).JSON(fiber.Map{
			"error": "Unable to shorten URL",
			"url":   req.Url,
			"id":    id,
		})
	}

	resp := response{
		Url:             req.Url,
		CustomShort:     os.Getenv("APP_HOST") + "/" + id,
		Expiry:          req.Expiry,
		XRateRemaining:  30,
		XRateLimitReset: 30 * 60 * time.Second,
	}

	xrr, _ := r2.Decr(db.Dbctx, c.IP()).Result()
	resp.XRateRemaining = int(xrr)
	resp.XRateLimitReset = r2.TTL(db.Dbctx, c.IP()).Val()

	return c.Status(200).JSON(resp)
}
