// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_makeByteID(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"0000000011111111", args{"0000000011111111"}, []byte{0, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeByteID(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeByteID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readIDsFromCSV(t *testing.T) {
	type args struct {
		inputFile string
	}
	tests := []struct {
		name string
		args args
		want *[][]byte
	}{
		{"SGTIN-96_3_3_458960468_102_1",
			args{os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/testdata/ids.csv"},
			&[][]byte{{48, 109, 181, 178, 229, 64, 25, 128, 0, 0, 0, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readIDsFromCSV(tt.args.inputFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readIDsFromCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
