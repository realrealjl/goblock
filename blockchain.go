package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

/*
	1. 打开一个数据库文件
    2。检查文件里面是否已经存储了一个区块链
    3。如果已经存储了一个区块链：
       创建一个新的Blockchain实例
       设置Blockchain实例的tip为数据库中存储的最后一个块的哈希
    4。如果没有区块链
       创建创始块
       存储到数据库
       将创世块哈希保存为最后一个块的哈希
       创建一个新的Blockchain实例，初始时tip指向创世块
*/

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// tip 指的是存储最后一个块的哈希，在链的末端可能出现暂时分叉的情况，db存储数据库的链接
type BlockChain struct {
	tip []byte
	db *bolt.DB
}

// NewBlockChain 创建一个有创始块的链
func NewBlockChain() *BlockChain{

	var tip []byte
	// 打开一个BoltDB文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 如果数据库中不存在区块链就创建一个，否则就直接读取最后一个哈希块
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		// 如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的哈希
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")

			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
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


// AddBlock 向链中加入一个新块
// data在实际中就是交易
func (bc *BlockChain) AddBlock(data string) {
	// 加入区块时 需要将区块持久化到数据库
	var lastHash []byte

	// 首先获取最后一个块的哈希用于生成新块的哈希
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
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

type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// 往回倒
	i.currentHash = block.PrevBlockHash

	return block
}