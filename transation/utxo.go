package transation

import "fmt"

type UnspentOutput struct {
	TransationId          string
	TransationOutputIndex int
	Address               string
	Amount                int
}

var unspentTxOuts = []*UnspentOutput{}

func GetUnspentTxOuts() []*UnspentOutput {
	return unspentTxOuts
}

type WalletInface interface {
	GetMyAddress() string
}

func UpdateUnspentTxOuts(newUnspentTxOuts []*UnspentOutput) {

	fmt.Printf("Updating unspent transaction outputs %#v", newUnspentTxOuts)
	unspentTxOuts = newUnspentTxOuts
}

func GetUnspentTxOut(txId string, txIndex int) *UnspentOutput {
	var result *UnspentOutput

	for _, tx := range unspentTxOuts {
		if tx.TransationId == txId && tx.TransationOutputIndex == txIndex {
			result = tx
		}
	}

	return result
}

func GetMyUnspentTxOuts(wallet WalletInface) []*UnspentOutput {
	var myUnspentTxOuts []*UnspentOutput

	for _, tx := range unspentTxOuts {
		if tx.Address == wallet.GetMyAddress() {
			myUnspentTxOuts = append(myUnspentTxOuts, tx)
		}
	}

	return myUnspentTxOuts
}

func FindUnspentTxOutput(wallet WalletInface, amount int) ([]*UnspentOutput, int) {
	myUnSpentedTxOutput := GetMyUnspentTxOuts(wallet)
	var unspentTxOutputs []*UnspentOutput
	totalOutputValue := 0
	restOutputBalance := 0
	for _, txOutput := range myUnSpentedTxOutput {
		totalOutputValue = totalOutputValue + txOutput.Amount
		unspentTxOutputs = append(unspentTxOutputs, txOutput)
		if totalOutputValue >= amount {
			break
		}
	}

	fmt.Println("findUnspentTxOutput", totalOutputValue, amount, len(unspentTxOutputs))
	fmt.Printf("findUnspentTxOutput value: %#v \n", unspentTxOutputs[0])

	if totalOutputValue < amount {
		return nil, 0
	}

	restOutputBalance = totalOutputValue - amount

	return unspentTxOutputs, restOutputBalance
}
