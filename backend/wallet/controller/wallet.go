package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	block "github.com/yagikota/blockchain_with_go/backend/blockchain/model"
	"github.com/yagikota/blockchain_with_go/backend/common"
	"github.com/yagikota/blockchain_with_go/backend/wallet/model"
	"github.com/yagikota/blockchain_with_go/wallet"
)

func createWallet(c *fiber.Ctx) error {
	myWallet := wallet.NewWallet()
	return c.JSON(myWallet)
}

func createTransaction(c *fiber.Ctx) error {
	var t model.TransactionRequest
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if err := t.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	publicKey := common.PublicKeyFromString(t.SenderPublicKey)
	privateKey := common.PrivateKeyFromString(t.SenderPrivateKey, publicKey)
	transaction := wallet.NewTransaction(privateKey, publicKey, t.SenderBlockchainAddress, t.RecipientBlockchainAddress, t.Value)
	signature := transaction.GenerateSignature()
	signatureStr := signature.String()

	// blockchain serverに投げる用
	bt := &block.TransactionRequest{
		SenderBlockchainAddress:    t.SenderBlockchainAddress,
		RecipientBlockchainAddress: t.RecipientBlockchainAddress,
		SenderPublicKey:            t.SenderPublicKey,
		Value:                      t.Value,
		Signature:                  signatureStr,
	}
	btByte, _ := json.Marshal(bt)
	buf := bytes.NewBuffer(btByte)
	url := "localhost:8001/transactions"
	resp, _ := http.Post(url, "application/json", buf)
	if resp.StatusCode == fiber.StatusCreated {
		return c.SendStatus(fiber.StatusCreated)
	}
	return c.SendStatus(fiber.StatusInternalServerError)
}
