// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

/*
import (
	"testing"

	"github.com/iomz/gosstrak/tdt"
)

func TestNewSplayTree(t *testing.T) {
	type args struct {
		sub ByteSubscriptions
	}
	tests := []struct {
		name string
		args args
		want *SplayTree
	}{
		{
			"",
			args{
				ByteSubscriptions{
					"0011":         &PartialSubscription{0, "3", ByteSubscriptions{}},
					"00110000":     &PartialSubscription{0, "3-0", ByteSubscriptions{}},
					"00110011":     &PartialSubscription{0, "3-3", ByteSubscriptions{}},
					"001100110000": &PartialSubscription{0, "3-3-0", ByteSubscriptions{}},
					"1111":         &PartialSubscription{0, "15", ByteSubscriptions{}},
				},
			},
			&SplayTree{
				&SplayTreeNode{
					"3",
					NewFilter("0011", 0),
					&SplayTreeNode{
						"3-0",
						NewFilter("0000", 4),
						nil,
						&SplayTreeNode{
							"3-3",
							NewFilter("0011", 4),
							&SplayTreeNode{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
							nil,
						},
					},
					&SplayTreeNode{
						"15",
						NewFilter("1111", 0),
						nil,
						nil,
					},
				},
				tdt.NewCore(),
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

func TestSplayTree_Name(t *testing.T) {
	type fields struct {
		root *SplayTreeNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"SplayTree.Name",
			fields{},
			"SplayTree",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			if got := spt.Name(); got != tt.want {
				t.Errorf("SplayTree.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
