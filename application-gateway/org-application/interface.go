package main

import (
	"bufio"
	"credit-evaluation/chaincode"
	"fmt"
	"os"
	"strings"
)

func main() {
	// todo: run a goroutine to receive docs from people

	reader := bufio.NewReader(os.Stdin)
	orgApplication, _ := NewOrgApplication()
	localDocuments := make([]chaincode.Document, 0)
	//signedDocuments := make([]chaincode.Document, 0)

	for {
		fmt.Println("\nHi..." +
			"\nwhat do you want to do?" +
			"\n1. upload a document" +
			"\n2. send a document to person" +
			"\n3. put a document on blockchain" +
			"\n4. get a document from blockchain" +
			"\n5. get all documents of a person from blockchain" +
			"\n6. get all documents from blockchain")

		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "1": // upload a document
			document, err := UploadDocument()
			if err != nil {
				fmt.Println(err)
				continue
			}
			localDocuments = append(localDocuments, document)

		case "2": // send a document to person
			SendDocumentToUser(localDocuments)

		case "3": // put a document on blockchain
			PutDocumentOnBlockchain(localDocuments, orgApplication)

		case "4": // get a document from blockchain
			document, err := GetDocumentById(orgApplication)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(document)

		case "5": // get all documents of a person
			fmt.Println("5")

		case "6": // get all documents from blockchain
			fmt.Println("6")

		default:
			fmt.Println("not a valid option!", text)
		}
	}
}

type localDocument struct {
	OrgId        string            `json:"OrgId"`
	OwnerId      string            `json:"OwnerId"`
	Title        string            `json:"Title"`
	Data         map[string]string `json:"Data"`
	OrgSignature string            `json:"OrgSignature"`
}
