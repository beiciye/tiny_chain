package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"tiny_chain/transation"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
}

var myWallet *Wallet = newWallet()

func newWallet() *Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ret := &Wallet{
		PrivateKey: privateKey,
	}
	ret.PrintMyAddress()
	return ret
}

func GetWallet() *Wallet {
	return myWallet
}

func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.PrivateKey
}

func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.PrivateKey.Public().(*ecdsa.PublicKey)
}

func (w *Wallet) GetMyAddress() string {
	publicKeyAddr, _ := x509.MarshalPKIXPublicKey(w.GetPublicKey())
	hexAddr := hex.EncodeToString(publicKeyAddr)
	return hexAddr
}

func (w *Wallet) PrintMyAddress() {
	fmt.Printf("My address: %s\n", w.GetMyAddress())
}

func (w *Wallet) GetBanlance() []*transation.UnspentOutput {
	allUnSpentedTxOutput := transation.GetUnspentTxOuts()

	myUnSpentedTxOutput := make([]*transation.UnspentOutput, 0)

	for _, txOutput := range allUnSpentedTxOutput {
		if txOutput.Address == w.GetMyAddress() {
			myUnSpentedTxOutput = append(myUnSpentedTxOutput, txOutput)
		}
	}
	return myUnSpentedTxOutput
}

func (w *Wallet) CreateTxInput(myUnSpentedTxOutput []*transation.UnspentOutput) []*transation.TransationInput {

	txInputs := make([]*transation.TransationInput, 0)

	for _, txOutput := range myUnSpentedTxOutput {
		txInputs = append(txInputs, &transation.TransationInput{TxOutputId: txOutput.TransationId, TxOutIndex: txOutput.TransationOutputIndex, Siguature: ""})
	}
	return txInputs
}

func (w *Wallet) CreateTxOutput(to string, amount int, restOutputBalance int) []*transation.TransationOutput {

	txOutputs := []*transation.TransationOutput{
		{Address: to, Amount: amount},
	}

	if restOutputBalance > 0 {
		txOutputs = append(txOutputs, &transation.TransationOutput{Address: w.GetMyAddress(), Amount: restOutputBalance})
	}

	return txOutputs
}

func (w *Wallet) CreateTransation(to string, amount int) *transation.Transation {
	unspentTxOutputs, restOutputBalance := transation.FindUnspentTxOutput(w, amount)
	if unspentTxOutputs == nil {
		fmt.Printf("Not enough balance")
		return nil
	}

	fmt.Printf("Unspent transaction outputs %#v", unspentTxOutputs)

	txInputs := w.CreateTxInput(unspentTxOutputs)
	txOutputs := w.CreateTxOutput(to, amount, restOutputBalance)

	for _, txInput := range txInputs {
		fmt.Println("tx input", txInput)
		fmt.Printf("tx input value %#v", txInput)
	}

	tx := &transation.Transation{Id: "", Output: txOutputs, Input: txInputs}
	tx.Id = tx.GetTransationId()

	fmt.Printf("Created transaction id: %s\n", tx.Id)

	sig := tx.Sign(w.GetPrivateKey())

	fmt.Printf("Created transaction sig: %s\n", sig)
	for _, txInput := range tx.Input {
		txInput.Siguature = sig
	}

	return tx
}
