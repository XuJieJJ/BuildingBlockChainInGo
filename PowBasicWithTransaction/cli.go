package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// 提供与程序交互的接口
type CLI struct {
	bc *BlockChain
}

func (cli *CLI) validateArgs() {
	//os.Args用于获取通过命令行传入的参数
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" addBlock -data BLOCK_DATA -add a block to the blockchain")
	fmt.Println(" printChain -print all the blocks of the blockchain")
}

/*
	func (cli *CLI) addBlock(data string) {
		cli.bc.AddBlock(data)
		fmt.Println("Success")
	}
*/
func (cli *CLI) createBLockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Done！")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBLockHash)
		fmt.Printf("Data: %s\n", block.Hash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBLockHash) == 0 {
			break

		}
	}
}

// 所有与命令行相关的操作都将由 CLI 结构处理：
func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send block reward to")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBLockchain(*createBlockchainAddress)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
