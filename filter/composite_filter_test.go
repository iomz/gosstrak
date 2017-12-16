package filter

import (
	"reflect"
	"testing"
)

func TestNewCompositeFilter(t *testing.T) {
	type args struct {
		filters []*Filter
	}
	filter1 := Filter{"001100110011", 12, 0, []byte{51, 63}, []byte{0, 15}, 0, 2}
	filter2 := Filter{"0011000011", 10, 0, []byte{48, 255}, []byte{0, 63}, 0, 2}
	filter3 := Filter{"0011001111", 10, 0, []byte{51, 255}, []byte{0, 63}, 0, 2}
	tests := []struct {
		name string
		args args
		want *CompositeFilter
	}{
		{"001100110011 0011001111", args{[]*Filter{&filter1, &filter3}}, &CompositeFilter{
			[]*Filter{{
				"00110011",
				8,
				0,
				[]byte{51},
				[]byte{0},
				0,
				1,
			}},
			ElementFilters{"001100110011": []*Filter{
				{
					"0011",
					4,
					8,
					[]byte{63},
					[]byte{15},
					1,
					1,
				}},
				"0011001111": []*Filter{
					{
						"11",
						2,
						8,
						[]byte{255},
						[]byte{63},
						1,
						1,
					}},
			},
		}},
		{"001100110011 0011000011", args{[]*Filter{&filter1, &filter2}}, &CompositeFilter{
			[]*Filter{{
				"001100xx",
				8,
				0,
				[]byte{51},
				[]byte{3},
				0,
				1,
			}},
			ElementFilters{"001100110011": []*Filter{
				{
					"001100110011",
					12,
					0,
					[]byte{51, 63},
					[]byte{0, 15},
					0,
					2,
				}},
				"0011000011": []*Filter{
					{
						"0011000011",
						10,
						0,
						[]byte{48, 255},
						[]byte{0, 63},
						0,
						2,
					}},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCompositeFilter(tt.args.filters)
			for i, f := range tt.want.Filters {
				if !reflect.DeepEqual(*got.Filters[i], *f) {
					t.Errorf("NewCompositeFilter().Filters[%v] = \n%v, want \n%v", i, *got.Filters[i], *f)
				}
			}
			for s, filters := range tt.want.Elements {
				for i, f := range filters {
					if !reflect.DeepEqual(*got.Elements[s][i], *f) {
						t.Errorf("NewCompositeFilter().Elements[%v][%v] = \n%v, want \n%v", s, i, *got.Elements[s][i], *f)
					}
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
