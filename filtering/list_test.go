// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
)

func TestList_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		list    *List
		want    []byte
		wantErr bool
	}{
	/*
		{
			"simple marshal",
			&List{
				FilterLists{
					&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
					&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				},
				tdt.NewCore(),
			},
			[]byte{},
			false,
		},
	*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("List.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.MarshalBinary() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestList_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		list    *List
		args    args
		wantErr bool
	}{
	/*
		{
			"simple unmarshal",
			&List{
				FilterLists{
					&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
					&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				},
				tdt.NewCore(),
			},
			args{
				[]byte{},
			},
			false,
		},
	*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.list.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("List.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterLists_IndexOf(t *testing.T) {
	type args struct {
		em *ExactMatch
	}
	tests := []struct {
		name string
		lf   ListFilters
		args args
		want int
	}{
		{
			"Contains true",
			ListFilters{
				&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
				&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				&ExactMatch{NewFilter("001100110000", 0), "http://localhost:8888/3-3-0"},
				&ExactMatch{NewFilter("1111", 0), "http://localhost:8888/15"},
			},
			args{
				&ExactMatch{NewFilter("1111", 0), "http://localhost:8888/15"},
			},
			3,
		},
		{
			"Contains false",
			ListFilters{
				&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
				&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				&ExactMatch{NewFilter("001100110000", 0), "http://localhost:8888/3-3-0"},
				&ExactMatch{NewFilter("1111", 0), "http://localhost:8888/15"},
			},
			args{
				&ExactMatch{NewFilter("11", 0), "http://localhost:8888/3"},
			},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lf.IndexOf(tt.args.em); got != tt.want {
				t.Errorf("FilterLists.IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Name(t *testing.T) {
	tests := []struct {
		name string
		list *List
		want string
	}{
		{
			"List.Name",
			&List{},
			"List",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Name(); got != tt.want {
				t.Errorf("List.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func benchmarkFilterListNTagsNSubs(nTags int, nSubs int, b *testing.B) {
	// build the engine
	sub := LoadSubscriptionsFromCSVFile(os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-ecspec.csv", nSubs))
	listEngine := NewList(sub)

	// prepare the workload
	largeTagsGOB := os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-tags.gob", nSubs)
	var largeTags llrp.Tags
	binutil.Load(largeTagsGOB, &largeTags)

	var res []*llrp.ReadEvent
	perms := rand.Perm(len(largeTags))
	for count, i := range perms {
		if count < nTags {
			t := largeTags[i]
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.BigEndian, t.PCBits)
			if err != nil {
				b.Fatal(err)
			}
			res = append(res, &llrp.ReadEvent{PC: buf.Bytes(), ID: t.EPC})
		} else {
			break
		}
		if count == len(largeTags) {
			b.Skip("given tag size is larger than the testdata available")
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, re := range res {
			pureIdentity, reportURIs, err := listEngine.Search(*re)
			if err != nil {
				b.Error(err)
			}
			if len(reportURIs) == 0 {
				b.Errorf("no match found for %v", pureIdentity)
			}
		}
	}
}

// Impact from n_{E}
func BenchmarkFilterList100Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 100, b) }
func BenchmarkFilterList200Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(200, 100, b) }
func BenchmarkFilterList300Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(300, 100, b) }
func BenchmarkFilterList400Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(400, 100, b) }
func BenchmarkFilterList500Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(500, 100, b) }
func BenchmarkFilterList600Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(600, 100, b) }
func BenchmarkFilterList700Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(700, 100, b) }
func BenchmarkFilterList800Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(800, 100, b) }
func BenchmarkFilterList900Tags100Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(900, 100, b) }
func BenchmarkFilterList1000Tags100Subs(b *testing.B) { benchmarkFilterListNTagsNSubs(1000, 100, b) }

// Impact from n_{S}
func BenchmarkFilterList100Tags200Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 200, b) }
func BenchmarkFilterList100Tags300Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 300, b) }
func BenchmarkFilterList100Tags400Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 400, b) }
func BenchmarkFilterList100Tags500Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 500, b) }
func BenchmarkFilterList100Tags600Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 600, b) }
func BenchmarkFilterList100Tags700Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 700, b) }
func BenchmarkFilterList100Tags800Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 800, b) }
func BenchmarkFilterList100Tags900Subs(b *testing.B)  { benchmarkFilterListNTagsNSubs(100, 900, b) }
func BenchmarkFilterList100Tags1000Subs(b *testing.B) { benchmarkFilterListNTagsNSubs(100, 1000, b) }
