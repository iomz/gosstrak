// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/iomz/gosstrak-fc/filtering"
)

func Benchmark_BuildPatriciaTrie(b *testing.B) {
	fm := loadFiltersFromCSVFile("filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(fm)
}

func Benchmark_runDumb(b *testing.B) {
	fm := loadFiltersFromCSVFile("filters.csv")
	b.ResetTimer()
	runDumb("ids.gob", fm)
}

func Benchmark_runPatricia(b *testing.B) {
	head := loadPatriciaTrie("filters.csv", "", true)
	b.ResetTimer()
	execute("ids.gob", head, "out.gob")
}

func Test_loadFiltersFromCSVFile(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want filtering.Subscriptions
	}{
		{"SGTIN-96_3_3_458960468_102",
			args{getPackagePath() + "/testdata/filters.csv"},
			filtering.Subscriptions{
				"0011000001101101101101011011001011100101010000000001100110": &filtering.Info{"SGTIN-96_3_3_458960468_102", 0},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadFiltersFromCSVFile(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFiltersFromCSVFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPackagePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"os.Getenv(\"GOPATH\") + \"/src/github.com/iomz/gosstrak-fc\"", os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak-fc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPackagePath(); got != tt.want {
				t.Errorf("getPackagePath() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_runAnalyzePatricia(t *testing.T) {
	type args struct {
		head    *filtering.PatriciaTrie
		inFile  string
		outFile string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runAnalyzePatricia(tt.args.head, tt.args.inFile, tt.args.outFile)
		})
	}
}

func Test_runDumb(t *testing.T) {
	type args struct {
		idFile string
		sub    filtering.Subscriptions
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDumb(tt.args.idFile, tt.args.sub)
		})
	}
}

func Test_execute(t *testing.T) {
	type args struct {
		idFile  string
		head    filtering.Engine
		outFile string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execute(tt.args.idFile, tt.args.head, tt.args.outFile)
		})
	}
}

func Test_loadHuffmanTree(t *testing.T) {
	type args struct {
		filterFile   string
		engineFile   string
		isRebuilding bool
	}
	tests := []struct {
		name string
		args args
		want *filtering.HuffmanTree
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadHuffmanTree(tt.args.filterFile, tt.args.engineFile, tt.args.isRebuilding); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadHuffmanTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadPatriciaTrie(t *testing.T) {
	type args struct {
		filterFile   string
		engineFile   string
		isRebuilding bool
	}
	tests := []struct {
		name string
		args args
		want *filtering.PatriciaTrie
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadPatriciaTrie(tt.args.filterFile, tt.args.engineFile, tt.args.isRebuilding); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadPatriciaTrie() = %v, want %v", got, tt.want)
			}
		})
	}
}
