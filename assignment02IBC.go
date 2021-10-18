package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const miningReward = 100
const rootUser = "Satoshi"

// BlockData is a structure containing attributes of a block
type BlockData struct {
	Title    string
	Sender   string
	Receiver string
	Amount   int
}

// Block is a structure containing information of blockchain
type Block struct {
	Data        []BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

// CalculateBalance is used to check the balance of the sender
func CalculateBalance(userName string, chainHead *Block) int {
	if strings.ToLower(userName) == "system" {
		return 100000 // System will have infinite balance
	}

	balance := 0 // Initially balance will be 0
	for c := chainHead; c != nil; c = c.PrevPointer {
		for i := range c.Data {
			if c.Data[i].Sender == userName {
				balance -= c.Data[i].Amount
			}
			if c.Data[i].Receiver == userName {
				balance += c.Data[i].Amount
			}
		}
	}

	// BONUS TASK IMPLEMENTED
	//------------------------------------------------------------------------------
	if balance < 0 {
		fmt.Println()
		fmt.Println("Balance cannot be negative. Balance would become ", balance)
		fmt.Println("As you punishment, your balance is now 0 >.<")
		fmt.Println()
		balance = 0
		return balance
	}
	//------------------------------------------------------------------------------
	return balance
}

// CalculateHash is used to compute the hash of the inputBlock
func CalculateHash(inputBlock *Block) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", *inputBlock))))
}

// VerifyTransaction is used to see if the transaction are valid or not
func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
	if CalculateBalance(transaction.Sender, chainHead) < 0 {
		return false
	}
	if CalculateBalance(transaction.Sender, chainHead) >= transaction.Amount {
		return true
	}

	fmt.Println()
	fmt.Println("\t\t\t\t!!!!!!!!!!!ERROR!!!!!!!!!!!")
	fmt.Println(transaction.Sender, "has current balance = ", CalculateBalance(transaction.Sender, chainHead), ", whereas he/she requested the amount = ", transaction.Amount)
	fmt.Println()
	return false
}

// InsertBlock is used to insert a block into the blockchain
func InsertBlock(blockData []BlockData, chainHead *Block) *Block {

	coinbaseTrans := BlockData{Title: "Coinbase", Sender: "System", Receiver: rootUser, Amount: miningReward}
	blockData = append(blockData, coinbaseTrans)
	for i := range blockData {
		transaction := blockData[i]
		if !VerifyTransaction(&transaction, chainHead) {
			fmt.Println("Invalid transaction detected = ", transaction)
			return chainHead
		}
	}

	if chainHead == nil {
		var newBlock Block
		chainHead = &newBlock
		chainHead.PrevPointer = nil
		chainHead.Data = blockData
		chainHead.CurrentHash = CalculateHash(chainHead)
		chainHead.PrevHash = ""
	} else {
		var newBlock Block
		newBlock.PrevPointer = chainHead
		newBlock.PrevHash = chainHead.CurrentHash
		newBlock.Data = blockData
		newBlock.CurrentHash = CalculateHash(&newBlock)
		chainHead = &newBlock
	}
	return chainHead
}

// ListBlocks is used to display the list of blocks in a chain
func ListBlocks(chainHead *Block) {
	fmt.Println("\t\t\t\t LIST OF BLOCKS")

	newHead := chainHead
	i := 1
	for newHead != nil {
		fmt.Println()
		fmt.Println("-------------------")
		fmt.Println("Block Number = ", i)
		fmt.Println("-------------------")
		fmt.Println("\tTransactions = ")
		for i := range newHead.Data {
			fmt.Println("\t\tTitle: ", newHead.Data[i].Title)
			fmt.Println("\t\tSender: ", newHead.Data[i].Sender)
			fmt.Println("\t\tReceiver: ", newHead.Data[i].Receiver)
			fmt.Println("\t\tAmount: ", newHead.Data[i].Amount)
			fmt.Println()
		}
		fmt.Print("\tCurrent Hash = ")
		fmt.Println(newHead.CurrentHash)
		fmt.Print("\tPrevious Hash = ")
		fmt.Println(newHead.PrevHash)
		i++
		fmt.Println()
		fmt.Println("............................................................................................")
		newHead = newHead.PrevPointer
	}
}

