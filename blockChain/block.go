package blockChain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"
	"encoding/gob"
	"log"
)

type Block struct {
	//1.版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkel根
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数
	Nonce uint64

	//当前区块哈希
	Hash []byte
	//数据
	Data []byte
}

func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}

	//block.SetHash()

	pow := NewProofOfWork(&block)
	blockHash, nonce := pow.RunPow()
	block.Hash = blockHash
	block.Nonce = nonce

	return &block

}

//uint64转[]byte
func Uint64ToBytes(num uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (block *Block) SetHash() {
	tmp := [][]byte{
		Uint64ToBytes(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToBytes(block.TimeStamp),
		Uint64ToBytes(block.Difficulty),
		Uint64ToBytes(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})

	blockHash := sha256.Sum256(blockInfo)
	block.Hash = blockHash[:]
}

//序列化
func (block *Block)serialization() []byte {

	//
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("block serialization err")
	}

	return buffer.Bytes()
}

func Deserialization(data []byte) Block  {

	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("data deserialization to block err")
	}

	return block

}
