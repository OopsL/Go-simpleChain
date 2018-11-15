package main

import (
			"fmt"
	"simpleChain/blockChain"
)

func main()  {
	block := blockChain.NewBlock("test block", []byte{})

	fmt.Println(block)
}

