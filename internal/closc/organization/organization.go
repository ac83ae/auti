package organization

import (
	"crypto/sha256"
	"encoding/hex"
)

var sha256Func = sha256.New()

type TypeID string

type Organization struct {
	ID     TypeID
	IDHash string
}

func New(id string) *Organization {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	idHash := hex.EncodeToString(sha256Func.Sum(nil))
	org := &Organization{
		ID:     TypeID(id),
		IDHash: idHash,
	}
	return org
}
