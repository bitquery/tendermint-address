package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func multisigHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	port := flag.Int("port", 8080, "an int")
	bind := flag.String("bind", "", "a string")

	flag.Parse()

	fmt.Printf("Listening on port %d binding to %s ", *port, *bind)

	http.HandleFunc("/multisig", multisigHandler)
	var err = http.ListenAndServe(*bind+":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Panicln("Server failed starting. Error:", err)
	}

}
