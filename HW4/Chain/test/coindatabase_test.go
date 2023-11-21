package test

import (
	"Chain/pkg/block"
	"Chain/pkg/blockchain/chainwriter"
	"Chain/pkg/blockchain/coindatabase"
	"reflect"
	"testing"
)

func TestValidateValidBlock(t *testing.T) {
	defer cleanUp()
	genBlock := GenesisBlock()
	coinDB := coindatabase.New(coindatabase.DefaultConfig())
	coinDB.StoreBlock(genBlock.Transactions, true)
	block1 := MakeBlockFromPrev(genBlock)
	if !coinDB.ValidateBlock(block1.Transactions) {
		t.Errorf("block1 should have validated")
	}
}

func TestValidateInvalidBlock(t *testing.T) {
	defer cleanUp()
	genBlock := GenesisBlock()
	coinDB := coindatabase.New(coindatabase.DefaultConfig())
	coinDB.StoreBlock(genBlock.Transactions, true)
	block1 := MakeBlockFromPrev(genBlock)
	block2 := MakeBlockFromPrev(block1)
	if coinDB.ValidateBlock(block2.Transactions) {
		t.Errorf("block2 should not have validated")
	}
}

func TestUndoCoins(t *testing.T) {
	defer cleanUp()
	genBlock := GenesisBlock()
	coinDB := coindatabase.New(coindatabase.DefaultConfig())
	coinDB.StoreBlock(genBlock.Transactions, true)
	block1 := MakeBlockFromPrev(genBlock)
	coinDB.StoreBlock(block1.Transactions, true)
	block2 := MakeBlockFromPrev(block1)
	ub2 := UndoBlockFromBlock(block2)
	coinDB.StoreBlock(block2.Transactions, true)
	coinDB.UndoCoins([]*block.Block{block2}, []*chainwriter.UndoBlock{ub2})
	// make sure coins from undo block are put back
	for i := 0; i < len(ub2.TransactionInputHashes); i++ {
		cl := coindatabase.CoinLocator{
			ReferenceTransactionHash: ub2.TransactionInputHashes[i],
			OutputIndex:              ub2.OutputIndexes[i],
		}
		coin := coinDB.GetCoin(cl)
		if coin == nil {
			t.Errorf("coin should exist")
		} else {
			if coin.IsSpent {
				t.Errorf("coin should not be spent")
			}
		}
	}
	// make sure coins from block are deleted
	for _, tx := range block2.Transactions {
		txHash := tx.Hash()
		for i := 0; i < len(tx.Outputs); i++ {
			cl := coindatabase.CoinLocator{
				ReferenceTransactionHash: txHash,
				OutputIndex:              uint32(i),
			}
			if coin := coinDB.GetCoin(cl); coin != nil {
				t.Errorf("Coin should not exist")
			}
		}
	}
}

func TestGetCoin(t *testing.T) {
	defer cleanUp()
	genBlock := GenesisBlock()
	coinDB := coindatabase.New(coindatabase.DefaultConfig())
	coinDB.StoreBlock(genBlock.Transactions, true)
	txHash := genBlock.Transactions[0].Hash()
	cl := coindatabase.CoinLocator{
		ReferenceTransactionHash: txHash,
		OutputIndex:              0,
	}
	coin := coinDB.GetCoin(cl)
	if coin.IsSpent {
		t.Errorf("Expected coin.IsSpent: %v\nActual coin.IsSpent:%v", false, coin.IsSpent)
	}
	if !reflect.DeepEqual(coin.TransactionOutput, genBlock.Transactions[0].Outputs[0]) {
		t.Errorf(
			"Expected transaction output: %v\nActual transactionoutput:%v",
			genBlock.Transactions[0].Outputs[0],
			coin.TransactionOutput,
		)
	}
}
