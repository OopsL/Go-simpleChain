package main

import (
	"fmt"
	"os"
)

const Usage = `
	addBlock --data DATA     "add data to blockchain"
	printChain               "print all blockchain data" 
	getBalance --address address
`

func (cli *CLI) Run() {

	args := os.Args
	if len(args) < 2 {
		fmt.Println("111" + Usage)
		return
	}

	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Println("add block")
		if len(args) == 4 && args[2] == "--data" {
			//TODO
			//data := args[3]
			//cli.AddBlock(data)
		} else {
			fmt.Printf("添加区块参数使用不当，请检查")
			fmt.Printf(Usage)
		}
	case "printChain":
		fmt.Println("打印区块")
		cli.PrintBlockChain()
	case "getBalance":
		address := args[3]
		cli.getBalance(address)

	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)

	}

}
