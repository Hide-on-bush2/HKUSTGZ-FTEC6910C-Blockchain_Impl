package test

import (
	"Chain/pkg/blockchain"
	"Chain/pkg/utils"
	"os"
	"testing"
)

func TestNewChain(t *testing.T) {
	defer cleanUp()
	bc := blockchain.New(blockchain.DefaultConfig())
	if bc.Length != 1 {
		t.Errorf("Expected chain length: %v\n Actual chain length: %v", 1, bc.Length)
	}
	if _, err := os.Stat("./blockinfodata"); os.IsNotExist(err) {
		t.Errorf("Did not create leveldb blockinfodata")
	}
	if _, err := os.Stat("./coindata"); os.IsNotExist(err) {
		t.Errorf("Did not create leveldb coindata")
	}
	if len(bc.UnsafeHashes) != 0 {
		t.Errorf("unsafe hashes not initialized properly")
	}
	if bc.LastBlock == nil {
		t.Errorf("Did not initialize last block")
	}
}

func TestHandleAppendingBlock(t *testing.T) {
	defer cleanUp()
	bc := blockchain.New(blockchain.DefaultConfig())
	lastBlock := bc.LastBlock
	newBlock := MakeBlockFromPrev(lastBlock)
	bc.HandleBlock(newBlock)
	if bc.Length != 2 {
		t.Errorf("Expected chain length: %v\n Actual chain length: %v", 2, bc.Length)
	}
	if bc.LastHash != newBlock.Hash() {
		t.Errorf("Expected last hash: %v\nActual last hash: %v", newBlock.Hash(), bc.LastHash)
	}
	if bc.LastBlock != newBlock {
		t.Errorf("Expected block: %v\n Actual block: %v", newBlock, bc.LastBlock)
	}
}

func TestHandleForkingBlock(t *testing.T) {
	defer cleanUp()
	bc := blockchain.New(blockchain.DefaultConfig())
	currBlock := bc.LastBlock
	currForkingBlock := bc.LastBlock

	for i := 0; i < 4; i++ {
		newBlock := MakeBlockFromPrev(currBlock)
		newForkingBlock := MakeBlockFromPrev(currForkingBlock)
		newForkingBlock.Header.Version = 2
		if newBlock.Hash() == newForkingBlock.Hash() {
			t.Errorf("Hashes should not be the same")
		}
		if i < 3 {
			bc.HandleBlock(newBlock)
		}
		bc.HandleBlock(newForkingBlock)
		currBlock = newBlock
		currForkingBlock = newForkingBlock
	}
	if bc.Length != 5 {
		t.Errorf("Expected chain length: %v\n Actual chain length: %v", 5, bc.Length)
	}
	if bc.LastHash != currForkingBlock.Hash() {
		t.Errorf("Expected last hash: %v\nActual last hash: %v", currForkingBlock.Hash(), bc.LastHash)
	}
	if bc.LastBlock != currForkingBlock {
		t.Errorf("Expected block: %v\n Actual block: %v", currForkingBlock, bc.LastBlock)
	}
}

func TestHandleInvalidBlock(t *testing.T) {
	defer cleanUp()
	bc := blockchain.New(blockchain.DefaultConfig())
	lastBlock := bc.LastBlock
	block1 := MakeBlockFromPrev(lastBlock)
	block2 := MakeBlockFromPrev(block1)
	bc.HandleBlock(block2)
	if bc.Length != 1 {
		t.Errorf("Expected chain length: %v\n Actual chain length: %v", 1, bc.Length)
	}
	if bc.LastHash != lastBlock.Hash() {
		t.Errorf("Expected last hash: %v\nActual last hash: %v", lastBlock.Hash(), bc.LastHash)
	}
	if bc.LastBlock != lastBlock {
		t.Errorf("Expected block: %v\n Actual block: %v", lastBlock, bc.LastBlock)
	}
}

func TestHandle50Blocks(t *testing.T) {
	defer cleanUp()
	bc := blockchain.New(blockchain.DefaultConfig())
	currBlock := bc.LastBlock
	for i := 0; i < 50; i++ {
		newBlock := MakeBlockFromPrev(currBlock)
		bc.HandleBlock(newBlock)
		currBlock = bc.LastBlock
		utils.Debug.Printf("iteration: %v/49", i)
	}
	if bc.Length != 51 {
		t.Errorf("Expected chain length: %v\n Actual chain length: %v", 51, bc.Length)
	}
	if bc.LastHash != currBlock.Hash() {
		t.Errorf("Expected last hash: %v\nActual last hash: %v", currBlock.Hash(), bc.LastHash)
	}
	if bc.LastBlock != currBlock {
		t.Errorf("Expected block: %v\n Actual block: %v", currBlock, bc.LastBlock)
	}
}

func TestHandle2Forks(t *testing.T) {
	defer cleanUp()
	bc := blockchain.New(blockchain.DefaultConfig())
	currBlock := bc.LastBlock
	currForkingBlock := bc.LastBlock

	for i := 0; i < 3; i++ {
		newBlock := MakeBlockFromPrev(currBlock)
		newForkingBlock := MakeBlockFromPrev(currForkingBlock)
		newForkingBlock.Header.Version = 2
		if newBlock.Hash() == newForkingBlock.Hash() {
			t.Errorf("Hashes should not be the same")
		}
		if i != 1 {
			bc.HandleBlock(newBlock)
			bc.HandleBlock(newForkingBlock)
		} else {
			bc.HandleBlock(newForkingBlock)
			bc.HandleBlock(newBlock)
		}
		currBlock = newBlock
		currForkingBlock = newForkingBlock
	}
	if bc.Length != 4 {
		t.Errorf("Expected chain length: %v\n Actual chain length: %v", 4, bc.Length)
	}
	if bc.LastHash != currBlock.Hash() {
		t.Errorf("Expected last hash: %v\nActual last hash: %v", currBlock.Hash(), bc.LastHash)
	}
	if bc.LastBlock != currBlock {
		t.Errorf("Expected block: %v\n Actual block: %v", currBlock, bc.LastBlock)
	}
}
