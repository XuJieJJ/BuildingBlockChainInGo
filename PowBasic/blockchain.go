package main

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBLock(data string) {
	preBLock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, preBLock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
