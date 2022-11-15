package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	blockchainAddressParam = "blockchain_address"
)

func InitRouter() *fiber.App {
	app := fiber.New()
	v1 := app.Group("/v1")
	// v1/health_check
	v1.Get("/health_check", healthCheck)
	// v1/wallet
	v1.Post("/wallet", createWallet)
	v1.Get(fmt.Sprintf("/wallet/amount/:%s?", blockchainAddressParam), getAmount)
	v1.Post("/transactions", createTransaction)

	return app
}
