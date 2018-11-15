package main

import (
	"simpleChain/blockChain"
	"time"
	"fmt"
)

func main()  {
	block := blockChain.Block{
		Version: 1,
		TimeStamp: uint64(time.Now().Unix()),
	}

	fmt.Println(block)
}

