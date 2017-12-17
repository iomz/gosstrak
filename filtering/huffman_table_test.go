// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func TestHuffmanTable_sort(t *testing.T) {
	tests := []struct {
		name string
		ht   HuffmanTable
		want HuffmanTable
	}{
		{
			"0011,0000,1111,1100 -> 0000,0011,1100,1111",
			HuffmanTable{
				HuffmanCode{&group{[]*group{}, "0011"}, 3},
				HuffmanCode{&group{[]*group{}, "0000"}, 0},
				HuffmanCode{&group{[]*group{}, "1111"}, 15},
				HuffmanCode{&group{[]*group{}, "1100"}, 12},
			},
			HuffmanTable{
				HuffmanCode{&group{[]*group{}, "0000"}, 0},
				HuffmanCode{&group{[]*group{}, "0011"}, 3},
				HuffmanCode{&group{[]*group{}, "1100"}, 12},
				HuffmanCode{&group{[]*group{}, "1111"}, 15},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ht.sort()
			if !reflect.DeepEqual(tt.ht, tt.want) {
				t.Errorf("HuffmanTable.sort() -> \n%v, want \n%v", tt.ht, tt.want)
			}
		})
	}
}

func TestHuffmanTable_autoencode(t *testing.T) {
	tests := []struct {
		name string
		ht   HuffmanTable
		want *HuffmanTable
	}{
		{
			"Test autoencoding HuffmanTable",
			HuffmanTable{
				HuffmanCode{&group{[]*group{}, "0000"}, 0},
				HuffmanCode{&group{[]*group{}, "0011"}, 3},
				HuffmanCode{&group{[]*group{}, "1100"}, 12},
				HuffmanCode{&group{[]*group{}, "1111"}, 15},
			},
			&HuffmanTable{
				HuffmanCode{&group{[]*group{}, "1111"}, 15},
				HuffmanCode{&group{[]*group{
					{[]*group{}, "1100"},
					{[]*group{
						{[]*group{}, "0011"},
						{[]*group{}, "0000"},
					}, ""}}, ""}, 15,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ht.autoencode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HuffmanTable.sort() -> \n%v, want \n%v", tt.ht, tt.want)
			}
		})
	}
}
