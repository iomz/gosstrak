// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

import (
	"reflect"
	"testing"
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

func TestNewCore(t *testing.T) {
	tests := []struct {
		name string
		want *core
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_core_Translate(t *testing.T) {
	type fields struct {
		epcTDSVersion string
	}
	type args struct {
		id []byte
		pc []byte
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
			args{[]byte{48, 112, 94, 48, 167, 0, 0, 64, 0, 0, 0, 1}, []byte{48, 0}},
			"urn:epc:id:sgtin:3.12345678.1.1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &core{
				epcTDSVersion: tt.fields.epcTDSVersion,
			}
			got, err := c.Translate(tt.args.id, tt.args.pc)
			if (err != nil) != tt.wantErr {
				t.Errorf("core.Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("core.Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_core_buildEPC(t *testing.T) {
	type fields struct {
		epcTDSVersion string
	}
	type args struct {
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
			"SGTIN-96_3_0_123456789012_1_1",
			fields{""},
			args{[]byte{48, 96, 114, 250, 100, 104, 80, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.123456789012.1.1",
			false,
		},
		{
			"SGTIN-96_3_1_12345678901_1_1",
			fields{""},
			args{[]byte{48, 100, 91, 251, 131, 134, 160, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.12345678901.1.1",
			false,
		},
		{
			"SGTIN-96_3_1_1234567890_1_1",
			fields{""},
			args{[]byte{48, 104, 73, 150, 2, 210, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.1234567890.1.1",
			false,
		},
		{
			"SGTIN-96_3_1_123456789_1_1",
			fields{""},
			args{[]byte{48, 108, 117, 188, 209, 80, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.123456789.1.1",
			false,
		},
		{
			"SGTIN-96_3_1_12345678_1_1",
			fields{""},
			args{[]byte{48, 112, 94, 48, 167, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.12345678.1.1",
			false,
		},
		{
			"SGTIN-96_3_1_1234567_1_1",
			fields{""},
			args{[]byte{48, 116, 75, 90, 28, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.1234567.1.1",
			false,
		},
		{
			"SGTIN-96_3_1_123456_1_1",
			fields{""},
			args{[]byte{48, 120, 120, 144, 0, 0, 0, 64, 0, 0, 0, 1}},
			"urn:epc:id:sgtin:3.123456.1.1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &core{
				epcTDSVersion: tt.fields.epcTDSVersion,
			}
			got, err := c.buildEPC(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("core.buildEPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("core.buildEPC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_core_LoadEPCTagDataTranslation(t *testing.T) {
	type fields struct {
		epcTDSVersion string
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &core{
				epcTDSVersion: tt.fields.epcTDSVersion,
			}
			c.LoadEPCTagDataTranslation()
		})
	}
}

func Test_core_buildUII(t *testing.T) {
	type fields struct {
		epcTDSVersion string
	}
	type args struct {
		id  []byte
		afi byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &core{
				epcTDSVersion: tt.fields.epcTDSVersion,
			}
			got, err := c.buildUII(tt.args.id, tt.args.afi)
			if (err != nil) != tt.wantErr {
				t.Errorf("core.buildUII() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("core.buildUII() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_core_buildProprietary(t *testing.T) {
	type fields struct {
		epcTDSVersion string
	}
	type args struct {
		id []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &core{
				epcTDSVersion: tt.fields.epcTDSVersion,
			}
			got, err := c.buildProprietary(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("core.buildProprietary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("core.buildProprietary() = %v, want %v", got, tt.want)
			}
		})
	}
}
