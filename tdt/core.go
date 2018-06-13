// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

import (
	"errors"
	//"io/ioutil"
	//"log"
	"math/big"
	//"os"
	"strconv"
	//"xml"
)

type core struct {
	//schemePrefixMap map[schemePrefix]string
	epcTDSVersion string
}

func NewCore() *core {
	c := new(core)
	//c.loadEPCTagDataTranslation()
	return c
}

func (c *core) LoadEPCTagDataTranslation() {
	//schemaDir := os.Getenv("GOPATH") + "/src/github.com/iomz/gosstrak/vendor/schemes/"
	//files, err := ioutil.ReadDir(schemaDir)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//for _, f := range files {
	//}
}

/*
type schemePrefix struct {
	prefixMatch     []byte
	prefixBitLength int
}

type schemePatternMap map[string][]schemePattern

type schemePattern struct {
	pattern string
	fields  []schemeField
}

type schemeFiled struct {
	decimalMinimum int
	decimalMaximum int
	bitLength      int
	name           string
}
*/

func (c *core) Translate(id []byte, pc []byte) (string, error) {
	if len(pc) != 2 {
		return "", errors.New("Invalid PC bits")
	}

	// Check the NSI toggle
	// 00000001 & pc[0]
	switch 1 & pc[0] {
	case 0: // GS1
		return c.buildEPC(id)
	case 1: // ISO
		return c.buildUII(id, pc[1])
	}
	// Proprietary
	return c.buildProprietary(id)
}

func (c *core) buildEPC(id []byte) (string, error) {
	urn := ""

	// EPC Header
	switch id[0] {
	case 48: // SGTIN-96 00110000
		if len(id) != 12 {
			return "", errors.New("Invalid ID")
		}
		urn = "urn:epc:id:sgtin:"
		// FILTER
		urn += strconv.Itoa(int((id[1]&224)>>5)) + "." // 224: 11100000
		// PARTITION
		partition := int((id[1] & 28) >> 2) // 28: 00011100
		ptm := map[PartitionTableKey]int{}
		for _, v := range SGTIN96PartitionTable {
			if v[PValue] == partition {
				ptm = v
				break
			}
		}
		// COMPANY_PREFIX and ITEM_REFERENCE
		z := new(big.Int)
		switch ptm[CPBits] {
		case 40:
			cp := make([]byte, 5)
			remainder := id[1] & 3
			cp[0] = remainder<<6 | id[2]>>2
			remainder = id[2] & 3
			cp[1] = remainder<<6 | id[3]>>2
			remainder = id[3] & 3
			cp[2] = remainder<<6 | id[4]>>2
			remainder = id[4] & 3
			cp[3] = remainder<<6 | id[5]>>2
			remainder = id[5] & 3
			cp[4] = remainder<<6 | id[6]>>2
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 1) // 4 bits
			remainder = id[6] & 3
			ir[0] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		case 37:
			cp := make([]byte, 5)
			remainder := id[1] & 3
			cp[0] = remainder<<3 | id[2]>>5
			remainder = id[2] & 31
			cp[1] = remainder<<3 | id[3]>>5
			remainder = id[3] & 31
			cp[2] = remainder<<3 | id[4]>>5
			remainder = id[4] & 31
			cp[3] = remainder<<3 | id[5]>>5
			remainder = id[5] & 31
			cp[4] = remainder<<3 | id[6]>>5
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 1) // 7 bits
			remainder = id[6] & 31
			ir[0] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		case 34:
			cp := make([]byte, 5)
			cp[0] = id[1] & 3
			cp[1] = id[2]
			cp[2] = id[3]
			cp[3] = id[4]
			cp[4] = id[5]
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 2) // 10 bits
			ir[0] = id[6] >> 2
			remainder := id[6] & 3
			ir[1] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		case 30:
			cp := make([]byte, 4)
			remainder := id[1] & 3
			cp[0] = remainder<<4 | id[2]>>4
			remainder = id[2] & 15
			cp[1] = remainder<<4 | id[3]>>4
			remainder = id[3] & 15
			cp[2] = remainder<<4 | id[4]>>4
			remainder = id[4] & 15
			cp[3] = remainder<<4 | id[5]>>4
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 2) // 14 bits
			remainder = id[5] & 15
			ir[0] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[1] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		case 27:
			cp := make([]byte, 4)
			remainder := id[1] & 3
			cp[0] = remainder<<1 | id[2]>>7
			remainder = id[2] & 127
			cp[1] = remainder<<1 | id[3]>>7
			remainder = id[3] & 127
			cp[2] = remainder<<1 | id[4]>>7
			remainder = id[4] & 127
			cp[3] = remainder<<1 | id[5]>>7
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 3) // 17 bits
			remainder = id[5] & 127
			ir[0] = remainder >> 6
			remainder = remainder & 63
			ir[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[2] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		case 24:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<6 | id[2]>>2
			remainder = id[2] & 3
			cp[1] = remainder<<6 | id[3]>>2
			remainder = id[3] & 3
			cp[2] = remainder<<6 | id[4]>>2
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 3) // 20 bits
			remainder = id[4] & 3
			ir[0] = remainder<<2 | id[5]>>6
			remainder = id[5] & 63
			ir[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[2] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		case 20:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<2 | id[2]>>6
			remainder = id[2] & 63
			cp[1] = remainder<<2 | id[3]>>6
			remainder = id[3] & 63
			cp[2] = remainder<<2 | id[4]>>6
			z.SetBytes(cp)
			urn += z.String() + "."
			ir := make([]byte, 3) // 24 bits
			remainder = id[4] & 63
			ir[0] = remainder<<2 | id[5]>>6
			remainder = id[5] & 63
			ir[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[2] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			urn += z.String() + "."
		}
		// SERIAL
		ser := make([]byte, 5)
		ser[0] = id[7] & 63
		ser[1] = id[8]
		ser[2] = id[9]
		ser[3] = id[10]
		ser[4] = id[11]
		z.SetBytes(ser)
		urn += z.String()
	case 49: // SSCC-96  00110001
		urn = "urn:epc:id:sscc:"
	case 51: // GRAI-96  00110011
		urn = "urn:epc:id:grai:"
	case 52: // GIAI-96  00110100
		urn = "urn:epc:id:giai:"
	}
	return urn, nil
}

