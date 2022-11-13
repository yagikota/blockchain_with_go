package main

import (
	"log"

	"github.com/yagikota/blockchain_with_go/backend/wallet/controller"
)

// https://docs.gofiber.io/api/app#group
func main() {
	app1 := controller.InitRouter()
	log.Fatal(app1.Listen(":8000"))
}
