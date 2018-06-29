// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func TestSplayTree_AnalyzeLocality(t *testing.T) {
	type fields struct {
		root *OptimalBST
	}
	type args struct {
		id   []byte
		path string
		lm   *LocalityMap
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			spt.AnalyzeLocality(tt.args.id, tt.args.path, tt.args.lm)
		})
	}
}

func TestSplayTree_Dump(t *testing.T) {
	type fields struct {
		root *OptimalBST
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			if got := spt.Dump(); got != tt.want {
				t.Errorf("SplayTree.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplayTree_MarshalBinary(t *testing.T) {
	type fields struct {
		root *OptimalBST
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			got, err := spt.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("SplayTree.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplayTree.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplayTree_Search(t *testing.T) {
	type fields struct {
		root *OptimalBST
	}
	type args struct {
		id []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			if got := spt.Search(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplayTree.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplayTree_UnmarshalBinary(t *testing.T) {
	type fields struct {
		root *OptimalBST
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			if err := spt.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SplayTree.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptimalBST_splaySearch(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
	}
	type args struct {
		spt    *SplayTree
		parent *OptimalBST
		id     []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if got := obst.splaySearch(tt.args.spt, tt.args.parent, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OptimalBST.splaySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSplayTree(t *testing.T) {
	type args struct {
		sub *Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *SplayTree
	}{
		{
			"",
			args{
				&Subscriptions{
					"0011":         &Info{0, "3", 10, nil},
					"00110000":     &Info{0, "3-0", 5, nil},
					"00110011":     &Info{0, "3-3", 5, nil},
					"001100110000": &Info{0, "3-3-0", 5, nil},
					"1111":         &Info{0, "15", 2, nil},
				},
			},
			&SplayTree{
				&OptimalBST{
					"3",
					NewFilter("0011", 0),
					&OptimalBST{
						"3-0",
						NewFilter("0000", 4),
						nil,
						&OptimalBST{
							"3-3",
							NewFilter("0011", 4),
							&OptimalBST{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
							nil,
						},
					},
					&OptimalBST{
						"15",
						NewFilter("1111", 0),
						nil,
						nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSplayTree(tt.args.sub); got.Dump() != tt.want.Dump() {
				t.Errorf("NewSplayTree() = \n%v, want \n%v", got.Dump(), tt.want.Dump())
			}
		})
	}
}
