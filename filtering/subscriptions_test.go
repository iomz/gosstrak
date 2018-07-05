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

func TestSubscriptions_keys(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want []string
	}{
		{
			"0,8",
			Subscriptions{
				"0000": &Info{0, "0", 100, Subscriptions{}},
				"1000": &Info{0, "8", 10, Subscriptions{}},
			},
			[]string{"0000", "1000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscriptions.keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriptions_Dump(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want string
	}{
		{
			"Test Dump Subscriptions",
			Subscriptions{
				"0011":         &Info{0, "3", 10, Subscriptions{}},
				"00110011":     &Info{0, "3-3", 5, Subscriptions{}},
				"1111":         &Info{0, "15", 2, Subscriptions{}},
				"00110000":     &Info{0, "3-0", 5, Subscriptions{}},
				"001100110000": &Info{0, "3-3-0", 5, Subscriptions{}},
			},
			"--0011 10.000000\n" +
				"--00110000 5.000000\n" +
				"--00110011 5.000000\n" +
				"--001100110000 5.000000\n" +
				"--1111 2.000000\n",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Dump(); got != tt.want {
				t.Errorf("Subscriptions.Dump() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestSubscriptions_linkSubset(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want Subscriptions
	}{
		{
			"Subset linking test for Subscriptions",
			Subscriptions{
				"0011":         &Info{0, "3", 10, Subscriptions{}},
				"00110011":     &Info{0, "3-3", 5, Subscriptions{}},
				"1111":         &Info{0, "15", 2, Subscriptions{}},
				"00110000":     &Info{0, "3-0", 5, Subscriptions{}},
				"001100110000": &Info{0, "3-3-0", 5, Subscriptions{}},
			},
			Subscriptions{
				"0011": &Info{0, "3", 10, Subscriptions{
					"0000": &Info{4, "3-0", 5, Subscriptions{}},
					"0011": &Info{4, "3-3", 5, Subscriptions{
						"0000": &Info{8, "3-3-0", 5, Subscriptions{}},
					}},
				}},
				"1111": &Info{0, "15", 2, Subscriptions{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sub.linkSubset()
			if tt.sub.Dump() != tt.want.Dump() {
				t.Errorf("Subscriptions.linkSubset() -> \n%v, want \n%v", tt.sub.Dump(), tt.want.Dump())
			}
		})
	}
}

func TestSubscriptions_print(t *testing.T) {
	type args struct {
		indent int
	}
	tests := []struct {
		name       string
		sub        Subscriptions
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
				t.Errorf("Subscriptions.print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func Test_recalculateEntropyValue(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"15",
			args{
				Subscriptions{
					"00110000": &Info{0, "3-0", 5, Subscriptions{}},
					"00110011": &Info{0, "3-3", 5, Subscriptions{
						"0000": &Info{8, "3-3-0", 5, Subscriptions{}},
					}},
				},
			},
			float64(15),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := recalculateEntropyValue(tt.args.sub); got != tt.want {
				t.Errorf("recalculateEntropyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_MarshalBinary(t *testing.T) {
	type fields struct {
		Offset          int
		NotificationURI string
		EntropyValue    float64
		Subset          Subscriptions
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &Info{
				Offset:          tt.fields.Offset,
				NotificationURI: tt.fields.NotificationURI,
				EntropyValue:    tt.fields.EntropyValue,
				Subset:          tt.fields.Subset,
			}
			got, err := info.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("Info.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Info.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_UnmarshalBinary(t *testing.T) {
	type fields struct {
		Offset          int
		NotificationURI string
		EntropyValue    float64
		Subset          Subscriptions
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &Info{
				Offset:          tt.fields.Offset,
				NotificationURI: tt.fields.NotificationURI,
				EntropyValue:    tt.fields.EntropyValue,
				Subset:          tt.fields.Subset,
			}
			if err := info.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Info.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscriptions_Clone(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want Subscriptions
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscriptions.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriptions_Keys(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscriptions.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriptions_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		sub     *Subscriptions
		want    []byte
		wantErr bool
	}{
		{
			"nil marshal sub",
			&Subscriptions{},
			[]byte{3, 4, 0, 0},
			false,
		},
		{
			"simple marshal sub",
			&Subscriptions{
				"010101": &Info{Offset: 0, NotificationURI: "hoge", EntropyValue: 0.3, Subset: Subscriptions{}},
				"1010": &Info{Offset: 0, NotificationURI: "foo", EntropyValue: 0.7, Subset: Subscriptions{
					"11": &Info{Offset: 4, NotificationURI: "bar", EntropyValue: 0.2, Subset: Subscriptions{}}}},
			},
			[]byte{3, 4, 0, 4, 9, 12, 0, 6, 48, 49, 48, 49, 48, 49, 10, 255, 139, 6, 1, 2, 255, 142, 0, 0, 0, 25, 255, 143, 6, 1, 1, 13, 83, 117, 98, 115, 99, 114, 105, 112, 116, 105, 111, 110, 115, 1, 255, 144, 0, 0, 0, 65, 255, 140, 0, 61, 3, 4, 0, 0, 7, 12, 0, 4, 104, 111, 103, 101, 11, 8, 0, 248, 51, 51, 51, 51, 51, 51, 211, 63, 25, 255, 143, 6, 1, 1, 13, 83, 117, 98, 115, 99, 114, 105, 112, 116, 105, 111, 110, 115, 1, 255, 144, 0, 0, 0, 10, 255, 139, 6, 1, 2, 255, 142, 0, 0, 0, 7, 12, 0, 4, 49, 48, 49, 48, 64, 255, 140, 0, 60, 3, 4, 0, 0, 6, 12, 0, 3, 102, 111, 111, 11, 8, 0, 248, 102, 102, 102, 102, 102, 102, 230, 63, 25, 255, 143, 6, 1, 1, 13, 83, 117, 98, 115, 99, 114, 105, 112, 116, 105, 111, 110, 115, 1, 255, 144, 0, 0, 0, 10, 255, 139, 6, 1, 2, 255, 142, 0, 0, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sub.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscriptions.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscriptions.MarshalBinary() = %v, want %v", got, tt.want)
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
		want Subscriptions
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
