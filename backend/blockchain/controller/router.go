package controller

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	app := fiber.New()
	v1 := app.Group("/v1")
	v1.Get("/health_check", healthCheck)
	v1.Get("/chain", getChainHandler)
	// TODO: implement
	v1.Get("/transactions", getTransactions)
	v1.Post("/transactions", createTransactions)
	v1.Get("/mine", mine)

	return app
}
