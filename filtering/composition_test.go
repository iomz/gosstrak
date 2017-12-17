package filtering

import (
	"reflect"
	"testing"
)

func TestNewComposition(t *testing.T) {
	type args struct {
		filters []*Filter
	}
	tests := []struct {
		name string
		args args
		want *Composition
	}{
		{"[8]0011000011 | [6]1100110011",
			args{[]*Filter{
				{"0011000011", 10, 8, []byte{48, 255}, []byte{0, 63}, 1, 2},
				{"1100110011", 10, 6, []byte{255, 51}, []byte{252, 0}, 0, 2},
			}},
			&Composition{
				&Filter{"001100xx", 8, 8, []byte{51}, []byte{3}, 1, 1},
				ChildFilters{
					"0011000011": &Filter{"0011000011xxxxxx", 16, 8, []byte{48, 255}, []byte{0, 63}, 1, 2},
					"1100110011": &Filter{"xxxxxx1100110011", 16, 0, []byte{255, 51}, []byte{252, 0}, 0, 2},
				},
			},
		},
		{"[0]0011000011 | [0]0011001111",
			args{[]*Filter{
				{"0011000011", 10, 0, []byte{48, 255}, []byte{0, 63}, 0, 2},
				{"0011001111", 10, 0, []byte{51, 255}, []byte{0, 63}, 0, 2},
			}},
			&Composition{
				&Filter{"001100xx11xxxxxx", 16, 0, []byte{51, 255}, []byte{3, 63}, 0, 2},
				ChildFilters{
					"0011000011": &Filter{"00110000xxxxxxxx", 16, 0, []byte{48, 255}, []byte{0, 255}, 0, 2},
					"0011001111": &Filter{"00110011xxxxxxxx", 16, 0, []byte{51, 255}, []byte{0, 255}, 0, 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewComposition(tt.args.filters)
			if !reflect.DeepEqual(*got.Filter, *tt.want.Filter) {
				t.Errorf("*NewComposition().Filter = \n%v, want \n%v", *got.Filter, *tt.want.Filter)
			}
			for k := range tt.want.Children {
				t.Logf("Look up key %v in NewComposition().Children", k)
				if _, ok := got.Children[k]; !ok || got.Children[k] == nil {
					t.Errorf("NewComposition().Children wants %v key", k)
				}
				if !reflect.DeepEqual(*got.Children[k], *tt.want.Children[k]) {
					t.Errorf("*NewComposition().Children[%v] = \n%v, want \n%v", k, *got.Children[k], *tt.want.Children[k])
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
