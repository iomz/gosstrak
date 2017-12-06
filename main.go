package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/iomz/gosstrak-fc/filter"
)

var (
	// Current Version
	version = "0.1.0"

	// kingpin app
	app = kingpin.New("gosstrak-fc", "An RFID middleware to replace Fosstrak F&C.")
	// kingpin verbose mode flag
	verbose = app.Flag("debug", "Enable verbose mode.").Short('v').Default("false").Bool()

	// kingpin patricia command
	patricia   = app.Command("patricia", "Run in Patricia Trie filtering mode.")
	filterFile = patricia.Flag("filterFile", "A CSV file contains filter and notify.").Default("filters.csv").String()
)

func runPatricia(f string) {
	fm := loadFiltersFromCSVFile(f)
	head := filter.BuildPatriciaTrie(fm)
	fmt.Println(head.Dump())
	return
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
		runPatricia(*filterFile)
	}
}
