package filtering

import (
	"testing"
)

func TestMap_keys(t *testing.T) {
	tests := []struct {
		name string
		fm   Map
		want []string
	}{
		{"0,8", Map{"0000": "0", "1000": "8"}, []string{"0000", "1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.keys()
			for i := 0; i < len(got); i++ {
				if _, ok := tt.fm[got[i]]; !ok {
					t.Errorf("Map.keys() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
