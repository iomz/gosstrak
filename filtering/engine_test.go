// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"math/rand"
	"os"
	"testing"
	//"unsafe"
	//"github.com/iomz/go-llrp"
	//"github.com/iomz/go-llrp/binutil"
	//"github.com/iomz/golemu"
)

func benchmarkEngineGeneration(size int, constructor EngineConstructor, b *testing.B) {
	// Load up the subs from the file
	largeSubsCSV := os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/testdata/large-subs.csv"
	largeSubs := LoadFiltersFromCSVFile(largeSubsCSV)

	// cap the Subscriptions with the given size
	var limitedSubs Subscriptions
	keys := largeSubs.Keys()
	perms := rand.Perm(len(keys))
	for n, i := range perms {
		if n < size {
			limitedSubs[keys[i]] = largeSubs[keys[i]]
		} else {
			break
		}
		if n+1 == len(keys) {
			b.Fatal("given subscription size is larger than the testdata available")
		}
	}

	b.ResetTimer()
	engine := constructor(limitedSubs)
	b.StopTimer()

	// measure the size of the generated engine
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(engine)
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("the allocated memory size is: %v", len(buf.Bytes()))
}

func BenchmarkListGen100(b *testing.B) {
	benchmarkEngineGeneration(100, AvailableEngines["List"], b)
}
