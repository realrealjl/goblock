package main

// BlockChain BlockChain是一个Block指针数组
type BlockChain struct {
	blocks []*Block
}

// NewBlockChain 创建一个有创始块的链
func NewBlockChain() *BlockChain{
	return &BlockChain{
		blocks: []*Block{NewGenesisBlock()} ,
	}
}


// AddBlock 向链中加入一个新块
// data在实际中就是交易
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks) - 1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}