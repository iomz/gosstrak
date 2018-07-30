// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestByteSubscriptions_keys(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want []string
	}{
		{
			"0,8",
			ByteSubscriptions{
				"0000": &PartialSubscription{0, "0", ByteSubscriptions{}},
				"1000": &PartialSubscription{0, "8", ByteSubscriptions{}},
			},
			[]string{"0000", "1000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteSubscriptions.keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSubscriptions_Dump(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want string
	}{
		{
			"Test Dump ByteSubscriptions",
			ByteSubscriptions{
				"0011":         &PartialSubscription{0, "3", ByteSubscriptions{}},
				"00110011":     &PartialSubscription{0, "3-3", ByteSubscriptions{}},
				"1111":         &PartialSubscription{0, "15", ByteSubscriptions{}},
				"00110000":     &PartialSubscription{0, "3-0", ByteSubscriptions{}},
				"001100110000": &PartialSubscription{0, "3-3-0", ByteSubscriptions{}},
			},
			"--0011 0 3\n" +
				"--00110000 0 3-0\n" +
				"--00110011 0 3-3\n" +
				"--001100110000 0 3-3-0\n" +
				"--1111 0 15\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sub.Dump(); got != tt.want {
				t.Errorf("ByteSubscriptions.Dump() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestByteSubscriptions_linkSubset(t *testing.T) {
	tests := []struct {
		name string
		sub  ByteSubscriptions
		want ByteSubscriptions
	}{
		{
			"Subset linking test for ByteSubscriptions",
			ByteSubscriptions{
				"0011":         &PartialSubscription{0, "3", ByteSubscriptions{}},
				"00110011":     &PartialSubscription{0, "3-3", ByteSubscriptions{}},
				"1111":         &PartialSubscription{0, "15", ByteSubscriptions{}},
				"00110000":     &PartialSubscription{0, "3-0", ByteSubscriptions{}},
				"001100110000": &PartialSubscription{0, "3-3-0", ByteSubscriptions{}},
			},
			ByteSubscriptions{
				"0011": &PartialSubscription{0, "3", ByteSubscriptions{
					"0000": &PartialSubscription{4, "3-0", ByteSubscriptions{}},
					"0011": &PartialSubscription{4, "3-3", ByteSubscriptions{
						"0000": &PartialSubscription{8, "3-3-0", ByteSubscriptions{}},
					}},
				}},
				"1111": &PartialSubscription{0, "15", ByteSubscriptions{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sub.linkSubset()
			if tt.sub.Dump() != tt.want.Dump() {
				t.Errorf("ByteSubscriptions.linkSubset() -> \n%v, want \n%v", tt.sub.Dump(), tt.want.Dump())
			}
		})
	}
}

func TestSubscriptions_ToByteSubscriptions(t *testing.T) {
	tests := []struct {
		name string
		sub  Subscriptions
		want ByteSubscriptions
	}{
		{
			"ecspec_sample.csv -> ByteSubscriptions",
			Subscriptions{
				"http://localhost:8888/grai":  []string{"urn:epc:pat:grai-96:3.123456.1.1"},
				"http://localhost:8888/17365": []string{"urn:epc:pat:iso17365:25S.UN.ABC.0THANK0YOU0FOR0READING0THIS1"},
				"http://localhost:8888/giai":  []string{"urn:epc:pat:giai-96:3.02283922192.45325296932379"},
				"http://localhost:8888/17363": []string{"urn:epc:pat:iso17363:7B.MTR"},
				"http://localhost:8888/sgtin": []string{"urn:epc:pat:sgtin-96:3.999203.7757355"},
				"http://localhost:8888/sscc":  []string{"urn:epc:pat:sscc-96:3.00039579721"},
			},
			ByteSubscriptions{
				"0011000001111011110011111100100011011101100101111000101011":                                                                                                                                                               &PartialSubscription{Offset: 0, ReportURI: "http://localhost:8888/sgtin"},
				"001100010110010000000000010010110111111000001001001":                                                                                                                                                                      &PartialSubscription{Offset: 0, ReportURI: "http://localhost:8888/sscc"},
				"001100110111100001111000100100000000000000000000000000000100000000000000000000000000000000000001":                                                                                                                         &PartialSubscription{Offset: 0, ReportURI: "http://localhost:8888/grai"},
				"0011010001100100000100010000010000111100011000100001010010011100100011110001110010001011000011011":                                                                                                                        &PartialSubscription{Offset: 0, ReportURI: "http://localhost:8888/giai"},
				"110010110101010011010101001110000001000010000011110000010100001000000001001110001011110000011001001111010101110000000110001111010010110000010010000101000001000100001001001110000111110000010100001000001001010011110001": &PartialSubscription{Offset: 0, ReportURI: "http://localhost:8888/17365"},
				"110111000010001101010100010010": &PartialSubscription{Offset: 0, ReportURI: "http://localhost:8888/17363"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sub.ToByteSubscriptions()
			for pfs, psub := range tt.want {
				if gotPsub, ok := got[pfs]; !ok {
					t.Errorf("Subscriptions.ToByteSubscriptions() =  want %v", pfs)
				} else if gotPsub.Offset != psub.Offset {
					t.Errorf("Subscriptions.ToByteSubscriptions() = %q, want %q", gotPsub, psub)
				} else if gotPsub.ReportURI != psub.ReportURI {
					t.Errorf("Subscriptions.ToByteSubscriptions() = %q, want %q", gotPsub, psub)
				}
			}
		})
	}
}

func TestLoadSubscriptionsFromCSVFile(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want Subscriptions
	}{
		{
			"ecspec_sample.csv",
			args{os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/test/data/ecspec_sample.csv"},
			Subscriptions{
				"http://localhost:8888/grai":  []string{"urn:epc:pat:grai-96:3.123456.1.1"},
				"http://localhost:8888/17365": []string{"urn:epc:pat:iso17365:25S.UN.ABC.0THANK0YOU0FOR0READING0THIS1"},
				"http://localhost:8888/giai":  []string{"urn:epc:pat:giai-96:3.02283922192.45325296932379"},
				"http://localhost:8888/17363": []string{"urn:epc:pat:iso17363:7B.MTR"},
				"http://localhost:8888/sgtin": []string{"urn:epc:pat:sgtin-96:3.999203.7757355"},
				"http://localhost:8888/sscc":  []string{"urn:epc:pat:sscc-96:3.00039579721"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadSubscriptionsFromCSVFile(tt.args.f)
			for reportURI, patterns := range got {
				if wantedPatterns, ok := tt.want[reportURI]; !ok {
					t.Errorf("LoadSubscriptionsFromCSVFile() = unknown key %v", reportURI)
				} else if !reflect.DeepEqual(patterns, wantedPatterns) {
					t.Errorf("LoadSubscriptionsFromCSVFile() = %q, want %q", patterns, wantedPatterns)
				}
			}
		})
	}
}

func benchmarkLoadNSubs(nSubs int, b *testing.B) {
	var sub Subscriptions
	for i := 0; i < b.N; i++ {
		sub = LoadSubscriptionsFromCSVFile(os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-ecspec.csv", nSubs))
	}
	b.StopTimer()

	// measure the size of the generated engine
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(sub)
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("the resulting engine size: %v bytes", buf.Len())
}

func BenchmarkEngineGenLegacy100Subs(b *testing.B)  { benchmarkLoadNSubs(100, b) }
func BenchmarkEngineGenLegacy200Subs(b *testing.B)  { benchmarkLoadNSubs(200, b) }
func BenchmarkEngineGenLegacy300Subs(b *testing.B)  { benchmarkLoadNSubs(300, b) }
func BenchmarkEngineGenLegacy400Subs(b *testing.B)  { benchmarkLoadNSubs(400, b) }
func BenchmarkEngineGenLegacy500Subs(b *testing.B)  { benchmarkLoadNSubs(500, b) }
func BenchmarkEngineGenLegacy600Subs(b *testing.B)  { benchmarkLoadNSubs(600, b) }
func BenchmarkEngineGenLegacy700Subs(b *testing.B)  { benchmarkLoadNSubs(700, b) }
func BenchmarkEngineGenLegacy800Subs(b *testing.B)  { benchmarkLoadNSubs(800, b) }
func BenchmarkEngineGenLegacy900Subs(b *testing.B)  { benchmarkLoadNSubs(900, b) }
func BenchmarkEngineGenLegacy1000Subs(b *testing.B) { benchmarkLoadNSubs(1000, b) }
