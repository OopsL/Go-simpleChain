package blockChain

import (
	"github.com/boltdb/bolt"
	"log"
)

const blockchainDB = "blockchain.db"
const blockBucket = "blockBucket"
const lastBlockHashKey  = "lastBlockHash"

type BlockChain struct {
	//Blocks []*Block

	DB *bolt.DB
	tail []byte
}

func NewBlockChain() *BlockChain {


	//return &BlockChain{
	//	[]*Block{genesisBlock},
	//}

	db, err := bolt.Open(blockchainDB, 0600, nil)
	if err != nil {
		log.Panic("bolt open err")
	}

	var lastHash []byte

	//操作数据库
	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//创建bucket
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("create bucket err")
			}

			//存入初始化数据
			genesisBlock := GenesisBlock()
			bucket.Put(genesisBlock.Hash, genesisBlock.serialization())
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
			lastHash = genesisBlock.Hash

		}else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	return &BlockChain{db, lastHash}

}

//创世区块
func GenesisBlock() *Block  {
	return NewBlock("这是一个genesis block",[]byte{})
}

//添加区块
func (bc *BlockChain)AddBlock(data string)  {
	//lastBlock := bc.Blocks[len(bc.Blocks) - 1]
	//newBlock := NewBlock(data, lastBlock.Hash)
	//bc.Blocks = append(bc.Blocks, newBlock)

	db := bc.DB
	if db == nil {
		log.Panic("addblock db is nil")
	}

	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			log.Panic("addblock bucket is nil err")
		}

		block := NewBlock(data, bc.tail)
		bucket.Put(block.Hash, block.serialization())
		bucket.Put([]byte(lastBlockHashKey), block.Hash)
		bc.tail = block.Hash

		return nil
	})


}