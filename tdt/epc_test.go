// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

import (
	"reflect"
	"testing"
)

func Test_getAssetType(t *testing.T) {
	type args struct {
		at string
		pr map[PartitionTableKey]int
	}
	tests := []struct {
		name          string
		args          args
		wantAssetType []rune
	}{
		{
			"1234-8",
			args{at: "1234", pr: GRAI96PartitionTable[8]},
			[]rune{48, 48, 48, 48, 48, 48, 49, 48, 48, 49, 49, 48, 49, 48, 48, 49, 48},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAssetType := getAssetType(tt.args.at, tt.args.pr); !reflect.DeepEqual(gotAssetType, tt.wantAssetType) {
				t.Errorf("getAssetType() = %v, want %v", gotAssetType, tt.wantAssetType)
			}
		})
	}
}

func Test_getCompanyPrefix(t *testing.T) {
	type args struct {
		cp string
		pt PartitionTable
	}
	tests := []struct {
		name              string
		args              args
		wantCompanyPrefix []rune
	}{
		{
			"1234-7",
			args{cp: "1234", pt: SGTIN96PartitionTable},
			[]rune{49, 48, 48, 49, 49, 48, 49, 48, 48, 49, 48},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCompanyPrefix := getCompanyPrefix(tt.args.cp, tt.args.pt); !reflect.DeepEqual(gotCompanyPrefix, tt.wantCompanyPrefix) {
				t.Errorf("getCompanyPrefix() = %v, want %v", gotCompanyPrefix, tt.wantCompanyPrefix)
			}
		})
	}
}

func Test_getExtension(t *testing.T) {
	type args struct {
		e  string
		pr map[PartitionTableKey]int
	}
	tests := []struct {
		name          string
		args          args
		wantExtension []rune
	}{
		{
			"1234-8",
			args{e: "1234", pr: SSCC96PartitionTable[11]},
			[]rune{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 48, 48, 49, 49, 48, 49, 48, 48, 49, 48},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExtension := getExtension(tt.args.e, tt.args.pr); !reflect.DeepEqual(gotExtension, tt.wantExtension) {
				t.Errorf("getExtension() = %v, want %v", gotExtension, tt.wantExtension)
			}
		})
	}
}

func Test_getFilter(t *testing.T) {
	type args struct {
		fv string
	}
	tests := []struct {
		name       string
		args       args
		wantFilter []rune
	}{
		{
			"3",
			args{"3"},
			[]rune{48, 49, 49},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFilter := getFilter(tt.args.fv); !reflect.DeepEqual(gotFilter, tt.wantFilter) {
				t.Errorf("getFilter() = %v, want %v", gotFilter, tt.wantFilter)
			}
		})
	}
}

func Test_getIndivisualAssetReference(t *testing.T) {
	type args struct {
		iar string
		pr  map[PartitionTableKey]int
	}
	tests := []struct {
		name                         string
		args                         args
		wantIndivisualAssetReference []rune
	}{
		{
			"123456-6",
			args{iar: "123456", pr: GIAI96PartitionTable[6]},
			[]rune{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 49, 49, 49, 48, 48, 48, 49, 48, 48, 49, 48, 48, 48, 48, 48, 48},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIndivisualAssetReference := getIndivisualAssetReference(tt.args.iar, tt.args.pr); !reflect.DeepEqual(gotIndivisualAssetReference, tt.wantIndivisualAssetReference) {
				t.Errorf("getIndivisualAssetReference() = %v, want %v", gotIndivisualAssetReference, tt.wantIndivisualAssetReference)
			}
		})
	}
}

func Test_getItemReference(t *testing.T) {
	type args struct {
		ir string
		pr map[PartitionTableKey]int
	}
	tests := []struct {
		name              string
		args              args
		wantItemReference []rune
	}{
		{
			"123-9",
			args{ir: "123", pr: SGTIN96PartitionTable[9]},
			[]rune{48, 48, 48, 48, 48, 48, 48, 49, 49, 49, 49, 48, 49, 49},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotItemReference := getItemReference(tt.args.ir, tt.args.pr); !reflect.DeepEqual(gotItemReference, tt.wantItemReference) {
				t.Errorf("getItemReference() = %v, want %v", gotItemReference, tt.wantItemReference)
			}
		})
	}
}

func Test_getSerial(t *testing.T) {
	type args struct {
		s            string
		serialLength int
	}
	tests := []struct {
		name       string
		args       args
		wantSerial []rune
	}{
		{
			"123456789",
			args{s: "123456789", serialLength: 38},
			[]rune{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 49, 49, 48, 49, 48, 49, 49, 48, 49, 49, 49, 49, 48, 48, 49, 49, 48, 49, 48, 48, 48, 49, 48, 49, 48, 49},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSerial := getSerial(tt.args.s, tt.args.serialLength); !reflect.DeepEqual(gotSerial, tt.wantSerial) {
				t.Errorf("getSerial() = %v, want %v", gotSerial, tt.wantSerial)
			}
		})
	}
}

func TestNewPrefixFilterGIAI96(t *testing.T) {
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
			"GIAI-96_3_1_02283922192_45325296932379",
			args{[]string{"3", "02283922192", "45325296932379"}},
			"0011010001100100000100010000010000111100011000100001010010011100100011110001110010001011000011011",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFilterGIAI96(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFilterGIAI96() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPrefixFilterGIAI96() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPrefixFilterGRAI96(t *testing.T) {
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
			"GRAI-96_3_6_123456_1_1",
			args{[]string{"3", "123456", "1", "1"}},
			"001100110111100001111000100100000000000000000000000000000100000000000000000000000000000000000001",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFilterGRAI96(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFilterGRAI96() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPrefixFilterGRAI96() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPrefixFilterSGTIN96(t *testing.T) {
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
			"SGTIN-96_3_1_00182565127_02",
			args{[]string{"3", "00182565127", "02"}},
			"0011000001100100000000010101110000110111001000001110000010",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFilterSGTIN96(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFilterSGTIN96() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPrefixFilterSGTIN96() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPrefixFilterSSCC96(t *testing.T) {
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
			"SSCC-96_3_1_00039579721",
			args{[]string{"3", "00039579721"}},
			"001100010110010000000000010010110111111000001001001",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFilterSSCC96(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFilterSSCC96() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPrefixFilterSSCC96() = %v, want %v", got, tt.want)
			}
		})
	}
}
