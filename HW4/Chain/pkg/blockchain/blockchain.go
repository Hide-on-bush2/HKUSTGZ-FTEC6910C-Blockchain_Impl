package blockchain

import (
	"Chain/pkg/block"
	"Chain/pkg/blockchain/blockinfodatabase"
	"Chain/pkg/blockchain/chainwriter"
	"Chain/pkg/blockchain/coindatabase"
	"Chain/pkg/utils"
	"math"
)

// BlockChain is the main type of this project.
// Length is the length of the active chain.
// LastBlock is the last block of the active chain.
// LastHash is the hash of the last block of the active chain.
// UnsafeHashes are the hashes of the "unsafe" blocks on the
// active chain. These "unsafe" blocks may be reverted during a
// fork.
// maxHashes is the number of unsafe hashes that the chain keeps track of.
// BlockInfoDB is a pointer to a block info database
// ChainWriter is a pointer to a chain writer.
// CoinDB is a pointer to a coin database.
type BlockChain struct {
	Length       uint32
	LastBlock    *block.Block
	LastHash     string
	UnsafeHashes []string
	maxHashes    int

	BlockInfoDB *blockinfodatabase.BlockInfoDatabase
	ChainWriter *chainwriter.ChainWriter
	CoinDB      *coindatabase.CoinDatabase
}

// New returns a blockchain given a Config.
func New(config *Config) *BlockChain {
	genBlock := GenesisBlock(config)
	hash := genBlock.Hash()
	bc := &BlockChain{
		Length:       1,
		LastBlock:    genBlock,
		LastHash:     hash,
		UnsafeHashes: []string{},
		maxHashes:    6,
		BlockInfoDB:  blockinfodatabase.New(blockinfodatabase.DefaultConfig()),
		ChainWriter:  chainwriter.New(chainwriter.DefaultConfig()),
		CoinDB:       coindatabase.New(coindatabase.DefaultConfig()),
	}
	// have to store the genesis block
	bc.CoinDB.StoreBlock(genBlock.Transactions, true)
	ub := &chainwriter.UndoBlock{}
	br := bc.ChainWriter.StoreBlock(genBlock, ub, 1)
	bc.BlockInfoDB.StoreBlockRecord(hash, br)
	return bc
}

// GenesisBlock creates the genesis Block, using the Config's
// InitialSubsidy and GenesisPublicKey.
func GenesisBlock(config *Config) *block.Block {
	txo := &block.TransactionOutput{
		Amount:        config.InitialSubsidy,
		LockingScript: config.GenesisPublicKey,
	}
	genTx := &block.Transaction{
		Version:  0,
		Inputs:   nil,
		Outputs:  []*block.TransactionOutput{txo},
		LockTime: 0,
	}
	return &block.Block{
		Header: &block.Header{
			Version:          0,
			PreviousHash:     "",
			MerkleRoot:       "",
			DifficultyTarget: "",
			Nonce:            0,
			Timestamp:        0,
		},
		Transactions: []*block.Transaction{genTx},
	}
}

// Task 7
// HandleBlock handles a new Block. At a high level, it:
// (1) Validates and stores the Block.
// (2) Stores the Block and resulting Undoblock to Disk.
// (3) Stores the BlockRecord in the BlockInfoDatabase.
// (4) Handles a fork, if necessary.
// (5) Updates the BlockChain's fields.
func (bc *BlockChain) HandleBlock(b *block.Block) {
	appends := bc.appendsToActiveChain(b)
	blockHash := b.Hash()

	//1. Validate Block

	//2. Make Undo Block

	//3. Store Block in CoinDatabase

	//4. Get BlockRecord for previous Block

	//5. Store UndoBlock and Block to Disk,
  //and get corresponding block record
	

	//6. Store BlockRecord to BlockInfoDatabase

	if appends {
		//7. Handle appending Block
    //update bc.Length, bc.LastBlock, bc.LastHash
		bc.Length++
		bc.LastBlock = b
		bc.LastHash = blockHash
    //8. update bc.UnsafeHashes if there are new block being confirmed
		
	} else if height > bc.Length {
		//9. Handle fork
	
  }
}

// Task 6
// handleFork updates the BlockChain when a fork occurs. First, it
// finds the Blocks the BlockChain must revert. Once found, it uses
// those Blocks to update the CoinDatabase. Lastly, it updates the
// BlockChain's fields to reflect the fork.
func (bc *BlockChain) handleFork(b *block.Block, blockHash string, height uint32) {
  //1. get all blocks to be reverted
	
  //2. as the confirmed block can not be reverted,
  //thus we need to get all hash values of unsafe reverted blocks
	
	//3. update the unsafe hashs of bc
	
  //4. revert the block via bc.CoinDB.UndoCoins
	
  //5. verify and store new blocks to the blockchain
	
  //6. update bc.LastBlock and bc.LastHash
	
}

// Task 2
// makeUndoBlock returns an UndoBlock given a slice of Transactions.
func (bc *BlockChain) makeUndoBlock(txs []*block.Transaction) *chainwriter.UndoBlock {
  //1. initilize multiple arrays corresponding to the fields in UndoBlock for later appending
	var transactionHashes []string
	var outputIndexes []uint32
	var amounts []uint32
	var lockingScripts []string
  //2. iterate each tx in txs
	for _, tx := range txs {
		for _, txi := range tx.Inputs {
      //3. initilize a coin locator, which consists of 
      //the corresponding hash value and output index
			
      //4. get the coin from bc.CoinDB
			
			//5. if the coin is nil it means this isn't even a possible fork,
      //return an empty UndoBlock
			
      //6. append the corresponding information to transactionHashes,
      //outputIndexes, amounts, and lockingScripts
			
		}
	}
  //7. consturct the UndoBlock and return the result


  return nil
}

