// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"reflect"
	"testing"
)

func TestHuffmanCodes_Dump(t *testing.T) {
	tests := []struct {
		name string
		ht   *HuffmanCodes
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ht.Dump(); got != tt.want {
				t.Errorf("HuffmanCodes.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_print(t *testing.T) {
	type fields struct {
		filter string
		p      float64
		bit    rune
		code   int
		group  []*entry
	}
	type args struct {
		indent int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ent := &entry{
				filter: tt.fields.filter,
				p:      tt.fields.p,
				bit:    tt.fields.bit,
				code:   tt.fields.code,
				group:  tt.fields.group,
			}
			writer := &bytes.Buffer{}
			ent.print(writer, tt.args.indent)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("entry.print() = %v, want %v", gotWriter, tt.wantWriter)
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
				&entry{"0011", 3, 0, -1, []*entry{}},
				&entry{"0001", 1, 0, -1, []*entry{}},
				&entry{"1111", 15, 0, -1, []*entry{}},
				&entry{"1100", 12, 0, -1, []*entry{}},
			},
			&HuffmanCodes{
				&entry{"0001", 1, 0, -1, []*entry{}},
				&entry{"0011", 3, 0, -1, []*entry{}},
				&entry{"1100", 12, 0, -1, []*entry{}},
				&entry{"1111", 15, 0, -1, []*entry{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hc.sortByP()
		})
	}
}

func TestHuffmanCodes_sortByCode(t *testing.T) {
	tests := []struct {
		name string
		hc   *HuffmanCodes
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hc.sortByCode()
		})
	}
}

func TestHuffmanCodes_autoencode(t *testing.T) {
	tests := []struct {
		name string
		hc   HuffmanCodes
		want *HuffmanCodes
	}{ /*
		{
			"Test autoencoding HuffmanTable",
			HuffmanCodes{
				&entry{"0000", 0,  0, -1,[]*entry{}},
				&entry{"0011", 3,  0, -1,[]*entry{}},
				&entry{"1100", 12, 0, -1,[]*entry{}},
				&entry{"1111", 15, 0, -1,[]*entry{}},
			},
			&HuffmanCodes{
				&entry{"0000", 0, 0, -1, []*entry{
					&entry{"",[]*entry{
						&entry{},
						&entry{},
					}, 12,0, -1},

				HuffmanCode{&group{[]*group{}, "1111"}, 15},
				HuffmanCode{&group{[]*group{
					{[]*group{}, "1100"},
					{[]*group{
						{[]*group{}, "0011"},
						{[]*group{}, "0000"},
					}, ""}}, ""}, 15,
				},
			},
		},
	*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hc.autoencode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HuffmanCodes.autoencode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHuffmanCodes(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *HuffmanCodes
	}{
		{
			"0011,0000,1111,1100",
			args{
				Subscriptions{
					"0011": &Info{"3", 3, &Subscriptions{}},
					"0001": &Info{"1", 1, &Subscriptions{}},
					"1111": &Info{"15", 15, &Subscriptions{}},
					"1100": &Info{"12", 12, &Subscriptions{}},
				},
			},
			&HuffmanCodes{
				&entry{"0001", 1, 0, -1, []*entry{}},
				&entry{"0011", 3, 0, -1, []*entry{}},
				&entry{"1100", 12, 0, -1, []*entry{}},
				&entry{"1111", 15, 0, -1, []*entry{}},
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
