package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBLockHash []byte
	Hash          []byte
	Nonce         int
}

//实现 Block 的序列化

func (b *Block) Serialize() []byte {
	var result bytes.Buffer            //声明缓冲区
	encoder := gob.NewEncoder(&result) //初始化gob编码器

	err := encoder.Encode(b) //对数据块编码
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 实现反序列化 接受字节数组的输入，返回数据块
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBLockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]

}

func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), preBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
