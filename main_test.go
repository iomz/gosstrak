package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFilterMap_keys(t *testing.T) {
	tests := []struct {
		name string
		fm   FilterMap
		want []string
	}{
		{"0,8", FilterMap{"0000": "0", "1000": "8"}, []string{"0000", "1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fm.keys()
			for i := 0; i < len(got); i++ {
				if _, ok := tt.fm[got[i]]; !ok {
					t.Errorf("FilterMap.keys() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestPatriciaTrie_constructTrie(t *testing.T) {
	type fields struct {
		prefix string
		one    *PatriciaTrie
		zero   *PatriciaTrie
		notify string
	}
	type args struct {
		prefix string
		fm     FilterMap
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
				prefix: tt.fields.prefix,
				one:    tt.fields.one,
				zero:   tt.fields.zero,
				notify: tt.fields.notify,
			}
			pt.constructTrie(tt.args.prefix, tt.args.fm)
		})
	}
}

func TestPatriciaTrie_dump(t *testing.T) {
	type fields struct {
		prefix string
		one    *PatriciaTrie
		zero   *PatriciaTrie
		notify string
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
			pt := &PatriciaTrie{
				prefix: tt.fields.prefix,
				one:    tt.fields.one,
				zero:   tt.fields.zero,
				notify: tt.fields.notify,
			}
			if got := pt.dump(); got != tt.want {
				t.Errorf("PatriciaTrie.dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatriciaTrie_print(t *testing.T) {
	type fields struct {
		prefix string
		one    *PatriciaTrie
		zero   *PatriciaTrie
		notify string
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
				prefix: tt.fields.prefix,
				one:    tt.fields.one,
				zero:   tt.fields.zero,
				notify: tt.fields.notify,
			}
			writer := &bytes.Buffer{}
			pt.print(writer, tt.args.indent)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("PatriciaTrie.print() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func Test_loadFiltersFromCSVFile(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want FilterMap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadFiltersFromCSVFile(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFiltersFromCSVFile() = %v, want %v", got, tt.want)
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
		{"0001,0010,0011,0000", args{[]string{"0001","0010","0011","0000"}}, "00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcp(tt.args.l); got != tt.want {
				t.Errorf("lcp() = %v, want %v", got, tt.want)
			}
		})
	}
}
