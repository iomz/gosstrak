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
	"github.com/iomz/gosstrak-fc/filtering"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Current Version
	version = "0.1.1"

	// app
	app = kingpin.New("gosstrak-fc", "An RFID middleware to replace Fosstrak F&C.")
	// common flag
	verbose    = app.Flag("debug", "Enable verbose mode.").Short('v').Default("false").Bool()
	filterFile = app.Flag("filter-file", "A CSV file contains filter and notify.").Short('f').Default("filters.csv").String()
	idFile     = app.Flag("id-file", "A gob file contains ids.").Short('i').Default("ids.gob").String()
	treeFile   = app.Flag("tree-file", "Indicate the filename for the Trie.").Short('t').Default("tree.gob").String()
	outFile    = app.Flag("out-file", "Indicate the filename for whatever output.").Short('o').Default("out.gob").String()

	// patricia command
	patricia         = app.Command("patricia", "Run in Patricia Trie filtering mode.")
	patriciaShowTrie = patricia.Flag("print", "Show the Patricia Trie.").Short('p').Default("false").Bool()
	patriciaRebuild  = patricia.Flag("rebuild", "Rebuild the Patricia Trie.").Short('r').Default("false").Bool()

	// dumb command
	dumb = app.Command("dumb", "Run in dumb filter mode.")

	// analyze command
	analyze       = app.Command("analyze", "Analyze the tree node reference locality.")
	analyzeEngine = analyze.Flag("engine", "Used filtering engine for the target.").Default("patricia").String()
	analyzeInput  = analyze.Flag("analyze-input", "A gob file contains results in NotifyMap.").Default("out.gob").String()
	analyzeOutput = analyze.Flag("analyze-output", "A JSON file for d3.js.").Default("public/patricia/locality.json").String()
)

func runAnalyzePatricia(head *filtering.PatriciaTrie, inFile string, outFile string) {
	matches := new(filtering.NotifyMap)
	if err := binutil.Load(inFile, matches); err != nil {
		panic(err)
	}
	log.Printf("Loaded %v notifies from %s\n", len(*matches), inFile)
	lm := filtering.LocalityMap{}
	for _, ids := range *matches {
		for _, id := range ids {
			head.AnalyzeLocality(id, "", &lm)
		}
	}
	// Save to file
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal("file:", err)
	}
	file.Write(lm.ToJSON())
	file.Close()
	log.Print("Saved the Patricia Trie locality to ", outFile)
}

func runDumb(idFile string, fm filtering.Map) {
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

func runPatricia(idFile string, head *filtering.PatriciaTrie, outFile string) {
	if *patriciaShowTrie {
		fmt.Println(head.Dump())
	}
	ids := new([][]byte)
	if err := binutil.Load(idFile, ids); err != nil {
		panic(err)
	}
	notifies := filtering.NotifyMap{}
	for _, id := range *ids {
		matches := head.Search(id)
		for _, n := range matches {
			if _, ok := notifies[n]; !ok {
				notifies[n] = [][]byte{}
			}
			notifies[n] = append(notifies[n], id)
		}
	}
	binutil.Save(outFile, notifies)
	log.Print("Saved the result to ", outFile)
}

func loadFiltersFromCSVFile(f string) filtering.Map {
	fm := filtering.Map{}
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

func loadPatriciaTrie(filterFile string, treeFile string, isRebuilding bool) *filtering.PatriciaTrie {
	var head *filtering.PatriciaTrie
	// Tree encode
	_, err := os.Stat(treeFile)
	if isRebuilding || os.IsNotExist(err) {
		fm := loadFiltersFromCSVFile(filterFile)
		log.Printf("Loaded %v filters from %s\n", len(fm), filterFile)
		var tree bytes.Buffer
		head = filtering.BuildPatriciaTrie(fm)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal("encode:", err)
		}
		// Save to file
		file, err := os.Create(treeFile)
		if err != nil {
			log.Fatal("file:", err)
		}
		file.Write(tree.Bytes())
		file.Close()
		log.Print("Saved the Patricia Trie to ", treeFile)
	} else {
		// Tree decode
		binutil.Load(treeFile, &head)
		log.Print("Loaded the Patricia Trie from ", treeFile)
	}
	return head
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parse {
	case patricia.FullCommand():
		head := loadPatriciaTrie(*filterFile, *treeFile, *patriciaRebuild)
		runPatricia(*idFile, head, *outFile)
	case dumb.FullCommand():
		fm := loadFiltersFromCSVFile(*filterFile)
		log.Printf("Loaded %v filters from %s\n", len(fm), *filterFile)
		runDumb(*idFile, fm)
	case analyze.FullCommand():
		switch strings.ToLower(*analyzeEngine) {
		case "patricia":
			// Force analyze mode to need the tree file
			head := loadPatriciaTrie(*filterFile, *treeFile, false)
			runAnalyzePatricia(head, *analyzeInput, *analyzeOutput)
		}
	}
}
