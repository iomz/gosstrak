// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/gob"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak-fc/filtering"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Constant Values
const (
	// BufferSize is a general size for a buffer
	BufferSize = 10240
)

// Environmental variables
var (
	// Current Version
	version = "0.2.0"

	// Data cache directory
	dataCacheDir = "/var/tmp/gosstrak-fc-cache"

	// app
	app = kingpin.
		New("gosstrak-fc", "An RFID middleware to replace Fosstrak F&C.")

	// common flag
	verbose = app.
		Flag("debug", "Enable verbose mode.").
		Short('v').
		Default("false").
		Bool()
	filterFile = app.
			Flag("filter-file", "A CSV file contains filter and notify.").
			Short('f').
			Default("filters.csv").
			String()
	idFile = app.
		Flag("id-file", "A gob file contains ids.").
		Short('i').
		Default("ids.gob").
		String()
	engineFile = app.
			Flag("engine-file", "Indicate the filename storing the filtering engine.").
			Short('e').
			Default(dataCacheDir + "/engine.gob").
			String()
	outFile = app.
		Flag("out-file", "Indicate the filename for saving the result.").
		Short('o').
		Default(dataCacheDir + "/out.gob").String()
	isRebuilding = app.
			Flag("rebuild", "Rebuild the filtering engine.").
			Short('r').
			Default("false").
			Bool()
	// LLRP related values
	initialMessageID = app.
				Flag("initialMessageID", "The initial messageID to start from.").
				Short('m').
				Default("1000").
				Int()
	ip = app.
		Flag("ip", "LLRP listening address.").
		Short('a').
		Default("0.0.0.0").
		IP()
	port = app.
		Flag("port", "LLRP listening port.").
		Short('p').
		Default("5084").
		Int()

	// analyze command
	cmdAnalyze = app.
			Command("analyze", "Analyze the tree node reference locality.")
	analyzeEngine = cmdAnalyze.
			Flag("engine", "Used filtering engine for the target.").
			Default("patricia").
			String()
	analyzeInput = cmdAnalyze.
			Flag("analyze-input", "A gob file contains results in NotifyMap.").
			Default(dataCacheDir + "/out.gob").
			String()

	// dumb command
	cmdDumb = app.
		Command("dumb", "Run in dumb filter mode.")

	// list command
	cmdList = app.
		Command("list", "Use list of byte filtering engine.")

	// obst command
	cmdOBST = app.
		Command("obst", "Use Optimal Binary Search Tree filtering engine.")

	// patricia command
	cmdPatricia = app.
			Command("patricia", "Use Patricia Trie filtering engine.")

	// obst command
	cmdSplay = app.
			Command("splay", "Use Splay Tree filtering engine.")

	// run command
	cmdRun = app.
		Command("run", "Run the F&C middleware.")

	// Current messageID
	messageID = uint32(*initialMessageID)
)

func analyze(head filtering.Engine, inFile string, outFile string) {
	matches := new(filtering.NotifyMap)
	if err := binutil.Load(inFile, matches); err != nil {
		panic(err)
	}
	log.Printf("Loaded %v notifications from %s\n", len(*matches), inFile)
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
	log.Print("Saved the result to ", outFile)
}

