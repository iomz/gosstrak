package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak-fc/filter"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Current Version
	version = "0.1.1"

	// kingpin app
	app = kingpin.New("gosstrak-fc", "An RFID middleware to replace Fosstrak F&C.")
	// kingpin verbose mode flag
	verbose    = app.Flag("debug", "Enable verbose mode.").Short('v').Default("false").Bool()
	filterFile = app.Flag("filterFile", "A CSV file contains filter and notify.").Default("filters.csv").String()

	// kingpin patricia command
	patricia         = app.Command("patricia", "Run in Patricia Trie filtering mode.")
	patriciaShowTrie = patricia.Flag("patriciaShowTrie", "Show the Patricia Trie.").Short('p').Default("false").Bool()

	// kingpin dumb command
	dumb = app.Command("dumb", "Run in dumb filter mode.")
)

func runDumb(fm filter.FilterMap) {
	ids := new([][]byte)
	if err := binutil.Load("ids.gob", ids); err != nil {
		panic(err)
	}
	matched := make([]string, 0, len(*ids))
	for _, id := range *ids {
		i := binutil.ParseByteSliceToBinString(id)
		for f, n := range fm {
			if strings.HasPrefix(i, f) {
				matched = append(matched, n)
			}
		}
	}
	//fmt.Printf("Matched ids: %v\n", len(matched))
}

func runPatricia(head *filter.PatriciaTrie) {
	if *patriciaShowTrie {
		fmt.Println(head.Dump())
	}
	ids := new([][]byte)
	if err := binutil.Load("ids.gob", ids); err != nil {
		panic(err)
	}
	matched := make([]string, 0, len(*ids))
	for _, id := range *ids {
		head.Match(id, &matched)
	}
	//fmt.Printf("Matched ids: %v\n", len(matched))
}

func loadFiltersFromCSVFile(f string) filter.FilterMap {
	fm := filter.FilterMap{}
	fp, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fm[record[0]] = record[1]
	}
	return fm
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parse {
	case patricia.FullCommand():
		fm := loadFiltersFromCSVFile(*filterFile)
		head := filter.BuildPatriciaTrie(fm)
		runPatricia(head)
	case dumb.FullCommand():
		fm := loadFiltersFromCSVFile(*filterFile)
		runDumb(fm)
	}
}
