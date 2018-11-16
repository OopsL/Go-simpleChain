package blockChain

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	DB *bolt.DB
	currentHashPointer []byte
}

func (bc *BlockChain)NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		DB:                 bc.DB,
		currentHashPointer: bc.tail,
	}
}

func (iter *BlockChainIterator)Next() *Block {
	var block Block

	iter.DB.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			 log.Panic("next bucket is nil err")
		}

		blockBytes := bucket.Get(iter.currentHashPointer)
		block = Deserialization(blockBytes)
		iter.currentHashPointer = block.PrevHash

		return nil
	})


	return &block
}
