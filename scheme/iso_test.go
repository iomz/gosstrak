package scheme

import (
	"reflect"
	"testing"
)

func TestGetISO6346CD(t *testing.T) {
	type args struct {
		cn string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"CSQU305438", args{"CSQU305438"}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetISO6346CD(tt.args.cn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetISO6346CD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetISO6346CD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeISO17363(t *testing.T) {
	type args struct {
		pf  bool
		oc  string
		ei  string
		csn string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   int
		want2   string
		want3   string
		wantErr bool
	}{
		{"A97BCSQU3054383", args{false, "CSQ", "U", "305438"}, []byte{220, 32, 211, 69, 92, 240, 215, 76, 248, 206}, 80, "", "", false},
		{"A97BCSQU3054383", args{true, "CSQ", "U", "305438"}, []byte{}, 0, "110111000010000011010011010001010101110011110000110101110100110011111000110011", "urn:epc:pat:iso17363:7B.CSQ.U.305438", false},
		{"A97BCSQU", args{true, "CSQ", "U", ""}, []byte{}, 0, "110111000010000011010011010001010101", "urn:epc:pat:iso17363:7B.CSQ.U", false},
		{"A97BCSQ", args{true, "CSQ", "", ""}, []byte{}, 0, "110111000010000011010011010001", "urn:epc:pat:iso17363:7B.CSQ", false},
		{"A97B", args{true, "", "", ""}, []byte{}, 0, "110111000010", "urn:epc:pat:iso17363:7B", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, err := MakeISO17363(tt.args.pf, tt.args.oc, tt.args.ei, tt.args.csn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeISO17363() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeISO17363() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MakeISO17363() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("MakeISO17363() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("MakeISO17363() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestMakeISO17365(t *testing.T) {
	type args struct {
		pf  bool
		di  string
		iac string
		cin string
		sn  string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   int
		want2   string
		want3   string
		wantErr bool
	}{
		{"25SUN043325711MH8031200000000001", args{false, "25S", "UN", "043325711", "MH8031200000000001"}, []byte{203, 84, 213, 59, 13, 51, 207, 45, 119, 199, 19, 72, 227, 12, 241, 203, 12, 48, 195, 12, 48, 195, 12, 49}, 192, "", "", false},
		{"25SUN043325711MH8031200000000001", args{true, "25S", "UN", "043325711", "MH8031200000000001"}, []byte{}, 0, "110010110101010011010101001110110000110100110011110011110010110101110111110001110001001101001000111000110000110011110001110010110000110000110000110000110000110000110000110000110000110000110001", "urn:epc:pat:iso17365:25S.UN.043325711.MH8031200000000001", false},
		{"25SUN043325711", args{true, "25S", "UN", "043325711", ""}, []byte{}, 0, "110010110101010011010101001110110000110100110011110011110010110101110111110001110001", "urn:epc:pat:iso17365:25S.UN.043325711", false},
		{"25SUN", args{true, "25S", "UN", "", ""}, []byte{}, 0, "110010110101010011010101001110", "urn:epc:pat:iso17365:25S.UN", false},
		{"25SUN", args{true, "25S", "", "", ""}, []byte{}, 0, "110010110101010011", "urn:epc:pat:iso17365:25S", false},
		{"25SODCIN10000000RTIA1B2C3DOSN12345", args{false, "25S", "OD", "CIN1", "0000000RTIA1B2C3DOSN12345"}, []byte{203, 84, 207, 16, 50, 78, 199, 12, 48, 195, 12, 48, 73, 66, 65, 196, 44, 131, 204, 67, 211, 59, 28, 179, 211, 88}, 208, "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, err := MakeISO17365(tt.args.pf, tt.args.di, tt.args.iac, tt.args.cin, tt.args.sn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeISO17365() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeISO17365() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MakeISO17365() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("MakeISO17365() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("MakeISO17365() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestPad6BitEncodingRuneSlice(t *testing.T) {
	type args struct {
		bs []rune
	}
	tests := []struct {
		name  string
		args  args
		want  []rune
		want1 int
	}{
		{"0000", args{[]rune("0000")}, []rune("0000100000100000"), 16},
		{"0000000000000000", args{[]rune("0000000000000000")}, []rune("0000000000000000"), 16},
		//{"", args{""}, []rune{""}, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Pad6BitEncodingRuneSlice(tt.args.bs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pad6BitEncodingRuneSlice() got = %v, want %v", string(got), string(tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("Pad6BitEncodingRuneSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
