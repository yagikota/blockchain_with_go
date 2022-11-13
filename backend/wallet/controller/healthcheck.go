package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Wallet Server! You've requested: %s", c.Path()))
}
