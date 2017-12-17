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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lcp(tt.args.l); got != tt.want {
				t.Errorf("lcp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatriciaTrie_MarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
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

func TestPatriciaTrie_UnmarshalBinary(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			if err := pt.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PatriciaTrie.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPatriciaTrie_AnalyzeLocality(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.AnalyzeLocality(tt.args.id, tt.args.prefix, tt.args.lm)
		})
	}
}

func TestPatriciaTrie_Search(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			if gotMatches := pt.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("PatriciaTrie.Search() = %v, want %v", gotMatches, tt.wantMatches)
			}
		})
	}
}

func TestPatriciaTrie_build(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			pt.build(tt.args.prefix, tt.args.sub)
		})
	}
}

func TestPatriciaTrie_Dump(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
		one             *PatriciaTrie
		zero            *PatriciaTrie
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
				notificationURI: tt.fields.notificationURI,
				filter:          tt.fields.filter,
				one:             tt.fields.one,
				zero:            tt.fields.zero,
			}
			if got := pt.Dump(); got != tt.want {
				t.Errorf("PatriciaTrie.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatriciaTrie_print(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildPatriciaTrie(tt.args.sub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildPatriciaTrie() = %v, want %v", got, tt.want)
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
	// TODO: Add test cases.
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
