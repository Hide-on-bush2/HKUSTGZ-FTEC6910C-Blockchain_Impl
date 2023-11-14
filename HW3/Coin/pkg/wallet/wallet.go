package wallet

import (
	"Coin/pkg/block"
	"Coin/pkg/blockchain/chainwriter"
	"Coin/pkg/id"
	"Coin/pkg/utils"
)

// CoinInfo holds the information about a TransactionOutput
// necessary for making a TransactionInput.
// ReferenceTransactionHash is the hash of the transaction that the
// output is from.
// OutputIndex is the index into the Outputs array of the
// Transaction that the TransactionOutput is from.
// TransactionOutput is the actual TransactionOutput
type CoinInfo struct {
	ReferenceTransactionHash string
	OutputIndex              uint32
	TransactionOutput        *block.TransactionOutput
}

// Wallet handles keeping track of the owner's coins
//
// CoinCollection is the owner of this wallet's set of coins
//
// UnseenSpentCoins is a mapping of transaction hashes (which are strings)
// to a slice of coinInfos. It's used for keeping track of coins that we've
// used in a transaction but haven't yet seen in a block.
//
// UnconfirmedSpentCoins is a mapping of Coins to number of confirmations
// (which are integers). We can't confirm that a Coin has been spent until
// we've seen enough POW on top the block containing our sent transaction.
//
// UnconfirmedReceivedCoins is a mapping of CoinInfos to number of confirmations
// (which are integers). We can't confirm we've received a Coin until
// we've seen enough POW on top the block containing our received transaction.
type Wallet struct {
	Config              *Config
	Id                  id.ID
	TransactionRequests chan *block.Transaction
	Address             string
	Balance             uint32

	// All coins
	CoinCollection map[*block.TransactionOutput]*CoinInfo

	// Not yet seen
	UnseenSpentCoins map[string][]*CoinInfo

	// Seen but not confirmed
	UnconfirmedSpentCoins    map[*CoinInfo]uint32
	UnconfirmedReceivedCoins map[*CoinInfo]uint32
}

// SetAddress sets the address
// of the node in the wallet.
func (w *Wallet) SetAddress(a string) {
	w.Address = a
}

// New creates a wallet object
func New(config *Config, id id.ID) *Wallet {
	if !config.HasWallet {
		return nil
	}
	return &Wallet{
		Config:                   config,
		Id:                       id,
		TransactionRequests:      make(chan *block.Transaction),
		Balance:                  0,
		CoinCollection:           make(map[*block.TransactionOutput]*CoinInfo),
		UnseenSpentCoins:         make(map[string][]*CoinInfo),
		UnconfirmedSpentCoins:    make(map[*CoinInfo]uint32),
		UnconfirmedReceivedCoins: make(map[*CoinInfo]uint32),
	}
}

// generateTransactionInputs creates the transaction inputs required to make a transaction.
// In addition to the inputs, it returns the amount of change the wallet holder should
// return to themselves, and the coinInfos used
func (w *Wallet) generateTransactionInputs(amount uint32, fee uint32) (uint32, []*block.TransactionInput, []*CoinInfo) {
	// the inputs that we will eventually be returning
	var inputs []*block.TransactionInput
	// the coinInfos that we're using
	var coinInfos []*CoinInfo
	// the total amount of the coins that we've used so far for our inputs
	total := uint32(0)
	// Now that we know our balance is enough, we can loop through our coins until we've reached
	// a large enough total to meet our amount and fee
	for coin, coinInfo := range w.CoinCollection {
		if total >= amount+fee {
			break
		}
		// have to generate the unlockingScripts so that we can prove we have the ability to spend
		// this coin
		unlockingScript, err := coin.MakeSignature(w.Id)
		if err != nil {
			utils.Debug.Printf("[generateTransactionInputs] Error: failed to create unlockingScript\n")
		}
		// actually create the transaction input
		txi := &block.TransactionInput{
			ReferenceTransactionHash: coinInfo.ReferenceTransactionHash,
			OutputIndex:              coinInfo.OutputIndex,
			UnlockingScript:          unlockingScript,
		}
		coinInfos = append(coinInfos, coinInfo)
		inputs = append(inputs, txi)
		total += coin.Amount
	}
	change := total - (amount + fee)
	return change, inputs, coinInfos
}

