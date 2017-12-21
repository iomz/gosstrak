// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func TestNewComposition(t *testing.T) {
	type args struct {
		filters []*FilterObject
	}
	tests := []struct {
		name string
		args args
		want *Composition
	}{
		{"[8]0011000011 | [6]1100110011",
			args{[]*FilterObject{
				{"0011000011", 10, 8, []byte{48, 255}, []byte{0, 63}, 1, 2},
				{"1100110011", 10, 6, []byte{255, 51}, []byte{252, 0}, 0, 2},
			}},
			&Composition{
				"001100xx",
				8,
				ChildFilters{
					"0011000011": &FilterObject{"0011000011xxxxxx", 16, 8, []byte{48, 255}, []byte{0, 63}, 1, 2},
					"1100110011": &FilterObject{"xxxxxx1100110011", 16, 0, []byte{255, 51}, []byte{252, 0}, 0, 2},
				},
			},
		},
		{"[0]0011000011 | [0]0011001111",
			args{[]*FilterObject{
				{"0011000011", 10, 0, []byte{48, 255}, []byte{0, 63}, 0, 2},
				{"0011001111", 10, 0, []byte{51, 255}, []byte{0, 63}, 0, 2},
			}},
			&Composition{
				"001100xx11xxxxxx",
				0,
				ChildFilters{
					"0011000011": &FilterObject{"00110000xxxxxxxx", 16, 0, []byte{48, 255}, []byte{0, 255}, 0, 2},
					"0011001111": &FilterObject{"00110011xxxxxxxx", 16, 0, []byte{51, 255}, []byte{0, 255}, 0, 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewComposition(tt.args.filters)
			if got.filter != tt.want.filter {
				t.Errorf("*NewComposition().filter = \n%v, want \n%v", got.filter, tt.want.filter)
			}
			if got.offset != tt.want.offset {
				t.Errorf("*NewComposition().offset = \n%v, want \n%v", got.offset, tt.want.offset)
			}
			for k := range tt.want.children {
				t.Logf("Look up key %v in NewComposition().children", k)
				if _, ok := got.children[k]; !ok || got.children[k] == nil {
					t.Errorf("NewComposition().children wants %v key", k)
				}
				if !reflect.DeepEqual(*got.children[k], *tt.want.children[k]) {
					t.Errorf("*NewComposition().children[%v] = \n%v, want \n%v", k, *got.children[k], *tt.want.children[k])
				}
			}
		})
	}
}

func Test_groupSequentialSlice(t *testing.T) {
	type args struct {
		rs []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{"(blank)", args{[]int{}}, [][]int{{}}},
		{"1", args{[]int{1}}, [][]int{{1}}},
		{"1 2", args{[]int{1, 2}}, [][]int{{1, 2}}},
		{"1 2 5", args{[]int{1, 2, 5}}, [][]int{{1, 2}, {5}}},
		{"1 5 6 ", args{[]int{1, 5, 6}}, [][]int{{1}, {5, 6}}},
		{"1 3 4 7 8 10", args{[]int{1, 3, 4, 7, 8, 10}}, [][]int{{1}, {3, 4}, {7, 8}, {10}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := groupSequentialSlice(tt.args.rs); !reflect.DeepEqual(got, tt.want) {
				if len(got[0]) == 0 && len(tt.want[0]) == 0 {
					// pass
				} else {
					t.Errorf("groupSequentialSlice() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_processWildcard(t *testing.T) {
	type args struct {
		b byte
		m byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"11111111 00001111 -> 1111xxxx", args{255, 15}, "1111xxxx"},
		{"11111111 11110000 -> xxxx1111", args{255, 240}, "xxxx1111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processWildcard(tt.args.b, tt.args.m); got != tt.want {
				t.Errorf("processWildcard() = %v, want %v", got, tt.want)
			}
		})
	}
}
