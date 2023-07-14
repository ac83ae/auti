package orgchain

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/ac83ae/auti/benchmark/clolc/internal/constant"
	"github.com/ac83ae/auti/benchmark/timecounter"
	"github.com/ac83ae/auti/internal/clolc/transaction"
)

const (
	channelName    = "mychannel"
	contractType   = "auti-org-chain"
	orgWalletPath  = "wallet"
	orgWalletLabel = "appUser"
	org1MSPid      = "Org1MSP"
)

var (
	fabloFilePath string
	org1CCPPath   string
	org1CREDPath  string
	org1CertPath  string
	org1KeyDir    string
)

func init() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}
	err = os.RemoveAll(orgWalletPath)
	if err != nil {
		log.Fatalf("Error removing wallet directory: %v", err)
	}

	autiOrgChainDir := os.Getenv("AUTI_ORG_CHAIN_DIR")

	fabloFilePath = filepath.Join(autiOrgChainDir, "fablo-target", "fabric-config")

	org1CCPPath = filepath.Join(fabloFilePath, "connection-profiles", "connection-profile-org1.yaml")
	org1CREDPath = filepath.Join(
		fabloFilePath,
		"crypto-config",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)
	org1CertPath = filepath.Join(org1CREDPath, "signcerts", "User1@org1.example.com-cert.pem")
	org1KeyDir = filepath.Join(org1CREDPath, "keystore")
}

func SubmitTX(numTXs int) ([]string, error) {
	lc, err := NewController()
	if err != nil {
		return nil, err
	}
	defer lc.Close()
	dummyTXs := DummyOnChainTransactions(numTXs)
	var txIDs []string
	startTime := time.Now()
	for batch := 0; batch < numTXs; batch += constant.SubmitTXBatchSize {
		right := batch + constant.SubmitTXBatchSize
		if right > numTXs {
			right = numTXs
		}
		for trail := 0; trail < constant.SubmitTXMaxRetries; trail++ {
			batchTXIDs, err := lc.SubmitBatchTXs(dummyTXs[batch:right])
			if err == nil {
				txIDs = append(txIDs, batchTXIDs...)
				break
			}
			log.Printf("Failed to submit batch TXs: %v\n", err)
			log.Printf("Retrying in %v seconds\n", constant.SubmitTXRetryDelaySeconds)
			time.Sleep(constant.SubmitTXRetryDelaySeconds * time.Second)
		}
		if err != nil {
			return nil, err
		}
	}
	elapsedTime := time.Since(startTime)
	timecounter.Print(elapsedTime)
	return txIDs, nil
}

func ReadTX() error {
	f, err := os.Open(constant.OrgChainTXIDLogPath)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var txIDList []string
	for fileScanner.Scan() {
		txIDList = append(txIDList, fileScanner.Text())
	}
	err = f.Close()
	if err != nil {
		return err
	}
	lc, err := NewController()
	if err != nil {
		return err
	}
	defer lc.Close()
	idx := rand.Int() % len(txIDList)
	startTime := time.Now()
	_, err = lc.ReadTX(txIDList[idx])
	elapsedTime := time.Since(startTime)
	timecounter.Print(elapsedTime)
	return err
}

func ReadAllTXsByPage() error {
	lc, err := NewController()
	if err != nil {
		return err
	}
	defer lc.Close()
	var (
		bookmark string
		txList   []*transaction.OrgOnChain
	)
	startTime := time.Now()
	for {
		var (
			pageTXList []*transaction.OrgOnChain
			err        error
		)
		pageTXList, bookmark, err = lc.ReadAllTXsByPage(bookmark)
		if err != nil {
			return err
		}
		txList = append(txList, pageTXList...)
		if bookmark == "" {
			break
		}
	}
	elapsedTime := time.Since(startTime)
	timecounter.Print(elapsedTime)
	return err
}

func SaveTXIDs(txIDs []string) error {
	f, err := os.OpenFile(constant.OrgChainTXIDLogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	for _, id := range txIDs {
		if _, err = f.WriteString(id + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func populateOrgWallet(wallet *gateway.Wallet) error {
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(org1CertPath))
	if err != nil {
		return err
	}

	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(org1KeyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(org1KeyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(org1MSPid, string(cert), string(key))

	return wallet.Put(orgWalletLabel, identity)
}
