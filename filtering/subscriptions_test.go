package filtering

import (
	"testing"
)

func TestMap_keys(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want []string
	}{
		{"0,8", Subscriptions{"0000": &Info{"0", 100}, "1000": &Info{"8", 10}}, []string{"0000", "1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sub.keys()
			for i := 0; i < len(got); i++ {
				if _, ok := tt.sub[got[i]]; !ok {
					t.Errorf("Map.keys() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
