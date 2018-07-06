// Copyright (c) 2017 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/docker/libchan/spdy"
	"github.com/iomz/go-llrp"
	"github.com/iomz/gosstrak/adapter"
	"github.com/iomz/gosstrak/filtering"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Notification is the struct to send/receive captured ID
type Notification struct {
	ID []byte
}

// Constant Values
const (
	// BufferSize is a general size for a buffer
	BufferSize = 64 * 1024 // 64 KiB
	QueueSize  = 128
)

// Environmental variables
var (
	// Current Version
	version = "0.3.0"

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

	// LLRP related values
	initialMessageID = app.
				Flag("initialMessageID", "The initial messageID to start from.").
				Short('m').
				Default("1000").
				Int()
	ip = app.
		Flag("ip", "LLRP emulator address.").
		Short('a').
		Default("127.0.0.1").
		IP()
	port = app.
		Flag("port", "LLRP emulator port.").
		Short('p').
		Default("5084").
		Int()

	// start command
	cmdStart = app.Command("start", "Start the gosstrak-fc.")

	// Current messageID
	currentMessageID = uint32(*initialMessageID)
)

func getPackagePath() string {
	// Determine the package dir
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Dir(filename)
}

func run(f string) {
	log.Println("initializing gosstrak-fc for master mode...")

	// Load existing subscriptions from file
	log.Println("loading subscriptions from file")
	sub := filtering.LoadFiltersFromCSVFile(f)

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

	// Set up an EngineFactory with a management channel
	log.Println("setting up an engine factory")
	engineFactory := filtering.NewEngineFactory(sub, mc)
	go engineFactory.Run()
	// Wait until the first engine becomes available
	for currentEngine == nil {
		time.Sleep(time.Second)
	}

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
				break
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

	// Receive incoming IDs
	log.Println("setting up an incoming LLRPReadEvent channel")
	var rq = make(chan []*adapter.LLRPReadEvent)
	go func() {
		for {
			res, ok := <-rq
			if !ok {
				break
			}

			log.Printf("%v tags received", len(res))
			matches := map[string][]byte{}
			for _, re := range res {
				dests := currentEngine.Search(re.ID)
				for _, d := range dests {
					matches[d] = re.ID
				}
			}
			log.Printf("%v matches", len(matches))
			//notify matches
			/*
				for _, m := range matches {
					log.Printf("match: %s <- %v,%v\n", m, re.PC, re.ID)
				}
			*/
		}
		log.Fatalln("LLRPReadEvent listener exited in gosstrak-fc")
	}()

	// Establish a connection to the llrp client
	log.Println("starting a connection to an interrogator")
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(*port))
	for err != nil {
		log.Print(err)
		log.Println("wait 5 seconds for the interrogator to becom online...")
		time.Sleep(5 * time.Second)
		conn, err = net.Dial("tcp", ip.String()+":"+strconv.Itoa(*port))
	}

	// Prepare LLRP header storage
	header := make([]byte, 2)
	length := make([]byte, 4)
	messageID := make([]byte, 4)
	for {
		_, err = io.ReadFull(conn, header)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.ReadFull(conn, length)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.ReadFull(conn, messageID)
		if err != nil {
			log.Fatal(err)
		}
		// length containts the size of the entire message in octets
		// starting from bit offset 0, hence, the message size is
		// length - 10 bytes
		var messageValue []byte
		if messageSize := binary.BigEndian.Uint32(length) - 10; messageSize != 0 {
			messageValue = make([]byte, binary.BigEndian.Uint32(length)-10)
			_, err = io.ReadFull(conn, messageValue)
			if err != nil {
				log.Fatal(err)
			}
		}

		h := binary.BigEndian.Uint16(header)
		switch h {
		case llrp.ReaderEventNotificationHeader:
			log.Println(">>> READER_EVENT_NOTIFICATION")
			conn.Write(llrp.SetReaderConfig(currentMessageID))
		case llrp.KeepaliveHeader:
			log.Println(">>> KEEP_ALIVE")
			conn.Write(llrp.KeepaliveAck())
		case llrp.SetReaderConfigResponseHeader:
			log.Println(">>> SET_READER_CONFIG_RESPONSE")
		case llrp.ROAccessReportHeader:
			log.Println(">>> RO_ACCESS_REPORT")
			rq <- adapter.UnmarshalROAccessReportBody(messageValue)
		default:
			log.Fatalf("Unknown LLRP Message Header: %v\n", h)
		}
	}
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Create cache directory if not exists
	// TODO: set OS specific dataCacheDir
	dataCacheDir := "/var/tmp/gosstrak-fc-cache"
	if _, err := os.Stat(dataCacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataCacheDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	switch parse {
	case cmdStart.FullCommand():
		run(*filterFile)
	}
}
