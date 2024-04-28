package network

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tiny_chain/blockchain"
	"tiny_chain/transation"
	"tiny_chain/wallet"
)

func getBlockChain() []*blockchain.Block {
	return blockchain.GetBlockChain()
}

func httpGetBlockChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(getBlockChain())
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}

	w.Write(data)
}

func httpGetPeers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(getPerers())
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	w.Write(data)
}

func httpAddPeer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	var newPeer *PeerPoint
	json.Unmarshal(resData, &newPeer)

	peers = append(peers, &PeerPoint{Addr: newPeer.Addr})
	data, err := json.Marshal(peers)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	w.Write(data)
}

func httpMineBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nextBlock := blockchain.MineBlock()
	data, err := json.Marshal(nextBlock)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	address := wallet.GetWallet().GetMyAddress()
	tx := transation.CreateCoinbaseTx(address, nextBlock.Index)

	tx.ExecuteCoinbaseTransation()
	broadcastNewNode()
	broadcastTx(tx)
	w.Write(data)
}

func httpOnBlockChainReceived(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Received blockchain")
	resData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	var receivedChain []*blockchain.Block
	json.Unmarshal(resData, &receivedChain)
	if receivedChain == nil {
		return
	}
	if len(receivedChain) > len(getBlockChain()) {
		if blockchain.ValidateBlockChain(receivedChain) {
			blockchain.ReplaceChain(receivedChain)
		}
	}
}

func httpOnTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}

	fmt.Printf("Received transaction %s\n", string(resData))

	var tx transation.Transation
	json.Unmarshal(resData, &tx)

	tx.ExecuteTransation()
	w.Write([]byte("Transaction executed"))
}

type PerformTx struct {
	To     string
	Amount int
}

func httpPerformTransation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	var userTx *PerformTx

	json.Unmarshal(resData, &userTx)

	wallet := wallet.GetWallet()
	tx := wallet.CreateTransation(userTx.To, userTx.Amount)
	tx.ExecuteTransation()
	broadcastTx(tx)

	w.Write([]byte("Transaction executed"))
}

func httpGetMyBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	balance := wallet.GetWallet().GetBanlance()
	data, err := json.Marshal(balance)
	if err != nil {
		log.Default().Fatalf("Error: %s", err)
	}
	w.Write(data)
}

func InitSever() {
	http.HandleFunc("/blockchain", httpGetBlockChain)

	http.HandleFunc("/mineBlock", httpMineBlock)
	http.HandleFunc("/peers", httpGetPeers)
	http.HandleFunc("/addPeer", httpAddPeer)

	http.HandleFunc("/syncBlockchain", httpOnBlockChainReceived)
	http.HandleFunc("/syncTransaction", httpOnTransaction)

	http.HandleFunc("/performTransaction", httpPerformTransation)
	http.HandleFunc("/getMyBalance", httpGetMyBalance)

	args := os.Args

	port := fmt.Sprintf(":%s", args[2])

	fmt.Printf("Listening on port %v\n", port)

	http.ListenAndServe(port, nil)
}
