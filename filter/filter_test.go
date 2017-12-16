package filter

import (
	"reflect"
	"testing"
)

func Test_makeFilter(t *testing.T) {
	type args struct {
		bs     []rune
		offset int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 []rune
		want2 []rune
	}{
		{"0000 from 0", args{[]rune("0000"), 0}, 0, []rune("00001111"), []rune("00001111")},
		{"10000000 from 0", args{[]rune("10000000"), 0}, 0, []rune("10000000"), []rune("00000000")},
		{"100000001000 from 0", args{[]rune("100000001000"), 0}, 0, []rune("1000000010001111"), []rune("0000000000001111")},
		{"1000000010000000 from 0", args{[]rune("1000000010000000"), 0}, 0, []rune("1000000010000000"), []rune("0000000000000000")},
		{"100000001000000010 from 0", args{[]rune("100000001000000010"), 0}, 0, []rune("100000001000000010111111"), []rune("000000000000000000111111")},
		{"0 from 6", args{[]rune("0"), 6}, 0, []rune("11111101"), []rune("11111101")},
		{"0000 from 4", args{[]rune("0000"), 4}, 0, []rune("11110000"), []rune("11110000")},
		{"0010 from 6", args{[]rune("0010"), 6}, 0, []rune("1111110010111111"), []rune("1111110000111111")},
		{"0000 from 8", args{[]rune("0000"), 8}, 1, []rune("00001111"), []rune("00001111")},
		{"xx000000 from 0", args{[]rune("xx000000"), 0}, 0, []rune("11000000"), []rune("11000000")},
		{"000xx000 from 0", args{[]rune("000xx000"), 0}, 0, []rune("00011000"), []rune("00011000")},
		{"000000xx from 0", args{[]rune("000000xx"), 0}, 0, []rune("00000011"), []rune("00000011")},
		{"xx00xx from 1", args{[]rune("xx00xx"), 1}, 0, []rune("11100111"), []rune("11100111")},
		{"xxxx from 0", args{[]rune("xxxx"), 0}, 0, []rune("11111111"), []rune("11111111")},
		{"xxxxxxxx from 0", args{[]rune("xxxxxxxx"), 0}, 0, []rune("11111111"), []rune("11111111")},
		{"xxx from 5", args{[]rune("xxx"), 5}, 0, []rune("11111111"), []rune("11111111")},
		{"xx from 7", args{[]rune("xx"), 7}, 0, []rune("1111111111111111"), []rune("1111111111111111")},
		//{" from ", args{[]rune(""), 0}, 0, []rune(""), []rune("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := makeFilter(tt.args.bs, tt.args.offset)
			if got != tt.want {
				t.Errorf("makeFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("makeFilter() got1 = \n%v, want \n%v", string(got1), string(tt.want1))
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("makeFilter() got2 = \n%v, want \n%v", string(got2), string(tt.want2))
			}
		})
	}
}

