package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/crypto"
)

const amountAmplifier = 100

type LocalPlain struct {
	CounterParty string
	Amount       int64
	Timestamp    int64
}

func NewLocalPlain(counterParty string, amount float64, timestamp int64) *LocalPlain {
	amountInt := int64(amount * amountAmplifier)
	return &LocalPlain{
		CounterParty: counterParty,
		Amount:       amountInt,
		Timestamp:    timestamp,
	}
}

func NewPairLocalPlain(
	fromID, toID string,
	amount float64,
	timestamp int64,
) (*LocalPlain, *LocalPlain) {
	amountInt := int64(amount * amountAmplifier)
	return &LocalPlain{
			CounterParty: toID,
			Amount:       amountInt,
			Timestamp:    timestamp,
		}, &LocalPlain{
			CounterParty: fromID,
			Amount:       -amountInt,
			Timestamp:    timestamp,
		}
}

func (l *LocalPlain) Hide() (hiddenTX *LocalHidden,
	commitment kyber.Point, randScalar kyber.Scalar, err error) {
	sha256Func := sha256.New()
	sha256Func.Write([]byte(l.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, randScalar, err = crypto.PedersenCommit(l.Amount)
	if err != nil {
		return nil, nil, nil, err
	}
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return nil, nil, nil, err
	}
	hiddenTX = NewLocalHidden(
		counterPartyHash,
		commitmentBytes,
		l.Timestamp,
	)
	return
}

type LocalHidden struct {
	CounterParty []byte
	Commitment   []byte
	Timestamp    int64
}

func NewLocalHidden(counterParty, commitment []byte, timestamp int64) *LocalHidden {
	return &LocalHidden{
		CounterParty: counterParty,
		Commitment:   commitment,
		Timestamp:    timestamp,
	}
}

func (h *LocalHidden) ToOnChain() *LocalOnChain {
	timestampStr := strconv.FormatInt(h.Timestamp, 10)
	return NewLocalOnChain(
		hex.EncodeToString(h.CounterParty),
		hex.EncodeToString(h.Commitment),
		timestampStr,
	)
}

type LocalOnChain struct {
	CounterParty string `json:"counter_party"`
	Commitment   string `json:"commitment"`
	Timestamp    string `json:"timestamp"`
}

func NewLocalOnChain(counterParty, commitment, timestamp string) *LocalOnChain {
	return &LocalOnChain{
		CounterParty: counterParty,
		Commitment:   commitment,
		Timestamp:    timestamp,
	}
}

func (l *LocalOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(l)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}

func (l *LocalOnChain) ToHidden() (*LocalHidden, error) {
	counterParty, err := hex.DecodeString(l.CounterParty)
	if err != nil {
		return nil, err
	}
	commitment, err := hex.DecodeString(l.Commitment)
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(l.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewLocalHidden(counterParty, commitment, timestamp), nil
}
