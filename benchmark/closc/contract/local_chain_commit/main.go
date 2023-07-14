package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/ac83ae/auti/closc/contract/local_chain_commit/chaincode"
)

func main() {
	sc := new(chaincode.SmartContract)
	cc, err := contractapi.NewChaincode(sc)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
