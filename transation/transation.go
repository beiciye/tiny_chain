package transation

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"slices"
	"tiny_chain/blockchain"
)

var coinbaseAmount = 50

type TransationInput struct {
	TxOutputId string
	TxOutIndex int
	Siguature  string
}

type TransationOutput struct {
	Address string
	Amount  int
}

type Transation struct {
	Id     string
	Output []*TransationOutput
	Input  []*TransationInput
}

func (txin *TransationInput) GetTransationInputContent() string {
	// Add your code here

	return fmt.Sprintf("%s%d", txin.TxOutputId, txin.TxOutIndex)
}

func (txout *TransationOutput) GetTransationOutputContent() string {
	// Add your code here

	return fmt.Sprintf("%s%d", txout.Address, txout.Amount)
}

func (t *Transation) GetTransationId() string {
	// Add your code here
	txContent := &bytes.Buffer{}
	for _, txin := range t.Input {
		txContent.WriteString(txin.GetTransationInputContent())
	}

	for _, txout := range t.Output {
		txContent.WriteString(txout.GetTransationOutputContent())
	}

	result := sha256.Sum256(txContent.Bytes())
	return hex.EncodeToString(result[:])
}

func (t *Transation) Sign(privateKey *ecdsa.PrivateKey) string {
	// Add your code here
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, []byte(t.GetTransationId()))
	if err != nil {
		log.Default().Fatalf("Failed to sign the transaction: %v", err)
	}

	return hex.EncodeToString(sig)
}

func (t *Transation) ExecuteTransation() {
	// Add your code here

	// validate coinbase transaction
	if t.IsCoinbaseTx() && t.ValidateCoinbaseTx(blockchain.GetLatestBlock().Index) {
		t.ExecuteCoinbaseTransation()
		return
	}

	valid, msg := t.IsValid()

	if !valid {
		fmt.Printf("Invalid transaction, %s\n", msg)
		return
	}

	aUnspentTxout := GetUnspentTxOuts()
	newUnspentTxOuts := []*UnspentOutput{}

	for idx, txout := range t.Output {
		newUnspentTxOuts = append(newUnspentTxOuts, &UnspentOutput{TransationId: t.Id, TransationOutputIndex: idx, Amount: txout.Amount, Address: txout.Address})
	}

	for _, at := range aUnspentTxout {
		isSpent := slices.ContainsFunc(t.Input, func(txin *TransationInput) bool {
			return at.TransationId == txin.TxOutputId && at.TransationOutputIndex == txin.TxOutIndex
		})
		if !isSpent {
			newUnspentTxOuts = append(newUnspentTxOuts, at)
		}
	}
	UpdateUnspentTxOuts(newUnspentTxOuts)
}

func (t *Transation) ExecuteCoinbaseTransation() {
	// Add your code here
	aUnspentTxout := GetUnspentTxOuts()
	newUnspentTxOuts := []*UnspentOutput{}

	for idx, txout := range t.Output {
		newUnspentTxOuts = append(newUnspentTxOuts, &UnspentOutput{TransationOutputIndex: idx, TransationId: t.Id, Amount: txout.Amount, Address: txout.Address})
	}

	newUnspentTxOuts = append(aUnspentTxout, newUnspentTxOuts...)

	UpdateUnspentTxOuts(newUnspentTxOuts)
}

func (txin *TransationInput) IsValid(tx *Transation) bool {

	fmt.Printf(" TransationInput txin: %#v\n", txin)

	consumedOutput := GetUnspentTxOut(txin.TxOutputId, txin.TxOutIndex)

	fmt.Printf("consumedOutput: %#v\n", consumedOutput)

	data, _ := hex.DecodeString(consumedOutput.Address)

	publicKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		log.Default().Fatalf("Failed to parse public key: %v", err)
	}

	sigData, _ := hex.DecodeString(txin.Siguature)

	return ecdsa.VerifyASN1(publicKey.(*ecdsa.PublicKey), []byte(tx.Id), sigData)
}

func (tx *Transation) IsValid() (bool, string) {

	// validate tx id
	if tx.Id != tx.GetTransationId() {
		return false, "Invalid transaction id"
	}

	// validate tx input
	// {
	// 	for _, txin := range tx.Input {
	// 		if !txin.IsValid(tx) {
	// 			return false, fmt.Sprintf("Invalid transaction input: %s", txin.TxOutputId)
	// 		}
	// 	}
	// }

	// validate tx output
	{
		totalInputAmount := 0
		totalOuptutAmount := 0

		for _, txin := range tx.Input {
			consumedOutput := GetUnspentTxOut(txin.TxOutputId, txin.TxOutIndex)
			if consumedOutput != nil {
				totalInputAmount += consumedOutput.Amount
			}
		}

		for _, txout := range tx.Output {
			totalOuptutAmount += txout.Amount
		}

		if totalInputAmount != totalOuptutAmount {
			return false, "Invalid transaction amount"
		}
	}

	return true, ""
}

func (tx *Transation) IsCoinbaseTx() bool {
	// Add your code here

	if len(tx.Input) == 1 && len(tx.Output) == 1 {
		if tx.Input[0].TxOutputId == "" {
			return true
		}
	}

	return false
}

func (tx *Transation) ValidateCoinbaseTx(blockIndex int) bool {
	// Add your code here

	if tx.Id != tx.GetTransationId() {
		return false
	}

	if len(tx.Input) != 1 {
		return false
	}

	if len(tx.Output) != 1 {
		return false
	}

	if tx.Input[0].TxOutIndex != blockIndex {
		return false
	}

	if tx.Output[0].Amount != coinbaseAmount {
		return false
	}

	return true
}

func CreateCoinbaseTx(address string, blockIndex int) *Transation {
	// Add your code here

	txin := &TransationInput{TxOutputId: "", TxOutIndex: blockIndex, Siguature: ""}
	if (blockIndex % 10) == 0 {
		coinbaseAmount = coinbaseAmount / 2
	}
	txout := &TransationOutput{Address: address, Amount: coinbaseAmount}
	tx := &Transation{Id: "", Output: []*TransationOutput{txout}, Input: []*TransationInput{txin}}
	tx.Id = tx.GetTransationId()

	fmt.Printf("Created coinbase transaction: %s\n", tx.Id)
	return tx
}
