package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yagikota/blockchain_with_go/backend/common"
	"github.com/yagikota/blockchain_with_go/backend/wallet/model"
)

const (
	gateWayURL = "http://localhost:8001/v1"
)

func createWallet(c *fiber.Ctx) error {
	myWallet := model.NewWallet()
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
	transaction := model.NewTransaction(privateKey, publicKey, t.SenderBlockchainAddress, t.RecipientBlockchainAddress, t.Value)
	signature := transaction.GenerateSignature()
	signatureStr := signature.String()

	// blockchain serverに投げる用
	bt := &model.BlockchainTransactionRequest{
		SenderBlockchainAddress:    t.SenderBlockchainAddress,
		RecipientBlockchainAddress: t.RecipientBlockchainAddress,
		SenderPublicKey:            t.SenderPublicKey,
		Value:                      t.Value,
		Signature:                  signatureStr,
	}
	btByte, _ := json.Marshal(bt)
	buf := bytes.NewBuffer(btByte)

	resp, err := http.Post(gateWayURL+"/transactions", "application/json", buf)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	if resp.StatusCode == fiber.StatusCreated {
		return c.SendStatus(fiber.StatusCreated)
	}
	return c.SendStatus(fiber.StatusInternalServerError)
}

func getAmount(c *fiber.Ctx) error {
	bcAddress := c.Query("blockchain_address")
	whereAddress := fmt.Sprintf("?blockchain_address=%s", bcAddress)
	bcResp, err := http.Get(gateWayURL + fmt.Sprintf("/amount%s", whereAddress))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	decorder := json.NewDecoder(bcResp.Body)
	var resp model.AmountResponse
	if err := decorder.Decode(&resp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewResponse(err.Error()))
	}
	return c.JSON(resp)
}
