package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	Threshold string   `json:"threshold"`
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

func multisigAddress(request) string {
	return "qqq"
}

func error(w http.ResponseWriter, message string) {
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
		error(w, fmt.Sprintf("Query parameter pubkey not specified"))
		return
	}

	request := request{}
	err := json.Unmarshal([]byte(pubkey[0]), &request)
	if err != nil {
		error(w, fmt.Sprintf("Error parsing request JSON: %s", err))
		return
	}

	address := multisigAddress(request)
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
