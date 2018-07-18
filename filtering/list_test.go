// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package filtering

import (
	"reflect"
	"testing"
	//"github.com/iomz/gosstrak/tdt"
)

func TestList_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		list    *List
		want    []byte
		wantErr bool
	}{
	/*
		{
			"simple marshal",
			&List{
				ListFilters{
					&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
					&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				},
				tdt.NewCore(),
			},
			[]byte{},
			false,
		},
	*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("List.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.MarshalBinary() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestList_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		list    *List
		args    args
		wantErr bool
	}{
	/*
		{
			"simple unmarshal",
			&List{
				ListFilters{
					&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
					&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				},
				tdt.NewCore(),
			},
			args{
				[]byte{},
			},
			false,
		},
	*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.list.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("List.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListFilters_IndexOf(t *testing.T) {
	type args struct {
		em *ExactMatch
	}
	tests := []struct {
		name string
		lf   ListFilters
		args args
		want int
	}{
		{
			"Contains true",
			ListFilters{
				&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
				&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				&ExactMatch{NewFilter("001100110000", 0), "http://localhost:8888/3-3-0"},
				&ExactMatch{NewFilter("1111", 0), "http://localhost:8888/15"},
			},
			args{
				&ExactMatch{NewFilter("1111", 0), "http://localhost:8888/15"},
			},
			3,
		},
		{
			"Contains false",
			ListFilters{
				&ExactMatch{NewFilter("0011", 0), "http://localhost:8888/3"},
				&ExactMatch{NewFilter("00110000", 0), "http://localhost:8888/3-0"},
				&ExactMatch{NewFilter("001100110000", 0), "http://localhost:8888/3-3-0"},
				&ExactMatch{NewFilter("1111", 0), "http://localhost:8888/15"},
			},
			args{
				&ExactMatch{NewFilter("11", 0), "http://localhost:8888/3"},
			},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lf.IndexOf(tt.args.em); got != tt.want {
				t.Errorf("ListFilters.IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Name(t *testing.T) {
	tests := []struct {
		name string
		list *List
		want string
	}{
		{
			"List.Name",
			&List{},
			"List",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Name(); got != tt.want {
				t.Errorf("List.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
