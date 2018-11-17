package blockChain

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"log"
	"os"
)

const walletsFileName = "wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.LoadFile()

	return &wallets
}

func (ws *Wallets) CreateWallet() string {

	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet
	ws.SaveToFile()

	return address
}

func (ws *Wallets) SaveToFile() {

	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile(walletsFileName, buffer.Bytes(), 0600)

}

func (ws *Wallets) LoadFile() {

	_, err := os.Stat(walletsFileName)
	if os.IsNotExist(err) {
		//err:  stat wallets.dat: no such file or directory
		//fmt.Println(walletsFileName, " not exists err : ", err )
		return
	}

	content, err := ioutil.ReadFile(walletsFileName)
	if err != nil {
		log.Panic(err)
		//return
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wsLocal Wallets
	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}

	ws.WalletsMap = wsLocal.WalletsMap
}

//获取所有地址
func (ws *Wallets) ListAllAddress() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}

func GetPubKeyFromAddress(address string) []byte {
	base58De := base58.Decode(address)
	//base58De由 1byte version + 20byte pubkeyhash + 4byte checkcode组成
	//截取pubkeyhash
	fmt.Println("len = ", len(base58De), "address = ", address)
	pubkeyHash := base58De[1 : len(base58De)-4]
	return pubkeyHash
}
