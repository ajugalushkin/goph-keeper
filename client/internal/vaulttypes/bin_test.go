package vaulttypes

import "testing"

func TestBin_String(t *testing.T) {
	type fields struct {
		FileName string
		Size     int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bin{
				FileName: tt.fields.FileName,
				Size:     tt.fields.Size,
			}
			if got := b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBin_Type(t *testing.T) {
	type fields struct {
		FileName string
		Size     int64
	}
	tests := []struct {
		name   string
		fields fields
		want   VaultType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bin{
				FileName: tt.fields.FileName,
				Size:     tt.fields.Size,
			}
			if got := b.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}
