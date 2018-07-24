// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"os"
	"testing"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
)

func Test_parse6BitEncodedByteSliceToString(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"11000010 00001000: 0 + padding (10 00001000)", args{[]byte{194, 8}}, "0", false},
		{"10010110 01101010: %& + padding (10)", args{[]byte{150, 106}}, "%&", false},
		{"00000100 00100000 11100000 10000010: ABC + padding (100000 10000010)", args{[]byte{4, 32, 224, 130}}, "ABC", false},
		{"11000111 00101100 11110100 10000010: 1234 + padding (10000010)", args{[]byte{199, 44, 244, 130}}, "1234", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse6BitEncodedByteSliceToString(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse6BitEncodedByteSliceToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parse6BitEncodedByteSliceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_core_Translate(t *testing.T) {
	type fields struct {
		epcTDSVersion string
	}
	type args struct {
		pc []byte
		id []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"SGTIN-96_3_1_12345678_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 112, 94, 48, 167, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:12345678.00001.1",
			false,
		},
		{
			"SGTIN-96_3_1_12345678901_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 100, 91, 251, 131, 134, 160, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:12345678901.01.1",
			false,
		},
		{
			"SGTIN-96_3_1_1234567890_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 104, 73, 150, 2, 210, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:1234567890.001.1",
			false,
		},
		{
			"SGTIN-96_3_1_123456789_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 108, 117, 188, 209, 80, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:123456789.0001.1",
			false,
		},
		{
			"SGTIN-96_3_1_12345678_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 112, 94, 48, 167, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:12345678.00001.1",
			false,
		},
		{
			"SGTIN-96_3_1_1234567_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 116, 75, 90, 28, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:1234567.000001.1",
			false,
		},
		{
			"SGTIN-96_3_1_123456_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{48, 120, 120, 144, 0, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:123456.0000001.1",
			false,
		},
		{
			"SSCC-96_3_0_123456789012_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 96, 114, 250, 100, 104, 80, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:123456789012.00001",
			false,
		},
		{
			"SSCC-96_3_0_12345678901_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 100, 91, 251, 131, 134, 160, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:12345678901.000001",
			false,
		}, {
			"SSCC-96_3_0_1234567890_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 104, 73, 150, 2, 210, 0, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:1234567890.0000001",
			false,
		}, {
			"SSCC-96_3_0_123456789_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 108, 117, 188, 209, 80, 0, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:123456789.00000001",
			false,
		}, {
			"SSCC-96_3_0_12345678_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 112, 94, 48, 167, 0, 0, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:12345678.000000001",
			false,
		}, {
			"SSCC-96_3_0_1234567_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 116, 75, 90, 28, 0, 0, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:1234567.0000000001",
			false,
		}, {
			"SSCC-96_3_0_123456_1",
			fields{""},
			args{[]byte{48, 0}, []byte{49, 120, 120, 144, 0, 0, 0, 0, 1, 0, 0, 0}},
			"urn:epc:id:sscc:123456.00000000001",
			false,
		}, {
			"GIAI-96_3_0_123456789012_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 96, 114, 250, 100, 104, 80, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:123456789012.12345",
			false,
		}, {
			"GIAI-96_3_0_12345678901_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 100, 91, 251, 131, 134, 160, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:12345678901.12345",
			false,
		}, {
			"GIAI-96_3_0_1234567890_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 104, 73, 150, 2, 210, 0, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:1234567890.12345",
			false,
		}, {
			"GIAI-96_3_0_123456789_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 108, 117, 188, 209, 80, 0, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:123456789.12345",
			false,
		}, {
			"GIAI-96_3_0_12345678_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 112, 94, 48, 167, 0, 0, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:12345678.12345",
			false,
		}, {
			"GIAI-96_3_0_1234567_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 116, 75, 90, 28, 0, 0, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:1234567.12345",
			false,
		}, {
			"GIAI-96_3_0_123456_12345",
			fields{""},
			args{[]byte{48, 0}, []byte{52, 120, 120, 144, 0, 0, 0, 0, 0, 0, 48, 57}},
			"urn:epc:id:giai:123456.12345",
			false,
		}, {
			"GRAI-96_3_0_123456789012_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 96, 114, 250, 100, 104, 80, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:123456789012.1",
			false,
		}, {
			"GRAI-96_3_0_12345678901_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 100, 91, 251, 131, 134, 160, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:12345678901.1.1",
			false,
		}, {
			"GRAI-96_3_0_1234567890_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 104, 73, 150, 2, 210, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:1234567890.01.1",
			false,
		}, {
			"GRAI-96_3_0_123456789_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 108, 117, 188, 209, 80, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:123456789.001.1",
			false,
		}, {
			"GRAI-96_3_0_12345678_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 112, 94, 48, 167, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:12345678.0001.1",
			false,
		}, {
			"GRAI-96_3_0_1234567_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 116, 75, 90, 28, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:1234567.00001.1",
			false,
		}, {
			"GRAI-96_3_0_123456_1_1",
			fields{""},
			args{[]byte{48, 0}, []byte{51, 120, 120, 144, 0, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:grai:123456.000001.1",
			false,
		}, {
			"ISO17363_7B_ABC_U_1234560",
			fields{""},
			args{[]byte{41, 169}, []byte{220, 32, 66, 13, 92, 114, 207, 77, 118, 194}},
			"urn:epc:id:iso17363:7BABCU1234560",
			false,
		},
		{
			"ISO17365_25S_UN_ABC_0THANK0YOU0FOR0READING0THIS1",
			fields{""},
			args{[]byte{113, 162}, []byte{203, 84, 213, 56, 16, 131, 193, 66, 1, 56, 188, 25, 61, 92, 6, 61, 44, 18, 20, 17, 9, 56, 124, 20, 32, 148, 241, 130}},
			"urn:epc:id:iso17365:25SUNABC0THANK0YOU0FOR0READING0THIS1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Core{
				epcTDSVersion: tt.fields.epcTDSVersion,
			}
			got, err := c.Translate(tt.args.pc, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("core.Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("core.Translate() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func benchmarkTranslateNTags(nTags int, b *testing.B) {
	largeTagsGOB := os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/test/data/bench-100subs-tags.gob"
	// load up the tags from the file
	var largeTags llrp.Tags
	binutil.Load(largeTagsGOB, &largeTags)
	tdtCore := NewCore()

	var limitedTags []*llrp.ReadEvent
	perms := rand.Perm(len(largeTags))
	for count, i := range perms {
		if count < nTags {
			t := largeTags[i]
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.BigEndian, t.PCBits)
			if err != nil {
				b.Fatal(err)
			}
			limitedTags = append(limitedTags, &llrp.ReadEvent{PC: buf.Bytes(), ID: t.EPC})
		} else {
			break
		}
		if count == len(largeTags) {
			b.Skip("given tag size is larger than the testdata available")
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tag := range limitedTags {
			pureIdentity, err := tdtCore.Translate(tag.PC, tag.ID)
			if err != nil {
				b.Fatal(err)
			}
			_ = pureIdentity
		}
	}
}

func BenchmarkTranslate100Tags(b *testing.B)  { benchmarkTranslateNTags(100, b) }
func BenchmarkTranslate200Tags(b *testing.B)  { benchmarkTranslateNTags(200, b) }
func BenchmarkTranslate300Tags(b *testing.B)  { benchmarkTranslateNTags(300, b) }
func BenchmarkTranslate400Tags(b *testing.B)  { benchmarkTranslateNTags(400, b) }
func BenchmarkTranslate500Tags(b *testing.B)  { benchmarkTranslateNTags(500, b) }
func BenchmarkTranslate600Tags(b *testing.B)  { benchmarkTranslateNTags(600, b) }
func BenchmarkTranslate700Tags(b *testing.B)  { benchmarkTranslateNTags(700, b) }
func BenchmarkTranslate800Tags(b *testing.B)  { benchmarkTranslateNTags(800, b) }
func BenchmarkTranslate900Tags(b *testing.B)  { benchmarkTranslateNTags(900, b) }
func BenchmarkTranslate1000Tags(b *testing.B) { benchmarkTranslateNTags(1000, b) }
