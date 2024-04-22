package http_network

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"tiny_chain/block"
)

type PeerPoint struct {
	Addr string
}

var peers []*PeerPoint

func getBlockChain() []*block.Block {
	return block.GetBlockChain()
}

func getPerers() []*PeerPoint {
	return peers
}

func httpGetBlockChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(getBlockChain())
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	w.Write(data)
}

func httpGetPeers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(getPerers())
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	w.Write(data)
}

func httpAddPeer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newPeer, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	peers = append(peers, &PeerPoint{Addr: string(newPeer)})
	data, err := json.Marshal(peers)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	w.Write(data)
}

func httpMineBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nextBlock := block.MineBlock()
	data, err := json.Marshal(nextBlock)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	w.Write(data)
}

func InitSever() {
	http.HandleFunc("/blockchain", httpGetBlockChain)

	http.HandleFunc("/mineBlock", httpMineBlock)
	http.HandleFunc("/peers", httpGetPeers)
	http.HandleFunc("/addPeer", httpAddPeer)

	http.ListenAndServe(":12001", nil)
}
