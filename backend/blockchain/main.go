package main

import (
	"log"

	"github.com/yagikota/blockchain_with_go/backend/blockchain/controller"
)

// https://docs.gofiber.io/api/app#group
func main() {
	app1 := controller.InitRouter()
	// go func() {
	// 	log.Fatal(app1.Listen(":8002"))
	// }()
	log.Fatal(app1.Listen(":8001"))
}
