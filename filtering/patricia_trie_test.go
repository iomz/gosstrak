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
			pt.AnalyzeLocality(tt.args.id, tt.args.prefix, tt.args.lm)
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
			if gotMatches := pt.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("PatriciaTrie.Search() = %v, want %v", gotMatches, tt.wantMatches)
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
			if err := pt.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PatriciaTrie.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatriciaTrie_build(t *testing.T) {
	type fields struct {
		notificationURI string
		filterObject    *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
	}
	type args struct {
		prefix string
		sub    Subscriptions
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
			pt.build(tt.args.prefix, tt.args.sub)
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

func TestBuildPatriciaTrie(t *testing.T) {
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
					"0011":         &Info{0, "3", 10, nil},
					"00110011":     &Info{0, "3-3", 5, nil},
					"1111":         &Info{0, "15", 2, nil},
					"00110000":     &Info{0, "3-0", 5, nil},
					"001100110000": &Info{0, "3-3-0", 5, nil},
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
			got := BuildPatriciaTrie(tt.args.sub)
			if ok, g, w := got.equal(tt.want); !ok {
				t.Errorf("BuildPatriciaTrie() = \n%v, want \n%v", g.Dump(), w.Dump())
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
