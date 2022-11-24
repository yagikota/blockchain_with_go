package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/yagikota/blockchain_with_go/backend/blockchain/controller"
)

// https://docs.gofiber.io/api/app#group
func main() {
	app1 := controller.InitRouter()
	port := flag.Int("port", 8001, "TCP Port Number of Blockchain Server")
	flag.Parse()
	fmt.Println(*port)
	log.Fatal(app1.Listen(net.JoinHostPort("localhost", strconv.Itoa(*port))))
}
