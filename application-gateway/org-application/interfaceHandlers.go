package main

import (
	"bufio"
	"bytes"
	"credit-evaluation/chaincode"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func UploadDocument() (chaincode.Document, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("alright. let's input the document. put the data in a json file and give me the address of the file on your computer." +
		"\n(use the template to create the document file.)")
	address, _ := reader.ReadString('\n')
	address = strings.Replace(address, "\n", "", -1)
	data, err := os.ReadFile(address)
	if err != nil {
		fmt.Println("can not read the file", err)
		return chaincode.Document{}, err
	}
	var document chaincode.Document
	err = json.Unmarshal(data, &document)
	if err != nil {
		fmt.Println("document format is invalid.", err)
		return chaincode.Document{}, err
	}

	fmt.Println(string(data))
	// todo: encrypt document
	// todo: sign document

	return document, nil
}

func SendDocumentToUser(localDocuments []chaincode.Document) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("enter the number of the document that you want to put on the blockchain")
	for i := 0; i < len(localDocuments); i++ {
		print(i+1, "- ", "owner: ", localDocuments[i].OwnerID, ", title: ", localDocuments[i].Title, "\n")
	}
	localDocId, _ := reader.ReadString('\n')
	localDocId = strings.Replace(localDocId, "\n", "", -1)
	localDocNumber, err := strconv.Atoi(localDocId)
	if err != nil || localDocNumber > len(localDocuments) || localDocNumber < 1 {
		fmt.Println("number not acceptable.", err)
		return false
	}

	jsonValue, _ := json.Marshal(localDocuments[localDocNumber-1])
	resp, err := http.Post("http://localhost:8082/document", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("could not send document to person", err)
		return false
	}
	defer resp.Body.Close()

	print("document send to person successfully.")
	return true
}

func PutDocumentOnBlockchain(localDocuments []chaincode.Document, orgApplication *OrgApplication) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("enter the number of the document that you want to put on the blockchain")
	for i := 0; i < len(localDocuments); i++ {
		print(i+1, "- ", "owner: ", localDocuments[i].OwnerID, ", title: ", localDocuments[i].Title, "\n")
	}
	localDocId, _ := reader.ReadString('\n')
	localDocId = strings.Replace(localDocId, "\n", "", -1)
	localDocNumber, err := strconv.Atoi(localDocId)
	if err != nil || localDocNumber > len(localDocuments) || localDocNumber < 1 {
		fmt.Println("number not acceptable.", err)
		return false
	}
	docId := orgApplication.CreateDocument(localDocuments[localDocNumber-1])
	fmt.Println("saved doc id", docId)
	return true
}

func GetDocumentById(application *OrgApplication) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("alright. let's input the id of the document.")
	docId, _ := reader.ReadString('\n')
	docId = strings.Replace(docId, "\n", "", -1)

	document, err := application.ReadDocumentByID(docId)
	if err != nil {
		return "", err
	}
	return document, nil
}
