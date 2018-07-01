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

func TestPatriciaTrie_AnalyzeLocality(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
			"simple analyze",
			fields{
				"3",
				NewFilter("0011", 0),
				nil,
				&PatriciaTrie{
					"",
					NewFilter("00", 4),
					&PatriciaTrie{
						"3-3",
						NewFilter("11", 6),
						nil,
						&PatriciaTrie{
							"3-3-0",
							NewFilter("0000", 8),
							nil,
							nil,
						},
					},
					&PatriciaTrie{
						"3-0",
						NewFilter("00", 6),
						nil,
						nil,
					},
				},
			},
			args{
				[]byte{51}, // "b00110011"
				"",
				&LocalityMap{},
			},
			&LocalityMap{
				",0011":       1,
				",0011,00":    1,
				",0011,00,11": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.AnalyzeLocality(tt.args.id, tt.args.prefix, tt.args.lm)
			if !reflect.DeepEqual(*tt.args.lm, *tt.want) {
				t.Errorf("PatriciaTrie.AnalyzeLocality() = \n%v, want \n%v", *tt.args.lm, *tt.want)
			}
		})
	}
}

func TestPatriciaTrie_Dump(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"simple patricia",
			fields{
				"3",
				NewFilter("0011", 0),
				nil,
				&PatriciaTrie{
					"",
					NewFilter("00", 4),
					&PatriciaTrie{
						"3-3",
						NewFilter("11", 6),
						nil,
						&PatriciaTrie{
							"3-3-0",
							NewFilter("0000", 8),
							nil,
							nil,
						},
					},
					&PatriciaTrie{
						"3-0",
						NewFilter("00", 6),
						nil,
						nil,
					},
				},
			},
			"--0011(0 4) -> 3\n  --00(4 2) \n    --11(6 2) -> 3-3\n      --0000(8 4) -> 3-3-0\n    --00(6 2) -> 3-0\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			if got := pt.Dump(); got != tt.want {
				t.Errorf("PatriciaTrie.Dump() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestPatriciaTrie_MarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"simple marshal patricia",
			fields{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			[]byte{32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 3, 12, 0, 0, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 11, 255, 130, 4, 1, 255, 1, 1, 255, 2, 2, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 255, 190, 255, 136, 0, 255, 185, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 5, 12, 0, 2, 49, 53, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 49, 49, 49, 49, 1, 8, 2, 1, 255, 1, 1, 15, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 254, 3, 234, 255, 136, 0, 254, 3, 228, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 4, 12, 0, 1, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 48, 48, 49, 49, 1, 8, 2, 1, 63, 1, 1, 15, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 254, 3, 30, 255, 136, 0, 254, 3, 24, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 3, 12, 0, 0, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 2, 48, 48, 1, 4, 1, 8, 1, 1, 243, 1, 1, 243, 2, 2, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 254, 1, 146, 255, 136, 0, 254, 1, 140, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 6, 12, 0, 3, 51, 45, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 2, 49, 49, 1, 4, 1, 12, 1, 1, 255, 1, 1, 252, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 255, 197, 255, 136, 0, 255, 192, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 8, 12, 0, 5, 51, 45, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 23, 255, 130, 1, 4, 48, 48, 48, 48, 1, 8, 1, 16, 1, 1, 15, 1, 1, 15, 1, 2, 1, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 255, 191, 255, 136, 0, 255, 186, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 6, 12, 0, 3, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 2, 48, 48, 1, 4, 1, 12, 1, 1, 252, 1, 1, 252, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			got, err := pt.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("PatriciaTrie.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatriciaTrie.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatriciaTrie_Search(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{
				[]byte{64},
			},
			[]string{},
		},
		{
			"1 match",
			fields{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{
				[]byte{240},
			},
			[]string{"15"},
		},
		{
			"2 matches",
			fields{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{
				[]byte{48},
			},
			[]string{"3", "3-0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			if gotMatches := pt.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				if len(gotMatches) != 0 && len(tt.wantMatches) != 0 {
					t.Errorf("PatriciaTrie.Search() = %v, want %v", gotMatches, tt.wantMatches)
				}
			}
		})
	}
}

func TestPatriciaTrie_UnmarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
			"simple unmarshal patricia",
			fields{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{
				[]byte{32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 3, 12, 0, 0, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 11, 255, 130, 4, 1, 255, 1, 1, 255, 2, 2, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 255, 190, 255, 136, 0, 255, 185, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 5, 12, 0, 2, 49, 53, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 49, 49, 49, 49, 1, 8, 2, 1, 255, 1, 1, 15, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 254, 3, 234, 255, 136, 0, 254, 3, 228, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 4, 12, 0, 1, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 4, 48, 48, 49, 49, 1, 8, 2, 1, 63, 1, 1, 15, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 254, 3, 30, 255, 136, 0, 254, 3, 24, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 3, 12, 0, 0, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 2, 48, 48, 1, 4, 1, 8, 1, 1, 243, 1, 1, 243, 2, 2, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 254, 1, 146, 255, 136, 0, 254, 1, 140, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 6, 12, 0, 3, 51, 45, 51, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 2, 49, 49, 1, 4, 1, 12, 1, 1, 255, 1, 1, 252, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 1, 10, 255, 135, 6, 1, 2, 255, 138, 0, 0, 0, 255, 197, 255, 136, 0, 255, 192, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 8, 12, 0, 5, 51, 45, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 23, 255, 130, 1, 4, 48, 48, 48, 48, 1, 8, 1, 16, 1, 1, 15, 1, 1, 15, 1, 2, 1, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0, 3, 2, 0, 1, 255, 191, 255, 136, 0, 255, 186, 32, 12, 0, 29, 69, 110, 103, 105, 110, 101, 58, 102, 105, 108, 116, 101, 114, 105, 110, 103, 46, 80, 97, 116, 114, 105, 99, 105, 97, 84, 114, 105, 101, 6, 12, 0, 3, 51, 45, 48, 3, 2, 0, 1, 113, 255, 129, 3, 1, 1, 12, 70, 105, 108, 116, 101, 114, 79, 98, 106, 101, 99, 116, 1, 255, 130, 0, 1, 7, 1, 6, 83, 116, 114, 105, 110, 103, 1, 12, 0, 1, 4, 83, 105, 122, 101, 1, 4, 0, 1, 6, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 10, 66, 121, 116, 101, 70, 105, 108, 116, 101, 114, 1, 10, 0, 1, 8, 66, 121, 116, 101, 77, 97, 115, 107, 1, 10, 0, 1, 10, 66, 121, 116, 101, 79, 102, 102, 115, 101, 116, 1, 4, 0, 1, 8, 66, 121, 116, 101, 83, 105, 122, 101, 1, 4, 0, 0, 0, 19, 255, 130, 1, 2, 48, 48, 1, 4, 1, 12, 1, 1, 252, 1, 1, 252, 2, 2, 0, 3, 2, 0, 0, 3, 2, 0, 0},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			if err := pt.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PatriciaTrie.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatriciaTrie_equal(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
	}
	type args struct {
		want *PatriciaTrie
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantGot    *PatriciaTrie
		wantWanted *PatriciaTrie
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
				&PatriciaTrie{
					"wanted",
					NewFilter("00000000", 0),
					nil,
					nil,
				},
			},
			false,
			&PatriciaTrie{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			&PatriciaTrie{
				"wanted",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
		},
		{
			"no one",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			args{
				&PatriciaTrie{
					"got",
					NewFilter("00000000", 0),
					&PatriciaTrie{
						"wanted",
						NewFilter("00000000", 0),
						nil,
						nil,
					},
					nil,
				},
			},
			false,
			nil,
			&PatriciaTrie{
				"wanted",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
		},
		{
			"wrong one",
			fields{
				"got",
				NewFilter("00000000", 0),
				&PatriciaTrie{
					"wanted",
					NewFilter("11110000", 0),
					nil,
					nil,
				},
				nil,
			},
			args{
				&PatriciaTrie{
					"got",
					NewFilter("00000000", 0),
					&PatriciaTrie{
						"wanted",
						NewFilter("00000000", 0),
						nil,
						nil,
					},
					nil,
				},
			},
			false,
			&PatriciaTrie{
				"wanted",
				NewFilter("11110000", 0),
				nil,
				nil,
			},
			&PatriciaTrie{
				"wanted",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
		},
		{
			"no zero",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
			args{
				&PatriciaTrie{
					"got",
					NewFilter("00000000", 0),
					nil,
					&PatriciaTrie{
						"wanted",
						NewFilter("00000000", 0),
						nil,
						nil,
					},
				},
			},
			false,
			nil,
			&PatriciaTrie{
				"wanted",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
		},
		{
			"wrong zero",
			fields{
				"got",
				NewFilter("00000000", 0),
				nil,
				&PatriciaTrie{
					"wanted",
					NewFilter("11110000", 0),
					nil,
					nil,
				},
			},
			args{
				&PatriciaTrie{
					"got",
					NewFilter("00000000", 0),
					nil,
					&PatriciaTrie{
						"wanted",
						NewFilter("00000000", 0),
						nil,
						nil,
					},
				},
			},
			false,
			&PatriciaTrie{
				"wanted",
				NewFilter("11110000", 0),
				nil,
				nil,
			},
			&PatriciaTrie{
				"wanted",
				NewFilter("00000000", 0),
				nil,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			gotOk, _, _ := pt.equal(tt.args.want)
			if gotOk != tt.wantOk {
				t.Errorf("PatriciaTrie.equal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestPatriciaTrie_print(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			writer := &bytes.Buffer{}
			pt.print(writer, tt.args.indent)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("PatriciaTrie.print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestNewPatriciaTrie(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *PatriciaTrie
	}{
		{
			"simple patricia",
			args{
				Subscriptions{
					m: SubMap{
						"0011":         &Info{0, "3", 10, Subscriptions{}},
						"00110011":     &Info{0, "3-3", 5, Subscriptions{}},
						"1111":         &Info{0, "15", 2, Subscriptions{}},
						"00110000":     &Info{0, "3-0", 5, Subscriptions{}},
						"001100110000": &Info{0, "3-3-0", 5, Subscriptions{}},
					},
				},
			},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPatriciaTrie(tt.args.sub)
			if ok, g, w := got.(*PatriciaTrie).equal(tt.want); !ok {
				t.Errorf("NewPatriciaTrie() = \n%v, want \n%v", g.Dump(), w.Dump())
			}
		})
	}
}

func Test_getNextBit(t *testing.T) {
	type args struct {
		id  []byte
		nbo int
	}
	tests := []struct {
		name    string
		args    args
		want    rune
		wantErr bool
	}{
		{"0000111100001111, 4", args{[]byte{15, 15}, 4}, '1', false},
		{"00001111, 8", args{[]byte{15}, 8}, 'x', false},
		{"00001111, 9", args{[]byte{15}, 9}, '?', true},
		{"00001111, 3", args{[]byte{15}, 3}, '0', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNextBit(tt.args.id, tt.args.nbo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNextBit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNextBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lcp(t *testing.T) {
	type args struct {
		l []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0001,0010,0011,0000", args{[]string{"0001", "0010", "0011", "0000"}}, "00"},
		{"0001", args{[]string{"0001"}}, "0001"},
		{"", args{[]string{}}, ""},
		{"0011,00", args{[]string{"0011", "00"}}, "00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcp(tt.args.l); got != tt.want {
				t.Errorf("lcp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatriciaTrie_add(t *testing.T) {
	type args struct {
		fs              string
		notificationURI string
	}
	tests := []struct {
		name   string
		fields *PatriciaTrie
		args   args
		want   *PatriciaTrie
	}{
		{
			"add a node to the edge",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"111110", "15-2"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					&PatriciaTrie{
						"15-2",
						NewFilter("10", 4),
						nil,
						nil,
					},
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"add a node to the top",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"1100", "12"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"",
					NewFilter("11", 0),
					&PatriciaTrie{
						"15",
						NewFilter("11", 2),
						nil,
						nil,
					},
					&PatriciaTrie{
						"12",
						NewFilter("00", 2),
						nil,
						nil,
					},
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"add a node in the middle 1",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"00110010", "3-0-2"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"",
							NewFilter("1", 6),
							&PatriciaTrie{
								"3-3",
								NewFilter("1", 7),
								nil,
								&PatriciaTrie{
									"3-3-0",
									NewFilter("0000", 8),
									nil,
									nil,
								},
							},
							&PatriciaTrie{
								"3-0-2",
								NewFilter("0", 7),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"add a node in the middle 2",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"001100110011", "3-3-3"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"",
								NewFilter("00", 8),
								&PatriciaTrie{
									"3-3-3",
									NewFilter("11", 10),
									nil,
									nil,
								},
								&PatriciaTrie{
									"3-3-0",
									NewFilter("00", 10),
									nil,
									nil,
								},
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"add a node in the middle 3",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"111101", "15-1"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"add a node that is already there",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"001100110000", "3-3-0_new"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0_new",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.add(tt.args.fs, tt.args.notificationURI)
			if ok, got, wanted := pt.equal(tt.want); !ok {
				t.Errorf("add() = \n%v, want \n%v", got.Dump(), wanted.Dump())
			}
		})
	}
}

func TestPatriciaTrie_AddSubscription(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.AddSubscription(tt.args.sub)
		})
	}
}

func TestPatriciaTrie_DeleteSubscription(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.DeleteSubscription(tt.args.sub)
		})
	}
}

func TestPatriciaTrie_delete(t *testing.T) {
	type args struct {
		fs              string
		notificationURI string
	}
	tests := []struct {
		name   string
		fields *PatriciaTrie
		args   args
		want   *PatriciaTrie
	}{
		{
			"delete a node from the edge",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"001100110000", "3-3-0"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							nil,
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"delete a node from the middle 1",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"00110011", "3-3"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3-0",
							NewFilter("110000", 6),
							nil,
							nil,
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"delete a node from the middle 3",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"3-0",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							nil,
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"001100", "3-0"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							nil,
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
		{
			"delete a node from the middle 4",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"3-0",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							nil,
						},
						nil,
					},
				},
			},
			args{"001100", "3-0"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"3-3",
						NewFilter("0011", 4),
						nil,
						nil,
					},
				},
			},
		},
		{
			"delete a node from the top",
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				&PatriciaTrie{
					"15",
					NewFilter("1111", 0),
					nil,
					nil,
				},
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
			args{"1111", "15"},
			&PatriciaTrie{
				"",
				NewFilter("", 0),
				nil,
				&PatriciaTrie{
					"3",
					NewFilter("0011", 0),
					nil,
					&PatriciaTrie{
						"",
						NewFilter("00", 4),
						&PatriciaTrie{
							"3-3",
							NewFilter("11", 6),
							nil,
							&PatriciaTrie{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
						},
						&PatriciaTrie{
							"3-0",
							NewFilter("00", 6),
							nil,
							nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PatriciaTrie{
				notificationURI: tt.fields.notificationURI,
				filterObject:    tt.fields.filterObject,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.delete(tt.args.fs, tt.args.notificationURI)
			if pt.Dump() != tt.want.Dump() {
				t.Errorf("delete() = \n%v, want \n%v", pt.Dump(), tt.want.Dump())
				//t.Errorf("delete() = \n%v, want \n%v", got, wanted)
			}
		})
	}
}
