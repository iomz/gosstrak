package filtering

import (
	"reflect"
	"testing"
)

func TestList_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		list    *List
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("List.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.MarshalBinary() = %v, want %v", got, tt.want)
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.list.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("List.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestList_AnalyzeLocality(t *testing.T) {
	type args struct {
		id     []byte
		prefix string
		lm     *LocalityMap
	}
	tests := []struct {
		name string
		list *List
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.AnalyzeLocality(tt.args.id, tt.args.prefix, tt.args.lm)
		})
	}
}

func TestList_Search(t *testing.T) {
	type args struct {
		id []byte
	}
	tests := []struct {
		name        string
		list        *List
		args        args
		wantMatches []string
	}{
		{
			"0011xxxx on []byte{60, 128}",
			&List{
				&ExactMatch{"3", NewFilter("0011", 0)},
			},
			args{[]byte{60, 128}},
			[]string{"3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMatches := tt.list.Search(tt.args.id); !reflect.DeepEqual(gotMatches, tt.wantMatches) {
				t.Errorf("List.Search() = %v, want %v", gotMatches, tt.wantMatches)
			}
		})
	}
}

func TestBuildList(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			"BuildList testing...",
			args{
				Subscriptions{
					"0011":         &Info{"3", 10, nil},
					"1111":         &Info{"15", 2, nil},
					"00110000":     &Info{"3-0", 5, nil},
					"001100110000": &Info{"3-3-0", 5, nil},
				},
			},
			&List{
				&ExactMatch{"3", NewFilter("0011", 0)},
				&ExactMatch{"3-0", NewFilter("00110000", 0)},
				&ExactMatch{"3-3-0", NewFilter("001100110000", 0)},
				&ExactMatch{"15", NewFilter("1111", 0)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildList(tt.args.sub); !reflect.DeepEqual(got, tt.want) {
				for i, em := range *got {
					if !reflect.DeepEqual(em.filter, (*tt.want)[i].filter) {
						t.Errorf("(*BuildList())[%v].filter = \n%v, want \n%v", i, em.filter, (*tt.want)[i].filter)
					} else if em.notificationURI != (*tt.want)[i].notificationURI {
						t.Errorf("(*BuildList())[%v].notificationURI = \n%v, want \n%v", i, em.notificationURI, (*tt.want)[i].notificationURI)
					}
				}
			}
		})
	}
}
