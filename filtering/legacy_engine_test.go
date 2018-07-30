package filtering

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak/tdt"
)

func TestLegacyEngine_AddSubscription(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
	}
	type args struct {
		sub Subscriptions
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
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			le.AddSubscription(tt.args.sub)
		})
	}
}

func TestLegacyEngine_DeleteSubscription(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
	}
	type args struct {
		sub Subscriptions
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
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			le.DeleteSubscription(tt.args.sub)
		})
	}
}

func TestLegacyEngine_Dump(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			if got := le.Dump(); got != tt.want {
				t.Errorf("LegacyEngine.Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyEngine_MarshalBinary(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
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
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			got, err := le.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("LegacyEngine.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LegacyEngine.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyEngine_Name(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			if got := le.Name(); got != tt.want {
				t.Errorf("LegacyEngine.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyEngine_Search(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
	}
	type args struct {
		re llrp.ReadEvent
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantPureIdentity string
		wantReportURIs   []string
		wantErr          bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			gotPureIdentity, gotReportURIs, err := le.Search(tt.args.re)
			if (err != nil) != tt.wantErr {
				t.Errorf("LegacyEngine.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPureIdentity != tt.wantPureIdentity {
				t.Errorf("LegacyEngine.Search() gotPureIdentity = %v, want %v", gotPureIdentity, tt.wantPureIdentity)
			}
			if !reflect.DeepEqual(gotReportURIs, tt.wantReportURIs) {
				t.Errorf("LegacyEngine.Search() gotReportURIs = %v, want %v", gotReportURIs, tt.wantReportURIs)
			}
		})
	}
}

func TestLegacyEngine_UnmarshalBinary(t *testing.T) {
	type fields struct {
		filters Subscriptions
		tdtCore *tdt.Core
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			le := &LegacyEngine{
				filters: tt.fields.filters,
				tdtCore: tt.fields.tdtCore,
			}
			if err := le.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("LegacyEngine.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewLegacyEngine(t *testing.T) {
	type args struct {
		sub Subscriptions
	}
	tests := []struct {
		name string
		args args
		want Engine
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLegacyEngine(tt.args.sub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLegacyEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringIndexInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringIndexInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("stringIndexInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func benchmarkFilterLegacyNTagsNSubs(nTags int, nSubs int, b *testing.B) {
	sub := LoadSubscriptionsFromCSVFile(os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-ecspec.csv", nSubs))
	largeTagsGOB := os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/bench-%vsubs-tags.gob", nSubs)
	var largeTags llrp.Tags
	binutil.Load(largeTagsGOB, &largeTags)
	tdtCore := tdt.NewCore()

	var res []*llrp.ReadEvent
	perms := rand.Perm(len(largeTags))
	for count, i := range perms {
		if count < nTags {
			t := largeTags[i]
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.BigEndian, t.PCBits)
			if err != nil {
				b.Fatal(err)
			}
			res = append(res, &llrp.ReadEvent{PC: buf.Bytes(), ID: t.EPC})
		} else {
			break
		}
		if count == len(largeTags) {
			b.Skip("given tag size is larger than the testdata available")
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, re := range res {
			// search start
			// translate the readevent to a PureIdentity
			b.StopTimer()
			pureIdentity, err := tdtCore.Translate(re.PC, re.ID)
			if err != nil {
				b.Error(err)
			}

			b.StartTimer()
			var reportURIs []string
			for reportURI, patterns := range sub {
				for _, pattern := range patterns {
					seq := strings.Split(pattern, ":")
					if len(seq) != 5 {
						continue
					}
					patternType := seq[3]
					pattern := seq[4]

					switch patternType {
					case "giai-96":
						fields := strings.Split(seq[4], ".")
						// remove filter value in tag uri to match with the received PureIdentity
						pattern = "giai:" + strings.Join(fields[1:], ".")
					case "grai-96":
						fields := strings.Split(seq[4], ".")
						// remove filter value in tag uri to match with the received PureIdentity
						pattern = "grai:" + strings.Join(fields[1:], ".")
					case "sgtin-96":
						fields := strings.Split(seq[4], ".")
						// remove filter value in tag uri to match with the received PureIdentity
						pattern = "sgtin:" + strings.Join(fields[1:], ".")
					case "sscc-96":
						fields := strings.Split(seq[4], ".")
						// remove filter value in tag uri to match with the received PureIdentity
						pattern = "sscc:" + strings.Join(fields[1:], ".")
					case "iso17363":
						pattern = patternType + ":" + strings.Replace(pattern, ".", "", -1)
					case "iso17365":
						pattern = patternType + ":" + strings.Replace(pattern, ".", "", -1)
					}
					if strings.HasPrefix(strings.TrimPrefix(pureIdentity, "urn:epc:id:"), pattern) {
						reportURIs = append(reportURIs, reportURI)
					}
				}
			}
			if len(reportURIs) == 0 {
				b.Errorf("no match found for %v", pureIdentity)
			}
		}
	}
}

// Impact from n_{E}
func BenchmarkFilterLegacy100Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 100, b) }
func BenchmarkFilterLegacy200Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(200, 100, b) }
func BenchmarkFilterLegacy300Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(300, 100, b) }
func BenchmarkFilterLegacy400Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(400, 100, b) }
func BenchmarkFilterLegacy500Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(500, 100, b) }
func BenchmarkFilterLegacy600Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(600, 100, b) }
func BenchmarkFilterLegacy700Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(700, 100, b) }
func BenchmarkFilterLegacy800Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(800, 100, b) }
func BenchmarkFilterLegacy900Tags100Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(900, 100, b) }
func BenchmarkFilterLegacy1000Tags100Subs(b *testing.B) { benchmarkFilterLegacyNTagsNSubs(1000, 100, b) }

// Impact from n_{S}
func BenchmarkFilterLegacy100Tags200Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 200, b) }
func BenchmarkFilterLegacy100Tags300Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 300, b) }
func BenchmarkFilterLegacy100Tags400Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 400, b) }
func BenchmarkFilterLegacy100Tags500Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 500, b) }
func BenchmarkFilterLegacy100Tags600Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 600, b) }
func BenchmarkFilterLegacy100Tags700Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 700, b) }
func BenchmarkFilterLegacy100Tags800Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 800, b) }
func BenchmarkFilterLegacy100Tags900Subs(b *testing.B)  { benchmarkFilterLegacyNTagsNSubs(100, 900, b) }
func BenchmarkFilterLegacy100Tags1000Subs(b *testing.B) { benchmarkFilterLegacyNTagsNSubs(100, 1000, b) }
