package main

func main() {
	bc := NewBlockChain()

	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()

	//bc.AddBlock("Send 1 BTC to Ivan")
	//bc.AddBlock("Send 2 BTC to Ivan")
	//
	//for _, block := range bc.blocks {
	//	fmt.Printf("Prev hash:%x\n", block.PrevBlockHash)
	//	fmt.Printf("Data:%x\n", block.Data)
	//	fmt.Printf("Hash:%x\n", block.Hash)
	//	pow := NewProofOfWork(block)
	//	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//}
}