// VerifyChain is used to verify the chain
func VerifyChain(chainHead *Block) {
	for c := chainHead; c != nil; c = c.PrevPointer {
		hashc := CalculateHash(c)
		if c.PrevPointer != nil {
			hashp := CalculateHash(c.PrevPointer)
			if hashp != c.PrevHash || hashc != c.CurrentHash {
				fmt.Println("Blockchain is compromised")
				return
			}
		}
		if hashc != c.CurrentHash {
			fmt.Println("Blockchain is compromised")
			return
		}
	}
	fmt.Println("Blockchain Verified")
	return
}

// PremineChain is used to premine the block
func PremineChain(chainHead *Block, numBlocks int) *Block {
	nilTrans := BlockData{Title: "Premined", Sender: "nil", Receiver: "nil", Amount: 0}
	coinbaseTrans := BlockData{Title: "Coinbase", Sender: "System", Receiver: rootUser, Amount: miningReward}

	for i := 0; i < numBlocks; i++ {
		chainHead = InsertBlock([]BlockData{nilTrans, coinbaseTrans}, chainHead)
	}

	return chainHead
}

/// MAIN STARTS HERE ///
// package main
//
// import (
// 	"fmt"
//
// 	a2 "github.com/HxnDev/assignment02IBC"
// )
//
// func main() {
// 	var chainHead *a2.Block
// 	//This insertion is invalid as Alice is neither miner nor has enough coins for the transaction, pay 50 from Alice to Bob
// 	aliceToBob := []a2.BlockData{{Title: "ALice2Bob", Sender: "Alice", Receiver: "Bob", Amount: 50}}
// 	chainHead = a2.InsertBlock(aliceToBob, chainHead)
// 	fmt.Println(chainHead)
//
// 	//Lets mine some blocks to start the chain and check Satoshi's balance
// 	chainHead = a2.PremineChain(chainHead, 2)
// 	fmt.Printf("Satoshi's balance: %v\n", a2.CalculateBalance("Satoshi", chainHead))
//
// 	//Now Satoshi can send some coins to Alice
// 	SatoshiToAlice := []a2.BlockData{{Title: "SatoshiToAlice", Sender: "Satoshi", Receiver: "Alice", Amount: 50}}
// 	chainHead = a2.InsertBlock(SatoshiToAlice, chainHead)
//
// 	//We can verify this by checking balances once again and listing the chain
// 	fmt.Printf("Satoshi's balance: %v\n", a2.CalculateBalance("Satoshi", chainHead))
// 	fmt.Printf("Alice's balance: %v\n", a2.CalculateBalance("Alice", chainHead))
// 	a2.ListBlocks(chainHead)
//
// 	//Alice can then make the transactions using her coins, She can make multiple
// 	//transactions at once, notice that field Data has type []BlockData in Block Struct
// 	AliceToBobCharlie := []a2.BlockData{{Title: "ALice2Bob", Sender: "Alice", Receiver: "Bob", Amount: 20}, {Title: "ALice2Charlie", Sender: "Alice", Receiver: "Charlie", Amount: 10}}
// 	chainHead = a2.InsertBlock(AliceToBobCharlie, chainHead)
//
// 	//We can verify this by checking balances once again and listing the chain
// 	a2.ListBlocks(chainHead)
//
// 	fmt.Printf("Satoshi's balance: %v ", a2.CalculateBalance("Satoshi", chainHead))
// 	fmt.Printf("Alice's balance: %v ", a2.CalculateBalance("Alice", chainHead))
// 	fmt.Printf("Charlie's balance: %v\n", a2.CalculateBalance("Charlie", chainHead))
//
// 	//Finally the transaction verification fails if any of the transaction is invalid
// 	oneInvalidoneValid := []a2.BlockData{{Title: "ALice2EZ", Sender: "Alice", Receiver: "Bob", Amount: 100}, {Title: "Satoshi2EZ", Sender: "Satoshi", Receiver: "EZ", Amount: 200}}
// 	chainHead = a2.InsertBlock(oneInvalidoneValid, chainHead)
//
// 	//Bonus (2 absolutes) - Fix the erroneous behavior below
// 	//The transactions are valid individually but when applied to chain Alice's balance
// 	//become negative :(
// 	bonusTransactions := []a2.BlockData{{Title: "ALice2Bob", Sender: "Alice", Receiver: "Bob", Amount: 15}, {Title: "AliceToEZ", Sender: "Alice", Receiver: "EZ", Amount: 15}}
// 	chainHead = a2.InsertBlock(bonusTransactions, chainHead)
// 	fmt.Printf("Alice's balance: %v\n", a2.CalculateBalance("Alice", chainHead))
//
// }
