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
	v1.Get("/health_check", healthCheck)
	v1.Get("/chain", getChainHandler)
	v1.Get("/transactions", getTransactions)
	v1.Post("/transactions", createTransactions)
	v1.Get("/mine", mine)
	v1.Get("/mine/start", startMine)
	v1.Get(fmt.Sprintf("/amount/:%s?", blockchainAddressParam), amount)

	return app
}
