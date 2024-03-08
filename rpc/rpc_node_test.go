package rpc

import (
	"testing"
)

func TestRand2Num(t *testing.T) {
	tests := []struct {
		n       int
		wantErr bool
	}{
		{n: 1, wantErr: false},
		{n: 2, wantErr: false},
		{n: 3, wantErr: false},
		{n: 4, wantErr: false},
		{n: 5, wantErr: false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got1, got2 := Rand2Num(tt.n)
			if tt.wantErr {
				t.Errorf("Rand2Num(%d) = (%d, %d), want error", tt.n, got1, got2)
			} else {
				// Validate that the returned values are within the expected range
				if got1 < 0 || got1 >= tt.n {
					t.Errorf("Rand2Num(%d) = (%d, %d), got1 out of range", tt.n, got1, got2)
				}
				if got2 < 0 || got2 >= tt.n {
					t.Errorf("Rand2Num(%d) = (%d, %d), got2 out of range", tt.n, got1, got2)
				}
			}
		})
	}
}
