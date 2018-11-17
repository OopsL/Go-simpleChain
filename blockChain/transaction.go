package blockChain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
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
	Sig string
}

type TXOutput struct {
	Value float64
	//锁定脚本
	PubkeyHash string
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

	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
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

	utxos, resVal := bc.FindNeedUTXOs(from, amount)

	if resVal < amount {
		fmt.Println("余额不足, 请充值~")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	//创建交易
	for keyID, indexArr := range utxos {
		for _, val := range indexArr {
			input := TXInput{[]byte(keyID), int64(val), from}
			inputs = append(inputs, input)
		}
	}

	//创建output
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	if resVal > amount {
		outputs = append(outputs, TXOutput{resVal - amount, from})
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()

	return &tx
}