func TestFilter_GetByteAt(t *testing.T) {
	type fields struct {
		String     string
		Size       int
		Offset     int
		ByteFilter []byte
		ByteMask   []byte
		ByteOffset int
		ByteSize   int
	}
	type args struct {
		bo int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    byte
		want1   byte
		wantErr bool
	}{
		{"00110011 from 0 at 0", fields{"00110011", 8, 0, []byte{51}, []byte{0}, 0, 1}, args{0}, byte(51), byte(0), false},
		{"00110011 from 0 at 1", fields{"00110011", 8, 0, []byte{51}, []byte{0}, 0, 1}, args{1}, byte(0), byte(0), true},
		{"00110011 from 4 at 0", fields{"00110011", 8, 4, []byte{243, 63}, []byte{240, 15}, 0, 2}, args{0}, byte(243), byte(240), false},
		{"00110011 from 4 at 1", fields{"00110011", 8, 4, []byte{243, 63}, []byte{240, 15}, 0, 2}, args{1}, byte(63), byte(15), false},
		{"00110011 from 4 at 3", fields{"00110011", 8, 4, []byte{243, 63}, []byte{240, 15}, 0, 2}, args{3}, byte(0), byte(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				String:     tt.fields.String,
				Size:       tt.fields.Size,
				Offset:     tt.fields.Offset,
				ByteFilter: tt.fields.ByteFilter,
				ByteMask:   tt.fields.ByteMask,
				ByteOffset: tt.fields.ByteOffset,
				ByteSize:   tt.fields.ByteSize,
			}
			got, got1, err := f.GetByteAt(tt.args.bo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter.GetByteAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Filter.GetByteAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Filter.GetByteAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFilter_HasByteAt(t *testing.T) {
	type fields struct {
		String     string
		Size       int
		Offset     int
		ByteFilter []byte
		ByteMask   []byte
		ByteOffset int
		ByteSize   int
	}
	type args struct {
		bo int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"00110011 from 0 at 0", fields{"00110011", 8, 0, []byte{51}, []byte{0}, 0, 1}, args{0}, true},
		{"00110011 from 4 at 1", fields{"00110011", 8, 4, []byte{243, 63}, []byte{240, 15}, 0, 2}, args{1}, true},
		{"00110011 from 8 at 0", fields{"00110011", 8, 8, []byte{51}, []byte{0}, 1, 1}, args{0}, false},
		{"00110011 from 8 at 1", fields{"00110011", 8, 8, []byte{51}, []byte{0}, 1, 1}, args{1}, true},
		{"001100110011 from 4 at 1", fields{"001100110011", 12, 4, []byte{243, 51}, []byte{240, 0}, 0, 2}, args{1}, true},
		{"0000 from 10 at 1", fields{"0000", 4, 10, []byte{195}, []byte{195}, 1, 1}, args{1}, true},
		{"0000 from 10 at 2", fields{"0000", 4, 10, []byte{195}, []byte{195}, 1, 1}, args{2}, false},
		{"11 from 15 at 0", fields{"11", 2, 15, []byte{255, 255}, []byte{254, 127}, 1, 2}, args{0}, false},
		{"11 from 15 at 1", fields{"11", 2, 15, []byte{255, 255}, []byte{254, 127}, 1, 2}, args{1}, true},
		{"11 from 15 at 2", fields{"11", 2, 15, []byte{255, 255}, []byte{254, 127}, 1, 2}, args{2}, true},
		{"11 from 15 at 3", fields{"11", 2, 15, []byte{255, 255}, []byte{254, 127}, 1, 2}, args{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				String:     tt.fields.String,
				Size:       tt.fields.Size,
				Offset:     tt.fields.Offset,
				ByteFilter: tt.fields.ByteFilter,
				ByteMask:   tt.fields.ByteMask,
				ByteOffset: tt.fields.ByteOffset,
				ByteSize:   tt.fields.ByteSize,
			}
			if got := f.HasByteAt(tt.args.bo); got != tt.want {
				t.Errorf("Filter.HasByteAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Match(t *testing.T) {
	type fields struct {
		String     string
		Size       int
		Offset     int
		ByteFilter []byte
		ByteMask   []byte
		ByteOffset int
		ByteSize   int
	}
	type args struct {
		id []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"01010101", fields{"0", 1, 2, []byte{223}, []byte{223}, 0, 1}, args{[]byte{85}}, true},
		{"01010101", fields{"1100", 4, 0, []byte{207}, []byte{15}, 0, 1}, args{[]byte{85}}, false},
		{"0000000011110000", fields{"0000000011111111", 16, 0, []byte{0, 255}, []byte{0, 15}, 0, 2}, args{[]byte{0, 240}}, true},
		{"000000001111000000000000", fields{"00111100", 8, 6, []byte{252, 243}, []byte{252, 3}, 0, 2}, args{[]byte{0, 240, 0}}, true},
		{"000000001111111100000000", fields{"0000", 4, 19, []byte{15}, []byte{15}, 2, 1}, args{[]byte{0, 255, 0}}, true},
		{"001100000111011000011110100011011101010000000000", fields{"1100001111010001101110101", 25, 13, []byte{254, 30, 141, 215}, []byte{248, 0, 0, 3}, 1, 4}, args{[]byte{48, 118, 30, 141, 212, 0}}, true},
		//{"", fields{"", 0, []byte{}, []byte{}, 0, 1}, args{[]byte{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				String:     tt.fields.String,
				Size:       tt.fields.Size,
				Offset:     tt.fields.Offset,
				ByteFilter: tt.fields.ByteFilter,
				ByteMask:   tt.fields.ByteMask,
				ByteOffset: tt.fields.ByteOffset,
				ByteSize:   tt.fields.ByteSize,
			}
			if got := f.Match(tt.args.id); got != tt.want {
				t.Errorf("Filter.match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_ToString(t *testing.T) {
	type fields struct {
		String     string
		Size       int
		Offset     int
		ByteFilter []byte
		ByteMask   []byte
		ByteOffset int
		ByteSize   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"0000 from 0", fields{"0000", 4, 0, []byte{15}, []byte{15}, 0, 1}, "0000(0 4)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				String:     tt.fields.String,
				Size:       tt.fields.Size,
				Offset:     tt.fields.Offset,
				ByteFilter: tt.fields.ByteFilter,
				ByteMask:   tt.fields.ByteMask,
				ByteOffset: tt.fields.ByteOffset,
				ByteSize:   tt.fields.ByteSize,
			}
			if got := f.ToString(); got != tt.want {
				t.Errorf("Filter.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFilter(t *testing.T) {
	type args struct {
		s string
		o int
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		{"00000000 from 0", args{"00000000", 0}, &Filter{"00000000", 8, 0, []byte{0}, []byte{0}, 0, 1}},
		{"0000xxxx from 0", args{"0000xxxx", 0}, &Filter{"0000xxxx", 8, 0, []byte{15}, []byte{15}, 0, 1}},
		{"0000 from 0", args{"0000", 0}, &Filter{"0000", 4, 0, []byte{15}, []byte{15}, 0, 1}},
		{"0000 from 4", args{"0000", 4}, &Filter{"0000", 4, 4, []byte{240}, []byte{240}, 0, 1}},
		{"0000000000000000 from 12", args{"0000000000000000", 12}, &Filter{"0000000000000000", 16, 12, []byte{240, 0, 15}, []byte{240, 0, 15}, 1, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFilter(tt.args.s, tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilter() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
