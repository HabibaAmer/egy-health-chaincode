package main

import (
	"fmt"

	"chaincode/chaincode"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	// "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	err := shim.Start(new(chaincode.SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	// assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	// if err != nil {
	// 	log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	// }

	// if err := assetChaincode.Start(); err != nil {
	// 	log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	// }
}
