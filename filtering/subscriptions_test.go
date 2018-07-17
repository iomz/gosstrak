// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"reflect"
	"testing"
)

func TestByteSubscriptions_keys(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want []string
	}{
		{
			"0,8",
			ByteSubscriptions{
				"0000": &PartialSubscription{0, "0", ByteSubscriptions{}},
				"1000": &PartialSubscription{0, "8", ByteSubscriptions{}},
			},
			[]string{"0000", "1000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteSubscriptions.keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSubscriptions_Dump(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want string
	}{
		{
			"Test Dump ByteSubscriptions",
			ByteSubscriptions{
				"0011":         &PartialSubscription{0, "3", ByteSubscriptions{}},
				"00110011":     &PartialSubscription{0, "3-3", ByteSubscriptions{}},
				"1111":         &PartialSubscription{0, "15", ByteSubscriptions{}},
				"00110000":     &PartialSubscription{0, "3-0", ByteSubscriptions{}},
				"001100110000": &PartialSubscription{0, "3-3-0", ByteSubscriptions{}},
			},
			"--0011\n" +
				"--00110000\n" +
				"--00110011\n" +
				"--001100110000\n" +
				"--1111\n",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Dump(); got != tt.want {
				t.Errorf("ByteSubscriptions.Dump() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestByteSubscriptions_linkSubset(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want ByteSubscriptions
	}{
		{
			"Subset linking test for ByteSubscriptions",
			ByteSubscriptions{
				"0011":         &PartialSubscription{0, "3", ByteSubscriptions{}},
				"00110011":     &PartialSubscription{0, "3-3", ByteSubscriptions{}},
				"1111":         &PartialSubscription{0, "15", ByteSubscriptions{}},
				"00110000":     &PartialSubscription{0, "3-0", ByteSubscriptions{}},
				"001100110000": &PartialSubscription{0, "3-3-0", ByteSubscriptions{}},
			},
			ByteSubscriptions{
				"0011": &PartialSubscription{0, "3", ByteSubscriptions{
					"0000": &PartialSubscription{4, "3-0", ByteSubscriptions{}},
					"0011": &PartialSubscription{4, "3-3", ByteSubscriptions{
						"0000": &PartialSubscription{8, "3-3-0", ByteSubscriptions{}},
					}},
				}},
				"1111": &PartialSubscription{0, "15", ByteSubscriptions{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sub.linkSubset()
			if tt.sub.Dump() != tt.want.Dump() {
				t.Errorf("ByteSubscriptions.linkSubset() -> \n%v, want \n%v", tt.sub.Dump(), tt.want.Dump())
			}
		})
	}
}

func TestByteSubscriptions_print(t *testing.T) {
	type args struct {
		indent int
	}
	tests := []struct {
		name       string
		sub        ByteSubscriptions
		args       args
		wantWriter string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			tt.sub.print(writer, tt.args.indent)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("ByteSubscriptions.print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestByteSubscriptions_Clone(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want ByteSubscriptions
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteSubscriptions.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSubscriptions_Keys(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteSubscriptions.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSubscriptions_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		sub     *ByteSubscriptions
		want    []byte
		wantErr bool
	}{
		{
			"nil marshal sub",
			&ByteSubscriptions{},
			[]byte{3, 4, 0, 0},
			false,
		},
		{
			"simple marshal sub",
			&ByteSubscriptions{
				"010101": &PartialSubscription{Offset: 0, ReportURI: "hoge", Subset: ByteSubscriptions{}},
				"1010": &PartialSubscription{Offset: 0, ReportURI: "foo", Subset: ByteSubscriptions{
					"11": &PartialSubscription{Offset: 4, ReportURI: "bar", Subset: ByteSubscriptions{}}}},
			},
			[]byte{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sub.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("ByteSubscriptions.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("ByteSubscriptions.MarshalBinary() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestLoadFiltersFromCSVFile(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want ByteSubscriptions
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadFiltersFromCSVFile(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFiltersFromCSVFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
