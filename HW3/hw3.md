# HW2-Coin

## Introduction

In Chain, we learned how a cryptocurrency might implement a blockchain for block storage and validation. Now it’s time to flesh out the rest of the system! You’ve heard about nodes, miners, and wallets in class (or will soon). Now’s your chance to dive into the details of how these parts come together.

This project builds on the blockchain that you implemented for your first project. Don’t worry if you didn’t pass all of our tests on Chain – we’ve provided our solution for the blockchain.

>**Before you begin**: After reviewing your feedback from Chain, we wanted to supply you with as much useful information up front as we could. While this handout is admittedly quite long, we expect Coin to take much less time than Chain.

## Components

Coin consists of several main components. Here’s a quick overview:

* **Blockchain**: See [Chain](https://hackmd.io/@cs1951l-2023/chain).
* Miner: The Miner handles our proof of work (POW) consensus mechanism. It’s in charge of finding a winning nonce for the blocks it forms. The miner keeps track of a Transaction Pool, which contains transactions passed to it by the node. Off to the races!
* **Node**: The node handles all top level logic. Among its various duties, it validates incoming blocks (using our blockchain), broadcasts transaction requests from the wallet, and tells the miner (if it has one) when to stop and start mining. Whenever it sees a valid block that appends to its chain, the node should tell its miner to get busy on a new block.
* **Wallet**: The wallet keeps track of our coins and creates transactions, which it passes along to its node. The node can then broadcast these transactions to the network. Miners of nodes that hear about these transactions will add them to their transaction pools. Given a large enough fee (or enough space on a block), the miner will include the transaction in their block.
* **Server**: If we wanted to, we could start a server and run our very own cryptocurrency out of Brown. We’re skeptical this would be a good idea, since there are plenty of bugs lurking around in the code base. But we still could! You don’t need to worry about servers for this assignment, but we still wanted to mention it here.

## Assignments

For this assignment, you will implement the functions listed below. We recommend that you tackle them in the following order:

* Wallet
    * HandleBlock
    * RequestTransaction
* Node
    * BroadcastTransaction
    * HandleMinerBlock
* Miner
    * CalculateNonce
    * GenerateCoinbaseTransaction
    * Mine


## Wallet

The vast majority of cryptocurrency holders trust a third-party to take care of their coins (e.g. services like Coinbase and [Crypto.com](http://crypto.com/)). Even financial services like Venmo and Paypal are starting to join the party. Nodes, on the other hand, nearly always keep track of their own holdings. You’ll learn more about wallets in a future lecture.

Here’s an overview of the relevant `Wallet` fields:

> * `Config`: our wallet’s configuration options. For this project, you only need to worry about `Config.SafeBlockAmount`.
> * `Id`: our wallet’s `SimpleID` keeps track of the node’s public and private keys. More complicated configurations–such as **hierarachical deterministic wallets**–are both more secure and more complex. We opted for keeping things simple: our wallet can only process transactions sent to our sole `PublicKey`.
> * `TransactionRequests`: the channel that our wallet uses to send transaction requests to the `Node`, which can then broadcast this transaction to its peers.
> * `CoinCollection`: this mapping contains all of our coins! This only contains coins that have been confirmed enough times (see `config.SafeBlockAmount`).
> * `UnseenSpentCoins`: this mapping contains all of the coins that we’ve spent (by requesting a transaction) but haven’t yet seen in a Block. Until we see our transaction actually included in a Block on the Chain, it hasn’t happened!
> * `UnconfirmedSpentCoins`: Once we’ve seen one of our spent coins on the chain, we should move it from `UnseenSpentCoins` to `UnconfirmedSpentCoins`. Even though we’ve seen our coin included in a block, we still have to wait until enough confirmations (new blocks added to the chain) have occurred until we can safely spend our coin.
> * `UnconfirmedReceivedCoins`: Since we’re an active cryto-user, we need a way to handle coins sent to us from our friends and peers, as well as the change from our own transactions! Like with `UnconfirmedSpentCoins`, we’ll have to wait until enough confirmations have occured until we can spend any Coin that we receive.

Ok, I think I understand how wallets work. But what are `CoinInfos`?

`CoinInfos` are convenient structs that hold the information about TransactionOutputs necessary for turning them into TransactionInputs. You’ll see them as either keys or values in all of the wallet’s mappings.

> **Note**: HandleBlock and RequestTransaction will likely be the two most time-consuming functions that you have to write (given the amount of steps involved). But we’re pretty confident you’ll feel like a superhuman when you do finish them :smiley:. And, you’ll see how Bitcoin-esque wallets vary drastically from traditional bank accounts.

### Implementation task 1

In file `pkg/wallet/wallet.go`:

```
// HandleBlock handles the transactions of a new block. It:
// (1) sees if any of the inputs are ones that we've spent
// (2) sees if any of the incoming outputs on the block are ours
// (3) updates our unconfirmed coins, since we've just gotten
// another confirmation!
func (w *Wallet) HandleBlock(txs []*block.Transaction)
```

Overview:
* Simply put, `HandleBlock` moves coins (via coinInfos) between our various mappings. This function looks at all of the transactions in a block and sees if any of them involve the wallet’s owner.
* When we see a coin that we’ve spent in a transaction, we should move its coinInfo from `UnseenSpentCoins` to `UnconfirmedSpentCoins`.
* When we see a coin sent to us, we should add it to our `UnconfirmedReceivedCoins`.
* Since we’ve just seen a new Block, we should update the number of confirmations for all of our coins in `UnconfirmedSpentCoins` and `UnconfirmedReceivedCoins`. If any of the coins contained in those mappings have had enough confirmations, we can safely add them to our `CoinCollection` (if received) or remove them (if spent).

### Implementation task 2

In file `pkg/wallet/wallet.go`:

```
// RequestTransaction allows the wallet to send a transaction to the node,
// which will propagate the transaction along the P2P network.
func (w *Wallet) RequestTransaction(amount uint32, fee uint32, recipientPK []byte) *block.Transaction
```

Overview:
* When the wallet’s owner wants to create a transaction, they must have both a large enough balance and tangible coins to send.
* Once you’re confident you have the funds to cover the transaction’s amount and fee, you’ll have to actually build it. Think about which order makes more sense: inputs or outputs first?
* You’ll want to temporarily remove any coins used in the transaction from the `CoinCollection` and move them to `UnseenSpentCoins`. **Note**: Our implementation doesn’t allow users to (1) send multiple transactions with the same coin or (2) rebroadcast transactions that they haven’t seen on the Chain in a while.
* Make sure that you actually send your finalized transaction to the `TransactionRequests` channel. Otherwise, the node won’t be able to broadcast it!
* Lastly, update your balance. We temporarily reduce our balance by the total of the inputs when we have a pending transaction. This means we’ll have to wait until we’ve seen our transaction in the Chain (with enough confirmations) to get our change back!

Helpful functions:

* `func (txo *TransactionOutput) MakeSignature(id id.ID) (string, error)`
* `func (id *SimpleID) GetPublicKeyString() string`
* While we did not provide our implementations of `generateTransactionInputs` and `generateTransactionOutputs`, we did leave you the function headers as a guide to how you might want to approach this task.

### Bonus implementation task

Currently, our wallet can’t recover from a fork. But it wouldn’t take much additional functionality to do so! In file `pkg/wallet/wallet.go`:

```
// HandleFork handles a fork, updating the wallet's relevant fields.
func (w *Wallet) HandleFork(blocks []*block.Block, undoBlocks []*chainwriter.UndoBlock)
```

Helpful functions:
* You’re on your own for this one! Similar to our blockchain’s `handleFork`––and its use of coinDB’s `undoCoins`––what fields need to be updated when reverting?
* This name of the unit test is `TestHandleFork`

## Node

Nodes are the workhouses of cryptocurrencies like Bitcoin. They send information to their peers, run miners, and generally do a lot of validation.

Our node should take advantage of non-blocking `goroutines` to delegate work to our miner/wallet/blockchain and continue on with its duties. That way, the nodes won’t have to wait for the miner/wallet/blockchain to finish before we can continue. Head back to HW0 for a refresher on goroutines.

Here are `Node` fields that you have to be aware of for this project:

* `Config`: the node’s configuration. You’ll want to know whether this node has a miner
* `Blockchain`, `Wallet`, and `Miner`: self-explanatory
* `SeenTransactions`: the transactions that our node has either (1) heard about from other nodes or (2) broadcast itself
* `SeenBlocks`: same as `SeenTransactions`, but for blocks
* `PeerDb`: the database of peers that the Node is currently connected to.

### Implementation task 3

In file `pkg/node.go`:

```
// BroadcastTransaction broadcasts transactions created by the wallet
// to other peers in the network.
func (n *Node) BroadcastTransaction(tx *block.Transaction)
```

Overview:
* Our wallet has just requested that we broadcast one of its transactions to our peers, who will then broadcast to their peers, etc…
* If we have a miner, it should update its `TXPool` somehow, since it can now include this transaction in one of its own blocks.
* We should add this transaction to those that we’ve seen.
* We should send this transaction to every peer in our `PeerDb`. Make sure to wrap your RPC call in a goroutine!

Helpful functions:

* `func (m *Miner) HandleTransaction(t *block.Transaction)`
* `func (a *Address) ForwardTransactionRPC(request *pro.Transaction) (*pro.Empty, error)`: this function allows us to forward a transaction to other nodes.

### Implementation task 4

In file `pkg/node.go`:
```
// HandleMinerBlock handles a block
// that was just made by the miner. It does this
// by sending the block to the chain so that it can be
// added, to the wallet, and to the network to be
// broadcast.
func (n *Node) HandleMinerBlock(b *block.Block)
```

Overview:
* Woo! Our miner has succesfully mined a block!
* We should add it to the blocks that we’ve seen.
* Our blockchain needs to be udpated somehow.
* If we have a wallet, it should update its various mappings somehow.
* We need to broadcast this block to all the peers in our network, just like we did with the transaction in `BroadcastTransaction`

Helpful functions:
* `func (a *Address) ForwardBlockRPC(request *pro.Block) (*pro.Empty, error)`: this function allows us to forward a block to other nodes.

## Miner

The infamous miners. Greedy miners compile the transactions that give them the highest reward into blocks, which they then race to find a winning nonce (number only used once) for. As soon as they hear about another block, they have to start from scratch. It’s a rat race.

Not only does the miner get to keep all transaction fees, they also get the minting reward for each block. Recall that minting rewards [halve](https://www.coinmama.com/blog/the-bitcoin-halving-a-history/) every 10,000 blocks. Our halvings occur more frequently.

Here’s an overview of the relevant Miner fields:
* `Config`: our miner’s configuration options. For this project, you’ll only really have to worry about `Config.NonceLimit`
* `TxPool`: all of the transactions that the miner can pull from to create their block. `TxPool` is an enhanced priority queue that also keeps track of total priority, to know whether or not the miner should mine at all.
* `MiningPool`: the transactions that the miner is actively mining in the current block. This is a subset of the `TxPool`.
* `PreviousHash`: the hash of the current last block on the main chain, so that our miner knows where to point to when constructing its new block.
* `ChainLength`: used by the miner to calculate the current minting reward.
* `SendBlock`: the channel that the node monitors to know when the miner has successfully mined a block. See `node.Start` for how this is used.
* **PoolUpdated**: the channel used to send information about `TxPool` updates to the miner.
* **DifficultyTarget**: the level of difficulty that a winning nonce’s resulting hash must beat.

### Uhm… what’s `context.Context`?

* We use `context` to send cancelleation signals to goroutines, which take on lives of their own once spawned.
* If you’re still a bit fuzzy on context, here is the [documentation](https://pkg.go.dev/context) and some [example code](https://gobyexample.com/context).

### Why are we using `atomic` variables?

* Since our miner runs concurrently using goroutines, we must be mindful of any potential synchronization issues that could arise.
* The `atomic` package provides low-level atomic memory primitives useful for implementing synchronization algorithms – it basically allows us to read/write to shared memory without having to deal with mutexes.
* Ask us on Ed if you have questions about concurrency or synchronization (especially if you haven’t taken a systems course like cs0300 or cs0330).
* **Fun fact**: Multi-processor synchronization is one of Dr. Herlihy’s specialties, so if you geek about this kind of stuff, go talk to him about it at his office hours!

>Important! A note on `select`:
>* Go’s `select` statement allows you to wait for results from multiple channels. It’s really just a fancy version of the classic switch statement, where your cases are channel results.
>* For examples of how the select statement might be used, check out **this post** and `node.Start`. We usually combine select statements with infinite loops, so that a process just monitors relevant channels until it’s terminated.

### Implementation task 5

In file `pkg/miner/mine.go`:

```
// CalculateNonce finds a winning nonce for a block. It uses 
// context to know whether it should quit before it finds a nonce 
// (if another block was found). ASICSs are optimized for this task.
func (m *Miner) CalculateNonce(ctx context.Context, b *block.Block) bool
```

Overview:
* This is what miners spend the majority of their time doing: they update the block header’s Nonce until they find a hash smaller than the difficulty target.
* You should try all possible nonce values, altering the block’s header.Nonce value until your resulting hash is less than the difficulty target.
    * If the current difficulty target is `0x0003eab192` (decimal: 65,712,530), you’ll need to find a smaller hash, such as `0x0002ffffff` (decimal: 50,331,647). **Note**: our hashes are not technically hex, so we’ll be comparing byte values.
* Don’t forget about context! How can we use it to monitor whether our nonce calculations are still a worthwhile pursuit?
* Check that our resulting hash beats the difficulty target.

Some helpful functions:
* context’s `Done() <-chan struct{}`: we need to know whether we should stop mining.

### Implementation task 6

In file `pkg/miner/mine.go`:

```
// GenerateCoinbaseTransaction generates a coinbase
// transaction based off the transactions in the mining pool.
// It does this by adding the fee reward to the minting reward.
func (m *Miner) GenerateCoinbaseTransaction(txs []*block.Transaction) *block.Transaction
```

Overview:
* The miner works so hard in hopes of securing a coinbase transaction, in which it receives transaction fees and the minting reward.
* Transaction fees are implicit: they’re the sum of the inputs minus the sum of the outputs.
* You should aggregate all of the fees for the transactions, get the minting reward, and make sure this transaction is addressed to you! How can we get that from our Id?
* The coinbase transaction is an exception to the rule: it does NOT have any inputs

Some helpful functions:

* `func CalculateMintingReward(c *Config, chainLength uint32) uint32`: we need to know what the reward is for this chain length!
* `func (m *Miner) getInputSums(txs []*block.Transaction) ([]uint32, error)` : our miner uses its `GetInputSums` channel to request this information from the node. Once it receives the request, the node will get the input sums for each transaction from the blockchain and return it to the miner via the `InputSums` channel

### Implementation task 7

In file `pkg/miner/mine.go`:

```
// When asked to mine, the miner selects the transactions
// with the highest priority to add to the mining pool. 
func (m *Miner) Mine()
```
Overview:
* In our cryptocurrency, a miner needs to be sure it’s even worth mining a block. After all, CPU usage can get pricey. Thus, there better be enough transactions worth mining. Hint: check out the `TxPool`’s functions.
* If we’re set to mine, then we should set our `Mining` field to true.
* Now, our miner should select the very best transactions to mine. We suggest you look at the `MiningPool`.
* Once we figure out which transactions we want to mine, it’s time to construct the block–with our coinbase transaction at the very top, of course.
* Then, we do the dirty work: we try and find a winning nonce for our block, so that our node can broadcast it to rest of our network and get our sweet reward.
* After successfully(or unsuccessfully finding our nonce), we can safely set our `Mining` field to false. We’re done.
* If we were successful in finding a winning nonce, we should send the block so that our node can hear about it, then we should handle the block ourselves (in order to remove the block’s transactions from our `TxPool`)

Some helpful functions:
* atomic’s `func (x *Bool) Store(v bool)`: this is how we update an atomic boolean.
* You should be able to figure out which other functions to use based on the overview above!


## Testing 

We have provided all test functions in the files under `test`, and each testing file corresponds to a file to be completed. We have provided several helper functions in `test/testing_utils.go` and `test/mocking_utils.go`.

In the `test` directory, run the following command to perform unit test for each module:

```
go test -v {module}_test.go testing_utils.go mocking_utils.go
```

or run `go test -v` to test all functions totally.

## Grading

Your score is based on the number of testing functions your implementation passes. 

### Wallet

|Test Function	|Points|
|:-:|:-|
|TestHandleBlock	|15|
|TestBasicTransactionRequest	|10|
|TestMultipleTransactionRequests	|15|
|TestRequestTransactionFails	|15|
|Total	|55|

### Node

|Test Function	|Points|
|:-:|:-|
|TestHandleMinerBlock	|15|
|Total	|15|

### Miner

|Test Function	|Points|
|:-:|:-|
|TestGenerateCoinbaseTransaction	|5|
|TestCalculateNonce	|5|
|TestCalculateNonceContext	|5|
|TestMine	|15|
|Total	|30|

### Fork (Bonus)
|Test Function	|Points|
|:-:|:-|
|TestHandleFork	|10|
|Total	|10|