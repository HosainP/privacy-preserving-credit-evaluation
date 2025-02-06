package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"time"
)

// SmartContract provides functions for managing a Document
type SmartContract struct {
	contractapi.Contract
}

// Document describes details of what makes up a document
type Document struct {
	ID      string `json:"ID"`
	OrgID   string `json:"OrgID"`
	OwnerID string `json:"OwnerID"`

	Title string            `json:"Title"`
	Time  time.Time         `json:"Time"`
	Data  map[string]string `json:"Data"`

	OrgSignature   string `json:"OrgSignature"`
	OwnerSignature string `json:"OwnerSignature"`
}

// InitLedger adds the fist document to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	document := Document{
		ID:             "Genesis ID",
		OrgID:          "Genesis organization",
		OwnerID:        "Genesis owner",
		Title:          "Genesis block",
		Time:           time.Time{},
		Data:           make(map[string]string),
		OrgSignature:   "",
		OwnerSignature: "",
	}

	documentJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("genesis-block", documentJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return err
}

func (s *SmartContract) CreateDocument(ctx contractapi.TransactionContextInterface, orgID string, ownerID string, title string, time time.Time, data map[string]string, orgSignature string, ownerSignature string) (string, error) {
	document := Document{
		OrgID:          orgID,
		OwnerID:        ownerID,
		Title:          title,
		Time:           time,
		Data:           data,
		OrgSignature:   orgSignature,
		OwnerSignature: ownerSignature,
	}
	id, err := document.getID()
	if err != nil {
		return "", err
	}
	id = "123456"
	document.ID = id

	documentJSON, err := json.Marshal(document)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, documentJSON)
	if err != nil {
		return "", err
	}

	return document.ID, nil
}

// ReadDocument returns the document stored in the world state with given id.
func (s *SmartContract) ReadDocument(ctx contractapi.TransactionContextInterface, id string) (*Document, error) {
	documentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if documentJSON == nil {
		return nil, fmt.Errorf("the document %s does not exist", id)
	}

	var document Document
	err = json.Unmarshal(documentJSON, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// UpdateDocument updates an existing document in the world state with provided parameters.
func (s *SmartContract) UpdateDocument(ctx contractapi.TransactionContextInterface, id string, title string, time time.Time, data map[string]string, orgSignature string, ownerSignature string) error {
	exists, err := s.DocumentExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the document %s does not exist", id)
	}

	// overwriting original document with new document
	document := Document{
		ID:             id,
		Title:          title,
		Time:           time,
		Data:           data,
		OrgSignature:   orgSignature,
		OwnerSignature: ownerSignature,
	}
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, documentJSON)
}

// DeleteDocument deletes a given document from the world state.
func (s *SmartContract) DeleteDocument(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.DocumentExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the document %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// DocumentExists returns true when document with given ID exists in world state
func (s *SmartContract) DocumentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	documentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return documentJSON != nil, nil
}

// GetAllDocuments returns all documents found in world state
func (s *SmartContract) GetAllDocuments(ctx contractapi.TransactionContextInterface) ([]*Document, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all documents in the chaincode namespace.
	resultIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var documents []*Document
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var document Document
		err = json.Unmarshal(queryResponse.Value, &document)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	return documents, nil
}

// GetAllDocumentsByOwner returns all documents found in world state belonging to an owner
func (s *SmartContract) GetAllDocumentsByOwner(ctx contractapi.TransactionContextInterface, ownerId string) ([]*Document, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all documents in the chaincode namespace.
	resultIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()

	var documents []*Document
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		var document Document
		err = json.Unmarshal(queryResponse.Value, &document)
		if err != nil {
			return nil, err
		}
		if document.OwnerID == ownerId {
			documents = append(documents, &document)
		}
	}

	return documents, nil
}

// getID generates and returns a unique ID for a document
func (d *Document) getID() (string, error) {
	return d.OwnerID + d.Title + d.Time.String(), nil // todo
}

///////////////////////////////////// users /////////////////////////////////////

type User struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`

	CreatedAt   time.Time `json:"Time"`
	DateOfBirth time.Time `json:"DateOfBirth"`

	GovSignature string `json:"GovSignature"`
	PublicKey    string `json:"PublicKey"`
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, userID string, name string, govSignature string, publicKey string) (string, error) {
	user := User{
		ID:           userID,
		Name:         name,
		GovSignature: govSignature,
		PublicKey:    publicKey,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(userID, userJSON)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the user %s does not exist", userID)
	}
	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
