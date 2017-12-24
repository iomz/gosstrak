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

func TestHuffmanTree_AnalyzeLocality(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	type args struct {
		id     []byte
		prefix string
		lm     *LocalityMap
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
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			ht.AnalyzeLocality(tt.args.id, tt.args.prefix, tt.args.lm)
		})
	}
}

func TestHuffmanTree_Dump(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"simple dump",
			fields{
				"00",
				NewFilter("00", 0),
				&HuffmanTree{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&HuffmanTree{
					"11",
					NewFilter("11", 0),
					nil,
					nil,
				},
			},
			"--00(0 2) -> 00\n  --11(2 2) -> 0011\n  --11(0 2) -> 11\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if got := ht.Dump(); got != tt.want {
				t.Errorf("HuffmanTree.Dump() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestHuffmanTree_MarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
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
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			got, err := ht.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("HuffmanTree.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HuffmanTree.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuffmanTree_Search(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	type args struct {
		id []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMatches []string
	}{
		{
			"no match",
			fields{
				"00",
				NewFilter("00", 0),
				&HuffmanTree{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&HuffmanTree{
					"11",
					NewFilter("11", 0),
					nil,
					nil,
				},
			},
			args{[]byte{64}},
			[]string{""},
		},
		{
			"1 match",
			fields{
				"00",
				NewFilter("00", 0),
				&HuffmanTree{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&HuffmanTree{
					"11",
					NewFilter("11", 0),
					nil,
					nil,
				},
			},
			args{[]byte{240}},
			[]string{"11"},
		},
		{
			"2 matches",
			fields{
				"00",
				NewFilter("00", 0),
				&HuffmanTree{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&HuffmanTree{
					"11",
					NewFilter("11", 0),
					nil,
					nil,
				},
			},
			args{[]byte{48}},
			[]string{"00", "0011"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if gotMatches := ht.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				if len(gotMatches) != 0 && len(tt.wantMatches) != 0 {
					t.Errorf("HuffmanTree.Search() = %v, want %v", gotMatches, tt.wantMatches)
				}
			}
		})
	}
}

func TestHuffmanTree_UnmarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
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
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if err := ht.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("HuffmanTree.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHuffmanTree_build(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	type args struct {
		sub *Subscriptions
		hc  *HuffmanCodes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HuffmanTree
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if got := ht.build(tt.args.sub, tt.args.hc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HuffmanTree.build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuffmanTree_equal(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	type args struct {
		want *HuffmanTree
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantGot    *HuffmanTree
		wantWanted *HuffmanTree
	}{
		{
			"wrong node",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			args{
				&HuffmanTree{
					"wanted",
					NewFilter("00000000", 0),
					nil,
					nil,
				},
			},
			false,
			&HuffmanTree{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			&HuffmanTree{
				"wanted",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
		},
		{
			"no matchNext",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			args{
				&HuffmanTree{
					"got",
					NewFilter("00000000", 0),
					&HuffmanTree{
						"wanted",
						NewFilter("11111111", 0),
						nil,
						nil,
					},
					nil,
				},
			},
			false,
			nil,
			&HuffmanTree{
				"wanted",
				NewFilter("11111111", 0),
				nil,
				nil,
			},
		},
		{
			"wrong matchNext",
			fields{
				"got",
				NewFilter("00000000", 0),
				&HuffmanTree{
					"got",
					NewFilter("11111111", 0),
					nil,
					nil,
				},
				nil,
			},
			args{
				&HuffmanTree{
					"got",
					NewFilter("00000000", 0),
					&HuffmanTree{
						"wanted",
						NewFilter("11111111", 0),
						nil,
						nil,
					},
					nil,
				},
			},
			false,
			nil,
			&HuffmanTree{
				"wanted",
				NewFilter("11111111", 0),
				nil,
				nil,
			},
		},
		{
			"no mismatchNext",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			args{
				&HuffmanTree{
					"got",
					NewFilter("00000000", 0),
					nil,
					&HuffmanTree{
						"wanted",
						NewFilter("11111111", 0),
						nil,
						nil,
					},
				},
			},
			false,
			nil,
			&HuffmanTree{
				"wanted",
				NewFilter("11111111", 0),
				nil,
				nil,
			},
		},
		{
			"wrong mismatchNext",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				&HuffmanTree{
					"got",
					NewFilter("11111111", 0),
					nil,
					nil,
				},
			},
			args{
				&HuffmanTree{
					"got",
					NewFilter("00000000", 0),
					nil,
					&HuffmanTree{
						"wanted",
						NewFilter("11111111", 0),
						nil,
						nil,
					},
				},
			},
			false,
			nil,
			&HuffmanTree{
				"wanted",
				NewFilter("11111111", 0),
				nil,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			gotOk, _, _ := ht.equal(tt.args.want)
			if gotOk != tt.wantOk {
				t.Errorf("HuffmanTree.equal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestHuffmanTree_print(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
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
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			writer := &bytes.Buffer{}
			ht.print(writer, tt.args.indent)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("HuffmanTree.print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestBuildHuffmanTree(t *testing.T) {
	type args struct {
		sub *Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *HuffmanTree
	}{
		{
			"",
			args{
				&Subscriptions{
					"0011":         &Info{0, "3", 10, nil},
					"00110011":     &Info{0, "3-3", 5, nil},
					"1111":         &Info{0, "15", 2, nil},
					"00110000":     &Info{0, "3-0", 5, nil},
					"001100110000": &Info{0, "3-3-0", 5, nil},
				},
			},
			&HuffmanTree{
				"3",
				NewFilter("0011", 0),
				&HuffmanTree{
					"3-3",
					NewFilter("0011", 4),
					&HuffmanTree{
						"3-3-0",
						NewFilter("0000", 8),
						nil,
						nil,
					},
					&HuffmanTree{
						"3-0",
						NewFilter("0000", 4),
						nil,
						nil,
					},
				},
				&HuffmanTree{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildHuffmanTree(tt.args.sub)
			if ok, g, w := got.equal(tt.want); !ok {
				t.Errorf("BuildHuffmanTree() = \n%v, want \n%v", g.Dump(), w.Dump())
			}
		})
	}
}
