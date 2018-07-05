// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/iomz/go-llrp/binutil"
	"github.com/iomz/gosstrak/filtering"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// kingpin app
	app     = kingpin.New("gobsubs", "Read filters from csv and save them in a gob file.")
	inFile  = app.Flag("in", "A source csv file contains filters.").Short('i').Default("filters.csv").String()
	outFile = app.Flag("out", "A destination gob file contains filters.").Short('o').Default("filters.gob").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	sub := filtering.LoadFiltersFromCSVFile(*inFile)
	binutil.Save(*outFile, &sub)
	log.Printf("Saved %v filters in %v\n", len(sub.Keys()), *outFile)
}