// generateTransactionOutputs generates the transaction outputs required to create a transaction.
func (w *Wallet) generateTransactionOutputs(
	amount uint32,
	receiverPK []byte,
	change uint32,
) []*block.TransactionOutput {
	// make sure that the public key we're sending our amount to is valid
	if receiverPK == nil || len(receiverPK) == 0 {
		utils.Debug.Printf("[generateTransactionOutputs] Error: receiver's public key is invalid")
		return nil
	}
	// the outputs that we will eventually return
	var outputs []*block.TransactionOutput
	// the output for the person we're sending this transaction output to
	txoSending := &block.TransactionOutput{Amount: amount, LockingScript: string(receiverPK)}
	outputs = append(outputs, txoSending)
	// if there's change, we should send that back to ourselves.
	if change != 0 {
		txoChange := &block.TransactionOutput{Amount: change, LockingScript: string(w.Id.GetPublicKeyBytes())}
		outputs = append(outputs, txoChange)
	}
	return outputs
}


// Implementation task 2
// RequestTransaction allows the wallet to send a transaction to the node,
// which will propagate the transaction along the P2P network.
func (w *Wallet) RequestTransaction(amount uint32, fee uint32, recipientPK []byte) *block.Transaction {
	// 1. have to ensure that we have enough money to actually make this transaction
	
	// 2. now that we have the transaction, we can add the coinInfos to our UnseenSpentCoins
	// and temporarily remove from the CoinCollection

	// 3. send to the channel that the node monitors
	
	// 4. we do this here in case generateTransactionInputs doesn't work
	// have to make sure that the balance is decremented so that the wallet owner can't keep spamming their coin

  return nil
}

// Implementation task 1
// HandleBlock handles the transactions of a new block. It:
// (1) sees if any of the inputs are ones that we've spent
// (2) sees if any of the incoming outputs on the block are ours
// (3) updates our unconfirmed coins, since we've just gotten
// another confirmation!
func (w *Wallet) HandleBlock(txs []*block.Transaction) {
	// 1.most of the time, we will just be handling the transactions
	
  // 2. update the confirmation state

}

// addCoin adds a received coin to our UnconfirmedReceivedCoins
func (w *Wallet) addCoin(hash string, index uint32, output *block.TransactionOutput) {
	coinInfo := &CoinInfo{
		ReferenceTransactionHash: hash,
		OutputIndex:              index,
		TransactionOutput:        output,
	}
	w.UnconfirmedReceivedCoins[coinInfo] = 0
}

func (w *Wallet) updateConfirmations() {
	// update unconfirmed spent coins
	for coinInfo, numConfirmations := range w.UnconfirmedSpentCoins {
		if numConfirmations == w.Config.SafeBlockAmount {
			// if we've seen enough blocks, we can safely remove this
			// coin from our coin collection. It's been spent!
			delete(w.CoinCollection, coinInfo.TransactionOutput)
			delete(w.UnconfirmedSpentCoins, coinInfo)
		} else {
			// otherwise, we still have to wait :(
			w.UnconfirmedSpentCoins[coinInfo] = numConfirmations + 1
		}
	}
	// update unconfirmed received coins
	for coinInfo, numConfirmations := range w.UnconfirmedReceivedCoins {
		if numConfirmations == w.Config.SafeBlockAmount {
			// if we've seen enough blocks, we can safely add this
			// coin to our coin collection. It's spendable!
			w.CoinCollection[coinInfo.TransactionOutput] = coinInfo
			// Also need to update our balance
			w.Balance += coinInfo.TransactionOutput.Amount
			delete(w.UnconfirmedReceivedCoins, coinInfo)
		} else {
			// otherwise, we still have to wait :(
			w.UnconfirmedReceivedCoins[coinInfo] = numConfirmations + 1
		}
	}
}

// handleSeenCoins moves coins from UnseenSpentCoins to
// UnconfirmedSpentCoins
func (w *Wallet) handleSeenCoins(hash string) {
	seenCoins, _ := w.UnseenSpentCoins[hash]
	// remove from unseen, since we've now seen our
	// transaction in a block
	delete(w.UnseenSpentCoins, hash)
	// move the seen coins over to unconfirmed
	for _, coinInfo := range seenCoins {
		w.UnconfirmedSpentCoins[coinInfo] = 0
	}
}

// Bonus implementation task
// HandleFork handles a fork, updating the wallet's relevant fields.
func (w *Wallet) HandleFork(blocks []*block.Block, undoBlocks []*chainwriter.UndoBlock) {

}
