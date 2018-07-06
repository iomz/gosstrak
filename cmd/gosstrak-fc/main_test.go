// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import "testing"

func Test_getPackagePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPackagePath(); got != tt.want {
				t.Errorf("getPackagePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(tt.args.f)
		})
	}
}
