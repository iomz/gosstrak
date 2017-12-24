// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func Test_entry_equal(t *testing.T) {
	type fields struct {
		filter          string
		offset          int
		notificationURI string
		p               float64
	}
	type args struct {
		want *entry
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantGot    *entry
		wantWanted *entry
	}{
		{
			"entry.equal simple test true",
			fields{"1111", 0, "15", 15},
			args{
				&entry{"1111", 0, "15", 15},
			},
			true, nil, nil,
		},
		{
			"entry.equal simple test false",
			fields{"1010", 0, "10", 10},
			args{
				&entry{"1111", 0, "15", 15},
			},
			false,
			&entry{"1010", 0, "10", 10},
			&entry{"1111", 0, "15", 15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ent := &entry{
				filter:          tt.fields.filter,
				offset:          tt.fields.offset,
				notificationURI: tt.fields.notificationURI,
				p:               tt.fields.p,
			}
			gotOk, gotGot, gotWanted := ent.equal(tt.args.want)
			if gotOk != tt.wantOk {
				t.Errorf("entry.equal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotGot, tt.wantGot) {
				t.Errorf("entry.equal() gotGot = %v, want %v", gotGot, tt.wantGot)
			}
			if !reflect.DeepEqual(gotWanted, tt.wantWanted) {
				t.Errorf("entry.equal() gotWanted = %v, want %v", gotWanted, tt.wantWanted)
			}
		})
	}
}

func TestHuffmanCodes_sortByP(t *testing.T) {
	tests := []struct {
		name string
		hc   *HuffmanCodes
		want *HuffmanCodes
	}{
		{
			"0011,0001,1111,1100 -> 0001,0011,1100,1111",
			&HuffmanCodes{
				&entry{"0011", 0, "3", 3},
				&entry{"0001", 0, "1", 1},
				&entry{"1111", 0, "15", 15},
				&entry{"1100", 0, "12", 12},
			},
			&HuffmanCodes{
				&entry{"1111", 0, "15", 15},
				&entry{"1100", 0, "12", 12},
				&entry{"0011", 0, "3", 3},
				&entry{"0001", 0, "1", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hc.sortByP()
			for i, g := range *tt.hc {
				if ok, got, wanted := g.equal((*tt.want)[i]); !ok {
					t.Errorf("HuffmanCodes.sortByP() = \n%v, want \n%v", *got, *wanted)
				}
			}
		})
	}
}

func TestNewHuffmanCodes(t *testing.T) {
	type args struct {
		sub *Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *HuffmanCodes
	}{
		{
			"0011,0000,1111,1100",
			args{
				&Subscriptions{
					"0011": &Info{"3", 3, &Subscriptions{}},
					"0001": &Info{"1", 1, &Subscriptions{}},
					"1111": &Info{"15", 15, &Subscriptions{}},
					"1100": &Info{"12", 12, &Subscriptions{}},
				},
			},
			&HuffmanCodes{
				&entry{"1111", 0, "15", 15},
				&entry{"1100", 0, "12", 12},
				&entry{"0011", 0, "3", 3},
				&entry{"0001", 0, "1", 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHuffmanCodes(tt.args.sub)
			for i, w := range *tt.want {
				if !reflect.DeepEqual(*(*got)[i], *w) {
					t.Errorf("NewHuffmanCodes() = \n%v, want \n%v", *(*got)[i], *w)
				}
			}
		})
	}
}
