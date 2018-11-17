package blockChain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
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

//生成地址
func (wallet *Wallet) NewAddress() string {

	pubKey := wallet.PublicKey
	// sha356
	hash := sha256.Sum256(pubKey)
	//ripemd160
	rip160hash := ripemd160.New()
	_, err := rip160hash.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}

	rip160Value := rip160hash.Sum(nil)

	version := byte(00)
	//拼接版本号
	payload := append([]byte{version}, rip160Value...)
	//对拼接结果进行两次sha256
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])

	//拼接hash2的前4个byte
	checkCode := hash2[:4]
	payload = append(payload, checkCode...)

	//base58
	address := base58.Encode(payload)

	return address
}
