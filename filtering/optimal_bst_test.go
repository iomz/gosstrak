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

func TestOptimalBST_AnalyzeLocality(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
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
		want   *LocalityMap
	}{
		{
			"simple analyze huffman",
			fields{
				"00",
				NewFilter("00", 0),
				&OptimalBST{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&OptimalBST{
					"11",
					NewFilter("11", 0),
					nil,
					nil,
				},
			},
			args{
				[]byte{48},
				"",
				&LocalityMap{},
			},
			&LocalityMap{
				"00":            1,
				"00,0011":       1,
				"00,0011,Match": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			obst.AnalyzeLocality(tt.args.id, tt.args.prefix, tt.args.lm)
			if !reflect.DeepEqual(*tt.args.lm, *tt.want) {
				t.Errorf("PatriciaTrie.AnalyzeLocality() = \n%v, want \n%v", *tt.args.lm, *tt.want)
			}
		})
	}
}

func TestOptimalBST_Dump(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
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
				&OptimalBST{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&OptimalBST{
					"11",
					NewFilter("11", 0),
					nil,
					nil,
				},
			},
			"--00(0 2) -> 00\n  ok--11(2 2) -> 0011\n  ng--11(0 2) -> 11\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if got := obst.Dump(); got != tt.want {
				t.Errorf("OptimalBST.Dump() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestOptimalBST_MarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"simple marshal",
			fields{
				"3",
				NewFilter("0011", 0),
				&OptimalBST{
					"3-3",
					NewFilter("0011", 4),
					&OptimalBST{
						"3-3-0",
						NewFilter("0000", 8),
						nil,
						nil,
					},
					&OptimalBST{
						"3-0",
						NewFilter("0000", 4),
						nil,
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
			[]byte{30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 4, 12, 0, 1, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 48, 48, 49, 49, 1, 8, 2, 1, 63, 1, 1, 15, 2, 2, 0, 3, 2, 0, 1, 10, 255, 131, 6, 1, 2, 255, 134, 0, 0, 0, 254, 2, 81, 255, 132, 0, 254, 2, 75, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 6, 12, 0, 3, 51, 45, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 21, 255, 130, 1, 4, 48, 48, 49, 49, 1, 8, 1, 8, 1, 1, 243, 1, 1, 240, 2, 2, 0, 3, 2, 0, 1, 10, 255, 131, 6, 1, 2, 255, 134, 0, 0, 0, 255, 195, 255, 132, 0, 255, 190, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 8, 12, 0, 5, 51, 45, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 23, 255, 130, 1, 4, 48, 48, 48, 48, 1, 8, 1, 16, 1, 1, 15, 1, 1, 15, 1, 2, 1, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 255, 191, 255, 132, 0, 255, 186, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 6, 12, 0, 3, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 21, 255, 130, 1, 4, 48, 48, 48, 48, 1, 8, 1, 8, 1, 1, 240, 1, 1, 240, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 255, 188, 255, 132, 0, 255, 183, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 5, 12, 0, 2, 49, 53, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 49, 49, 49, 49, 1, 8, 2, 1, 255, 1, 1, 15, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			got, err := obst.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("OptimalBST.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OptimalBST.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptimalBST_Search(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
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
				&OptimalBST{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&OptimalBST{
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
				&OptimalBST{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&OptimalBST{
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
				&OptimalBST{
					"0011",
					NewFilter("11", 2),
					nil,
					nil,
				},
				&OptimalBST{
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
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if gotMatches := obst.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				if len(gotMatches) != 0 && len(tt.wantMatches) != 0 {
					t.Errorf("OptimalBST.Search() = %v, want %v", gotMatches, tt.wantMatches)
				}
			}
		})
	}
}

func TestOptimalBST_UnmarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
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
		{
			"simple unmarshal huffman",
			fields{
				"3",
				NewFilter("0011", 0),
				&OptimalBST{
					"3-3",
					NewFilter("0011", 4),
					&OptimalBST{
						"3-3-0",
						NewFilter("0000", 8),
						nil,
						nil,
					},
					&OptimalBST{
						"3-0",
						NewFilter("0000", 4),
						nil,
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
			args{
				[]byte{30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 4, 12, 0, 1, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 48, 48, 49, 49, 1, 8, 2, 1, 63, 1, 1, 15, 2, 2, 0, 3, 2, 0, 1, 10, 255, 131, 6, 1, 2, 255, 134, 0, 0, 0, 254, 2, 81, 255, 132, 0, 254, 2, 75, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 6, 12, 0, 3, 51, 45, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 21, 255, 130, 1, 4, 48, 48, 49, 49, 1, 8, 1, 8, 1, 1, 243, 1, 1, 240, 2, 2, 0, 3, 2, 0, 1, 10, 255, 131, 6, 1, 2, 255, 134, 0, 0, 0, 255, 195, 255, 132, 0, 255, 190, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 8, 12, 0, 5, 51, 45, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 23, 255, 130, 1, 4, 48, 48, 48, 48, 1, 8, 1, 16, 1, 1, 15, 1, 1, 15, 1, 2, 1, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 255, 191, 255, 132, 0, 255, 186, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 6, 12, 0, 3, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 21, 255, 130, 1, 4, 48, 48, 48, 48, 1, 8, 1, 8, 1, 1, 240, 1, 1, 240, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 255, 188, 255, 132, 0, 255, 183, 30, 12, 0, 27, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 79, 112, 116, 105, 109, 97, 108, 66, 83, 84, 5, 12, 0, 2, 49, 53, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 49, 49, 49, 49, 1, 8, 2, 1, 255, 1, 1, 15, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if err := obst.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("OptimalBST.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptimalBST_build(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
	}
	type args struct {
		sub *Subscriptions
		nds *Nodes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *OptimalBST
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
			if got := obst.build(tt.args.sub, tt.args.nds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OptimalBST.build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptimalBST_equal(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
	}
	type args struct {
		want *OptimalBST
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantGot    *OptimalBST
		wantWanted *OptimalBST
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
				&OptimalBST{
					"wanted",
					NewFilter("00000000", 0),
					nil,
					nil,
				},
			},
			false,
			&OptimalBST{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			&OptimalBST{
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
				&OptimalBST{
					"got",
					NewFilter("00000000", 0),
					&OptimalBST{
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
			&OptimalBST{
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
				&OptimalBST{
					"got",
					NewFilter("11111111", 0),
					nil,
					nil,
				},
				nil,
			},
			args{
				&OptimalBST{
					"got",
					NewFilter("00000000", 0),
					&OptimalBST{
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
			&OptimalBST{
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
				&OptimalBST{
					"got",
					NewFilter("00000000", 0),
					nil,
					&OptimalBST{
						"wanted",
						NewFilter("11111111", 0),
						nil,
						nil,
					},
				},
			},
			false,
			nil,
			&OptimalBST{
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
				&OptimalBST{
					"got",
					NewFilter("11111111", 0),
					nil,
					nil,
				},
			},
			args{
				&OptimalBST{
					"got",
					NewFilter("00000000", 0),
					nil,
					&OptimalBST{
						"wanted",
						NewFilter("11111111", 0),
						nil,
						nil,
					},
				},
			},
			false,
			nil,
			&OptimalBST{
				"wanted",
				NewFilter("11111111", 0),
				nil,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			gotOk, _, _ := obst.equal(tt.args.want)
			if gotOk != tt.wantOk {
				t.Errorf("OptimalBST.equal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestOptimalBST_print(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
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
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			writer := &bytes.Buffer{}
			obst.print(writer, tt.args.indent)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("OptimalBST.print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestBuildOptimalBST(t *testing.T) {
	type args struct {
		sub *Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *OptimalBST
	}{
	/*
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
			&OptimalBST{
				"3",
				NewFilter("0011", 0),
				&OptimalBST{
					"3-3",
					NewFilter("0011", 4),
					&OptimalBST{
						"3-3-0",
						NewFilter("0000", 8),
						nil,
						nil,
					},
					&OptimalBST{
						"3-0",
						NewFilter("0000", 4),
						nil,
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
	*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildOptimalBST(tt.args.sub)
			if ok, g, w := got.equal(tt.want); !ok {
				t.Errorf("BuildOptimalBST() = \n%v, want \n%v", g.Dump(), w.Dump())
			}
		})
	}
}

func TestOptimalBST_AddSubscription(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
	}
	type args struct {
		sub Subscriptions
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
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			obst.AddSubscription(tt.args.sub)
		})
	}
}

func TestOptimalBST_DeleteSubscription(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		matchNext       *OptimalBST
		mismatchNext    *OptimalBST
	}
	type args struct {
		sub Subscriptions
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
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			obst.DeleteSubscription(tt.args.sub)
		})
	}
}

func TestOptimalBST_add(t *testing.T) {
	type args struct {
		fs              string
		notificationURI string
	}
	tests := []struct {
		name   string
		fields *OptimalBST
		args   args
		want   *OptimalBST
	}{
		{
			"add a subset of subscription 1",
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
			args{"11111010", "15-10"},
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
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
		},
		{
			"add a subset of subscription 2",
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
			args{"001100111100", "3-3-8"},
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
							&OptimalBST{
								"3-3-8",
								NewFilter("1100", 8),
								nil,
								nil,
							},
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
		{
			"add totally different subscription",
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
			args{"1010", "10"},
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
					&OptimalBST{
						"10",
						NewFilter("1010", 0),
						nil,
						nil,
					},
				},
			},
		},
		{
			"add an existing subscription (and update)",
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
			args{"00110000", "3-0_new"},
			&OptimalBST{
				"3",
				NewFilter("0011", 0),
				&OptimalBST{
					"3-0_new",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			obst.add(tt.args.fs, tt.args.notificationURI)
			if obst.Dump() != tt.want.Dump() {
				t.Errorf("add() = \n%v, want \n%v", obst.Dump(), tt.want.Dump())
			}
		})
	}
}

func TestOptimalBST_delete(t *testing.T) {
	type args struct {
		fs              string
		notificationURI string
	}
	tests := []struct {
		name   string
		fields *OptimalBST
		args   args
		want   *OptimalBST
	}{
		{
			"delete a subscription from the top branch",
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
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
			args{"11111010", "15-10"},
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
		{
			"delete the top branch edge",
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
					&OptimalBST{
						"10",
						NewFilter("1010", 0),
						nil,
						nil,
					},
				},
			},
			args{"1010", "10"},
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
		{
			"delete a subset of subscription 2",
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
							&OptimalBST{
								"3-3-8",
								NewFilter("1100", 8),
								nil,
								nil,
							},
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
			args{"001100111100", "3-3-8"},
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
		{
			"delete a subset of subscription and replace it with the mismatchNext node",
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
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
			args{"00110000", "3-0"},
			&OptimalBST{
				"3",
				NewFilter("0011", 0),
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
				&OptimalBST{
					"15",
					NewFilter("1111", 0),
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
		},
		{
			"delete a subset of subscription and concatenate with the to-be-deleted node",
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
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
			args{"00110011", "3-3"},
			&OptimalBST{
				"3",
				NewFilter("0011", 0),
				&OptimalBST{
					"3-0",
					NewFilter("0000", 4),
					nil,
					&OptimalBST{
						"3-3-0",
						NewFilter("00110000", 4),
						nil,
						nil,
					},
				},
				&OptimalBST{
					"15",
					NewFilter("1111", 0),
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
		},
		{
			"delete a subset of subscription and to-be-deleted node becomes an aggregation node",
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
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
			args{"0011", "3"},
			&OptimalBST{
				"",
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
					&OptimalBST{
						"15-10",
						NewFilter("1010", 4),
						nil,
						nil,
					},
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obst := &OptimalBST{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			obst.delete(tt.args.fs, tt.args.notificationURI)
			if obst.Dump() != tt.want.Dump() {
				t.Errorf("delete() = \n%v, want \n%v", obst.Dump(), tt.want.Dump())
			}
		})
	}
}
