// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

// Package tdt contains Tag Data Translation module from binary to Pure Identity
package tdt

import (
	"errors"
	"fmt"
	//"io/ioutil"
	//"log"
	"math/big"
	"strconv"
	"strings"
	//"xml"
)

// Core is the TDT core
type Core struct {
	//schemePrefixMap map[schemePrefix]string
	epcTDSVersion string
}

// NewCore returns a new instance of TDT core
func NewCore() *Core {
	c := new(Core)
	//c.loadEPCTagDataTranslation()
	return c
}

// LoadEPCTagDataTranslation loads EPC scheme from scheme files
func (c *Core) LoadEPCTagDataTranslation() {
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

// Translate takes ID in binary ([]byte) and returns the corresponding PureIdentity
func (c *Core) Translate(pc []byte, id []byte) (string, error) {
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

func (c *Core) buildEPC(id []byte) (string, error) {
	urn := ""

	// EPC Header
	switch id[0] {
	case 48: /* ------------- SGTIN-96 00110000 ------------- */
		if len(id) != 12 {
			return "", errors.New("Invalid ID")
		}
		urn = "urn:epc:id:sgtin-96:"
		// FILTER
		urn += strconv.Itoa(int((id[1]&224)>>5)) + "." // 224: 11100000
		// PARTITION
		partition := int((id[1] & 28) >> 2) // 28: 00011100
		ptm := map[PartitionTableKey]int{}
		var cpLength int
		for k, v := range SGTIN96PartitionTable {
			if v[PValue] == partition {
				ptm = v
				cpLength = k
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 1) // 4 bits
			remainder = id[6] & 3
			ir[0] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 1) // 7 bits
			remainder = id[6] & 31
			ir[0] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
		case 34:
			cp := make([]byte, 5)
			cp[0] = id[1] & 3
			cp[1] = id[2]
			cp[2] = id[3]
			cp[3] = id[4]
			cp[4] = id[5]
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 2) // 10 bits
			ir[0] = id[6] >> 6
			remainder := id[6] & 63
			ir[1] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 2) // 14 bits
			remainder = id[5] & 15
			ir[0] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[1] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 3) // 17 bits
			remainder = id[5] & 127
			ir[0] = remainder >> 6
			remainder = remainder & 63
			ir[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[2] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
		case 24:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<6 | id[2]>>2
			remainder = id[2] & 3
			cp[1] = remainder<<6 | id[3]>>2
			remainder = id[3] & 3
			cp[2] = remainder<<6 | id[4]>>2
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 3) // 20 bits
			remainder = id[4] & 3
			ir[0] = remainder<<2 | id[5]>>6
			remainder = id[5] & 63
			ir[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[2] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
		case 20:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<2 | id[2]>>6
			remainder = id[2] & 63
			cp[1] = remainder<<2 | id[3]>>6
			remainder = id[3] & 63
			cp[2] = remainder<<2 | id[4]>>6
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ir := make([]byte, 3) // 24 bits
			remainder = id[4] & 63
			ir[0] = remainder<<2 | id[5]>>6
			remainder = id[5] & 63
			ir[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			ir[2] = remainder<<2 | id[7]>>6
			z.SetBytes(ir)
			itemReference := z.String()
			urn += strings.Repeat("0", ptm[IRDigits]-len(itemReference)) + itemReference + "."
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
	case 49: /* ------------- SSCC-96  00110001 ------------- */
		if len(id) != 12 {
			return "", errors.New("Invalid ID")
		}
		urn = "urn:epc:id:sscc-96:"
		// FILTER
		urn += strconv.Itoa(int((id[1]&224)>>5)) + "." // 224: 11100000
		// PARTITION
		partition := int((id[1] & 28) >> 2) // 28: 00011100
		ptm := map[PartitionTableKey]int{}
		var cpLength int
		for k, v := range SSCC96PartitionTable {
			if v[PValue] == partition {
				ptm = v
				cpLength = k
				break
			}
		}
		// COMPANY_PREFIX and Extension
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 3) // 18 bits
			remainder = id[6] & 3
			ext[0] = remainder
			ext[1] = id[7]
			ext[2] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 3) // 21 bits
			remainder = id[6] & 31
			ext[0] = remainder
			ext[1] = id[7]
			ext[2] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
		case 34:
			cp := make([]byte, 5)
			cp[0] = id[1] & 3
			cp[1] = id[2]
			cp[2] = id[3]
			cp[3] = id[4]
			cp[4] = id[5]
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 3) // 24 bits
			ext[0] = id[6]
			ext[1] = id[7]
			ext[2] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 4) // 28 bits
			remainder = id[5] & 15
			ext[0] = remainder
			ext[1] = id[6]
			ext[2] = id[7]
			ext[3] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 4) // 31 bits
			remainder = id[5] & 127
			ext[0] = remainder
			ext[1] = id[6]
			ext[2] = id[7]
			ext[3] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
		case 24:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<6 | id[2]>>2
			remainder = id[2] & 3
			cp[1] = remainder<<6 | id[3]>>2
			remainder = id[3] & 3
			cp[2] = remainder<<6 | id[4]>>2
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 5) // 34 bits
			remainder = id[4] & 3
			ext[0] = remainder
			ext[1] = id[5]
			ext[2] = id[6]
			ext[3] = id[7]
			ext[4] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
		case 20:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<2 | id[2]>>6
			remainder = id[2] & 63
			cp[1] = remainder<<2 | id[3]>>6
			remainder = id[3] & 63
			cp[2] = remainder<<2 | id[4]>>6
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			ext := make([]byte, 5) // 38 bits
			remainder = id[4] & 63
			ext[0] = remainder
			ext[1] = id[5]
			ext[2] = id[6]
			ext[3] = id[7]
			ext[4] = id[8]
			z.SetBytes(ext)
			extension := z.String()
			urn += strings.Repeat("0", ptm[EDigits]-len(extension)) + extension
		}
	case 51: /* ------------- GRAI-96  00110011------------- */
		if len(id) != 12 {
			return "", errors.New("Invalid ID")
		}
		urn = "urn:epc:id:grai-96:"
		// FILTER
		urn += strconv.Itoa(int((id[1]&224)>>5)) + "." // 224: 11100000
		// PARTITION
		partition := int((id[1] & 28) >> 2) // 28: 00011100
		ptm := map[PartitionTableKey]int{}
		var cpLength int
		for k, v := range GRAI96PartitionTable {
			if v[PValue] == partition {
				ptm = v
				cpLength = k
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 1) // 4 bits
			remainder = id[6] & 3
			at[0] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 1) // 7 bits
			remainder = id[6] & 31
			at[0] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
		case 34:
			cp := make([]byte, 5)
			cp[0] = id[1] & 3
			cp[1] = id[2]
			cp[2] = id[3]
			cp[3] = id[4]
			cp[4] = id[5]
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 2) // 10 bits
			at[0] = id[6] >> 6
			remainder := id[6] & 63
			at[1] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 2) // 14 bits
			remainder = id[5] & 15
			at[0] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			at[1] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 3) // 17 bits
			remainder = id[5] & 127
			at[0] = remainder >> 6
			remainder = remainder & 63
			at[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			at[2] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
		case 24:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<6 | id[2]>>2
			remainder = id[2] & 3
			cp[1] = remainder<<6 | id[3]>>2
			remainder = id[3] & 3
			cp[2] = remainder<<6 | id[4]>>2
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 3) // 20 bits
			remainder = id[4] & 3
			at[0] = remainder<<2 | id[5]>>6
			remainder = id[5] & 63
			at[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			at[2] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
		case 20:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<2 | id[2]>>6
			remainder = id[2] & 63
			cp[1] = remainder<<2 | id[3]>>6
			remainder = id[3] & 63
			cp[2] = remainder<<2 | id[4]>>6
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			at := make([]byte, 3) // 24 bits
			remainder = id[4] & 63
			at[0] = remainder<<2 | id[5]>>6
			remainder = id[5] & 63
			at[1] = remainder<<2 | id[6]>>6
			remainder = id[6] & 63
			at[2] = remainder<<2 | id[7]>>6
			z.SetBytes(at)
			assetType := z.String()
			urn += strings.Repeat("0", ptm[ATDigits]-len(assetType)) + assetType + "."
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
	case 52: /* ------------- GIAI-96  00110100 ------------- */
		if len(id) != 12 {
			return "", errors.New("Invalid ID")
		}
		urn = "urn:epc:id:giai-96:"
		// FILTER
		urn += strconv.Itoa(int((id[1]&224)>>5)) + "." // 224: 11100000
		// PARTITION
		partition := int((id[1] & 28) >> 2) // 28: 00011100
		ptm := map[PartitionTableKey]int{}
		var cpLength int
		for k, v := range GIAI96PartitionTable {
			if v[PValue] == partition {
				ptm = v
				cpLength = k
				break
			}
		}
		// COMPANY_PREFIX and Individual Asset Reference
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 6) // 42 bits
			remainder = id[6] & 3
			iar[0] = remainder
			iar[1] = id[7]
			iar[2] = id[8]
			iar[3] = id[9]
			iar[4] = id[10]
			iar[5] = id[11]
			z.SetBytes(iar)
			urn += z.String()
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 6) // 45 bits
			remainder = id[6] & 31
			iar[0] = remainder
			iar[1] = id[7]
			iar[2] = id[8]
			iar[3] = id[9]
			iar[4] = id[10]
			iar[5] = id[11]
			z.SetBytes(iar)
			urn += z.String()
		case 34:
			cp := make([]byte, 5)
			cp[0] = id[1] & 3
			cp[1] = id[2]
			cp[2] = id[3]
			cp[3] = id[4]
			cp[4] = id[5]
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 6) // 48 bits
			iar[0] = id[6]
			iar[1] = id[7]
			iar[2] = id[8]
			iar[3] = id[9]
			iar[4] = id[10]
			iar[5] = id[11]
			z.SetBytes(iar)
			urn += z.String()
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 7) // 52 bits
			remainder = id[5] & 15
			iar[0] = remainder
			iar[1] = id[6]
			iar[2] = id[7]
			iar[3] = id[8]
			iar[4] = id[9]
			iar[5] = id[10]
			iar[6] = id[11]
			z.SetBytes(iar)
			urn += z.String()
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
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 7) // 55 bits
			remainder = id[5] & 127
			iar[0] = remainder
			iar[1] = id[6]
			iar[2] = id[7]
			iar[3] = id[8]
			iar[4] = id[9]
			iar[5] = id[10]
			iar[6] = id[11]
			z.SetBytes(iar)
			urn += z.String()
		case 24:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<6 | id[2]>>2
			remainder = id[2] & 3
			cp[1] = remainder<<6 | id[3]>>2
			remainder = id[3] & 3
			cp[2] = remainder<<6 | id[4]>>2
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 8) // 58 bits
			remainder = id[4] & 3
			iar[0] = remainder
			iar[1] = id[5]
			iar[2] = id[6]
			iar[3] = id[7]
			iar[4] = id[8]
			iar[5] = id[9]
			iar[6] = id[10]
			iar[7] = id[11]
			z.SetBytes(iar)
			urn += z.String()
		case 20:
			cp := make([]byte, 3)
			remainder := id[1] & 3
			cp[0] = remainder<<2 | id[2]>>6
			remainder = id[2] & 63
			cp[1] = remainder<<2 | id[3]>>6
			remainder = id[3] & 63
			cp[2] = remainder<<2 | id[4]>>6
			z.SetBytes(cp)
			companyPrefix := z.String()
			urn += strings.Repeat("0", cpLength-len(companyPrefix)) + companyPrefix + "."
			iar := make([]byte, 8) // 62 bits
			remainder = id[4] & 63
			iar[0] = remainder
			iar[1] = id[5]
			iar[2] = id[6]
			iar[3] = id[7]
			iar[4] = id[8]
			iar[5] = id[9]
			iar[6] = id[10]
			iar[7] = id[11]
			z.SetBytes(iar)
			urn += z.String()
		}
	}
	return urn, nil
}

