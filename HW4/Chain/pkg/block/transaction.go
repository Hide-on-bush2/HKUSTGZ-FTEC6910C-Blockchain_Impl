package block

import (
	"Chain/pkg/pro"
	"crypto/sha256"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// TransactionInput is used as the input to create a TransactionOutput.
// Recall that TransactionInputs generate TransactionOutputs which in turn
// generate new TransactionInputs and so forth.
// ReferenceTransactionHash is the hash of the parent TransactionOutput's Transaction.
// OutputIndex is the index of the parent TransactionOutput's Transaction.
// Signature verifies that the payer can spend the referenced TransactionOutput.
type TransactionInput struct {
	ReferenceTransactionHash string
	OutputIndex              uint32
	UnlockingScript          string
}

// TransactionOutput is an output created from a TransactionInput.
// Recall that TransactionOutputs generate TransactionInputs which in turn
// generate new TransactionOutputs and so forth.
// Amount is how much this TransactionOutput is worth.
// PublicKey is used to verify the payee's signature.
type TransactionOutput struct {
	Amount        uint32
	LockingScript string
}

// Transaction contains information about a transaction.
// Version is the version of this transaction.
// Inputs is a slice of TransactionInputs.
// Outputs is a slice of TransactionOutputs.
// LockTime is the future time after which the Transaction is valid.
type Transaction struct {
	Version  uint32
	Inputs   []*TransactionInput
	Outputs  []*TransactionOutput
	LockTime uint32
}

// Task 1-1:
// EncodeTransaction returns a pro.Transaction given a Transaction.
func EncodeTransaction(tx *Transaction) *pro.Transaction {
  //1. initilize an array of pro.TransactionInput
	var ptxis []*pro.TransactionInput
  //2. encode each input of tx and append it to the ptxis

  //3. initilize an array of pro.TransactionOutput
	var ptxos []*pro.TransactionOutput
  //4. encode each output of tx and append it to the ptxos

  //5. construct the encoded tx, which is of type pro.Transaction,
  //and return the result
  
  return nil
}

// Task 1-2:
// DecodeTransaction returns a Transaction given a pro.Transaction.
func DecodeTransaction(ptx *pro.Transaction) *Transaction {
  //1. initilize an array of TransactionInput
	var txis []*TransactionInput
  //2. decode each input of ptx and append it to txis
	
  //3. initilize an array of TransactionOutput
	var txos []*TransactionOutput
  //4. decode each output of ptx and append it to txos
	
  //5. construct the decoded tx, which is of type Transaction, 
  //and return the result

  return nil
}

// Task 1-3:
// Hash returns the hash of the transaction
func (tx *Transaction) Hash() string {
  //1. initilize a sha256 instance
	
  //2. encode the tx
	
  //3. serialize the encoded tx
	
  //4. hash the serialized tx
	
  //5. return the result
  
  return ""
}
