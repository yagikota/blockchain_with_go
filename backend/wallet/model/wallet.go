package model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yagikota/blockchain_with_go/backend/common"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
func NewWallet() *Wallet {
	// 1. create ECDSA private key (32 bytes), public key (64 bytes).
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil
	}
	w := &Wallet{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}

	// create blockchainAddress from publicKey
	// 2.Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil) // digest3(20bytes)
	// 4. Add network ID byte in front of RIPEMD-160 hash (0x00 for Main Network).
	digest4 := make([]byte, len(digest3)+1)
	digest4[0] = 0x00
	copy(digest4[1:], digest3)
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(digest4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	checkSum := digest6[:4]
	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	digestCheckSum := make([]byte, len(digest4)+len(checkSum))
	copy(digestCheckSum[:len(digest4)], digest4)
	copy(digestCheckSum[len(digest4):], checkSum)
	// 9. Convert the result from a byte string into base58.
	address := base58.Encode(digestCheckSum)
	w.blockchainAddress = address
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X, w.publicKey.Y)
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.blockchainAddress,
	})
}

// 署名用transaction
// walletで生成したprivateKey, publicKey, BlockchainAddressなどの情報を使用する
// https://dev.classmethod.jp/articles/blockchain-basic/
type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey `json:"-"`
	senderPublicKey            *ecdsa.PublicKey  `json:"-"`
	SenderBlockchainAddress    string            `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string            `json:"recipient_blockchain_address"`
	Value                      float64           `json:"value"`
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value float64) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (t *Transaction) GenerateSignature() *common.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256(m)
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &common.Signature{
		R: r,
		S: s,
	}
}

// validater: https://zenn.dev/mattn/articles/893f28eff96129
type TransactionRequest struct {
	SenderPrivateKey           string  `json:"sender_private_key"`
	SenderPublicKey            string  `json:"sender_public_key"`
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float64 `json:"value"`
}

func (t TransactionRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.SenderPrivateKey, validation.Required, validation.Length(64, 64)), // 32 bytes(256bits): 16進数1桁が2進数4桁→4ビットで表現できる。 16進数で64文字,
		validation.Field(&t.SenderPublicKey, validation.Required, validation.Length(128, 128)),
		validation.Field(&t.SenderBlockchainAddress, validation.Required, validation.Length(26, 35)),
		validation.Field(&t.RecipientBlockchainAddress, validation.Required, validation.Length(26, 35)),
		validation.Field(&t.Value, validation.Required, validation.Min(1.0)), // TODO: validation効いてない
	)
}

type BlockchainTransactionRequest struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	SenderPublicKey            string  `json:"sender_public_key"`
	Value                      float64 `json:"value"`
	Signature                  string  `json:"signature"`
}

func (t BlockchainTransactionRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.SenderBlockchainAddress, validation.Required, validation.Length(26, 35)),
		validation.Field(&t.RecipientBlockchainAddress, validation.Required, validation.Length(26, 35)),
		validation.Field(&t.SenderPublicKey, validation.Required, validation.Length(128, 128)),
		validation.Field(&t.Value, validation.Required),
		validation.Field(&t.Signature, validation.Required),
	)
}

type AmountResponse struct {
	Amount float64 `json:"amount"`
}
