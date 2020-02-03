package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io"
	"log"
	"net/http"
	"strconv"
)

type pubkey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type value struct {
	Threshold int      `json:"threshold,string"`
	Pubkeys   []pubkey `json:"pubkeys"`
}

type request struct {
	Type  string `json:"type"`
	Value value  `json:"value"`
}

type response struct {
	Error   string `json:"error,omitempty"`
	Address string `json:"address,omitempty"`
}

func multisigAddress(request request) (string, string) {

	if request.Type != "tendermint/PubKeyMultisigThreshold" {
		return "", "Unexpected signature type, PubKeyMultisigThreshold expected"
	}

	var pubkeys = make([]crypto.PubKey, len(request.Value.Pubkeys))
	for index, element := range request.Value.Pubkeys {

		if element.Type != "tendermint/PubKeySecp256k1" {
			return "", "Unexpected signature type, tendermint/PubKeySecp256k1 expected"
		}

		data, _ := base64.StdEncoding.DecodeString(element.Value)
		var bytes [33]byte
		copy(bytes[:33], data)
		pubkeys[index] = secp256k1.PubKeySecp256k1(bytes)

	}

	var msig = multisig.NewPubKeyMultisigThreshold(request.Value.Threshold, pubkeys)
	var addr = msig.Address()
	return hex.EncodeToString(addr.Bytes()), ""
}

func errorResponse(w http.ResponseWriter, message string) {
	response := &response{
		Error: message,
	}
	res_string, _ := json.Marshal(response)
	io.WriteString(w, string(res_string))
	w.WriteHeader(400)
}

func multisigHandler(w http.ResponseWriter, req *http.Request) {
	var pubkey, ok = req.URL.Query()["pubkey"]

	if !ok || len(pubkey[0]) < 1 {
		errorResponse(w, fmt.Sprintf("Query parameter pubkey not specified"))
		return
	}

	request := request{}
	err := json.Unmarshal([]byte(pubkey[0]), &request)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Error parsing request JSON: %s", err))
		return
	}

	address, message := multisigAddress(request)
	if message != "" {
		errorResponse(w, fmt.Sprintf("Error parsing request JSON: %s", message))
		return
	}

	response := &response{
		Address: address,
	}

	res_string, _ := json.Marshal(response)
	io.WriteString(w, string(res_string))

}

func main() {
	port := flag.Int("port", 8080, "an int")
	bind := flag.String("bind", "", "a string")

	flag.Parse()

	fmt.Printf("Listening on %s port %d ...", *bind, *port)

	http.HandleFunc("/multisig", multisigHandler)
	var err = http.ListenAndServe(*bind+":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Panicln("Server failed starting. Error:", err)
	}

}
