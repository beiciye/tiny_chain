package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"tiny_chain/transation"
)

type PeerPoint struct {
	Addr string
}

var peers []*PeerPoint

func getPerers() []*PeerPoint {
	return peers
}

func InitPeerPoint(addr string) {

}

func (p *PeerPoint) Ping() error {
	// Add your code here

	url := fmt.Sprintf("%s/ping", p.Addr)

	_, err := http.Get(url)

	return err
}

func (p *PeerPoint) broadcastBlockChain() {
	blockchains := getBlockChain()

	data, _ := json.Marshal(blockchains)

	fmt.Printf("Broadcasting blockchain to %s\n", p.Addr)
	_, err := http.Post(p.Addr+"/syncBlockchain", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}

func (p *PeerPoint) broadcastTransaction(tx *transation.Transation) {
	// Add your code here

	data, _ := json.Marshal(tx)

	fmt.Printf("Broadcasting transaction to %s\n", p.Addr)

	_, err := http.Post(p.Addr+"/syncTransaction", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}

func broadcastNewNode() {
	// Add your code here
	for _, peer := range getPerers() {
		peer.broadcastBlockChain()
	}
}

func broadcastTx(tx *transation.Transation) {
	// Add your code here
	for _, peer := range getPerers() {
		peer.broadcastTransaction(tx)
	}
}
