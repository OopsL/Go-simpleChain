package main

import (
	"simpleChain/blockChain"
	"fmt"
)

func main()  {

	bc := blockChain.NewBlockChain()

	bc.AddBlock("这是第二个block")
	bc.AddBlock("这是第三个block")

	for index, block := range bc.Blocks{
		fmt.Println("当前区块", index)
		fmt.Printf("prevBlockHash: %x\n", block.PrevHash)
		fmt.Printf("currentBlockHash: %x\n", block.Hash)
		fmt.Printf("blockData: %s\n", block.Data)
	}
}

