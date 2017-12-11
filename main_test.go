package main

import (
	"reflect"
	"testing"

	"github.com/iomz/gosstrak-fc/filter"
)

func Benchmark_BuildPatriciaTrie(b *testing.B) {
	fm := loadFiltersFromCSVFile("filters.csv")
	b.ResetTimer()
	filter.BuildPatriciaTrie(fm)
}

func Benchmark_runDumb(b *testing.B) {
	fm := loadFiltersFromCSVFile("filters.csv")
	b.ResetTimer()
	runDumb(fm)
}

func Benchmark_runPatricia(b *testing.B) {
	fm := loadFiltersFromCSVFile("filters.csv")
	head := filter.BuildPatriciaTrie(fm)
	b.ResetTimer()
	runPatricia(head)
}

func Test_loadFiltersFromCSVFile(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want filter.FilterMap
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
