package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"go.dedis.ch/kyber/v3"
)

func PedersenCommit(amount int64) (kyber.Point, kyber.Scalar, error) {
	amountScalar, err := amountToScalar(amount)
	if err != nil {
		return nil, nil, err
	}
	commitment := KyberSuite.Point().Mul(amountScalar, PointG)
	randScalar := KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
	randPoint := KyberSuite.Point().Mul(randScalar, PointH)
	commitment.Add(commitment, randPoint)
	return commitment, randScalar, nil
}

func computeHashPoint(timestamp int64, receiverHash []byte, counter uint64) (kyber.Point, error) {
	// concatenated bytes for calculating the commitment
	timestampByte, err := int64ToBytes(timestamp)
	if err != nil {
		return nil, err
	}
	counterByte, err := uint64ToBytes(counter)
	if err != nil {
		return nil, err
	}
	concatBytes := append(timestampByte, receiverHash...)
	concatBytes = append(concatBytes, counterByte...)
	// calculate the hash of the concatenated bytes
	sha256Func := sha256.New()
	sha256Func.Write(concatBytes)
	concatByteHash := sha256Func.Sum(nil)
	hashScalar := KyberSuite.Scalar().SetBytes(concatByteHash)
	hashPoint := KyberSuite.Point().Mul(hashScalar, PointH)
	return hashPoint, nil
}

func PedersonCommitWithHash(amount, timestamp int64,
	receiverHash []byte, counter uint64) (kyber.Point, kyber.Point, error) {
	amountScalar, err := amountToScalar(amount)
	if err != nil {
		return nil, nil, err
	}
	commitment := KyberSuite.Point().Mul(amountScalar, PointG)
	hashPoint, err := computeHashPoint(timestamp, receiverHash, counter)
	if err != nil {
		return nil, nil, err
	}
	commitment.Add(commitment, hashPoint)
	return commitment, hashPoint, nil
}

func amountToScalar(amount int64) (kyber.Scalar, error) {
	positive := true
	if amount < 0 {
		amount = -amount
		positive = false
	}
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, err
	}
	amountScalar := KyberSuite.Scalar().SetBytes(amountBytes)
	if !positive {
		amountScalar.Neg(amountScalar)
	}
	return amountScalar, nil
}

func int64ToBytes(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, i); err != nil {
		return nil, fmt.Errorf("int64ToBytes: %v", err)
	}
	return buf.Bytes(), nil
}

func bytesToInt64(b []byte) (int64, error) {
	buf := bytes.NewReader(b)
	var i int64
	if err := binary.Read(buf, binary.BigEndian, &i); err != nil {
		return 0, fmt.Errorf("bytesToInt64: %v", err)
	}
	return i, nil
}

func uint64ToBytes(i uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, i); err != nil {
		return nil, fmt.Errorf("uint64ToBytes: %v", err)
	}
	return buf.Bytes(), nil
}

func bytesToUint64(b []byte) (uint64, error) {
	buf := bytes.NewReader(b)
	var i uint64
	if err := binary.Read(buf, binary.BigEndian, &i); err != nil {
		return 0, fmt.Errorf("bytesToUint64: %v", err)
	}
	return i, nil
}
