package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"log"
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
	filterFile = app.Flag("filter-file", "A CSV file contains filter and notify.").Short('f').Default("filters.csv").String()
	idFile     = app.Flag("id-file", "A gob file contains ids.").Short('i').Default("ids.gob").String()
	treeFile   = app.Flag("tree-file", "Indicate the filename for the Trie.").Short('t').Default("tree.gob").String()

	// kingpin patricia command
	patricia         = app.Command("patricia", "Run in Patricia Trie filtering mode.")
	patriciaShowTrie = patricia.Flag("print", "Show the Patricia Trie.").Short('p').Default("false").Bool()
	patriciaRebuild  = patricia.Flag("rebuild", "Rebuild the Patricia Trie.").Short('r').Default("false").Bool()

	// kingpin dumb command
	dumb = app.Command("dumb", "Run in dumb filter mode.")
)

func runDumb(idFile string, fm filter.FilterMap) {
	ids := new([][]byte)
	if err := binutil.Load(idFile, ids); err != nil {
		panic(err)
	}
	//fmt.Printf("Loaded %v ids from %v\n", len(*ids), idFile)
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
		var head *filter.PatriciaTrie

		// Tree encode
		_, err := os.Stat(*treeFile)
		if *patriciaRebuild || os.IsNotExist(err) {
			var tree bytes.Buffer
			head = filter.BuildPatriciaTrie(fm)
			enc := gob.NewEncoder(&tree)
			err = enc.Encode(head)
			if err != nil {
				log.Fatal("encode:", err)
			}
			// Save to file
			file, err := os.Create(*treeFile)
			if err != nil {
				log.Fatal("file:", err)
			}
			file.Write(tree.Bytes())
			file.Close()
			log.Print("Saved the Patricia Trie to ", *treeFile)
		} else {
			// Tree decode
			binutil.Load(*treeFile, &head)
			log.Print("Loaded the Patricia Trie from ", *treeFile)
		}

		runPatricia(*idFile, head)
	case dumb.FullCommand():
		fm := loadFiltersFromCSVFile(*filterFile)
		log.Printf("Loaded %v filters from %s\n", len(fm), *filterFile)
		runDumb(*idFile, fm)
	}
}
