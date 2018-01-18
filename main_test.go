// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak-fc/filtering"
)

func Benchmark_buildList15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildList(sub)
}

func Benchmark_buildPatricia15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(sub)
}

func Benchmark_buildOBST15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildOptimalBST(&sub)
}

func Benchmark_buildSplay15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildSplayTree(&sub)
}

func Benchmark_buildList148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildList(sub)
}

func Benchmark_buildPatricia148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(sub)
}

func Benchmark_buildOBST148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildOptimalBST(&sub)
}

func Benchmark_buildSplay148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildSplayTree(&sub)
}

func Benchmark_buildList44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildList(sub)
}

/* Run too long
func Benchmark_buildPatricia44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(sub)
}
*/

func Benchmark_buildOBST44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildOptimalBST(&sub)
}

func Benchmark_buildSplay44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildSplayTree(&sub)
}

func Benchmark_filterList15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildList(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterPatricia15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildPatriciaTrie(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterOBST15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildOptimalBST(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterSplay15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildSplayTree(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterList148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildList(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterPatricia148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildPatriciaTrie(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterOBST148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildOptimalBST(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterSplay148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildSplayTree(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterList44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildList(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

/* Run too long
func Benchmark_filterPatricia44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildPatriciaTrie(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}
*/

func Benchmark_filterOBST44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildOptimalBST(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func Benchmark_filterSplay44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildSplayTree(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	b.ResetTimer()
	for _, id := range *ids {
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

/*
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
*/

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
