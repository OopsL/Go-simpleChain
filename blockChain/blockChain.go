package blockChain

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const blockchainDB = "blockchain.db"
const blockBucket = "blockBucket"
const lastBlockHashKey = "lastBlockHash"

type BlockChain struct {
	//Blocks []*Block

	DB   *bolt.DB
	tail []byte
}

func NewBlockChain(address string) *BlockChain {

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
			genesisBlock := GenesisBlock(address)
			bucket.Put(genesisBlock.Hash, genesisBlock.serialization())
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
			lastHash = genesisBlock.Hash

		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	return &BlockChain{db, lastHash}

}

//创世区块
func GenesisBlock(address string) *Block {

	coinbase := NewCoinbaseTX(address, "genesis block")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//lastBlock := bc.Blocks[len(bc.Blocks) - 1]
	//newBlock := NewBlock(data, lastBlock.Hash)
	//bc.Blocks = append(bc.Blocks, newBlock)

	db := bc.DB
	if db == nil {
		log.Panic("addblock db is nil")
	}

	db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("addblock bucket is nil err")
		}

		block := NewBlock(txs, bc.tail)
		bucket.Put(block.Hash, block.serialization())
		bucket.Put([]byte(lastBlockHashKey), block.Hash)
		bc.tail = block.Hash

		return nil
	})

}

//查询未被消耗的所有的utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {

	var UTXO []TXOutput

	txs := bc.FindUTXOTransantions(address)

	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			//将未消耗的output添加进UTXO
			if output.PubkeyHash == address {
				UTXO = append(UTXO, output)
			}
		}
	}

	return UTXO
}

//查找未被消耗的并且满足交易额度的utxo

func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]int64, float64) {

	utxos := make(map[string][]int64)
	var calc float64

	txs := bc.FindUTXOTransantions(from)

	for _, tx := range txs {
		for i, output := range tx.TXOutputs {

			//将未消耗的output添加进UTXO
			if output.PubkeyHash == from {

				if calc < amount {
					//添加utxo
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], int64(i))
					calc += output.Value

					if calc >= amount {
						return utxos, calc
					}
				}

			}
		}
	}

	return utxos, calc

}

//提取代码
func (bc *BlockChain) FindUTXOTransantions(address string) []*Transaction {

	spentOutput := make(map[string][]int64)
	var txs []*Transaction

	//1. 遍历block
	iter := bc.NewIterator()
	for {

		block := iter.Next()
		//2. 遍历交易
		for _, tx := range block.Transactions {

			//3. 遍历output
		OUTPUT:
			for i, output := range tx.TXOutputs {

				//判断output是否被消耗
				if spentOutput[string(tx.TXID)] != nil {
					for _, value := range spentOutput[string(tx.TXID)] {
						if value == int64(i) {
							continue OUTPUT
						}
					}
				}

				//将未消耗的output添加进UTXO
				if output.PubkeyHash == address {
					//UTXO = append(UTXO, output)

					txs = append(txs, tx)

				}

			}

			//4. 遍历input
			//如果当前交易是挖矿交易,那么不做遍历
			if !tx.IsCoinbase() {
				for _, input := range tx.TXInputs {
					//先取出之前的值
					//spentIndexArr := spentOutput[string(input.TXID)]
					//将当前值存入
					if input.Sig == address {
						spentOutput[string(input.TXID)] = append(spentOutput[string(input.TXID)], input.Index)

					}
				}
			} else {
				//fmt.Println("这是coinbase")
			}
		}

		if block.PrevHash == nil {
			fmt.Println("FindUTXOs 区块遍历完成")
			break
		}
	}

	return txs
}
