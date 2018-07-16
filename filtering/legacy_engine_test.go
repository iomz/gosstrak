package filtering

import (
	"reflect"
	"testing"

	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/tdt"
)

func TestLegacyEngine_Search(t *testing.T) {
	type fields struct {
		subscriptions map[string][]string
		tdtCore       *tdt.Core
	}
	type args struct {
		re *llrp.LLRPReadEvent
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantMatched      []string
		wantPureIdentity string
		wantErr          bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			le := &LegacyEngine{
				subscriptions: tt.fields.subscriptions,
				tdtCore:       tt.fields.tdtCore,
			}
			gotMatched, gotPureIdentity, err := le.Search(tt.args.re)
			if (err != nil) != tt.wantErr {
				t.Errorf("LegacyEngine.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMatched, tt.wantMatched) {
				t.Errorf("LegacyEngine.Search() gotMatched = %v, want %v", gotMatched, tt.wantMatched)
			}
			if gotPureIdentity != tt.wantPureIdentity {
				t.Errorf("LegacyEngine.Search() gotPureIdentity = %v, want %v", gotPureIdentity, tt.wantPureIdentity)
			}
		})
	}
}

func Test_convertLabelToURNPrefix(t *testing.T) {
	type args struct {
		label string
	}
	tests := []struct {
		name          string
		args          args
		wantUrnPrefix string
		wantErr       bool
	}{
		{
			"SGTIN-96_3_3_160772053_8516",
			args{"SGTIN-96_3_3_160772053_8516"},
			"urn:epc:id:sgtin:3.160772053.8516",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrnPrefix, err := convertLabelToURNPrefix(tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertLabelToURNPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUrnPrefix != tt.wantUrnPrefix {
				t.Errorf("convertLabelToURNPrefix() = \n%v, want \n%v", gotUrnPrefix, tt.wantUrnPrefix)
			}
		})
	}
}
