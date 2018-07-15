// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/iomz/go-llrp"
	"github.com/iomz/go-llrp/binutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// kingpin app
	app     = kingpin.New("gobtags", "Read pcbits and ids from csv and save them in a gob file.")
	inFile  = app.Arg("in", "A source csv file contains tags.").Required().String()
	outFile = app.Flag("out", "A destination gob file contains tags.").Short('o').Default("out.gob").String()
)

func makeByteID(s string) ([]byte, error) {
	id, err := binutil.ParseBinRuneSliceToUint8Slice([]rune(s))
	if err != nil {
		return []byte{}, nil
	}
	return binutil.Pack([]interface{}{id}), err
}

func makeUint16PC(s string) (uint16, error) {
	pc64, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		return uint16(0), err
	}
	return uint16(pc64), err
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	tags := llrp.LoadTagsFromCSV(*inFile)
	binutil.Save(*outFile, &tags)
	/*
		var tagBuf bytes.Buffer
		enc := gob.NewEncoder(&tagBuf)
		err := enc.Encode(tags)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Write(tagBuf.Bytes())
		file.Close()
	*/
	log.Printf("Saved %v tags in %v\n", len(tags), *outFile)
}
