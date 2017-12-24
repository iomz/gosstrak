// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
)

func TestLocalityMap_ToJSON(t *testing.T) {
	tests := []struct {
		name string
		lm   LocalityMap
		want []byte
	}{
		{
			"simple locality map test",
			LocalityMap{
				"":            12,
				",0011":       4,
				",0011,00":    2,
				",0011,00,11": 1,
			},
			[]byte("[{\"name\":\"Entry Node\",\"value\":100,\"children\":[{\"name\":\"0011\",\"value\":33.333332,\"children\":[{\"name\":\"00\",\"value\":16.666666,\"children\":[{\"name\":\"11\",\"value\":8.333333,\"children\":null}]}]}]}]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lm.ToJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalityMap.ToJSON() = \n%v, want \n%v", string(got), string(tt.want))
			}
		})
	}
}

func TestLocalityData_MarshalJSON(t *testing.T) {
	type fields struct {
		name     string
		locality float32
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
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := &LocalityData{
				name:     tt.fields.name,
				locality: tt.fields.locality,
				children: tt.fields.children,
			}
			ld.InsertLocality(tt.args.path, tt.args.locality)
		})
	}
}