func (c *core) buildUII(id []byte, afi byte) (string, error) {
	return "", nil
}

func (c *core) buildProprietary(id []byte) (string, error) {
	return "", nil
}

func parse6BitEncodedByteSliceToString(in []byte) (string, error) {
	bitLength := len(in) * 8
	var buf []byte
	for offset := 0; offset+6 <= bitLength; offset += 6 {
		var c byte
		switch offset % 8 {
		/*
			i = 0 1  2  3  4
			o = 0 6 12 18 24
			x = 0 6  4  2  0
			y = 0 0  1  2  3
			where x = offset%8 and y = offset/8
		*/
		case 0:
			// (11111100 & in[y] ) >> 2
			c = (252 & in[offset/8]) >> 2
		case 2:
			// (00111111 & in[y])
			c = 63 & in[offset/8]
		case 4:
			// ((00001111 & in[y]) << 2 ) | ((11000000 & in[y+1]) >> 6)
			y := offset / 8
			c = ((15 & in[y]) << 2) | ((192 & in[y+1]) >> 6)
		case 6:
			// ((00000011 & in[y]) << 4 ) | ((11110000 & in[y+1]) >> 4)
			y := offset / 8
			c = ((3 & in[y]) << 4) | ((240 & in[y+1]) >> 4)
		}
		// (00100000 & c) != 00100000
		if (32 & c) != 32 {
			// c = 01000000 | c
			c |= 64
		}
		// if c is NOT SPACE(100000), append the c
		if c^32 != 0 {
			buf = append(buf, c)
		}
	}
	return string(buf), nil
}