// getBlock uses the ChainWriter to retrieve a Block from Disk
// given that Block's hash
func (bc *BlockChain) getBlock(blockHash string) *block.Block {
	br := bc.BlockInfoDB.GetBlockRecord(blockHash)
	fi := &chainwriter.FileInfo{
		FileName:    br.BlockFile,
		StartOffset: br.BlockStartOffset,
		EndOffset:   br.BlockEndOffset,
	}
	return bc.ChainWriter.ReadBlock(fi)
}

// getUndoBlock uses the ChainWriter to retrieve an UndoBlock
// from Disk given the corresponding Block's hash
func (bc *BlockChain) getUndoBlock(blockHash string) *chainwriter.UndoBlock {
	br := bc.BlockInfoDB.GetBlockRecord(blockHash)
	fi := &chainwriter.FileInfo{
		FileName:    br.UndoFile,
		StartOffset: br.UndoStartOffset,
		EndOffset:   br.UndoEndOffset,
	}
	return bc.ChainWriter.ReadUndoBlock(fi)
}

// GetBlocks retrieves a slice of blocks from the main chain given a
// starting and ending height, inclusive. Given a chain of length 50,
// GetBlocks(10, 20) returns blocks 10 through 20.
func (bc *BlockChain) GetBlocks(start, end uint32) []*block.Block {
	if start >= end || end <= 0 || start <= 0 || end > bc.Length {
		utils.Debug.Printf("cannot get chain blocks with values start: %v end: %v", start, end)
	}

	var blocks []*block.Block
	currentHeight := bc.Length
	nextHash := bc.LastBlock.Hash()

	for currentHeight >= start {
		br := bc.BlockInfoDB.GetBlockRecord(nextHash)
		fi := &chainwriter.FileInfo{
			FileName:    br.BlockFile,
			StartOffset: br.BlockStartOffset,
			EndOffset:   br.BlockEndOffset,
		}
		if currentHeight <= end {
			nextBlock := bc.ChainWriter.ReadBlock(fi)
			blocks = append(blocks, nextBlock)
		}
		nextHash = br.Header.PreviousHash
		currentHeight--
	}
	return reverseBlocks(blocks)
}

// GetHashes retrieves a slice of hashes from the main chain given a
// starting and ending height, inclusive. Given a BlockChain of length
// 50, GetHashes(10, 20) returns the hashes of Blocks 10 through 20.
func (bc *BlockChain) GetHashes(start, end uint32) []string {
	if start >= end || end <= 0 || start <= 0 || end > bc.Length {
		utils.Debug.Printf("cannot get chain blocks with values start: %v end: %v", start, end)
	}

	var hashes []string
	currentHeight := bc.Length
	nextHash := bc.LastBlock.Hash()

	for currentHeight >= start {
		br := bc.BlockInfoDB.GetBlockRecord(nextHash)
		if currentHeight <= end {
			hashes = append(hashes, nextHash)
		}
		nextHash = br.Header.PreviousHash
		currentHeight--
	}
	return reverseHashes(hashes)
}

// appendsToActiveChain returns whether a Block appends to the
// BlockChain's active chain or not.
func (bc *BlockChain) appendsToActiveChain(b *block.Block) bool {
	return bc.LastBlock.Hash() == b.Header.PreviousHash
}

// getForkedBlocks returns a slice of Blocks given a starting hash.
// It returns a maximum of maxHashes Blocks, where maxHashes is the
// BlockChain's maximum number of unsafe hashes.
func (bc *BlockChain) getForkedBlocks(startHash string) []*block.Block {
	unsafeHashes := make(map[string]bool)
	for _, h := range bc.UnsafeHashes {
		unsafeHashes[h] = true
	}
	var forkedBlocks []*block.Block
	nextHash := startHash
	for i := 0; i < len(bc.UnsafeHashes); i++ {
		forkedBlock := bc.getBlock(nextHash)
		forkedBlocks = append(forkedBlocks, forkedBlock)
		if _, ok := unsafeHashes[nextHash]; ok {
			return forkedBlocks
		}
		nextHash = forkedBlock.Header.PreviousHash
	}
	return forkedBlocks
}

// getBlocksAndUndoBlocks returns a slice of n Blocks with a
// corresponding slice of n UndoBlocks.
func (bc *BlockChain) getBlocksAndUndoBlocks(n int) ([]*block.Block, []*chainwriter.UndoBlock) {
	var blocks []*block.Block
	var undoBlocks []*chainwriter.UndoBlock
	nextHash := bc.LastHash
	for i := 0; i < n; i++ {
		b := bc.getBlock(nextHash)
		ub := bc.getUndoBlock(nextHash)
		blocks = append(blocks, b)
		undoBlocks = append(undoBlocks, ub)
		nextHash = b.Header.PreviousHash
	}
	return blocks, undoBlocks
}

// reverseBlocks returns a reversed slice of Blocks.
func reverseBlocks(s []*block.Block) []*block.Block {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// reverseHashes returns a reversed slice of hashes.
func reverseHashes(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
