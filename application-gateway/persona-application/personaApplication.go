package main

import (
	"bufio"
	"credit-evaluation/chaincode"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const privateKey = "privateKey" // todo
var receivedDos = make([]chaincode.Document, 0)
var signedDos = make([]chaincode.Document, 0)

func main() {
	http.HandleFunc("/document", handleDocument)
	go func() {
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nHi..." +
			"\nwhat do you want to do?" +
			"\n1. read and sign received documents" +
			"\n2. send signed document to organization" +
			"\n3. get my documents from blockchain")

		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "1": // read and sign document
			// todo: check org signature (VERIFY)
			readDocuments()
		case "2": // send document to organization
			fmt.Println("Hello world")
		case "3": // retrieve documents from blockchain
			fmt.Println("Hello world")
		default:
			fmt.Println("not a valid option!", text)
		}
	}
}

func handleDocument(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))

	var doc chaincode.Document
	err = json.Unmarshal(body, &doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"error":"invalid request"}`))
		if err != nil {
			return
		}
	}

	receivedDos = append(receivedDos, doc)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(calcSignature(body))
	if err != nil {
		return
	}
}

func readDocuments() {
	if len(receivedDos) == 0 {
		fmt.Println("No documents received")
		return
	}

	fmt.Println("which document do you want to read?")
	for i := 0; i < len(receivedDos); i++ {
		fmt.Println(fmt.Sprintf("%d- title: %s", i+1, receivedDos[i].Title))
	}

	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, "\n", "", -1)
	optionNum, err := strconv.Atoi(option)
	if err != nil || optionNum > len(receivedDos) || optionNum < 1 {
		fmt.Println("not a valid option", err)
		return
	}

	// todo: decrypt
	docJson, err := json.Marshal(receivedDos[optionNum-1])
	if err != nil {
		fmt.Println("doc format is malformed", err)
		return
	}
	fmt.Println(string(docJson))

	// todo: sign option

	fmt.Println("press enter to get back to the menu")
	_, _ = reader.ReadString('\n')
}

func calcSignature(document []byte) []byte {
	h := sha256.New()
	h.Write(document)
	return h.Sum(nil)
	// todo
}
