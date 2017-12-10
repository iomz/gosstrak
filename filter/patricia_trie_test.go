package filter

import (
	"testing"
)

func Test_lcp(t *testing.T) {
	type args struct {
		l []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0001,0010,0011,0000", args{[]string{"0001", "0010", "0011", "0000"}}, "00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcp(tt.args.l); got != tt.want {
				t.Errorf("lcp() = %v, want %v", got, tt.want)
			}
		})
	}
}
