package localchain

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/ac83ae/auti/internal/clolc/transaction"
)

const (
	channelName           = "mychannel"
	contractType          = "auti-local-chain"
	createTXFuncName      = "CreateTX"
	createBatchTXFuncName = "CreateBatchTXs"
	txExistsName          = "TXExists"
	readTXFuncName        = "ReadTX"
	readAllTXFuncName     = "ReadAllTXs"
	readAllTXsByPageName  = "ReadAllTXsByPage"
)

type Controller struct {
	gw *gateway.Gateway
	ct *gateway.Contract
}

// NewController starts a new service instance
func NewController(walletPath, walletLabel, ccpPath string) (*Controller, error) {
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return nil, err
	}
	if !wallet.Exists(walletLabel) {
		if walletLabel == orgWalletLabel {
			if err = populateOrgWallet(wallet); err != nil {
				return nil, err
			}
		} else {
			if err = populateAudWallet(wallet); err != nil {
				return nil, err
			}
		}
	}
	var gw *gateway.Gateway
	if gw, err = gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, walletLabel),
	); err != nil {
		return nil, err
	}
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		return nil, err
	}
	contract := network.GetContract(contractType)
	return &Controller{gw: gw, ct: contract}, nil
}

func (c *Controller) Close() {
	c.gw.Close()
}

func (c *Controller) SubmitTX(tx *transaction.LocalOnChain) (string, error) {
	// log.Println("--> Submit Transaction: Invoke, function that adds a new asset")
	txID, err := c.ct.SubmitTransaction(createTXFuncName,
		tx.CounterParty,
		tx.Commitment,
		tx.Timestamp,
	)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	return string(txID), nil
}

func (c *Controller) SubmitBatchTXs(txList []*transaction.LocalOnChain) ([]string, error) {
	txListJSON, err := json.Marshal(txList)
	if err != nil {
		return nil, err
	}
	txListJSONstr := hex.EncodeToString(txListJSON)
	resBytes, err := c.ct.SubmitTransaction(createBatchTXFuncName, txListJSONstr)
	if err != nil {
		return nil, err
	}
	var txIDList []string
	err = json.Unmarshal(resBytes, &txIDList)
	if err != nil {
		return nil, err
	}
	return txIDList, nil
}

func (c *Controller) TXExists(txID string) (bool, error) {
	resBytes, err := c.ct.EvaluateTransaction(txExistsName, txID)
	if err != nil {
		return false, err
	}
	var result bool
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *Controller) ReadTX(id string) (*transaction.LocalOnChain, error) {
	result, err := c.ct.EvaluateTransaction(readTXFuncName, id)
	if err != nil {
		return nil, err
	}
	var tx transaction.LocalOnChain
	err = json.Unmarshal(result, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (c *Controller) ReadAllTXs() ([]*transaction.LocalOnChain, error) {
	results, err := c.ct.EvaluateTransaction(readAllTXFuncName)
	if err != nil {
		return nil, err
	}
	var txList []*transaction.LocalOnChain
	err = json.Unmarshal(results, &txList)
	if err != nil {
		return nil, err
	}
	return txList, nil
}

type PageResponse struct {
	Bookmark string                      `json:"bookmark"`
	TXs      []*transaction.LocalOnChain `json:"txs"`
}

func (c *Controller) ReadAllTXsByPage(bookmark string) ([]*transaction.LocalOnChain, string, error) {
	results, err := c.ct.EvaluateTransaction(readAllTXsByPageName, bookmark)
	if err != nil {
		return nil, "", err
	}
	var pageResponse PageResponse
	err = json.Unmarshal(results, &pageResponse)
	if err != nil {
		return nil, "", err
	}
	return pageResponse.TXs, pageResponse.Bookmark, nil
}
