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
		{"0,8", Subscriptions{"0000": &Info{"0", 100}, "1000": &Info{"8", 10}}, []string{"0000", "1000"}},
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
				"0011": &Info{"3", 3},
				"0000": &Info{"0", 0},
				"1111": &Info{"15", 15},
				"1100": &Info{"12", 12},
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
