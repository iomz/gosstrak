package filter

import (
	"testing"
)

func TestFilterMap_keys(t *testing.T) {
	tests := []struct {
		name string
		fm   FilterMap
		want []string
	}{
		{"0,8", FilterMap{"0000": "0", "1000": "8"}, []string{"0000", "1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.keys()
			for i := 0; i < len(got); i++ {
				if _, ok := tt.fm[got[i]]; !ok {
					t.Errorf("FilterMap.keys() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
