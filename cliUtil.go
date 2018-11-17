package main

import (
	"fmt"
	"os"
	"strconv"
)

const Usage = `
	printChain               "print all blockchain data" 
	getBalance --address address
	send from to amount miner data
	newWallet "创建一个钱包实例"
	"listAllAddress" "遍历所有地址"
`

func (cli *CLI) Run() {

	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}

	cmd := args[1]
	switch cmd {
	case "printChain":
		fmt.Println("打印区块")
		cli.PrintBlockChain()
	case "getBalance":
		address := args[3]
		cli.getBalance(address)
	case "send":
		if len(args) != 7 {
			fmt.Println("传入的参数有误")
			return
		}
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		cli.SendTraction(from, to, amount, miner, data)

	case "newWallet":
		cli.NewWallet()

	case "listAllAddress":
		cli.ListAllAddress()

	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usage)

	}

}
