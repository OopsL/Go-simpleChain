package main

import (
	"simpleChain/blockChain"

	"fmt"
)

type CLI struct {
	bc *blockChain.BlockChain
}

func (cli *CLI) AddBlock(txs []*blockChain.Transaction) {
	cli.bc.AddBlock(txs)
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	iter := bc.NewIterator()

	for {
		//返回区块，左移
		block := iter.Next()

		fmt.Printf("===========================\n\n")
		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
		fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %d\n", block.TimeStamp)
		fmt.Printf("难度值: %d\n", block.Difficulty)
		fmt.Printf("随机数: %d\n", block.Nonce)
		fmt.Printf("当前区块哈希值: %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].Sig)

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束！")
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	utxos := cli.bc.FindUTXOs(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Println(address, " 的余额为: ", total)
}
