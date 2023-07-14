package transaction

import (
	"crypto/sha256"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/internal/crypto"
)

const amountAmplifier = 100

// Plain is the struct for plaintext transaction
type Plain struct {
	Sender    string
	Receiver  string
	Amount    int64
	Counter   uint64
	Timestamp int64
}

// NewPlain creates a new plaintext transaction
func NewPlain(sender, receiver string, amount float64, counter uint64, timestamp int64) *Plain {
	amountInt := int64(amount * amountAmplifier)
	return &Plain{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amountInt,
		Counter:   counter,
		Timestamp: timestamp,
	}
}

func NewPairPlain(sender, receiver string, amount float64, counter uint64, timestamp int64) (*Plain, *Plain) {
	amountInt := int64(amount * amountAmplifier)
	return &Plain{
			Sender:    sender,
			Receiver:  receiver,
			Amount:    amountInt,
			Counter:   counter,
			Timestamp: timestamp,
		}, &Plain{
			Sender:    receiver,
			Receiver:  sender,
			Amount:    -amountInt,
			Counter:   counter,
			Timestamp: timestamp,
		}
}

// Hidden is the struct for hidden transaction
type Hidden struct {
	Sender     []byte
	Receiver   []byte
	Commitment []byte
	Timestamp  int64
}

func (p *Plain) Hide() (*Hidden, kyber.Point, error) {
	// sender hash
	sha256Func := sha256.New()
	sha256Func.Write([]byte(p.Sender))
	senderHash := sha256Func.Sum(nil)
	// receiver hash
	sha256Func.Reset()
	sha256Func.Write([]byte(p.Receiver))
	receiverHash := sha256Func.Sum(nil)
	commitment, hashPoint, err := crypto.PedersonCommitWithHash(
		p.Amount, p.Timestamp, receiverHash, p.Counter,
	)
	if err != nil {
		return nil, nil, err
	}
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	return NewHidden(senderHash, receiverHash, commitmentBytes, p.Timestamp), hashPoint, nil
}

func NewHidden(sender, receiver, commitment []byte, timestamp int64) *Hidden {
	return &Hidden{
		Sender:     sender,
		Receiver:   receiver,
		Commitment: commitment,
		Timestamp:  timestamp,
	}
}
