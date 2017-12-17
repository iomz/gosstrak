package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/iomz/go-llrp/binutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// kingpin app
	app     = kingpin.New("gobencid", "Read ids from csv and save them in a gob file.")
	inFile  = app.Flag("in", "A source csv file contains ids.").Short('i').Default("ids.csv").String()
	outFile = app.Flag("out", "A destination gob file contains ids.").Short('o').Default("ids.gob").String()
)

func makeByteID(s string) []byte {
	id, _ := binutil.ParseBinRuneSliceToUint8Slice([]rune(s))
	return binutil.Pack([]interface{}{id})
}

func readIDsFromCSV(inputFile string) *[][]byte {
	// Check inputFile
	fp, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// Read CSV and store in [][]byte
	ids := [][]byte{}
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
		if len(record) == 2 {
			id := makeByteID(record[1])
			ids = append(ids, id)
		}
	}
	return &ids
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	ids := readIDsFromCSV(*inFile)
	binutil.Save(*outFile, ids)
	log.Printf("Saved %v ids in %v\n", len(*ids), *outFile)
}
