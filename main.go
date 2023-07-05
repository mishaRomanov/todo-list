package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func About(c *fiber.Ctx) error {
	return c.SendString("This is About page!")
}

func main() {
	app := fiber.New()

	app.Get("/about", About)

	logrus.Fatal(app.Listen(":8080"))
}
