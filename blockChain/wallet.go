package blockChain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {

	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKeyOri := privateKey.PublicKey
	pubKey := append(pubKeyOri.X.Bytes(), pubKeyOri.Y.Bytes()...)

	return &Wallet{privateKey, pubKey}

}
