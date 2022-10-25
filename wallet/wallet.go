package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
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
	return fmt.Sprintf("%x%x", w.publicKey.X, w.publicKey.Y)
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}
