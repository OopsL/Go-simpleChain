package blockChain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	CurrentBlock *Block
	target       *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		CurrentBlock: block,
	}

	targetStr := "0000f00000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr, 16)
	pow.target = &tmpInt
	return &pow
}

func (pow *ProofOfWork) RunPow() ([]byte, uint64) {

	//拼接block中的数据
	var blockHash [32]byte
	var nonce uint64

	var tmp [][]byte
	var blockInfo []byte

	block := pow.CurrentBlock
	for {
		tmp = [][]byte{
			Uint64ToBytes(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToBytes(block.TimeStamp),
			Uint64ToBytes(block.Difficulty),
			Uint64ToBytes(nonce),
			block.Data,
		}
		blockInfo = bytes.Join(tmp, []byte{})

		blockHash = sha256.Sum256(blockInfo)

		tmpInt := big.Int{}
		tmpInt.SetBytes(blockHash[:])

		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash : %x, nonce : %d\n", blockHash, nonce)
			break
		} else {
			nonce++
		}

	}

	return blockHash[:], nonce
}
