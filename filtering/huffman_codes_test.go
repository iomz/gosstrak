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
		code   string
		group  HuffmanCodes
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
				&entry{"0011", 3, 0, "", HuffmanCodes{}},
				&entry{"0001", 1, 0, "", HuffmanCodes{}},
				&entry{"1111", 15, 0, "", HuffmanCodes{}},
				&entry{"1100", 12, 0, "", HuffmanCodes{}},
			},
			&HuffmanCodes{
				&entry{"0001", 1, 0, "", HuffmanCodes{}},
				&entry{"0011", 3, 0, "", HuffmanCodes{}},
				&entry{"1100", 12, 0, "", HuffmanCodes{}},
				&entry{"1111", 15, 0, "", HuffmanCodes{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hc.sortByP()
		})
	}
}

func TestHuffmanCodes_autoencode(t *testing.T) {
	tests := []struct {
		name string
		hc   HuffmanCodes
		want *HuffmanCodes
	}{
		{
			"Test autoencoding HuffmanTable",
			HuffmanCodes{
				&entry{"0000", 0, 0, "", HuffmanCodes{}},
				&entry{"0011", 3, 0, "", HuffmanCodes{}},
				&entry{"1100", 12, 0, "", HuffmanCodes{}},
				&entry{"1111", 15, 0, "", HuffmanCodes{}},
			},
			&HuffmanCodes{
				&entry{"1111", 15, '0', "", HuffmanCodes{}},
				&entry{"", 15, '1', "", HuffmanCodes{
					&entry{"1100", 12, '1', "", HuffmanCodes{}},
					&entry{"", 3, '0', "", HuffmanCodes{
						&entry{"0011", 3, '1', "", HuffmanCodes{}},
						&entry{"0000", 0, '0', "", HuffmanCodes{}},
					}},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hc.autoencode()
			for i := range *got {
				if ok, g, w := (*got)[i].equal((*tt.want)[i]); !ok {
					t.Errorf("HuffmanCodes.autoencode() = \n%v, want \n%v", *g, *w)
				}
			}
		})
	}
}

func TestHuffmanCodes_gencode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		hc   *HuffmanCodes
		args args
		want HuffmanCodes
	}{
		{
			"gencode() nested HuffmanCodes test",
			&HuffmanCodes{
				&entry{"1111", 15, '0', "", HuffmanCodes{}},
				&entry{"", 15, '1', "", HuffmanCodes{
					&entry{"1100", 12, '1', "", HuffmanCodes{}},
					&entry{"", 3, '0', "", HuffmanCodes{
						&entry{"0011", 3, '1', "", HuffmanCodes{}},
						&entry{"0000", 0, '0', "", HuffmanCodes{}},
					}},
				}},
			},
			args{""},
			HuffmanCodes{
				&entry{"1111", 15, '0', "0", HuffmanCodes{}},
				&entry{"", 15, '1', "", HuffmanCodes{
					&entry{"1100", 12, '1', "11", HuffmanCodes{}},
					&entry{"", 3, '0', "", HuffmanCodes{
						&entry{"0011", 3, '1', "101", HuffmanCodes{}},
						&entry{"0000", 0, '0', "100", HuffmanCodes{}},
					}},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hc.gencode(tt.args.code)
			for i := range *tt.hc {
				if ok, g, w := (*tt.hc)[i].equal(tt.want[i]); !ok {
					t.Errorf("HuffmanCodes.autoencode() = \n%v, want \n%v", *g, *w)
				}
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
				&entry{"0001", 1, 0, "", HuffmanCodes{}},
				&entry{"0011", 3, 0, "", HuffmanCodes{}},
				&entry{"1100", 12, 0, "", HuffmanCodes{}},
				&entry{"1111", 15, 0, "", HuffmanCodes{}},
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

func Test_entry_equal(t *testing.T) {
	type fields struct {
		filter string
		p      float64
		bit    rune
		code   string
		group  HuffmanCodes
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
			fields{"", 15, '1', "", HuffmanCodes{}},
			args{
				&entry{"", 15, '1', "", HuffmanCodes{}},
			},
			true, nil, nil,
		},
		{
			"entry.equal simple test false",
			fields{"", 10, '0', "", HuffmanCodes{}},
			args{
				&entry{"", 15, '1', "", HuffmanCodes{}},
			},
			false,
			&entry{"", 10, '0', "", HuffmanCodes{}},
			&entry{"", 15, '1', "", HuffmanCodes{}},
		},
		{
			"entry.equal simple test invalid",
			fields{"", 15, '1', "", HuffmanCodes{
				&entry{"", 0, '0', "", HuffmanCodes{}},
			}},
			args{
				&entry{"", 15, '1', "", HuffmanCodes{
					&entry{"1010", 10, '0', "", HuffmanCodes{}},
					&entry{"0101", 5, '1', "", HuffmanCodes{}},
				}},
			},
			false,
			nil,
			&entry{"1010", 10, '0', "", HuffmanCodes{}},
		},
		{
			"entry.equal nest test true",
			fields{"", 15, '1', "", HuffmanCodes{
				&entry{"1100", 12, '1', "", HuffmanCodes{}},
				&entry{"", 3, '0', "", HuffmanCodes{
					&entry{"0011", 3, '1', "", HuffmanCodes{}},
					&entry{"0000", 0, '0', "", HuffmanCodes{}},
				}},
			},
			},
			args{
				&entry{"", 15, '1', "", HuffmanCodes{
					&entry{"1100", 12, '1', "", HuffmanCodes{}},
					&entry{"", 3, '0', "", HuffmanCodes{
						&entry{"0011", 3, '1', "", HuffmanCodes{}},
						&entry{"0000", 0, '0', "", HuffmanCodes{}},
					}},
				}},
			},
			true,
			nil,
			nil,
		},
		{
			"entry.equal nest test false",
			fields{"", 15, '1', "", HuffmanCodes{
				&entry{"1100", 12, '1', "", HuffmanCodes{}},
				&entry{"", 3, '0', "", HuffmanCodes{
					&entry{"0001", 1, '1', "", HuffmanCodes{}},
					&entry{"0000", 0, '0', "", HuffmanCodes{}},
				}},
			},
			},
			args{
				&entry{"", 15, '1', "", HuffmanCodes{
					&entry{"1100", 12, '1', "", HuffmanCodes{}},
					&entry{"", 3, '0', "", HuffmanCodes{
						&entry{"0011", 3, '1', "", HuffmanCodes{}},
						&entry{"0000", 0, '0', "", HuffmanCodes{}},
					}},
				}},
			},
			false,
			&entry{"0001", 1, '1', "", HuffmanCodes{}},
			&entry{"0011", 3, '1', "", HuffmanCodes{}},
		},
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
