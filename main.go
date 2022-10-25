package main

import (
	"fmt"

	"github.com/yagikota/blockchain_with_go/block"
	"github.com/yagikota/blockchain_with_go/wallet"
)

func main() {
	wallet := wallet.NewWallet()
	fmt.Println(wallet.PrivateKeyStr())
	fmt.Println(wallet.PublicKeyStr())
	fmt.Println(wallet.BlockchainAddress())

	blockChain := block.NewBlockchain(wallet.BlockchainAddress())

	blockChain.AddTransaction("A", "B", 1.0)
	blockChain.Mining() // blockが追加される
	// blockChain.Print()

	blockChain.AddTransaction("C", "D", 2.0)
	blockChain.AddTransaction("X", "Y", 3.0)
	blockChain.Mining()
	blockChain.Print()

	fmt.Printf("my %.1f\n", blockChain.CalculateTotalAmount("my_blockchain_address"))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount("C"))
	fmt.Printf("D %.1f\n", blockChain.CalculateTotalAmount("D"))
}
