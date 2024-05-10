package main

import (
	"log"
	"os"

	"noto/database"
	"noto/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	godotenv.Load()
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Expose-Headers", "Set-Cookie")
		return c.Next()
	})
	app.Static("/", "storage")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://noto.moe",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		ExposeHeaders:    "Set-Cookie",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	database.ConnectDB()

	router.Initalize(app)
	log.Fatal(app.Listen("127.0.0.1:80"))
}
