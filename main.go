package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/josepsoares/video-render-api/pkg/logger"
	"github.com/josepsoares/video-render-api/pkg/utils"
)

func init() {
	// initialize required log files
	logger.InitializeLogger("main")
	logger.InitializeLogger("render")
	logger.InitializeLogger("error")
}

func main() {
	// load .env file
	err := godotenv.Load()
	utils.FailOnError("couldn't read .env file", err)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
