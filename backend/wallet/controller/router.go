package controller

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter() *fiber.App {
	app := fiber.New()
	v1 := app.Group("/v1")
	// v1/health_check
	v1.Get("/health_check", healthCheck)
	// v1/wallet
	v1.Post("/wallet", createWallet)
	v1.Post("/transaction", createTransaction)

	return app
}
