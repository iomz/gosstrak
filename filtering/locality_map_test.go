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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lm.ToJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalityMap.ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
