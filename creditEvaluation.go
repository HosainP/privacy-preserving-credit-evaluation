/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"credit-evaluation/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	creditEvaluationChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating credit-evaluation chaincode: %v", err)
	}

	if err := creditEvaluationChaincode.Start(); err != nil {
		log.Panicf("Error starting credit-evaluation chaincode: %v", err)
	}
}
