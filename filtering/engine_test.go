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
	limitedSubs := Subscriptions{}
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

func BenchmarkListGen10000(b *testing.B) {
	benchmarkEngineGeneration(10000, AvailableEngines["List"], b)
}

func BenchmarkListGen20000(b *testing.B) {
	benchmarkEngineGeneration(20000, AvailableEngines["List"], b)
}

func BenchmarkListGen30000(b *testing.B) {
	benchmarkEngineGeneration(30000, AvailableEngines["List"], b)
}

func BenchmarkListGen40000(b *testing.B) {
	benchmarkEngineGeneration(40000, AvailableEngines["List"], b)
}

func BenchmarkListGen50000(b *testing.B) {
	benchmarkEngineGeneration(50000, AvailableEngines["List"], b)
}

func BenchmarkListGen60000(b *testing.B) {
	benchmarkEngineGeneration(60000, AvailableEngines["List"], b)
}

func BenchmarkListGen70000(b *testing.B) {
	benchmarkEngineGeneration(70000, AvailableEngines["List"], b)
}

func BenchmarkListGen80000(b *testing.B) {
	benchmarkEngineGeneration(80000, AvailableEngines["List"], b)
}

func BenchmarkListGen90000(b *testing.B) {
	benchmarkEngineGeneration(90000, AvailableEngines["List"], b)
}

func BenchmarkListGen100000(b *testing.B) {
	benchmarkEngineGeneration(100000, AvailableEngines["List"], b)
}

func BenchmarkPatriciaTrieGen1000(b *testing.B) {
	benchmarkEngineGeneration(1000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen2000(b *testing.B) {
	benchmarkEngineGeneration(2000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen3000(b *testing.B) {
	benchmarkEngineGeneration(3000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen4000(b *testing.B) {
	benchmarkEngineGeneration(4000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen5000(b *testing.B) {
	benchmarkEngineGeneration(5000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen6000(b *testing.B) {
	benchmarkEngineGeneration(6000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen7000(b *testing.B) {
	benchmarkEngineGeneration(7000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen8000(b *testing.B) {
	benchmarkEngineGeneration(8000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen9000(b *testing.B) {
	benchmarkEngineGeneration(9000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkPatriciaTrieGen10000(b *testing.B) {
	benchmarkEngineGeneration(10000, AvailableEngines["PatriciaTrie"], b)
}

func BenchmarkSplayTreeGen1000(b *testing.B) {
	benchmarkEngineGeneration(1000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen2000(b *testing.B) {
	benchmarkEngineGeneration(2000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen3000(b *testing.B) {
	benchmarkEngineGeneration(3000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen4000(b *testing.B) {
	benchmarkEngineGeneration(4000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen5000(b *testing.B) {
	benchmarkEngineGeneration(5000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen6000(b *testing.B) {
	benchmarkEngineGeneration(6000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen7000(b *testing.B) {
	benchmarkEngineGeneration(7000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen8000(b *testing.B) {
	benchmarkEngineGeneration(8000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen9000(b *testing.B) {
	benchmarkEngineGeneration(9000, AvailableEngines["SplayTree"], b)
}

func BenchmarkSplayTreeGen10000(b *testing.B) {
	benchmarkEngineGeneration(10000, AvailableEngines["SplayTree"], b)
}
