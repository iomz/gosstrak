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
	"testing"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
)

/*
func TestNewSplayTree(t *testing.T) {
	type args struct {
		sub ByteSubscriptions
	}
	tests := []struct {
		name string
		args args
		want *SplayTree
	}{
		{
			"",
			args{
				ByteSubscriptions{
					"0011":         &PartialSubscription{0, "3", ByteSubscriptions{}},
					"00110000":     &PartialSubscription{0, "3-0", ByteSubscriptions{}},
					"00110011":     &PartialSubscription{0, "3-3", ByteSubscriptions{}},
					"001100110000": &PartialSubscription{0, "3-3-0", ByteSubscriptions{}},
					"1111":         &PartialSubscription{0, "15", ByteSubscriptions{}},
				},
			},
			&SplayTree{
				&SplayTreeNode{
					"3",
					NewFilter("0011", 0),
					&SplayTreeNode{
						"3-0",
						NewFilter("0000", 4),
						nil,
						&SplayTreeNode{
							"3-3",
							NewFilter("0011", 4),
							&SplayTreeNode{
								"3-3-0",
								NewFilter("0000", 8),
								nil,
								nil,
							},
							nil,
						},
					},
					&SplayTreeNode{
						"15",
						NewFilter("1111", 0),
						nil,
						nil,
					},
				},
				tdt.NewCore(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSplayTree(tt.args.sub); got.Dump() != tt.want.Dump() {
				t.Errorf("NewSplayTree() = \n%v, want \n%v", got.Dump(), tt.want.Dump())
			}
		})
	}
}

func TestSplayTree_Name(t *testing.T) {
	type fields struct {
		root *SplayTreeNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"SplayTree.Name",
			fields{},
			"SplayTree",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spt := &SplayTree{
				root: tt.fields.root,
			}
			if got := spt.Name(); got != tt.want {
				t.Errorf("SplayTree.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func benchmarkFilterSplayNTagsNSubs(nTags int, nSubs int, b *testing.B) {
	// build the engine
	sub := LoadSubscriptionsFromCSVFile(os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-ecspec.csv", nSubs))
	splayEngine := NewSplayTree(sub)

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
			pureIdentity, reportURIs, err := splayEngine.Search(*re)
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
func BenchmarkFilterSplay100Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 100, b) }
func BenchmarkFilterSplay200Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(200, 100, b) }
func BenchmarkFilterSplay300Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(300, 100, b) }
func BenchmarkFilterSplay400Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(400, 100, b) }
func BenchmarkFilterSplay500Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(500, 100, b) }
func BenchmarkFilterSplay600Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(600, 100, b) }
func BenchmarkFilterSplay700Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(700, 100, b) }
func BenchmarkFilterSplay800Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(800, 100, b) }
func BenchmarkFilterSplay900Tags100Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(900, 100, b) }
func BenchmarkFilterSplay1000Tags100Subs(b *testing.B) { benchmarkFilterSplayNTagsNSubs(1000, 100, b) }

// Impact from n_{S}
func BenchmarkFilterSplay100Tags200Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 200, b) }
func BenchmarkFilterSplay100Tags300Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 300, b) }
func BenchmarkFilterSplay100Tags400Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 400, b) }
func BenchmarkFilterSplay100Tags500Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 500, b) }
func BenchmarkFilterSplay100Tags600Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 600, b) }
func BenchmarkFilterSplay100Tags700Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 700, b) }
func BenchmarkFilterSplay100Tags800Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 800, b) }
func BenchmarkFilterSplay100Tags900Subs(b *testing.B)  { benchmarkFilterSplayNTagsNSubs(100, 900, b) }
func BenchmarkFilterSplay100Tags1000Subs(b *testing.B) { benchmarkFilterSplayNTagsNSubs(100, 1000, b) }
