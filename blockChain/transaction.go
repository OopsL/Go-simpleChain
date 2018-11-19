package blockChain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
)

const reward = 50.0

//交易
type Transaction struct {
	TXID      []byte     //交易ID
	TXInputs  []TXInput  //交易输入
	TXOutputs []TXOutput //交易输出
}

type TXInput struct {
	TXID []byte
	//引用output的索引值
	Index int64
	//解锁脚本
	//Sig string

	//签名
	Signature []byte

	//公钥
	PubKey []byte
}

type TXOutput struct {
	Value float64
	//锁定脚本
	//PubkeyHash string
	//接收人公钥的hash
	PubKeyHash []byte
}

//构建新的output
// 1. 当前存储的是address, 根据address生成接收人的pubkeyHash
func (output *TXOutput) Lock(address string) {

	output.PubKeyHash = GetPubKeyFromAddress(address)
}

func NewTXOutput(amount float64, address string) *TXOutput {
	output := TXOutput{Value: amount}
	output.Lock(address)

	return &output
}

//设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic("tx encoder err")
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//coinbase
func NewCoinbaseTX(address string, data string) *Transaction {

	//挖矿交易 不需要填真实的pubkey
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	//output := TXOutput{reward, address}
	output := NewTXOutput(reward, address)
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	tx.SetHash()

	return &tx

}

//判断是否是genesis block
func (tx *Transaction) IsCoinbase() bool {
	input := tx.TXInputs[0]
	if len(tx.TXInputs) == 1 && input.Index == -1 && bytes.Equal(input.TXID, []byte{}) {
		return true
	}

	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	//根据address获取钱包中的私钥和公钥
	ws := NewWallets()
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Println("获取钱包失败")
		return nil
	}

	privateKey := wallet.PrivateKey
	pubKey := wallet.PublicKey

	//需要传入pubkeyHash
	//utxos, resVal := bc.FindNeedUTXOs(from, amount)
	pubKeyHash := HashPubKey(pubKey)
	utxos, resVal := bc.FindNeedUTXOs(pubKeyHash, amount)

	if resVal < amount {
		fmt.Println("余额不足, 请充值~")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	//创建交易
	for keyID, indexArr := range utxos {
		for _, val := range indexArr {
			input := TXInput{[]byte(keyID), int64(val), nil, pubKey}
			inputs = append(inputs, input)
		}
	}

	//创建output
	//output := TXOutput{amount, to}
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)

	if resVal > amount {
		outputSel := NewTXOutput(resVal-amount, from)
		outputs = append(outputs, *outputSel)
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()

	//签名
	bc.SignTransaction(&tx, privateKey)

	return &tx
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	//遍历copy的tx
	for i, input := range txCopy.TXInputs {
		//根据input的TXID找到prevTXs中对应的tx
		prevTx := prevTXs[string(input.TXID)]

		if len(prevTx.TXID) == 0 {
			log.Panic("无效的交易")
		}

		//对input的Pubkey进行赋值
		txCopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash

		//对txCopy进行hash运算
		txCopy.SetHash()
		//还原, 以免影响后面的input的签名
		txCopy.TXInputs[i].PubKey = nil
		//签名是对交易的hash进行签名 本例中交易的TXID的值是交易的hash值
		signDataHash := txCopy.TXID
		//签名
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic("签名错误", err)
		}

		//给原tx的input的signature赋值
		tx.TXInputs[i].Signature = append(r.Bytes(), s.Bytes()...)

	}

}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXID, input.Index, nil, nil})
	}
	for _, output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}

	return Transaction{tx.TXID, inputs, outputs}
}

//校验Transaction
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {

	if tx.IsCoinbase() {
		return true
	}
	txCopy := tx.TrimmedCopy()

	for i, input := range tx.TXInputs {
		prevTx := prevTXs[string(input.TXID)]
		if len(prevTx.TXID) == 0 {
			log.Panic("引用的交易无效")
		}
		//本地先生成tx的数据hash, 再校验
		txCopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash

		txCopy.SetHash()
		//交易的hash
		dataHash := txCopy.TXID
		//根据交易提供的签名和公钥 进行验证
		signature := input.Signature
		pubKey := input.PubKey

		//根据signature, 分解出r, s
		r := big.Int{}
		s := big.Int{}

		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//根据pubkey分解出X Y后, 再生成公钥
		X := big.Int{}
		Y := big.Int{}

		X.SetBytes(pubKey[:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])
		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), &X, &Y}

		if !ecdsa.Verify(&pubKeyOrigin, dataHash, &r, &s) {
			return false
		}
	}
	return true
}
