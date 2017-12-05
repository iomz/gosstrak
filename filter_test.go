package main

import (
	"reflect"
	"testing"
)

func Test_makeFilter(t *testing.T) {
	type args struct {
		bs     []rune
		offset int
	}
	tests := []struct {
		name  string
		args  args
		want  []rune
		want1 []rune
	}{
		{"00000000 from 0", args{[]rune("00000000"), 0}, []rune("00000000"), []rune("00000000")},
		{"0000 from 0", args{[]rune("0000"), 0}, []rune("00001111"), []rune("00001111")},
		{"0000 from 4", args{[]rune("0000"), 4}, []rune("11110000"), []rune("11110000")},
		{"0 from 6", args{[]rune("0"), 6}, []rune("11111101"), []rune("11111101")},
		{"000000000 from 1", args{[]rune("000000000"), 1}, []rune("1000000000111111"), []rune("1000000000111111")},
		{"00000000000000000000 from 4", args{[]rune("00000000000000000000"), 4}, []rune("111100000000000000000000"), []rune("111100000000000000000000")},
		{"000000000000000 from 7", args{[]rune("000000000000000"), 7}, []rune("111111100000000000000011"), []rune("111111100000000000000011")},
		//{" from ", args{[]rune(""), 0}, []rune(""), []rune("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := makeFilter(tt.args.bs, tt.args.offset)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("makeFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