func execute(idFile string, engine filtering.Engine, outFile string) {
	ids := new([][]byte)
	if err := binutil.Load(idFile, ids); err != nil {
		// if ids.csv exists, use gobencids command
		if _, err := os.Stat("ids.csv"); !os.IsNotExist(err) {
			log.Printf("%v not found, using ids.csv instead...", idFile)
			err = exec.Command("gobencid").Run()
			if err != nil {
				panic(err)
			}
			err = binutil.Load(idFile, ids)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	notifies := filtering.NotifyMap{}

	for _, id := range *ids {
		matches := engine.Search(id)
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

func getPackagePath() string {
	// Determine the package dir
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Dir(filename)
}

func loadFiltersFromCSVFile(f string) filtering.Subscriptions {
	sub := filtering.Subscriptions{}
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
		if len(record) < 3 {
			// Default case
			// prefix as key, *filtering.Info as value
			sub[record[1]] = &filtering.Info{Offset: 0, NotificationURI: record[0], EntropyValue: 0, Subset: nil}
		} else {
			// For OptimalBST, filter with EntropyValue
			// prefix as key, *filtering.Info as value
			pValue, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				panic(err)
			}
			fs := record[1]
			uri := record[0]
			sub[fs] = &filtering.Info{Offset: 0, NotificationURI: uri, EntropyValue: pValue, Subset: nil}
		}
	}
	return sub
}

func loadList(filterFile string, engineFile string, isRebuilding bool) *filtering.List {
	var list *filtering.List
	// List encode
	_, err := os.Stat(engineFile)
	if isRebuilding || os.IsNotExist(err) {
		sub := loadFiltersFromCSVFile(filterFile)
		log.Printf("Loaded %v filters from %s\n", len(sub), filterFile)
		var listBuf bytes.Buffer
		list = filtering.BuildList(sub)
		enc := gob.NewEncoder(&listBuf)
		err = enc.Encode(list)
		if err != nil {
			log.Fatal("encode:", err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal("file:", err)
		}
		file.Write(listBuf.Bytes())
		file.Close()
		log.Print("Saved the List filtering engine to ", engineFile)
	} else {
		// Tree decode
		err = binutil.Load(engineFile, &list)
		if err != nil {
			log.Fatal("engine: ", err)
		}
		log.Print("Loaded the List filtering engine from ", engineFile)
	}
	return list
}

func loadOptimalBST(filterFile string, engineFile string, isRebuilding bool) *filtering.OptimalBST {
	var head *filtering.OptimalBST
	// Tree encode
	_, err := os.Stat(engineFile)
	if isRebuilding || os.IsNotExist(err) {
		sub := loadFiltersFromCSVFile(filterFile)
		log.Printf("Loaded %v filters from %s\n", len(sub), filterFile)
		var tree bytes.Buffer
		head = filtering.BuildOptimalBST(&sub)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal("encode:", err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal("file:", err)
		}
		file.Write(tree.Bytes())
		file.Close()
		log.Print("Saved the Optimal BST filtering engine to ", engineFile)
	} else {
		// Tree decode
		binutil.Load(engineFile, &head)
		log.Print("Loaded the Optimal BST filtering engine from ", engineFile)
	}
	return head
}

func loadPatriciaTrie(filterFile string, engineFile string, isRebuilding bool) *filtering.PatriciaTrie {
	var head *filtering.PatriciaTrie
	// Tree encode
	_, err := os.Stat(engineFile)
	if isRebuilding || os.IsNotExist(err) {
		sub := loadFiltersFromCSVFile(filterFile)
		log.Printf("Loaded %v filters from %s\n", len(sub), filterFile)
		var tree bytes.Buffer
		head = filtering.BuildPatriciaTrie(sub)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal("encode:", err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal("file:", err)
		}
		file.Write(tree.Bytes())
		file.Close()
		log.Print("Saved the Patricia Trie filtering engine to ", engineFile)
	} else {
		// Tree decode
		err = binutil.Load(engineFile, &head)
		if err != nil {
			log.Fatal("engine: ", err)
		}
		log.Print("Loaded the Patricia Trie filtering engine from ", engineFile)
	}
	return head
}

func loadSplayTree(filterFile string, engineFile string, isRebuilding bool) *filtering.SplayTree {
	var head *filtering.SplayTree
	// Tree encode
	_, err := os.Stat(engineFile)
	if isRebuilding || os.IsNotExist(err) {
		sub := loadFiltersFromCSVFile(filterFile)
		log.Printf("Loaded %v filters from %s\n", len(sub), filterFile)
		var tree bytes.Buffer
		head = filtering.BuildSplayTree(&sub)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal("encode:", err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal("file:", err)
		}
		file.Write(tree.Bytes())
		file.Close()
		log.Print("Saved the Splay Tree filtering engine to ", engineFile)
	} else {
		// Tree decode
		err = binutil.Load(engineFile, &head)
		if err != nil {
			log.Fatal("engine: ", err)
		}
		log.Print("Loaded the Splay Tree filtering engine from ", engineFile)
	}
	return head
}

func run(filterFile string, engineFile string) {
	log.Println("Initializing gosstrak-fc middleware...")

	/* Build the filtering engine */
	engine := loadPatriciaTrie(filterFile, engineFile, true)

	/* Initialize the LLRP connection */
	// Establish a connection to the llrp client
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("connection:", err)
	}
	buf := make([]byte, BufferSize)
	for {
		// Read the incoming connection into the buffer.
		msgLength, err := conn.Read(buf)
		if err == io.EOF {
			// Close the connection when you're done with it.
			return
		} else if err != nil {
			log.Fatal(err.Error())
			conn.Close()
			break
		}

		header := binary.BigEndian.Uint16(buf[:2])
		if header == llrp.ReaderEventNotificationHeader {
			log.Println(">>> READER_EVENT_NOTIFICATION")
			conn.Write(llrp.SetReaderConfig(messageID))
		} else if header == llrp.KeepaliveHeader {
			log.Println(">>> KEEP_ALIVE")
			conn.Write(llrp.KeepaliveAck())
		} else if header == llrp.SetReaderConfigResponseHeader {
			log.Println(">>> SET_READER_CONFIG_RESPONSE")
		} else if header == llrp.ROAccessReportHeader {
			log.Println(">>> RO_ACCESS_REPORT")
			roarSize := uint16(binary.BigEndian.Uint32(buf[2:6])) // ROAR size
			//log.Println(roarSize)
			trds := buf[10:roarSize]                      // TRD stack
			trdSize := binary.BigEndian.Uint16(trds[2:4]) // First TRD size
			offset := uint16(0)
			for trdSize != 0 && int(offset) != len(trds) {
				log.Printf("trdSize: %v, len(trds): %v\n", trdSize, len(trds))
				var id []byte
				if trds[offset+4] == 141 { // EPC-96
					id = trds[offset+5 : offset+17]
					log.Printf("EPC: %v\n", id)
				} else if binary.BigEndian.Uint16(trds[offset+4:offset+6]) == 241 {
					epcDataSize := binary.BigEndian.Uint16(trds[offset+6 : offset+8])
					epcLengthBits := binary.BigEndian.Uint16(trds[offset+8 : offset+10])
					id = trds[offset+10 : offset+epcDataSize*2]
					id = id[0 : epcLengthBits/8]
					log.Printf("non-EPC: %v\n", id)
				}
				matches := engine.Search(id)
				//_ = engine.Search(id)

				offset += trdSize
				//log.Printf("offset: %v\n", offset)
				if offset != roarSize && int(trdSize) != len(trds) {
					trdSize = binary.BigEndian.Uint16(trds[offset+2 : offset+4])
				} else {
					trdSize = 0
				}
				log.Printf("%v matches", len(matches))
			}
		} else {
			log.Fatalf("Unknown header: %v\n", header)
			log.Fatalf("%v\n", buf[:msgLength])
		}
	}
}

func runDumb(idFile string, sub filtering.Subscriptions) {
	ids := new([][]byte)
	if err := binutil.Load(idFile, ids); err != nil {
		panic(err)
	}
	//fmt.Printf("Loaded %v ids from %v\n", len(*ids), idFile)
	matches := map[string][]string{}
	for _, id := range *ids {
		i := binutil.ParseByteSliceToBinString(id)
		for f, info := range sub {
			if strings.HasPrefix(i, f) {
				if _, ok := matches[info.NotificationURI]; !ok {
					matches[info.NotificationURI] = []string{}
				}
				matches[info.NotificationURI] = append(matches[info.NotificationURI], i)
			}
		}
	}

	//for n, m := range matches {
	//	fmt.Printf("%v: %v\n", n, len(m))
	//}
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Create cache directory if not exists
	if _, err := os.Stat(dataCacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataCacheDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	switch parse {
	case cmdAnalyze.FullCommand():
		switch strings.ToLower(*analyzeEngine) {
		case "obst":
			head := loadOptimalBST(*filterFile, *engineFile, false)
			aoFile := getPackagePath() + "/public/obst/locality.json"
			analyze(head, *analyzeInput, aoFile)
		case "patricia":
			head := loadPatriciaTrie(*filterFile, *engineFile, false)
			aoFile := getPackagePath() + "/public/patricia/locality.json"
			analyze(head, *analyzeInput, aoFile)
		}
	case cmdDumb.FullCommand():
		sub := loadFiltersFromCSVFile(*filterFile)
		log.Printf("Loaded %v filters from %s\n", len(sub), *filterFile)
		runDumb(*idFile, sub)
	case cmdList.FullCommand():
		list := loadList(*filterFile, *engineFile, *isRebuilding)
		execute(*idFile, list, *outFile)
	case cmdOBST.FullCommand():
		head := loadOptimalBST(*filterFile, *engineFile, *isRebuilding)
		execute(*idFile, head, *outFile)
	case cmdPatricia.FullCommand():
		head := loadPatriciaTrie(*filterFile, *engineFile, *isRebuilding)
		execute(*idFile, head, *outFile)
	case cmdSplay.FullCommand():
		head := loadSplayTree(*filterFile, *engineFile, *isRebuilding)
		execute(*idFile, head, *outFile)
	case cmdRun.FullCommand():
		run(*filterFile, *engineFile)
	}
}
