// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
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
				"0000": &Info{"0", 100, &Subscriptions{}},
				"1000": &Info{"8", 10, &Subscriptions{}},
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

func TestSubscriptions_HuffmanTable(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want HuffmanTable
	}{
		{
			"0011,0000,1111,1100",
			Subscriptions{
				"0011": &Info{"3", 3, &Subscriptions{}},
				"0000": &Info{"0", 0, &Subscriptions{}},
				"1111": &Info{"15", 15, &Subscriptions{}},
				"1100": &Info{"12", 12, &Subscriptions{}},
			},
			HuffmanTable{
				HuffmanCode{[]string{"0000"}, 0},
				HuffmanCode{[]string{"0011"}, 3},
				HuffmanCode{[]string{"1100"}, 12},
				HuffmanCode{[]string{"1111"}, 15},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.HuffmanTable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscriptions.HuffmanTable() = \n%v, want \n%v", got, tt.want)
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
				"0011":         &Info{"3", 10, nil},
				"00110011":     &Info{"3-3", 5, nil},
				"1111":         &Info{"15", 2, nil},
				"00110000":     &Info{"3-0", 5, nil},
				"001100110000": &Info{"3-3-0", 5, nil},
			},
			Subscriptions{
				"0011": &Info{"3", 25,
					&Subscriptions{
						"00110000": &Info{"3-0", 5, nil},
						"00110011": &Info{"3-3", 10, &Subscriptions{
							"001100110000": &Info{"3-3-0", 5, nil},
						},
						},
					},
				},
				"1111": &Info{"15", 2, nil},
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
