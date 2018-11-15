package blockChain

import "time"

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

	return &block
}
