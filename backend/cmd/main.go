package main

import (
	"fmt"

	"github.com/yagikota/blockchain_with_go/backend/common"
)

func main() {
	fmt.Println(common.FindNeighbors("localhost", 8000, 0, 3, 8000, 8003))
}
