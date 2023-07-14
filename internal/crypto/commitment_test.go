package crypto

import (
	"testing"
)

func TestPedersenCommit(t *testing.T) {
	tests := []struct {
		name    string
		amount  int64
		wantErr bool
	}{
		{
			name:    "test 10",
			amount:  1,
			wantErr: false,
		},
		{
			name:    "test 100",
			amount:  1,
			wantErr: false,
		},
		{
			name:    "test 1000",
			amount:  1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point1, randScalar1, err := PedersenCommit(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("PedersenCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			point2, randScalar2, err := PedersenCommit(-tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("PedersenCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			randPoint1 := KyberSuite.Point().Mul(randScalar1, PointH)
			randPoint2 := KyberSuite.Point().Mul(randScalar2, PointH)
			point1.Sub(point1, randPoint1)
			point2.Sub(point2, randPoint2)
			point1.Add(point1, point2)
			neutralPoint := KyberSuite.Point().Null()
			if !point1.Equal(neutralPoint) {
				t.Errorf("Amount - Amount = %v, want %v", point1, neutralPoint)
			}
		})
	}
}

func Test_amountToScalar(t *testing.T) {
	tests := []struct {
		name    string
		amount  int64
		wantErr bool
	}{
		{
			name:    "test_10",
			amount:  10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := amountToScalar(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("amountToScalar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got2, err := amountToScalar(-tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("amountToScalar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			point1 := KyberSuite.Point().Mul(got1, PointG)
			point2 := KyberSuite.Point().Mul(got2, PointG)
			point1.Add(point1, point2)
			if !point1.Equal(PointG) {
				t.Errorf("amount - amount = %v, want %v", point1, PointG)
			}
		})
	}
}
