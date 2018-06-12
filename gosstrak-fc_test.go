// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak/filtering"
)

// Build engines -----------------------------------------------------

func BenchmarkBuildList15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildList(sub)
}

func BenchmarkBuildPatricia15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(sub)
}

func BenchmarkBuildOBST15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildOptimalBST(&sub)
}

func BenchmarkBuildSplay15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	b.ResetTimer()
	filtering.BuildSplayTree(&sub)
}

func BenchmarkBuildList148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildList(sub)
}

func BenchmarkBuildPatricia148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(sub)
}

func BenchmarkBuildOBST148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildOptimalBST(&sub)
}

func BenchmarkBuildSplay148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	b.ResetTimer()
	filtering.BuildSplayTree(&sub)
}

func BenchmarkBuildList44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildList(sub)
}

/* Run too long
func BenchmarkBuildPatricia44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildPatriciaTrie(sub)
}
*/

func BenchmarkBuildOBST44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildOptimalBST(&sub)
}

func BenchmarkBuildSplay44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	b.ResetTimer()
	filtering.BuildSplayTree(&sub)
}

// Filter IDs -----------------------------------------------------

func BenchmarkFilterList15915_4763(b *testing.B) {
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

func BenchmarkFilterPatricia15915_4763(b *testing.B) {
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

func BenchmarkFilterOBST15915_4763(b *testing.B) {
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

func BenchmarkFilterSplay15915_4763(b *testing.B) {
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

func BenchmarkFilterList148825_721(b *testing.B) {
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

func BenchmarkFilterPatricia148825_721(b *testing.B) {
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

func BenchmarkFilterOBST148825_721(b *testing.B) {
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

func BenchmarkFilterSplay148825_721(b *testing.B) {
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

func BenchmarkFilterList44456_47561(b *testing.B) {
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
func BenchmarkFilterPatricia44456_47561(b *testing.B) {
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

func BenchmarkFilterOBST44456_47561(b *testing.B) {
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

func BenchmarkFilterSplay44456_47561(b *testing.B) {
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

// Filter IDs (Random) -----------------------------------------------------

func BenchmarkFilterListRandom15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildList(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterPatriciaRandom15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildPatriciaTrie(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterOBSTRandom15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildOptimalBST(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterSplayRandom15915_4763(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/15915-4763/filters.csv")
	engine := filtering.BuildSplayTree(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/15915-4763/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterListRandom148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildList(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterPatriciaRandom148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildPatriciaTrie(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterOBSTRandom148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildOptimalBST(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterSplayRandom148825_721(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/148825-721/filters.csv")
	engine := filtering.BuildSplayTree(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/148825-721/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterListRandom44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildList(sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterOBSTRandom44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildOptimalBST(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

func BenchmarkFilterSplayRandom44456_47561(b *testing.B) {
	sub := loadFiltersFromCSVFile("scenarios/44456-47561/filters.csv")
	engine := filtering.BuildSplayTree(&sub)
	ids := new([][]byte)
	binutil.Load("scenarios/44456-47561/ids.gob", ids)
	notifies := filtering.NotifyMap{}
	perms := rand.Perm(len(*ids))
	b.ResetTimer()
	for _, i := range perms {
		id := (*ids)[i]
		matches := engine.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
}

// Other tests for main.go -----------------------------------------------------

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
		{"os.Getenv(\"GOPATH\") + \"/src/github.com/iomz/gosstrak\"", os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPackagePath(); got != tt.want {
				t.Errorf("getPackagePath() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
