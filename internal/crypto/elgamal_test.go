package crypto

import (
	"testing"
)

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		amount  int64
		wantErr bool
	}{
		{
			name:    "test_1",
			amount:  1,
			wantErr: false,
		},
		{
			name:    "test_11",
			amount:  11,
			wantErr: false,
		},
		{
			name:    "test_123123123123",
			amount:  123123123123,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, publicKey, err := KeyGen()
			if err != nil {
				t.Errorf("KeyGen() error = %v", err)
				return
			}
			var cipherText *CipherText
			cipherText, err = Encrypt(publicKey, tt.amount)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}
			var plainText int64
			plainText, err = Decrypt(privateKey, cipherText)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if plainText != tt.amount {
				t.Errorf("Decrypt() plainText = %v, amount %v", plainText, tt.amount)
			}
		})
	}
}

func Test_int64ToBytes(t *testing.T) {
	tests := []struct {
		name    string
		i       int64
		wantErr bool
	}{
		{
			name:    "test_0",
			i:       0,
			wantErr: false,
		},
		{
			name:    "test_10",
			i:       10,
			wantErr: false,
		},
		{
			name:    "test_-10",
			i:       -10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := int64ToBytes(tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("int64ToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			reversed, err := bytesToInt64(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("bytesToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reversed != tt.i {
				t.Errorf("bytesToInt64() reversed = %v, i %v", reversed, tt.i)
			}
		})
	}
}
