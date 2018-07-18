// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

import (
	"reflect"
	"testing"
)

func Test_getISO6346CD(t *testing.T) {
	type args struct {
		cn string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getISO6346CD(tt.args.cn)
			if (err != nil) != tt.wantErr {
				t.Errorf("getISO6346CD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getISO6346CD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pad6BitEncodingRuneSlice(t *testing.T) {
	type args struct {
		bs []rune
	}
	tests := []struct {
		name  string
		args  args
		want  []rune
		want1 int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := pad6BitEncodingRuneSlice(tt.args.bs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pad6BitEncodingRuneSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("pad6BitEncodingRuneSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewPrefixFilterISO17363(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"ISO17363_7B_AAJ",
			args{[]string{"7B", "AAJ"}},
			"110111000010000001000001001010",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFilterISO17363(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFilterISO17363() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPrefixFilterISO17363() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestNewPrefixFilterISO17365(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"ISO17365_25S_UN_ABC_0THANK0YOU0FOR0READING0THIS1",
			args{[]string{"25S", "UN", "ABC", "0THANK0YOU0FOR0READING0THIS1"}},
			"110010110101010011010101001110000001000010000011110000010100001000000001001110001011110000011001001111010101110000000110001111010010110000010010000101000001000100001001001110000111110000010100001000001001010011110001",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFilterISO17365(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFilterISO17365() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPrefixFilterISO17365() = %v, want %v", got, tt.want)
			}
		})
	}
}
