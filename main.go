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
	filterFile = app.Flag("filterFile", "A CSV file contains filter and notify.").Short('f').Default("filters.csv").String()
	idFile     = app.Flag("idFile", "A gob file contains ids.").Short('i').Default("ids.gob").String()

	// kingpin patricia command
	patricia         = app.Command("patricia", "Run in Patricia Trie filtering mode.")
	patriciaShowTrie = patricia.Flag("patriciaShowTrie", "Show the Patricia Trie.").Short('p').Default("false").Bool()

	// kingpin dumb command
	dumb = app.Command("dumb", "Run in dumb filter mode.")
)

func runDumb(idFile string, fm filter.FilterMap) {
	ids := new([][]byte)
	if err := binutil.Load(idFile, ids); err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %v ids from %v\n", len(*ids), idFile)
	matches := map[string][]string{}
	for _, id := range *ids {
		i := binutil.ParseByteSliceToBinString(id)
		for f, n := range fm {
			if strings.HasPrefix(i, f) {
				if _, ok := matches[n]; !ok {
					matches[n] = []string{}
				}
				matches[n] = append(matches[n], i)
			}
		}
	}

	//for n, m := range matches {
	//	fmt.Printf("%v: %v\n", n, len(m))
	//}
}

func runPatricia(idFile string, head *filter.PatriciaTrie) {
	if *patriciaShowTrie {
		fmt.Println(head.Dump())
	}
	ids := new([][]byte)
	if err := binutil.Load(idFile, ids); err != nil {
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
		// prefix as key, notify string as value
		fm[record[1]] = record[0]
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
		runPatricia(*idFile, head)
	case dumb.FullCommand():
		fm := loadFiltersFromCSVFile(*filterFile)
		fmt.Printf("Loaded %v filters from %s\n", len(fm), *filterFile)
		runDumb(*idFile, fm)
	}
}
