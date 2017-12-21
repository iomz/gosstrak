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
			if gotMatches := ht.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("HuffmanTree.Search() = %v, want %v", gotMatches, tt.wantMatches)
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
			gotOk, gotGot, gotWanted := ht.equal(tt.args.want)
			if gotOk != tt.wantOk {
				t.Errorf("HuffmanTree.equal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotGot, tt.wantGot) {
				t.Errorf("HuffmanTree.equal() gotGot = %v, want %v", gotGot, tt.wantGot)
			}
			if !reflect.DeepEqual(gotWanted, tt.wantWanted) {
				t.Errorf("HuffmanTree.equal() gotWanted = %v, want %v", gotWanted, tt.wantWanted)
			}
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
			if got := ht.Dump(); got != tt.want {
				t.Errorf("HuffmanTree.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuffmanTree_makeBranch(t *testing.T) {
	type args struct {
		sub           *Subscriptions
		isComposition bool
		first         *entry
		second        *entry
	}
	tests := []struct {
		name string
		args args
		want *HuffmanTree
	}{
		{
			"test makeBranch 0 0 false",
			args{
				&Subscriptions{
					"0011": &Info{"3", 3, nil},
					"0001": &Info{"1", 1, nil},
					"1111": &Info{"15", 15, nil},
					"1100": &Info{"12", 12, nil},
				},
				false,
				&entry{"1111", 0, "15", 15, 0, HuffmanCodes{}, nil},
				&entry{"1100", 0, "12", 12, 0, HuffmanCodes{}, nil},
			},
			&HuffmanTree{
				"15",
				NewFilter("1111", 0),
				nil,
				&HuffmanTree{
					"12",
					NewFilter("1100", 0),
					nil,
					nil,
				},
			},
		},
		{
			"test makeBranch 0 0 true meaningless composition",
			args{
				&Subscriptions{
					"1111":     &Info{"240", 240, nil},
					"00001111": &Info{"15", 12, nil},
				},
				true,
				&entry{"00001111", 0, "240", 240, 0, HuffmanCodes{}, nil},
				&entry{"1111", 0, "15", 15, 0, HuffmanCodes{}, nil},
			},
			&HuffmanTree{
				"240",
				NewFilter("00001111", 0),
				nil,
				&HuffmanTree{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
			},
		},
		{
			"test makeBranch 0 0 true",
			args{
				&Subscriptions{
					"0011": &Info{"3", 3, nil},
					"0001": &Info{"1", 1, nil},
					"1111": &Info{"15", 15, nil},
					"1100": &Info{"12", 12, nil},
				},
				true,
				&entry{"1111", 0, "15", 15, 0, HuffmanCodes{}, nil},
				&entry{"1100", 0, "12", 12, 0, HuffmanCodes{}, nil},
			},
			&HuffmanTree{
				"",
				NewFilter("11xxxxxx", 0),
				&HuffmanTree{
					"15",
					NewFilter("1111xxxx", 0),
					nil,
					&HuffmanTree{
						"12",
						NewFilter("1100xxxx", 0),
						nil,
						nil,
					},
				},
				nil,
			},
		},
		{
			"test makeBranch 1 1 false",
			args{
				&Subscriptions{
					"0111": &Info{"7", 7, nil},
					"0110": &Info{"6", 6, nil},
					"0101": &Info{"5", 5, nil},
					"0100": &Info{"4", 4, nil},
				},
				false,
				&entry{"", 0, "", 13, 1, HuffmanCodes{}, &HuffmanTree{
					"",
					NewFilter("011xxxxx", 0),
					&HuffmanTree{
						"7",
						NewFilter("0111xxxx", 0),
						nil,
						&HuffmanTree{
							"6",
							NewFilter("0110xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				}},
				&entry{"", 0, "", 9, 1, HuffmanCodes{}, &HuffmanTree{
					"",
					NewFilter("010xxxxx", 0),
					&HuffmanTree{
						"5",
						NewFilter("0101xxxx", 0),
						nil,
						&HuffmanTree{
							"4",
							NewFilter("0100xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				}},
			},
			&HuffmanTree{
				"",
				NewFilter("011xxxxx", 0),
				&HuffmanTree{
					"7",
					NewFilter("0111xxxx", 0),
					nil,
					&HuffmanTree{
						"6",
						NewFilter("0110xxxx", 0),
						nil,
						nil,
					},
				},
				&HuffmanTree{
					"",
					NewFilter("010xxxxx", 0),
					&HuffmanTree{
						"5",
						NewFilter("0101xxxx", 0),
						nil,
						&HuffmanTree{
							"4",
							NewFilter("0100xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				},
			},
		},
		{
			"test makeBranch 1 0 true",
			args{
				&Subscriptions{
					"0111": &Info{"7", 7, nil},
					"0101": &Info{"5", 5, nil},
					"0100": &Info{"4", 4, nil},
				},
				true,
				&entry{"", 0, "", 9, 1, HuffmanCodes{}, &HuffmanTree{
					"",
					NewFilter("010xxxxx", 0),
					&HuffmanTree{
						"5",
						NewFilter("0101xxxx", 0),
						nil,
						&HuffmanTree{
							"4",
							NewFilter("0100xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				}},
				&entry{"0111", 0, "7", 7, 0, HuffmanCodes{}, nil},
			},
			&HuffmanTree{
				"",
				NewFilter("01xxxxxx", 0),
				&HuffmanTree{
					"",
					NewFilter("010xxxxx", 0),
					&HuffmanTree{
						"5",
						NewFilter("0101xxxx", 0),
						nil,
						&HuffmanTree{
							"4",
							NewFilter("0100xxxx", 0),
							nil,
							nil,
						},
					},
					&HuffmanTree{
						"7",
						NewFilter("0111xxxx", 0),
						nil,
						nil,
					},
				},
				nil,
			},
		},
		{
			"test makeBranch 1 0 true meaningless",
			args{
				&Subscriptions{
					"00001111": &Info{"15", 15, nil},
					"1110":     &Info{"224", 64, nil},
					"1101":     &Info{"208", 80, nil},
				},
				true,
				&entry{"", 0, "", 9, 1, HuffmanCodes{}, &HuffmanTree{
					"",
					NewFilter("11xxxxxx", 0),
					&HuffmanTree{
						"224",
						NewFilter("1110xxxx", 0),
						nil,
						&HuffmanTree{
							"208",
							NewFilter("1101xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				}},
				&entry{"00001111", 0, "15", 15, 0, HuffmanCodes{}, nil},
			},
			&HuffmanTree{
				"",
				NewFilter("11xxxxxx", 0),
				&HuffmanTree{
					"224",
					NewFilter("1110xxxx", 0),
					nil,
					&HuffmanTree{
						"208",
						NewFilter("1101xxxx", 0),
						nil,
						nil,
					},
				},
				&HuffmanTree{
					"15",
					NewFilter("00001111", 0),
					nil,
					nil,
				},
			},
		},
		{
			"test makeBranch 0 1 true",
			args{
				&Subscriptions{
					"1010": &Info{"10", 10, nil},
					"0101": &Info{"5", 5, nil},
					"0100": &Info{"4", 4, nil},
				},
				true,
				&entry{"1010", 0, "10", 10, 0, HuffmanCodes{}, nil},
				&entry{"", 0, "", 9, 1, HuffmanCodes{}, &HuffmanTree{
					"",
					NewFilter("010xxxxx", 0),
					&HuffmanTree{
						"5",
						NewFilter("0101xxxx", 0),
						nil,
						&HuffmanTree{
							"4",
							NewFilter("0100xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				}},
			},
			&HuffmanTree{
				"10",
				NewFilter("1010", 0),
				nil,
				&HuffmanTree{
					"",
					NewFilter("010xxxxx", 0),
					&HuffmanTree{
						"5",
						NewFilter("0101xxxx", 0),
						nil,
						&HuffmanTree{
							"4",
							NewFilter("0100xxxx", 0),
							nil,
							nil,
						},
					},
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeBranch(tt.args.sub, tt.args.isComposition, tt.args.first, tt.args.second)
			if ok, g, w := got.equal(tt.want); !ok {
				t.Errorf("HuffmanTree.makeBranch() = \n%v, want \n%v", g.Dump(), w.Dump())
			}
		})
	}
}

func TestHuffmanTree_makeSubsetBranch(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	type args struct {
		sub *Subscriptions
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
			if got := ht.makeSubsetBranch(tt.args.sub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HuffmanTree.makeSubsetBranch() = %v, want %v", got, tt.want)
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
		sub       Subscriptions
		compLimit int
	}
	tests := []struct {
		name string
		args args
		want *HuffmanTree
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildHuffmanTree(tt.args.sub, tt.args.compLimit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildHuffmanTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
