package test

import (
  "Chain/pkg/block"
  "testing"
  "math/rand"
  "time"
  "fmt"
)

//Generate a random tx
func GenRandomTransaction() *block.Transaction {
  //The seed is changed with varying time
  rand.Seed(time.Now().UnixNano())
  txi := &block.TransactionInput {
    ReferenceTransactionHash: "",
    OutputIndex: uint32(rand.Intn(10)),
    UnlockingScript: "",
  }
  inputs := []*block.TransactionInput{txi}

  txo := &block.TransactionOutput{
    Amount: uint32(rand.Intn(10)),
    LockingScript: "",
  }
  outputs := []*block.TransactionOutput{txo}

  tx := &block.Transaction{
    Version: uint32(rand.Intn(10)),
    Inputs: inputs,
    Outputs: outputs,
    LockTime: uint32(rand.Intn(10)),
  }
  return tx
}

//Comparation of two transactions
func IsTransactionsEqual(tx1 *block.Transaction, tx2 *block.Transaction) bool {
  if tx1.Version != tx2.Version {
    return false
  }
  if tx1.LockTime != tx2.LockTime {
    return false
  }
  tx1_in_len := len(tx1.Inputs)
  tx1_out_len := len(tx1.Outputs)

  tx2_in_len := len(tx2.Inputs)
  tx2_out_len := len(tx2.Outputs)

  if tx1_in_len != tx2_in_len {
    return false
  }
  if tx1_out_len != tx2_out_len {
    return false
  }

  for i :=0; i < tx1_in_len; i++ {
    tx1_input := tx1.Inputs[i]
    tx2_input := tx2.Inputs[i]

    if tx1_input.OutputIndex != tx2_input.OutputIndex {
      return false
    }
  }

  for i := 0; i < tx1_out_len; i++ {
    tx1_output := tx1.Outputs[i]
    tx2_output := tx2.Outputs[i]

    if tx1_output.Amount != tx2_output.Amount {
      return false
    }
  }
  return true
}

func TestTransactionCoding(t *testing.T) {
  //1. generate a random tx
  tx := GenRandomTransaction()
  //2. encode the tx
  encoded_tx := block.EncodeTransaction(tx)
  //3. decode the encoded tx
  decoded_tx := block.DecodeTransaction(encoded_tx)
  
  //4. check whether tx and decoded_tx are the same
//  if !IsTransactionsEqual(tx, decoded_tx) {
//    t.Errorf("The decoded tx is not equal to the original one")
//  }
  if tx.Hash() != decoded_tx.Hash() {
    t.Errorf("The decoded tx is not equal to the original one")
  }

}

//Simply test if the hash function is executed as expected
//however the correctness of it is not tested, as the implementation
//varies personally
func TestTransactionHash(t *testing.T) {
  //1. generate a random tx
  tx := GenRandomTransaction()
  //2. run the hash function
  tx_hash := tx.Hash()
  fmt.Println("the hash is ", tx_hash)
}

//Generate a random block
func GenRandomBlock() *block.Block {
  //The seed is changed with varying time
  rand.Seed(time.Now().UnixNano())
  
  //create a random header
  header := &block.Header{
    Version: uint32(rand.Intn(10)),
    PreviousHash: "",
    MerkleRoot: "",
    DifficultyTarget: "",
    Nonce: uint32(rand.Intn(10)),
    Timestamp: uint32(rand.Intn(10)),
  }

  //generate 3 random txs
  tx1 := GenRandomTransaction()
  tx2 := GenRandomTransaction()
  tx3 := GenRandomTransaction()
  txs := []*block.Transaction{tx1, tx2, tx3}

  block := &block.Block{
    Header: header,
    Transactions: txs,
  }

  return block
}

func TestBlockCoding(t *testing.T) {
  //1. generate a random block
  blk := GenRandomBlock()
  //2. encode the block
  encoded_blk := block.EncodeBlock(blk)
  //3. decode the block
  decoded_blk := block.DecodeBlock(encoded_blk)
  //4. check whether the block and decoded_block are the same

  if blk.Hash() != decoded_blk.Hash() {
    t.Errorf("The decoded block is not equal to the original one")
  }
}

//Simply test if the hash function is executed as expected
//however the correctness of it is not tested, as the implementation
//varies personally
func TestBlockHash(t *testing.T) {
  //1. generate a random block
  block := GenRandomBlock()
  //2. run the hash function
  block_hash := block.Hash()
  fmt.Println("the hash is ", block_hash)
}
