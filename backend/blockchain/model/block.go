package model

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yagikota/blockchain_with_go/backend/common"
)

const (
	MINING_DIFFICULTY = 3 // difficulty of mining.
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	Timestamp    int64          `json:"timestamp"`
	Nonce        int            `json:"nonce"`
	PreviousHash string         `json:"previous_hash"`
	Transactions []*Transaction `json:"transactions"`
}

func (b *Block) Print() {
	fmt.Printf("timestamp       %d\n", b.Timestamp)
	fmt.Printf("nonce           %d\n", b.Nonce)
	fmt.Printf("previous_hash   %x\n", b.PreviousHash)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func NewBlock(nonce int, previousHash string, transactions []*Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
}

// TODO: hashの方法調べる
func (b *Block) Hash() string {
	m, _ := json.Marshal(b)
	h := sha256.Sum256(m)
	return string(h[:])
}

type Blockchain struct {
	transactionPool   []*Transaction `json:"-"`
	Chain             []*Block       `json:"chains"`
	BlockchainAddress string         `json:"-"` // Use bitcoin address as blockchainAddress.
	port              int            `json:"-"`
}

func NewBlockchain(blockchainAddress string, port int) *Blockchain {
	b := new(Block)
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	bc.BlockchainAddress = blockchainAddress
	bc.port = port
	return bc
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// TODO: function name maybe incorrect.
func (bc *Blockchain) CreateBlock(nonce int, previousHash string) {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.transactionPool = []*Transaction{}
}

func (bc *Blockchain) CreateTransaction(sender, recipient string, value float64, senderPublicKey *ecdsa.PublicKey, s *common.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)
	// TODO: Sync
	return isTransacted
}

// AddTransaction add a transaction to pool.
func (bc *Blockchain) AddTransaction(sender, recipient string, value float64, senderPublicKey *ecdsa.PublicKey, s *common.Signature) bool {
	t := NewTransaction(sender, recipient, value)
	// マイニングの場合、送り手はBCになる。
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: Not enough balance in a wallet")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	log.Println("ERROR: Verify Transaction")
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *common.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256(m)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionFromPool() []*Transaction {
	transaction := make([]*Transaction, 0, len(bc.transactionPool))
	for _, t := range bc.transactionPool {
		transaction = append(transaction, NewTransaction(
			t.SenderBlockchainAddress,
			t.RecipientBlockchainAddress,
			t.Value,
		))
	}
	return transaction
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

// マイニング競争に勝った者がブロックを生成するコンセンサスアルゴリズムの1種
// TODO: 時間かかる
func (bc *Blockchain) ValidProof(nonce int, previousHash string, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		Timestamp:    0,
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	guessHashStr := guessBlock.Hash()
	return guessHashStr[:difficulty] == zeros
}

// ProofOfWork finds nonce.
func (bc *Blockchain) ProofOfWork() int {
	// transactions := bc.CopyTransactionFromPool() // bc.transactionPoolじゃため？
	transactions := bc.transactionPool
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.BlockchainAddress, MINING_REWARD, nil, nil)
	previousHash := bc.LastBlock().Hash()
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float64 {
	totalAmount := 0.0
	for _, block := range bc.Chain {
		for _, t := range block.Transactions {
			value := t.Value
			if blockchainAddress == t.RecipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.SenderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

// block内のtransaction
type Transaction struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float64 `json:"value"`
}

func NewTransaction(sender, recipient string, value float64) *Transaction {
	return &Transaction{
		SenderBlockchainAddress:    sender,
		RecipientBlockchainAddress: recipient,
		Value:                      value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("senderBlockchainAddress    %s\n", t.SenderBlockchainAddress)
	fmt.Printf("recipientBlockchainAddress %s\n", t.RecipientBlockchainAddress)
	fmt.Printf("value                      %v\n", t.Value)
}

type BlockchainTransactionRequest struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	SenderPublicKey            string  `json:"sender_public_key"`
	Value                      float64 `json:"value"`
	Signature                  string  `json:"signature"`
}

// TODO: 長さのvalidate
func (t BlockchainTransactionRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.SenderBlockchainAddress, validation.Required, validation.Length(26, 35)),
		validation.Field(&t.RecipientBlockchainAddress, validation.Required, validation.Length(26, 35)),
		validation.Field(&t.SenderPublicKey, validation.Required, validation.Length(128, 128)),
		validation.Field(&t.Value, validation.Required),
		validation.Field(&t.Signature, validation.Required, validation.Length(128, 128)),
	)
}

type BlockchainTransactionResponse struct {
	Transactions []*Transaction `json:"transactions"`
	Length       int            `json:"length"`
}
