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
	execute("ids.gob", head, "out.gob", false)
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
				"0011000001101101101101011011001011100101010000000001100110": &filtering.Info{
					NotificationURI: "SGTIN-96_3_3_458960468_102",
					EntropyValue:    0,
					Subset:          nil,
				},
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
