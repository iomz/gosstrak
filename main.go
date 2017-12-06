package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
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

type (
	FilterMap map[string]string
)

type PatriciaTrie struct {
	prefix string
	one    *PatriciaTrie
	zero   *PatriciaTrie
	notify string
}

func (fm FilterMap) keys() []string {
	ks := []string{}
	for k, _ := range fm {
		ks = append(ks, k)
	}
	return ks
}

func (pt *PatriciaTrie) constructTrie(prefix string, fm FilterMap) {
	onePrefixBranch := ""
	zeroPrefixBranch := ""
	fks := fm.keys()
	for i := 0; i < len(fks); i++ {
		if len(fks[i]) < len(prefix) {
			continue
		}
		if !strings.HasPrefix(fks[i], prefix) {
			//fmt.Printf("x%s\n", fks[i])
			continue
		}
		p := fks[i][len(prefix):]
		if len(p) == 0 {
			continue
		}
		if strings.HasPrefix(p, "1") {
			if len(onePrefixBranch) == 0 {
				onePrefixBranch = p
			} else {
				onePrefixBranch = lcp([]string{p, onePrefixBranch})
			}
		} else if strings.HasPrefix(p, "0") {
			if len(zeroPrefixBranch) == 0 {
				zeroPrefixBranch = p
			} else {
				zeroPrefixBranch = lcp([]string{p, zeroPrefixBranch})
			}
		}
	}
	//fmt.Printf("opb: %s(%d), zpb: %s(%d)\n", onePrefixBranch, len(onePrefixBranch), zeroPrefixBranch, len(zeroPrefixBranch))
	cumulativePrefix := ""
	if len(onePrefixBranch) != 0 {
		pt.one = &PatriciaTrie{}
		pt.one.prefix = onePrefixBranch
		cumulativePrefix = prefix + onePrefixBranch
		// check if the prefix matches whole filter
		if n, ok := fm[cumulativePrefix]; ok {
			pt.one.notify = n
		}
		pt.one.constructTrie(cumulativePrefix, fm)
	}
	if len(zeroPrefixBranch) != 0 {
		pt.zero = &PatriciaTrie{}
		pt.zero.prefix = zeroPrefixBranch
		cumulativePrefix = prefix + zeroPrefixBranch
		// check if the prefix matches whole filter
		if n, ok := fm[cumulativePrefix]; ok {
			pt.zero.notify = n
		}
		pt.zero.constructTrie(cumulativePrefix, fm)
	}
}

func (pt *PatriciaTrie) dump() string {
	writer := &bytes.Buffer{}
	pt.print(writer, 0)
	return writer.String()
}

func (pt *PatriciaTrie) print(writer io.Writer, indent int) {
	var n string
	if len(pt.notify) != 0 {
		n = "-> " + pt.notify
	}
	fmt.Fprintf(writer, "%s--%s %s\n", strings.Repeat(" ", indent), pt.prefix, n)
	if pt.one != nil {
		pt.one.print(writer, indent+2)
	}
	if pt.zero != nil {
		pt.zero.print(writer, indent+2)
	}
}

// lcp finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
func lcp(l []string) string {
	// Special cases first
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

func runPatricia(f string) {
	fm := loadFiltersFromCSVFile(f)
	p1 := lcp(fm.keys())
	if len(p1) == 0 {
		// do something if there's no common prefix
	}
	head := &PatriciaTrie{}
	head.prefix = p1
	head.constructTrie(p1, fm)
	fmt.Println(head.dump())
	return
}

func loadFiltersFromCSVFile(f string) FilterMap {
	filters := FilterMap{}
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
		filters[record[0]] = record[1]
	}
	return filters
}

func main() {
	app.Version(version)
	parse := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch parse {
	case patricia.FullCommand():
		runPatricia(*filterFile)
	}
}
