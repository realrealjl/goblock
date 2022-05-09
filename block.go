package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block 由区块头和交易两部分构成
// Timestamp : 当前时间戳
// PrevBlockHash：前一个块的哈希
// Hash ：当前块的哈希
// Data：区块实际存储的信息，比特币中也就是交易
type Block struct {
	Timestamp int64
	PrevBlockHash []byte
	Hash []byte
	Data []byte
	Nonce int
}

// NewBlock 用于生成新块，参数需要Data与PrevBlockHash
// 当前块的哈希会给予Data和PrevBlockHash计算得到
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash: []byte{},
		Data: []byte(data),
		Nonce: 0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// SetHash 设置当前块哈希
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// NewGenesisBlock 生成创世块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}


