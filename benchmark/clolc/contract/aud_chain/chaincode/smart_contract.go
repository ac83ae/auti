package chaincode

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const pageSize = 10000

// SmartContract provides functions for managing an Transaction.
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of digests to the Aud Chain.
func (s *SmartContract) InitLedger(tci contractapi.TransactionContextInterface) error {
	transactions := []Transaction{
		{
			ID:        "0",
			CipherRes: "0",
			CipherB:   "0",
			CipherC:   "0",
			CipherD:   "0",
		},
	}
	for _, tx := range transactions {
		key, val, err := tx.KeyVal()
		if err != nil {
			return err
		}
		if err := tci.GetStub().PutState(key, val); err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}
	return nil
}

// CreateTX issues a new transaction to the world state with given details.
func (s *SmartContract) CreateTX(ctx contractapi.TransactionContextInterface,
	id, cipherRes, cipherB, cipherC, cipherD string) (string, error) {
	tx := NewTransaction(id, cipherRes, cipherB, cipherC, cipherD)
	key, val, err := tx.KeyVal()
	if err != nil {
		return "", err
	}
	var exists bool
	exists, err = s.TXExists(ctx, key)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("the transaction %s already exists", key)
	}
	return key, ctx.GetStub().PutState(key, val)
}

func (s *SmartContract) CreateBatchTXs(ctx contractapi.TransactionContextInterface, txListJSONString string) ([]string, error) {
	digestListJSONBytes, err := hex.DecodeString(txListJSONString)
	if err != nil {
		return nil, err
	}
	var txList []*Transaction
	err = json.Unmarshal(digestListJSONBytes, &txList)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(txList))
	for i, tx := range txList {
		key, val, err := tx.KeyVal()
		if err != nil {
			return nil, err
		}
		var exists bool
		exists, err = s.TXExists(ctx, key)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("the transaction %s already exists", key)
		}
		if err := ctx.GetStub().PutState(key, val); err != nil {
			return nil, fmt.Errorf("failed to put to world state: %v", err)
		}
		keys[i] = key
	}
	return keys, nil
}

// ReadTX returns the transaction stored in the world state with given id.
func (s *SmartContract) ReadTX(ctx contractapi.TransactionContextInterface,
	key string) (*Transaction, error) {
	tx, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("the transaction %s does not exist", key)
	}
	if err != nil {
		return nil, err
	}
	var txObj Transaction
	err = json.Unmarshal(tx, &txObj)
	if err != nil {
		return nil, err
	}
	return &txObj, nil
}

// DeleteTX deletes a given transaction from the world state.
func (s *SmartContract) DeleteTX(ctx contractapi.TransactionContextInterface, key string) error {
	exists, err := s.TXExists(ctx, key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the transaction %s does not exist", key)
	}
	return ctx.GetStub().DelState(key)
}

// TXExists returns true when transaction with given ID exists in world state.
func (s *SmartContract) TXExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	transactionJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return transactionJSON != nil, nil
}

// ReadAllTXs returns all transactions found in world state.
func (s *SmartContract) ReadAllTXs(ctx contractapi.TransactionContextInterface) (txList []*Transaction, err error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all transactions in the chaincode namespace.
	var iter shim.StateQueryIteratorInterface
	iter, err = ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err = resultsIterator.Close()
	}(iter)

	var response *queryresult.KV
	for iter.HasNext() {
		response, err = iter.Next()
		if err != nil {
			return
		}
		var tx Transaction
		err = json.Unmarshal(response.Value, &tx)
		if err != nil {
			return
		}
		txList = append(txList, &tx)
	}
	return
}

type PageResponse struct {
	Bookmark string         `json:"bookmark"`
	TXs      []*Transaction `json:"txs"`
}

// ReadAllTXsByPage returns the transactions found in world state with pagination.
func (s *SmartContract) ReadAllTXsByPage(ctx contractapi.TransactionContextInterface,
	bookmarkStr string) (pageResponse PageResponse, err error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all transactions in the chaincode namespace.
	var (
		iter         shim.StateQueryIteratorInterface
		responseMeta *peer.QueryResponseMetadata
	)
	if iter, responseMeta, err = ctx.GetStub().GetStateByRangeWithPagination(
		"", "", pageSize, bookmarkStr,
	); err != nil {
		return
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err = resultsIterator.Close()
	}(iter)
	var qr *queryresult.KV
	for iter.HasNext() {
		qr, err = iter.Next()
		if err != nil {
			return
		}
		var tx Transaction
		err = json.Unmarshal(qr.Value, &tx)
		if err != nil {
			return
		}
		pageResponse.TXs = append(pageResponse.TXs, &tx)
	}
	pageResponse.Bookmark = responseMeta.Bookmark
	return
}
