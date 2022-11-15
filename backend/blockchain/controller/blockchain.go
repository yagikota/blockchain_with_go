package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yagikota/blockchain_with_go/backend/blockchain/model"
	"github.com/yagikota/blockchain_with_go/backend/common"
)

type key string

const (
	cacheKey key = "blockchain"
)

var cache map[key]*model.Blockchain = make(map[key]*model.Blockchain)

func getBlockchain(c *fiber.Ctx) *model.Blockchain {
	bc, ok := cache[cacheKey]
	if !ok {
		minersWallet := model.NewWallet()
		port, _ := strconv.Atoi(c.Port())
		bc = model.NewBlockchain(minersWallet.BlockchainAddress(), port)
		cache[cacheKey] = bc
		log.Printf("private_key %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

func getChainHandler(c *fiber.Ctx) error {
	bc := getBlockchain(c)
	return c.JSON(bc)
}

func getTransactions(c *fiber.Ctx) error {
	bc := getBlockchain(c)
	transactions := bc.TransactionPool()
	return c.JSON(model.BlockchainTransactionResponse{
		Transactions: transactions,
		Length:       len(transactions),
	})
}

func createTransactions(c *fiber.Ctx) error {
	var t model.BlockchainTransactionRequest
	if err := c.BodyParser(&t); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	publicKey := common.PublicKeyFromString(t.SenderPublicKey)
	signature := common.SignatureFromString(t.Signature)
	bc := getBlockchain(c)
	isCreated := bc.CreateTransaction(t.SenderBlockchainAddress, t.RecipientBlockchainAddress, t.Value, publicKey, signature)
	if !isCreated {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func mine(c *fiber.Ctx) error {
	bc := getBlockchain(c)
	isMined := bc.Mining()
	if !isMined {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewResponse("couldn't mine"))
	}
	return c.Status(fiber.StatusOK).JSON(common.NewResponse("a mining success"))
}

func startMine(c *fiber.Ctx) error {
	bc := getBlockchain(c)
	bc.StartMining()
	return c.Status(fiber.StatusInternalServerError).JSON(common.NewResponse("auto mining start"))
}
