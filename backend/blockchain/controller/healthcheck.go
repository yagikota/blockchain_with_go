package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Hello! You've requested: %s", c.Path()))
}
