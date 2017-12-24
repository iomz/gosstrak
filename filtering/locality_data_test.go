// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func TestLocalityData_MarshalJSON(t *testing.T) {
	type fields struct {
		name     string
		locality float32
		parent   *LocalityData
		children []*LocalityData
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := &LocalityData{
				name:     tt.fields.name,
				locality: tt.fields.locality,
				parent:   tt.fields.parent,
				children: tt.fields.children,
			}
			got, err := ld.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalityData.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalityData.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalityData_JSON(t *testing.T) {
	type fields struct {
		name     string
		locality float32
		parent   *LocalityData
		children []*LocalityData
	}
	tests := []struct {
		name   string
		fields fields
		want   LocalityDataJSON
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := &LocalityData{
				name:     tt.fields.name,
				locality: tt.fields.locality,
				parent:   tt.fields.parent,
				children: tt.fields.children,
			}
			if got := ld.JSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalityData.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalityData_InsertLocality(t *testing.T) {
	type fields struct {
		name     string
		locality float32
		parent   *LocalityData
		children []*LocalityData
	}
	type args struct {
		path     []string
		locality float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{ /*
		{
			"test insertLocality",
			fields{
			},
			args{
			},
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := &LocalityData{
				name:     tt.fields.name,
				locality: tt.fields.locality,
				parent:   tt.fields.parent,
				children: tt.fields.children,
			}
			ld.InsertLocality(tt.args.path, tt.args.locality)
		})
	}
}
