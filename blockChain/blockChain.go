package blockChain

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		[]*Block{genesisBlock},
	}
}

//创世区块
func GenesisBlock() *Block  {
	return NewBlock("这是一个genesis block",[]byte{})
}

//添加区块
func (bc *BlockChain)AddBlock(data string)  {
	lastBlock := bc.Blocks[len(bc.Blocks) - 1]
	newBlock := NewBlock(data, lastBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}