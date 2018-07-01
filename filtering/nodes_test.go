// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func Test_node_equal(t *testing.T) {
	type fields struct {
		filter          string
		offset          int
		notificationURI string
		p               float64
	}
	type args struct {
		want *node
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantGot    *node
		wantWanted *node
	}{
		{
			"node.equal simple test true",
			fields{"1111", 0, "15", 15},
			args{
				&node{"1111", 0, "15", 15},
			},
			true, nil, nil,
		},
		{
			"node.equal simple test false",
			fields{"1010", 0, "10", 10},
			args{
				&node{"1111", 0, "15", 15},
			},
			false,
			&node{"1010", 0, "10", 10},
			&node{"1111", 0, "15", 15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := &node{
				filter:          tt.fields.filter,
				offset:          tt.fields.offset,
				notificationURI: tt.fields.notificationURI,
				p:               tt.fields.p,
			}
			gotOk, gotGot, gotWanted := nd.equal(tt.args.want)
			if gotOk != tt.wantOk {
				t.Errorf("node.equal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotGot, tt.wantGot) {
				t.Errorf("node.equal() gotGot = %v, want %v", gotGot, tt.wantGot)
			}
			if !reflect.DeepEqual(gotWanted, tt.wantWanted) {
				t.Errorf("node.equal() gotWanted = %v, want %v", gotWanted, tt.wantWanted)
			}
		})
	}
}

func TestNodes_sortByP(t *testing.T) {
	tests := []struct {
		name string
		nds  *Nodes
		want *Nodes
	}{
		{
			"0011,0001,1111,1100 -> 0001,0011,1100,1111",
			&Nodes{
				&node{"0011", 0, "3", 3},
				&node{"0001", 0, "1", 1},
				&node{"1111", 0, "15", 15},
				&node{"1100", 0, "12", 12},
			},
			&Nodes{
				&node{"1111", 0, "15", 15},
				&node{"1100", 0, "12", 12},
				&node{"0011", 0, "3", 3},
				&node{"0001", 0, "1", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.nds.sortByP()
			for i, g := range *tt.nds {
				if ok, got, wanted := g.equal((*tt.want)[i]); !ok {
					t.Errorf("Nodes.sortByP() = \n%v, want \n%v", *got, *wanted)
				}
			}
		})
	}
}

func TestNewNodes(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want Nodes
	}{
		{
			"0011,0000,1111,1100",
			args{
				Subscriptions{
					"0011": &Info{0, "3", 3, Subscriptions{}},
					"0001": &Info{0, "1", 1, Subscriptions{}},
					"1111": &Info{0, "15", 15, Subscriptions{}},
					"1100": &Info{0, "12", 12, Subscriptions{}},
				},
			},
			Nodes{
				&node{"1111", 0, "15", 15},
				&node{"1100", 0, "12", 12},
				&node{"0011", 0, "3", 3},
				&node{"0001", 0, "1", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNodes(tt.args.sub)
			for i, w := range tt.want {
				if !reflect.DeepEqual(got[i], w) {
					t.Errorf("NewNodes() = \n%v, want \n%v", got[i], w)
				}
			}
		})
	}
}
