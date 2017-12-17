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
				HuffmanCode{[]string{"0011"}, 3},
				HuffmanCode{[]string{"0000"}, 0},
				HuffmanCode{[]string{"1111"}, 15},
				HuffmanCode{[]string{"1100"}, 12},
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
			tt.ht.sort()
			if !reflect.DeepEqual(tt.ht, tt.want) {
				t.Errorf("HuffmanTable.sort() -> \n%v, want \n%v", tt.ht, tt.want)
			}
		})
	}
}
