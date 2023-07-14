package transaction

import (
	"testing"

	"github.com/ac83ae/auti/internal/crypto"
)

func TestLocalPlain_Hide(t *testing.T) {
	type fields struct {
		CounterParty string
		Amount       int64
		Timestamp    int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test_1",
			fields: fields{
				CounterParty: "test_1",
				Amount:       1,
				Timestamp:    1,
			},
			wantErr: false,
		},
		{
			name: "test_11",
			fields: fields{
				CounterParty: "test_11",
				Amount:       11,
				Timestamp:    11,
			},
			wantErr: false,
		},
		{
			name: "test_123123123123",
			fields: fields{
				CounterParty: "test_123123123123",
				Amount:       123123123123,
				Timestamp:    123123123123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c1 := &LocalPlain{
				CounterParty: tt.fields.CounterParty,
				Amount:       tt.fields.Amount,
				Timestamp:    tt.fields.Timestamp,
			}
			c2 := &LocalPlain{
				CounterParty: tt.fields.CounterParty,
				Amount:       -tt.fields.Amount,
				Timestamp:    tt.fields.Timestamp,
			}
			_, com1, randScalar1, err := c1.Hide()
			if (err != nil) != tt.wantErr {
				t.Errorf("Hide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, com2, randScalar2, err := c2.Hide()
			if (err != nil) != tt.wantErr {
				t.Errorf("Hide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			point1 := crypto.KyberSuite.Point().Mul(randScalar1, crypto.PointH)
			point2 := crypto.KyberSuite.Point().Mul(randScalar2, crypto.PointH)
			com1.Sub(com1, point1)
			com2.Sub(com2, point2)
			com1.Add(com1, com2)
			neutralPoint := crypto.KyberSuite.Point().Null()
			if !com1.Equal(neutralPoint) {
				t.Errorf("Hide() com1 = %v, want %v", com1, neutralPoint)
			}
		})
	}
}
