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
	largeSubsCSV := os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/test/data/large-subs.csv"
	largeSubs := LoadFiltersFromCSVFile(largeSubsCSV)

	// cap the Subscriptions with the given size
	limitedSubs := Subscriptions{}
	keys := largeSubs.Keys()
	rand.Seed(time.Now().UTC().UnixNano())
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
	b.Logf("the resulting engine size: %v bytes", buf.Len())
}

func BenchmarkEngineGenList1000(b *testing.B) {
	benchmarkEngineGeneration(1000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList2000(b *testing.B) {
	benchmarkEngineGeneration(2000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList3000(b *testing.B) {
	benchmarkEngineGeneration(3000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList4000(b *testing.B) {
	benchmarkEngineGeneration(4000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList5000(b *testing.B) {
	benchmarkEngineGeneration(5000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList6000(b *testing.B) {
	benchmarkEngineGeneration(6000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList7000(b *testing.B) {
	benchmarkEngineGeneration(7000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList8000(b *testing.B) {
	benchmarkEngineGeneration(8000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList9000(b *testing.B) {
	benchmarkEngineGeneration(9000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenList10000(b *testing.B) {
	benchmarkEngineGeneration(10000, AvailableEngines["List"], b)
}

func BenchmarkEngineGenPatricia1000(b *testing.B) {
	benchmarkEngineGeneration(1000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia2000(b *testing.B) {
	benchmarkEngineGeneration(2000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia3000(b *testing.B) {
	benchmarkEngineGeneration(3000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia4000(b *testing.B) {
	benchmarkEngineGeneration(4000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia5000(b *testing.B) {
	benchmarkEngineGeneration(5000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia6000(b *testing.B) {
	benchmarkEngineGeneration(6000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia7000(b *testing.B) {
	benchmarkEngineGeneration(7000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia8000(b *testing.B) {
	benchmarkEngineGeneration(8000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia9000(b *testing.B) {
	benchmarkEngineGeneration(9000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenPatricia10000(b *testing.B) {
	benchmarkEngineGeneration(10000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkEngineGenSplay1000(b *testing.B) {
	benchmarkEngineGeneration(1000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay2000(b *testing.B) {
	benchmarkEngineGeneration(2000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay3000(b *testing.B) {
	benchmarkEngineGeneration(3000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay4000(b *testing.B) {
	benchmarkEngineGeneration(4000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay5000(b *testing.B) {
	benchmarkEngineGeneration(5000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay6000(b *testing.B) {
	benchmarkEngineGeneration(6000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay7000(b *testing.B) {
	benchmarkEngineGeneration(7000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay8000(b *testing.B) {
	benchmarkEngineGeneration(8000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay9000(b *testing.B) {
	benchmarkEngineGeneration(9000, AvailableEngines["SplayTree"], b)
}

func BenchmarkEngineGenSplay10000(b *testing.B) {
	benchmarkEngineGeneration(10000, AvailableEngines["SplayTree"], b)
}
