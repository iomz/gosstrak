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
	"time"

	"github.com/docker/libchan/spdy"
	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak/filtering"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Notification is the struct to send/receive captured ID
type Notification struct {
	ID []byte
}

// Event is the struct to hold data on RFTags
type ReadEvent struct {
	id []byte
	pc []byte
}

// Constant Values
const (
	// BufferSize is a general size for a buffer
	BufferSize = 64 * 1024 // 64 KiB
	QueueSize  = 128
	SGTINHOST  = "192.168.22.1:9323"
	//SGTINHOST = "localhost:9323"
	SSCCHOST = "192.168.22.4:9323"
	//SSCCHOST = "localhost:9324"
)

// Environmental variables
var (
	// Current Version
	version = "0.3.0"

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
		Flag("ip", "LLRP emulator address.").
		Short('a').
		Default("0.0.0.0").
		IP()
	port = app.
		Flag("port", "LLRP emulator port.").
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
	cmdDumb = app.Command("dumb", "Run in dumb filter mode.")

	// list command
	cmdList = app.Command("list", "Use list of byte filtering engine.")

	// obst command
	cmdOBST = app.Command("obst", "Use Optimal Binary Search Tree filtering engine.")

	// patricia command
	cmdPatricia = app.Command("patricia", "Use Patricia Trie filtering engine.")

	// obst command
	cmdSplay = app.Command("splay", "Use Splay Tree filtering engine.")

	// run command
	cmdRun = app.Command("run", "Run the F&C middleware.")

	// start command
	cmdStart = app.Command("start", "Start the gosstrak-fc.")

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
		log.Fatal(err)
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
			sub[record[1]] = &filtering.Info{
				Offset:          0,
				NotificationURI: record[0],
				EntropyValue:    0,
				Subset:          filtering.Subscriptions{},
			}
		} else {
			// For OptimalBST, filter with EntropyValue
			// prefix as key, *filtering.Info as value
			pValue, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				panic(err)
			}
			fs := record[1]
			uri := record[0]
			sub[fs] = &filtering.Info{
				Offset:          0,
				NotificationURI: uri,
				EntropyValue:    pValue,
				Subset:          filtering.Subscriptions{},
			}
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
		list = filtering.NewList(sub).(*filtering.List)
		enc := gob.NewEncoder(&listBuf)
		err = enc.Encode(list)
		if err != nil {
			log.Fatal(err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Write(listBuf.Bytes())
		file.Close()
		log.Print("Saved the List filtering engine to ", engineFile)
	} else {
		// Tree decode
		err = binutil.Load(engineFile, &list)
		if err != nil {
			log.Fatal(err)
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
		head = filtering.NewOptimalBST(sub).(*filtering.OptimalBST)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal(err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal(err)
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
		head = filtering.NewPatriciaTrie(sub).(*filtering.PatriciaTrie)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal(err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Write(tree.Bytes())
		file.Close()
		log.Print("Saved the Patricia Trie filtering engine to ", engineFile)
	} else {
		// Tree decode
		err = binutil.Load(engineFile, &head)
		if err != nil {
			log.Fatal(err)
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
		head = filtering.NewSplayTree(sub).(*filtering.SplayTree)
		enc := gob.NewEncoder(&tree)
		err = enc.Encode(head)
		if err != nil {
			log.Fatal(err)
		}
		// Save to file
		file, err := os.Create(engineFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Write(tree.Bytes())
		file.Close()
		log.Print("Saved the Splay Tree filtering engine to ", engineFile)
	} else {
		// Tree decode
		err = binutil.Load(engineFile, &head)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Loaded the Splay Tree filtering engine from ", engineFile)
	}
	return head
}

func run(filterFile string, engineFile string) {
	log.Println("Initializing gosstrak-fc middleware...")

	/* Build the filtering engine */
	engine := loadPatriciaTrie(filterFile, engineFile, true)

	/* Initialize notify client */
	// SGTIN
	sgtinConn, err := net.Dial("tcp", SGTINHOST)
	if err != nil {
		log.Fatal(err)
	}
	sgtinp, err := spdy.NewSpdyStreamProvider(sgtinConn, false)
	if err != nil {
		log.Fatal(err)
	}
	sgtintransport := spdy.NewTransport(sgtinp)
	sgtinNotify, err := sgtintransport.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}
	// SSCC
	ssccConn, err := net.Dial("tcp", SSCCHOST)
	if err != nil {
		log.Fatal(err)
	}
	ssccp, err := spdy.NewSpdyStreamProvider(ssccConn, false)
	if err != nil {
		log.Fatal(err)
	}
	sscctransport := spdy.NewTransport(ssccp)
	ssccNotify, err := sscctransport.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}

llrpinit:
	/* Initialize the LLRP connection */
	// Establish a connection to the llrp client
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, BufferSize)
	for {
		// Read the incoming connection into the buffer.
		msgLength, err := conn.Read(buf)
		if err == io.EOF {
			// Close the connection when you're done with it.
			return
		} else if err != nil {
			log.Fatal(err)
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
			log.Println("<<< KEEP_ALIVE_ACK")
		} else if header == llrp.SetReaderConfigResponseHeader {
			log.Println(">>> SET_READER_CONFIG_RESPONSE")
		} else if header == llrp.ROAccessReportHeader {
			log.Println(">>> RO_ACCESS_REPORT")
			notifies := filtering.NotifyMap{}
			roarSize := uint16(binary.BigEndian.Uint32(buf[2:6])) // ROAR size
			//log.Println(roarSize)
			trds := buf[10:roarSize]                      // TRD stack
			trdSize := binary.BigEndian.Uint16(trds[2:4]) // First TRD size
			offset := uint16(0)
			go func() {
				for trdSize != 0 && int(offset) != len(trds) {
					//log.Printf("trdSize: %v, len(trds): %v\n", trdSize, len(trds))
					var id []byte
					if trds[offset+4] == 141 { // EPC-96
						id = trds[offset+5 : offset+17]
						//log.Printf("EPC: %v\n", id)
					} else if binary.BigEndian.Uint16(trds[offset+4:offset+6]) == 241 {
						epcDataSize := binary.BigEndian.Uint16(trds[offset+6 : offset+8])
						epcLengthBits := binary.BigEndian.Uint16(trds[offset+8 : offset+10])
						id = trds[offset+10 : offset+epcDataSize*2]
						id = id[0 : epcLengthBits/8]
						//log.Printf("non-EPC: %v\n", id)
					}
					matches := engine.Search(id)
					for _, n := range matches {
						if _, ok := notifies[n]; !ok {
							notifies[n] = [][]byte{}
						}
						notifies[n] = append(notifies[n], id)
					}

					offset += trdSize
					//log.Printf("offset: %v\n", offset)
					if offset != roarSize && int(trdSize) != len(trds) {
						trdSize = binary.BigEndian.Uint16(trds[offset+2 : offset+4])
					} else {
						trdSize = 0
					}
				}
				for filter, ids := range notifies {
					switch filter {
					case "SGTIN-96":
						for _, id := range ids {
							err = sgtinNotify.Send(&Notification{id})
							if err != nil {
								log.Fatal(err)
							}
						}
					case "SSCC-96":
						for _, id := range ids {
							err = ssccNotify.Send(&Notification{id})
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				}
			}()
		} else {
			log.Printf("Unknown header: %v\n", header)
			log.Printf("%v\n", buf[:msgLength])
			goto llrpinit
		}
	}
}

func runMaster(f string) {
	log.Println("initializing gosstrak-fc for master mode...")

	// Load existing subscriptions from file
	log.Println("loading subscriptions from file")
	sub := loadFiltersFromCSVFile(f)
	log.Println("...complete")

	// Receive the engine instance
	log.Println("setting up a management channel")
	var currentEngine filtering.Engine
	mc := make(chan filtering.ManagementMessage, QueueSize)
	go func() {
		for {
			msg, ok := <-mc
			if !ok {
				break
			}
			switch msg.Type {
			case filtering.DeployEngine:
				currentEngine = msg.EngineGeneratorInstance.Engine
				log.Printf("currentEngine switched to %s\n", msg.EngineGeneratorInstance.Name)
			}
		}
		log.Fatalln("managementChannel listener exited in gosstrak-fc")
	}()
	log.Println("...complete")

	// Set up an EngineFactory with a management channel
	log.Println("setting up an engine factory")
	engineFactory := filtering.NewEngineFactory(sub, mc)
	go engineFactory.Run()
	// Wait until the first engine becomes available
	for currentEngine == nil {
		time.Sleep(time.Second)
	}
	log.Println("...complete")

	// Receive management access
	log.Println("setting up an management interface")
	go func() {
		managementListener, err := net.Listen("tcp", "0.0.0.0:2784")
		if err != nil {
			log.Fatal(err)
		}
		for {
			c, err := managementListener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			p, err := spdy.NewSpdyStreamProvider(c, true)
			if err != nil {
				log.Print(err)
				continue
			}
			t := spdy.NewTransport(p)

			receiver, err := t.WaitReceiveChannel()
			if err != nil {
				log.Print(err)
				continue
			}

			mm := &filtering.ManagementMessage{}
			err = receiver.Receive(mm)
			if err != nil {
				log.Print(err)
				continue
			}
			log.Print(mm)
			mc <- *mm
		}
		log.Fatalln("managementListener closed in gosstrak-fc")
	}()
	log.Println("...complete")

	// Receive incoming IDs
	log.Println("setting up an incoming ReadEvent channel")
	var rq = make(chan ReadEvent, QueueSize)
	go func() {
		for {
			re, ok := <-rq
			if !ok {
				break
			}
			matches := currentEngine.Search(re.id)
			//notify matches
			for _, m := range matches {
				log.Printf("match: %s <- %v,%v\n", m, re.pc, re.id)
			}
		}
		log.Fatalln("ic listener exited in gosstrak-fc")
	}()
	log.Println("...complete")

	// Establish a connection to the llrp client
	log.Println("starting a connection to an interrogator")
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}

	// Prepare LLRP header storage
	header := make([]byte, 2)
	length := make([]byte, 4)
	for {
		_, err = io.ReadFull(conn, header)
		h := binary.BigEndian.Uint16(header)
		if h == llrp.ReaderEventNotificationHeader {
			_, err = io.ReadFull(conn, length)
			message := make([]byte, binary.BigEndian.Uint32(length)-6)
			_, err = io.ReadFull(conn, message)
			log.Println(">>> READER_EVENT_NOTIFICATION")
			conn.Write(llrp.SetReaderConfig(messageID))
		} else if h == llrp.KeepaliveHeader {
			_, err = io.ReadFull(conn, length)
			message := make([]byte, binary.BigEndian.Uint32(length)-6)
			_, err = io.ReadFull(conn, message)
			log.Println(">>> KEEP_ALIVE")
			conn.Write(llrp.KeepaliveAck())
		} else if h == llrp.SetReaderConfigResponseHeader {
			_, err = io.ReadFull(conn, length)
			message := make([]byte, binary.BigEndian.Uint32(length)-6)
			_, err = io.ReadFull(conn, message)
			log.Println(">>> SET_READER_CONFIG_RESPONSE")
		} else if h == llrp.ROAccessReportHeader {
			_, err = io.ReadFull(conn, length)
			l := binary.BigEndian.Uint32(length)
			message := make([]byte, l-6)
			_, err = io.ReadFull(conn, message)
			log.Println(">>> RO_ACCESS_REPORT")
			decapsulateROAccessReport(l, message, rq)
		} else {
			log.Fatalf("Unknown header: %v\n", h)
		}
	}
}

// decapsulate the ROAccessReport and extract IDs
func decapsulateROAccessReport(roarLength uint32, buf []byte, rq chan ReadEvent) {
	//defer timeTrack(time.Now(), fmt.Sprintf("unpacking %v bytes", len(buf)))
	trds := buf[4 : roarLength-6] // TRD stack
	trdLength := uint16(0)        // First TRD size
	offset := uint32(0)           // the start of TRD
	//for trdLength != 0 && int(offset) != len(trds) {
	for {
		if uint32(10+offset) < roarLength {
			trdLength = binary.BigEndian.Uint16(trds[offset+2 : offset+4])
		} else {
			break
		}
		var id, pc []byte
		if trds[offset+4] == 141 { // EPC-96
			id = trds[offset+5 : offset+17]
			if trds[offset+17] == 140 { // C1G2-PC parameter
				pc = trds[offset+18 : offset+20]
			}
			//log.Printf("EPC: %v, (%x)\n", id, pc)
		} else if binary.BigEndian.Uint16(trds[offset+4:offset+6]) == 241 { // EPCData
			epcDataLength := binary.BigEndian.Uint16(trds[offset+6 : offset+8])  // length
			epcLengthBits := binary.BigEndian.Uint16(trds[offset+8 : offset+10]) // EPCLengthBits
			epcLengthBytes := uint32(epcLengthBits / 8)
			/*
				// ID length in byte = Length - (6 + 10 + 16 + 16)/8
				//id = trds[offset+6 : offset+epcDataSize-6]
				// trim the last 1 byte if it's not a multiple of a word
				//id = id[0 : epcLengthBits/8]
			*/
			id = trds[offset+10 : offset+10+epcLengthBytes]
			if 4+epcDataLength < trdLength && trds[offset+10+epcLengthBytes] == 140 { // C1G2-PC parameter
				pc = trds[offset+10+epcLengthBytes+1 : offset+10+epcLengthBytes+3]
			}
			//log.Printf("EPC: %v, (%x)\n", id, pc)
		}
		rq <- ReadEvent{id, pc}
		offset += uint32(trdLength) // move the offset at the end of this TRD
		//log.Printf("offset: %v, roarLength: %v\n", offset, roarLength)
		//log.Printf("trdLength: %v, len(trds): %v\n", trdLength, len(trds))
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
		for _, fs := range sub.Keys() {
			if strings.HasPrefix(i, fs) {
				info := sub[fs]
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
	case cmdStart.FullCommand():
		runMaster(*filterFile)
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
