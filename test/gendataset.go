// Generate Tag data sets
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak/scheme"
)

type UIIParam struct {
	Type                  string
	Scheme                string
	CompanyPrefix         string
	ItemReference         string
	AssetType             string
	OwnerCode             string
	DataIdentifier        string
	IssuingAgencyCode     string
	CompanyIdentification string
	ExtDigits             int
	IARMaxDigits          int
}

var (
	NumRepeat = 89
	NumSerial = 10
	DIR       = os.Getenv("GOPATH") + fmt.Sprintf("/src/github.com/iomz/gosstrak/test/data/benchset%v", NumRepeat*10/100*100)
)

func filterWriter(fs chan string) {
	f, err := os.OpenFile(path.Join(DIR, "filters.csv"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for {
		select {
		case s := <-fs:
			if len(s) != 0 {
				w.WriteString(s + "\n")
				w.Flush()
			}
		}
	}
}

func generateNonZeroDigits(len int) (s string) {
	ok := false
	for !ok {
		s = binutil.GenerateNLengthDigitString(len)
		if !strings.HasPrefix(s, "0") {
			ok = true
		}
	}
	return s
}

func generateUIISet(wg *sync.WaitGroup, q chan UIIParam, fq chan string) {
	defer wg.Done()
	for {
		param, ok := <-q
		if !ok {
			return
		}
		// do stuff
		switch param.Scheme {
		case "sgtin-96":
			schemeDir := path.Join(DIR, param.Scheme, param.CompanyPrefix)
			os.MkdirAll(schemeDir, 0755)
			fileName := path.Join(schemeDir, param.ItemReference)
			if _, err := os.Stat(fileName); !os.IsNotExist(err) {
				break
			}
			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}

			bs, opt := scheme.MakeEPC(true, param.Scheme, "3", param.CompanyPrefix, "", "", "", "", "")
			fq <- scheme.PrintID(bs, opt)
			bs, opt = scheme.MakeEPC(true, param.Scheme, "3", param.CompanyPrefix, param.ItemReference, "", "", "", "")
			fq <- scheme.PrintID(bs, opt)

			w := bufio.NewWriter(f)
			for ser := 0; ser < NumSerial; ser++ {
				bs, opt := scheme.MakeEPC(false, param.Scheme, "3", param.CompanyPrefix, param.ItemReference, "", strconv.Itoa(ser), "", "")
				w.WriteString(scheme.PrintID(bs, opt) + "\n")
			}
			w.Flush()
			f.Close()
		case "sscc-96":
			schemeDir := path.Join(DIR, param.Scheme)
			os.MkdirAll(schemeDir, 0755)
			fileName := path.Join(schemeDir, param.CompanyPrefix)
			if _, err := os.Stat(fileName); !os.IsNotExist(err) {
				break
			}
			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}

			bs, opt := scheme.MakeEPC(true, param.Scheme, "3", param.CompanyPrefix, "", "", "", "", "")
			fq <- scheme.PrintID(bs, opt)

			w := bufio.NewWriter(f)
			for ext := 0; ext < NumSerial; ext++ {
				bs, opt := scheme.MakeEPC(false, param.Scheme, "3", param.CompanyPrefix, "", generateNonZeroDigits(param.ExtDigits), "", "", "")
				w.WriteString(scheme.PrintID(bs, opt) + "\n")
			}
			w.Flush()
			f.Close()
		case "giai-96":
			schemeDir := path.Join(DIR, param.Scheme)
			os.MkdirAll(schemeDir, 0755)
			fileName := path.Join(schemeDir, param.CompanyPrefix)
			if _, err := os.Stat(fileName); !os.IsNotExist(err) {
				break
			}
			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}

			bs, opt := scheme.MakeEPC(true, param.Scheme, "3", param.CompanyPrefix, "", "", "", "", "")
			fq <- scheme.PrintID(bs, opt)

			w := bufio.NewWriter(f)
			for iar := 0; iar < NumSerial; iar++ {
				iarLen := binutil.GenerateRandomInt(1, param.IARMaxDigits)
				bs, opt := scheme.MakeEPC(false, param.Scheme, "3", param.CompanyPrefix, "", "", "", generateNonZeroDigits(iarLen), "")
				w.WriteString(scheme.PrintID(bs, opt) + "\n")
			}
			w.Flush()
			f.Close()
		case "grai-96":
			schemeDir := path.Join(DIR, param.Scheme)
			os.MkdirAll(schemeDir, 0755)
			fileName := path.Join(schemeDir, param.CompanyPrefix)
			if _, err := os.Stat(fileName); !os.IsNotExist(err) {
				log.Fatal(err)
			}
			f, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}

			bs, opt := scheme.MakeEPC(true, param.Scheme, "3", param.CompanyPrefix, "", "", "", "", "")
			fq <- scheme.PrintID(bs, opt)
			bs, opt = scheme.MakeEPC(true, param.Scheme, "3", param.CompanyPrefix, "", "", "", "", param.AssetType)
			fq <- scheme.PrintID(bs, opt)

			w := bufio.NewWriter(f)
			for ser := 0; ser < NumSerial; ser++ {
				bs, opt := scheme.MakeEPC(false, param.Scheme, "3", param.CompanyPrefix, "", "", strconv.Itoa(ser), "", param.AssetType)
				w.WriteString(scheme.PrintID(bs, opt) + "\n")
			}
			w.Flush()
			f.Close()
		case "17363":
			schemeDir := path.Join(DIR, param.Type+param.Scheme)
			os.MkdirAll(schemeDir, 0755)
			fileName := path.Join(schemeDir, param.OwnerCode)
			if _, err := os.Stat(fileName); !os.IsNotExist(err) {
				break
			}
			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}

			bs, opt := scheme.MakeISO(true, param.Scheme, param.OwnerCode, "", "", "", "", "", "")
			fq <- scheme.PrintID(bs, opt)

			w := bufio.NewWriter(f)
			for ser := 0; ser < NumSerial; ser++ {
				bs, opt := scheme.MakeISO(false, param.Scheme, param.OwnerCode, "", strconv.Itoa(ser), "", "", "", "")
				w.WriteString(scheme.PrintID(bs, opt) + "\n")
			}
			w.Flush()
			f.Close()
		case "17365":
			schemeDir := path.Join(DIR, param.Type+param.Scheme)
			os.MkdirAll(schemeDir, 0755)
			fileName := path.Join(schemeDir, param.CompanyIdentification)
			if _, err := os.Stat(fileName); !os.IsNotExist(err) {
				break
			}
			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}

			bs, opt := scheme.MakeISO(true, param.Scheme, "", "", "", param.DataIdentifier, param.IssuingAgencyCode, "", "")
			fq <- scheme.PrintID(bs, opt)
			bs, opt = scheme.MakeISO(true, param.Scheme, "", "", "", param.DataIdentifier, param.IssuingAgencyCode, param.CompanyIdentification, "")
			fq <- scheme.PrintID(bs, opt)

			w := bufio.NewWriter(f)
			for ser := 0; ser < NumSerial; ser++ {
				serLen := binutil.GenerateRandomInt(10, 30)
				bs, opt := scheme.MakeISO(false, param.Scheme, "", "", "", param.DataIdentifier, param.IssuingAgencyCode, param.CompanyIdentification, binutil.GenerateNLengthAlphanumericString(serLen))
				w.WriteString(scheme.PrintID(bs, opt) + "\n")
			}
			w.Flush()
			f.Close()
		}
	}
}

