// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"testing"
)

func benchmarkEngineGenerationFromNSubs(nSubs int, constructor EngineConstructor, b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		sub := LoadSubscriptionsFromCSVFile(os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-ecspec.csv", nSubs))
		b.StartTimer()
		engine := constructor(sub)
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
}

// List engine generation 100-1000
func BenchmarkEngineGenList100(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(100, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList200(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(200, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList300(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(300, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList400(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(400, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList500(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(500, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList600(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(600, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList700(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(700, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList800(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(800, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList900(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(900, AvailableEngines["List"], b)
}
func BenchmarkEngineGenList1000(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(1000, AvailableEngines["List"], b)
}

// Patricia engine generation 100-1000
func BenchmarkEngineGenPatricia100(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(100, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia200(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(200, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia300(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(300, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia400(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(400, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia500(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(500, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia600(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(600, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia700(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(700, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia800(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(800, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia900(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(900, AvailableEngines["PatriciaTrie"], b)
}
func BenchmarkEngineGenPatricia1000(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(1000, AvailableEngines["PatriciaTrie"], b)
}

// Splay engine generation 100-1000
func BenchmarkEngineGenSplay100(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(100, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay200(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(200, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay300(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(300, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay400(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(400, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay500(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(500, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay600(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(600, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay700(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(700, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay800(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(800, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay900(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(900, AvailableEngines["SplayTree"], b)
}
func BenchmarkEngineGenSplay1000(b *testing.B) {
	benchmarkEngineGenerationFromNSubs(1000, AvailableEngines["SplayTree"], b)
}
