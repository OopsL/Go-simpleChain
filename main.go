package main

import "simpleChain/blockChain"

func main()  {

	bc := blockChain.NewBlockChain()

	cli := CLI{bc:bc}
	cli.Run()


	AAA()

	//bc.AddBlock("这是第二个block")
	//bc.AddBlock("这是第三个block")
	//
	//iter := bc.NewIterator()
	//
	//for  {
	//	block := iter.Next()
	//	fmt.Println("=========================")
	//	fmt.Printf("prevBlockHash: %x\n", block.PrevHash)
	//	fmt.Printf("currentBlockHash: %x\n", block.Hash)
	//	fmt.Printf("blockData: %s\n", block.Data)
	//	fmt.Println("=========================")
	//
	//	if block.PrevHash == nil {
	//		break
	//	}
	//}

}

