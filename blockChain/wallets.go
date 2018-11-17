package blockChain

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
)

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.LoadFile()

	return &wallets
}

func (wallets *Wallets) CreateWallet() string {

	wallet := NewWallet()
	address := wallet.NewAddress()

	wallets.WalletsMap[address] = wallet
	wallets.SaveToFile()

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

	ioutil.WriteFile("wallets.dat", buffer.Bytes(), 0600)

}

func (ws *Wallets) LoadFile() {
	content, err := ioutil.ReadFile("wallets.dat")
	if err != nil {
		//log.Panic(err)
		return
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