func main() {
	// create the target dir
	os.MkdirAll(DIR, 0755)
	log.Printf("Create target directory: %s", DIR)

	// prepare the workers
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(runtime.NumCPU())

	// prepare the filter writer
	fq := make(chan string)
	go filterWriter(fq)

	q := make(chan UIIParam, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go generateUIISet(&wg, q, fq)
	}

	var cpLen int
	for i := 0; i < NumRepeat; i++ {
		log.Printf("Iteration: %v\n", i)
		cpLen = binutil.GenerateRandomInt(6, 12)
		q <- UIIParam{
			Type:          "epc",
			Scheme:        "sgtin-96",
			CompanyPrefix: binutil.GenerateNLengthDigitString(cpLen),
			ItemReference: binutil.GenerateNLengthDigitString(13 - cpLen),
		}
		cpLen = binutil.GenerateRandomInt(6, 12)
		q <- UIIParam{
			Type:          "epc",
			Scheme:        "sscc-96",
			CompanyPrefix: binutil.GenerateNLengthDigitString(cpLen),
			ExtDigits:     17 - cpLen,
		}
		cpLen = binutil.GenerateRandomInt(6, 12)
		q <- UIIParam{
			Type:          "epc",
			Scheme:        "sscc-96",
			CompanyPrefix: binutil.GenerateNLengthDigitString(cpLen),
			ExtDigits:     17 - cpLen,
		}
		cpLen = binutil.GenerateRandomInt(6, 12)
		q <- UIIParam{
			Type:          "epc",
			Scheme:        "giai-96",
			CompanyPrefix: binutil.GenerateNLengthDigitString(cpLen),
			IARMaxDigits:  25 - cpLen,
		}
		cpLen = binutil.GenerateRandomInt(6, 11)
		q <- UIIParam{
			Type:          "epc",
			Scheme:        "grai-96",
			CompanyPrefix: binutil.GenerateNLengthDigitString(cpLen),
			AssetType:     binutil.GenerateNLengthDigitString(12 - cpLen),
		}
		q <- UIIParam{
			Type:      "iso",
			Scheme:    "17363",
			OwnerCode: binutil.GenerateNLengthAlphabetString(3),
		}
		ciLen := binutil.GenerateRandomInt(3, 7)
		q <- UIIParam{
			Type:                  "iso",
			Scheme:                "17365",
			DataIdentifier:        "25S",
			IssuingAgencyCode:     "U",
			CompanyIdentification: binutil.GenerateNLengthAlphanumericString(ciLen),
		}
	}
	close(q)

	wg.Wait()
}