func (c *Core) buildUII(id []byte, afi byte) (string, error) {
	urn := "urn:epc:id:iso"
	switch afi {
	case 161:
		urn += "17367:"
	case 162:
		urn += "17365:"
	case 163:
		urn += "17364:"
	case 164:
		urn += "17367h:"
	case 165:
		urn += "17366:"
	case 166:
		urn += "17366h:"
	case 167:
		urn += "17365h:"
	case 168:
		urn += "17364h:"
	case 169:
		urn += "17363:"
	case 170:
		urn += "17363h:"
	default:
		return "", errors.New("invalid afi")
	}

	sid, err := parse6BitEncodedByteSliceToString(id)
	if err != nil {
		return "", err
	}

	urn += sid
	return urn, nil
}

func (c *Core) buildProprietary(id []byte) (string, error) {
	return "", nil
}

// MakePrefixFilterString takes a pattern type and a slice of fields
// return a binary reporesentation of the prefix filter in string
func MakePrefixFilterString(patternType string, fields []string) (string, error) {
	switch patternType { // type
	case "giai-96":
		return NewPrefixFilterGIAI96(fields)
	case "grai-96":
		return NewPrefixFilterGRAI96(fields)
	case "sgtin-96":
		return NewPrefixFilterSGTIN96(fields)
	case "sscc-96":
		return NewPrefixFilterSSCC96(fields)
	case "iso17363":
		return NewPrefixFilterISO17363(fields)
	case "iso17365":
		return NewPrefixFilterISO17365(fields)
	default:
		return "", fmt.Errorf("unknown patternType: %v", patternType)
	}
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
