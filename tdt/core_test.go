// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

import "testing"

func Test_translate(t *testing.T) {
	type args struct {
		afi byte
		uii []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := translate(tt.args.afi, tt.args.uii)
			if (err != nil) != tt.wantErr {
				t.Errorf("translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse6BitEncodedByteSliceToString(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"11000010 00001000: 0 + padding (10 00001000)", args{[]byte{194, 8}}, "0", false},
		{"10010110 01101010: %& + padding (10)", args{[]byte{150, 106}}, "%&", false},
		{"00000100 00100000 11100000 10000010: ABC + padding (100000 10000010)", args{[]byte{4, 32, 224, 130}}, "ABC", false},
		{"11000111 00101100 11110100 10000010: 1234 + padding (10000010)", args{[]byte{199, 44, 244, 130}}, "1234", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse6BitEncodedByteSliceToString(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse6BitEncodedByteSliceToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parse6BitEncodedByteSliceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
