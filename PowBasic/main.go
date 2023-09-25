package main

import "fmt"

func main() {
	bc := NewBlockchain()

	bc.AddBLock("Send 1 BTC to Alice")
	bc.AddBLock("Send 2 more BTC to Bob")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBLockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

	}
}
