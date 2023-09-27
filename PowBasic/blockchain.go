package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

// 区块链迭代器
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

/*func (bc *BlockChain) AddBLock(data string) {
	preBLock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, preBLock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}*/

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte
	//这是 BoltDB 事务的另一种（只读）类型。在这里，我们从数据库中获取最后一个区块的哈希值，并用它来挖掘一个新的区块哈希值
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lastHash)
	//挖掘出新区块后，我们会将其序列化表示保存到数据库中，并更新 l key，现在它存储了新区块的哈希值。
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})

}

func NewBlockchain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}
	/**
	函数的核心部分。在这里，我们获取存储区块的 "桶"：
	如果它存在，我们就从中读取 l 密钥；
	如果它不存在，我们就生成创世区块，创建 "桶"，
	将区块保存到其中，并更新存储链上最后一个区块哈希值的 l 密钥
	*/
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			fmt.Println("No exitsting blockchian found .Creationg a new one ...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil

	})
	if err != nil {
		log.Panic(err)
	}
	bc := BlockChain{tip, db}

	return &bc
}

//BoltDB 允许遍历一个桶中的所有密钥，但密钥是按字节排序存储的，而我们希望按区块链中的顺序打印区块。

func (bc *BlockChain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// BlockchainIterator 只做一件事：从区块链中返回下一个区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBLockHash
	return block
}
