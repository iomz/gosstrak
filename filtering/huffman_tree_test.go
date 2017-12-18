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
		filter          *FilterObject
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
				filter:          tt.fields.filter,
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
		filter          *FilterObject
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
				filter:          tt.fields.filter,
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
		filter          *FilterObject
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
				filter:          tt.fields.filter,
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
		filter          *FilterObject
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
				filter:          tt.fields.filter,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if gotMatches := ht.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("HuffmanTree.Search() = %v, want %v", gotMatches, tt.wantMatches)
			}
		})
	}
}

func TestHuffmanTree_build(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
		matchNext       *HuffmanTree
		mismatchNext    *HuffmanTree
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := &HuffmanTree{
				notificationURI: tt.fields.notificationURI,
				filter:          tt.fields.filter,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			ht.build()
		})
	}
}

func TestHuffmanTree_Dump(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
				matchNext:       tt.fields.matchNext,
				mismatchNext:    tt.fields.mismatchNext,
			}
			if got := ht.Dump(); got != tt.want {
				t.Errorf("HuffmanTree.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuffmanTree_print(t *testing.T) {
	type fields struct {
		notificationURI string
		filter          *FilterObject
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
				filter:          tt.fields.filter,
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

/*
func TestBuildHuffmanTree(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *HuffmanTree
	}{
		{
			"",
			args{
				Subscriptions{
					"0011": &Info{"3", 3, &Subscriptions{}},
					"0000": &Info{"0", 0, &Subscriptions{}},
					"1111": &Info{"15", 15, &Subscriptions{}},
					"1100": &Info{"12", 12, &Subscriptions{}},
				},
			},
			&HuffmanTree{
				"test",
				&FilterObject{},
				&HuffmanTree{},
				&HuffmanTree{},
				[]*HuffmanTree{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildHuffmanTree(tt.args.sub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildHuffmanTree() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
*/
