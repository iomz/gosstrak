// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func TestLocalityMap_ToJSON(t *testing.T) {
	tests := []struct {
		name string
		lm   LocalityMap
		want []byte
	}{
	//{"00: 20, 0011: 10, 0000: 10", LocalityMap{"0011": 10, "00": 20, "0000": 10}, []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lm.ToJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalityMap.ToJSON() = \n%v, want \n%v", string(got), string(tt.want))
			}
		})
	}
}
