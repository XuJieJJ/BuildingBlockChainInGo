package main

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"
const genesisCoinbaseData = "The Times 09/Dec/2023 Chancellor on brink of second bailout for banks"

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

func (bc *BlockChain) MineBlock(transaction []*Transaction) {
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
	newBlock := NewBlock(transaction, lastHash)
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

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewBlockchain(address string) *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing blockchain fount. Create one first.")
		os.Exit(1) //程序终止 返回特定退出码 0表示成功 非0表示出错
	}

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
		tip = b.Get([]byte("l"))

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

func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)

	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

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

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := BlockChain{tip, db}
	return &bc
}

// 遍历区块链中的所有区块和交易，查找与指定地址相关的未花费交易
func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspenrTXs []Transaction        //创建一个空的交易切片，用于存储找到的未花费交易。
	spentTXOs := make(map[string][]int) //创建一个空的交易切片，用于存储找到的未花费交易。 key--txID value--spentIndex
	bci := bc.Iterator()                //迭代器

	for { //无线循环
		block := bci.Next() //从迭代器中获得下一个区块

		for _, tx := range block.Transaction { //遍历当前区块所有交易
			txID := hex.EncodeToString(tx.ID) //将ID转换为十六进制字符串 以便在map中查找相关的已花费输出
		Outputs: //label标签，用于后续循环中进行标识
			for outIdx, out := range tx.Vout { //遍历当前交易的所有输出
				if spentTXOs[txID] != nil { //检查是否存在已花费的输出记录
					for _, spentOut := range spentTXOs[txID] { //遍历已经花费的输出索引
						if spentOut == outIdx { //如果当前输出索引等于已花费的输出索引
							continue Outputs
						}
					}
				} //如果当前输出未被花费，通过address进行解锁
				if out.CanBeUnlockedWith(address) {
					unspenrTXs = append(unspenrTXs, *tx) //将当前交易添加到未花费交易列表中
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin { //遍历当前交易所有输入
					if in.CanUnlockOutputWith(address) { //如果当前输入可以使用给定该地址解锁，则其标记为已花费
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBLockHash) == 0 { //已遍历到创世区块，跳出循环
			break
		}
	}
	return unspenrTXs
}

func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransaction := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransaction {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// 在创建新的输出之前，先找到所有的未花费出账并确认它们存了足够的币
func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address) //返回与指定地址相关的未花费交易列表
	accumulated := 0                                  //跟踪累计的金额

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount { //检查输出是否可以由给定的地址解锁，并且累积的金额小于需要的金额
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

//
