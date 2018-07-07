// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package adapter

import (
	"math/rand"
	"os"
	"testing"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/golemu"
)

func TestUnmarshalROAccessReportBody(t *testing.T) {
	largeTagsGOB := os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/testdata/1000-tags.gob"
	size := 100
	// load up the tags from the file
	var largeTags golemu.Tags
	binutil.Load(largeTagsGOB, &largeTags)

	// cap the tags with the given size
	var limitedTags golemu.Tags
	perms := rand.Perm(len(largeTags))
	for n, i := range perms {
		if n < size {
			limitedTags = append(limitedTags, largeTags[i])
		} else {
			break
		}
		if n+1 == len(largeTags) {
			t.Fatal("given tag size is larger than the testdata available")
		}
	}

	// build ROAR message
	pdu := int(1500)
	trds := limitedTags.BuildTagReportDataStack(pdu)
	if len(trds) == 0 {
		t.Fatal("TagReportDataStack generation failed")
	}

	var res []*LLRPReadEvent
	for i, trd := range trds {
		roar := llrp.ROAccessReport(trd.Parameter, uint32(i))
		res = append(res, UnmarshalROAccessReportBody(roar[10:])...)
	}

	if len(res) != size {
		t.Errorf("decapsulateROAccessReport() = %v", res)
	}
}

func BenchmarkUnmarshalROAccessReportBody(b *testing.B) {
	largeTagsGOB := os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/testdata/million-tags.gob"
	// load up the tags from the file
	var largeTags golemu.Tags
	binutil.Load(largeTagsGOB, &largeTags)

	cycle := b.N / len(largeTags)
	remaining := b.N % len(largeTags)

	// cap the tags with the given size
	var limitedTags golemu.Tags
	perms := rand.Perm(len(largeTags))
	for n, i := range perms {
		if n < remaining {
			limitedTags = append(limitedTags, largeTags[i])
		} else {
			break
		}
		if n == len(largeTags) {
			b.Skip("given tag size is larger than the testdata available")
		}
	}

	// build ROAR message
	pdu := int(1500)
	trds := largeTags.BuildTagReportDataStack(pdu)
	if len(trds) == 0 {
		b.Fatal("TagReportDataStack generation was failed")
	}
	limitedTRDs := limitedTags.BuildTagReportDataStack(pdu)
	if len(limitedTRDs) == 0 && remaining != 0 {
		b.Logf("len(limitedTags): %v, len(limitedTRDs: %v", len(limitedTags), len(limitedTRDs))
		b.Logf("b.N: %v, cycle: %v, remaining: %v", b.N, cycle, remaining)
		b.Fatal("TagReportDataStack generation failed")
	}

	var res []*LLRPReadEvent
	b.ResetTimer()
	for c := 0; c < cycle; c++ {
		for i, trd := range trds {
			b.StopTimer()
			roar := llrp.ROAccessReport(trd.Parameter, uint32(i))
			b.StartTimer()
			res = append(res, UnmarshalROAccessReportBody(roar[10:])...)
		}
	}

	for i, trd := range limitedTRDs {
		b.StopTimer()
		roar := llrp.ROAccessReport(trd.Parameter, uint32(i))
		b.StartTimer()
		res = append(res, UnmarshalROAccessReportBody(roar[10:])...)
	}
	b.StopTimer()
	if b.N != len(res) {
		b.Fatal("LLRP unmarshaller failed")
	}
}
