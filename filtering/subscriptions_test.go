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
				"0000": &Info{0, "0", 100, nil},
				"1000": &Info{0, "8", 10, nil},
			},
			[]string{"0000", "1000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.keys(); !reflect.DeepEqual(got, tt.want) {
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
				"0011":         &Info{0, "3", 10, nil},
				"00110011":     &Info{0, "3-3", 5, nil},
				"1111":         &Info{0, "15", 2, nil},
				"00110000":     &Info{0, "3-0", 5, nil},
				"001100110000": &Info{0, "3-3-0", 5, nil},
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
				"0011":         &Info{0, "3", 10, nil},
				"00110011":     &Info{0, "3-3", 5, nil},
				"1111":         &Info{0, "15", 2, nil},
				"00110000":     &Info{0, "3-0", 5, nil},
				"001100110000": &Info{0, "3-3-0", 5, nil},
			},
			Subscriptions{
				"0011": &Info{0, "3", 10,
					&Subscriptions{
						"0000": &Info{4, "3-0", 5, nil},
						"0011": &Info{4, "3-3", 5, &Subscriptions{
							"0000": &Info{8, "3-3-0", 5, nil},
						}},
					},
				},
				"1111": &Info{0, "15", 2, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sub.linkSubset()
			if !reflect.DeepEqual(tt.sub, tt.want) {
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
		sub *Subscriptions
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"15",
			args{
				&Subscriptions{
					"00110000": &Info{0, "3-0", 5, nil},
					"00110011": &Info{0, "3-3", 5, &Subscriptions{
						"0000": &Info{8, "3-3-0", 5, nil},
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
