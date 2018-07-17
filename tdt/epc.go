// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

// Package tdt contains Tag Data Translation module from binary to Pure Identity
package tdt

import (
	"fmt"
	"math"
	"strconv"

	"github.com/iomz/go-llrp/binutil"
)

// PartitionTableKey is used for PartitionTables
type PartitionTableKey int

// PartitionTable is used to get the related values for each coding scheme
type PartitionTable map[int]map[PartitionTableKey]int

// Key values for PartitionTables
const (
	PValue PartitionTableKey = iota
	CPBits
	IRBits
	IRDigits
	EBits
	EDigits
	ATBits
	ATDigits
	IARBits
	IARDigits
)

// GIAI96PartitionTable is PT for GIAI
var GIAI96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, IARBits: 42, IARDigits: 13},
	11: {PValue: 1, CPBits: 37, IARBits: 45, IARDigits: 14},
	10: {PValue: 2, CPBits: 34, IARBits: 48, IARDigits: 15},
	9:  {PValue: 3, CPBits: 30, IARBits: 52, IARDigits: 16},
	8:  {PValue: 4, CPBits: 27, IARBits: 55, IARDigits: 17},
	7:  {PValue: 5, CPBits: 24, IARBits: 58, IARDigits: 18},
	6:  {PValue: 6, CPBits: 20, IARBits: 62, IARDigits: 19},
}

// GRAI96PartitionTable is PT for GRAI
var GRAI96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, ATBits: 4, ATDigits: 0},
	11: {PValue: 1, CPBits: 37, ATBits: 7, ATDigits: 1},
	10: {PValue: 2, CPBits: 34, ATBits: 10, ATDigits: 2},
	9:  {PValue: 3, CPBits: 30, ATBits: 14, ATDigits: 3},
	8:  {PValue: 4, CPBits: 27, ATBits: 17, ATDigits: 4},
	7:  {PValue: 5, CPBits: 24, ATBits: 20, ATDigits: 5},
	6:  {PValue: 6, CPBits: 20, ATBits: 24, ATDigits: 6},
}

// SGTIN96PartitionTable is PT for SGTIN
var SGTIN96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, IRBits: 4, IRDigits: 1},
	11: {PValue: 1, CPBits: 37, IRBits: 7, IRDigits: 2},
	10: {PValue: 2, CPBits: 34, IRBits: 10, IRDigits: 3},
	9:  {PValue: 3, CPBits: 30, IRBits: 14, IRDigits: 4},
	8:  {PValue: 4, CPBits: 27, IRBits: 17, IRDigits: 5},
	7:  {PValue: 5, CPBits: 24, IRBits: 20, IRDigits: 6},
	6:  {PValue: 6, CPBits: 20, IRBits: 24, IRDigits: 7},
}

// SSCC96PartitionTable is PT for SSCC
var SSCC96PartitionTable = PartitionTable{
	12: {PValue: 0, CPBits: 40, EBits: 18, EDigits: 5},
	11: {PValue: 1, CPBits: 37, EBits: 21, EDigits: 6},
	10: {PValue: 2, CPBits: 34, EBits: 24, EDigits: 7},
	9:  {PValue: 3, CPBits: 30, EBits: 28, EDigits: 8},
	8:  {PValue: 4, CPBits: 27, EBits: 31, EDigits: 9},
	7:  {PValue: 5, CPBits: 24, EBits: 34, EDigits: 10},
	6:  {PValue: 6, CPBits: 20, EBits: 38, EDigits: 11},
}

// getAssetType returns Asset Type as rune slice
func getAssetType(at string, pr map[PartitionTableKey]int) (assetType []rune) {
	if at != "" {
		assetType = binutil.ParseDecimalStringToBinRuneSlice(at)
		if pr[ATBits] > len(assetType) {
			leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(pr[ATBits] - len(assetType))
			assetType = append(leftPadding, assetType...)
		}
	} else {
		assetType, _ = binutil.GenerateNLengthRandomBinRuneSlice(pr[ATBits], uint(math.Pow(float64(10), float64(pr[ATDigits]))))
	}
	return
}

// getCompanyPrefix returns Company Prefix as rune slice
func getCompanyPrefix(cp string, pt PartitionTable) (companyPrefix []rune) {
	if cp != "" {
		companyPrefix = binutil.ParseDecimalStringToBinRuneSlice(cp)
		if pt[len(cp)][CPBits] > len(companyPrefix) {
			leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(pt[len(cp)][CPBits] - len(companyPrefix))
			companyPrefix = append(leftPadding, companyPrefix...)
		}
	}
	return
}

// getExtension returns Extension digit and Serial Reference as rune slice
func getExtension(e string, pr map[PartitionTableKey]int) (extension []rune) {
	if e != "" {
		extension = binutil.ParseDecimalStringToBinRuneSlice(e)
		if pr[EBits] > len(extension) {
			leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(pr[EBits] - len(extension))
			extension = append(leftPadding, extension...)
		}
	} else {
		extension, _ = binutil.GenerateNLengthRandomBinRuneSlice(pr[EBits], uint(math.Pow(float64(10), float64(pr[EDigits]))))
	}
	return
}

// getFilter returns filter value as rune slice
func getFilter(fv string) (filter []rune) {
	if fv != "" {
		n, _ := strconv.ParseInt(fv, 10, 32)
		filter = []rune(fmt.Sprintf("%.3b", n))
	} else {
		filter, _ = binutil.GenerateNLengthRandomBinRuneSlice(3, 7)
	}
	return
}

// getIndivisualAssetReference returns iar as rune slice
func getIndivisualAssetReference(iar string, pr map[PartitionTableKey]int) (indivisualAssetReference []rune) {
	if iar != "" {
		indivisualAssetReference = binutil.ParseDecimalStringToBinRuneSlice(iar)
		if pr[IARBits] > len(indivisualAssetReference) {
			leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(pr[IARBits] - len(indivisualAssetReference))
			indivisualAssetReference = append(leftPadding, indivisualAssetReference...)
		}
	} else {
		indivisualAssetReference, _ = binutil.GenerateNLengthRandomBinRuneSlice(pr[IARBits], uint(math.Pow(float64(10), float64(pr[IARDigits]))))
	}
	return
}

// getItemReference converts ItemReference value to rune slice
func getItemReference(ir string, pr map[PartitionTableKey]int) (itemReference []rune) {
	if ir != "" {
		itemReference = binutil.ParseDecimalStringToBinRuneSlice(ir)
		// If the itemReference is short, pad zeroes to the left
		if pr[IRBits] > len(itemReference) {
			leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(pr[IRBits] - len(itemReference))
			itemReference = append(leftPadding, itemReference...)
		}
	} else {
		itemReference, _ = binutil.GenerateNLengthRandomBinRuneSlice(pr[IRBits], uint(math.Pow(float64(10), float64(pr[IRDigits]))))
	}
	return
}

// getSerial converts serial to rune slice
func getSerial(s string, serialLength int) (serial []rune) {
	if s != "" {
		serial = binutil.ParseDecimalStringToBinRuneSlice(s)
		if serialLength > len(serial) {
			leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(serialLength - len(serial))
			serial = append(leftPadding, serial...)
		}
	} else {
		serial, _ = binutil.GenerateNLengthRandomBinRuneSlice(serialLength, uint(math.Pow(float64(2), float64(serialLength))))
	}
	return serial
}

// NewPrefixFilterGIAI96 takes field values in a slice and return a prefix filter string
func NewPrefixFilterGIAI96(fields []string) (string, error) {
	nFields := len(fields) // filter, companyPrefix, indivisualAssetReference
	if nFields == 0 {
		return "", fmt.Errorf("wrong fields: %q", fields)
	}

	// filter
	filter := getFilter(fields[0])
	if nFields == 1 {
		return "00110100" + string(filter), nil
	}

	// companyPrefix
	companyPrefix := getCompanyPrefix(fields[1], GIAI96PartitionTable)
	partition := []rune(fmt.Sprintf("%.3b", GIAI96PartitionTable[len(fields[1])][PValue]))
	if nFields == 2 {
		return "00110100" + string(filter) + string(partition) + string(companyPrefix), nil
	}

	// indivisualAssetReference
	indivisualAssetReference := getIndivisualAssetReference(fields[2], GIAI96PartitionTable[len(fields[1])])
	if nFields == 3 {
		return "00110100" + string(filter) + string(partition) + string(companyPrefix) + string(indivisualAssetReference), nil
	}

	return "", fmt.Errorf("unknown fields provided %q", fields)
}

// NewPrefixFilterGRAI96 takes field values in a slice and return a prefix filter string
func NewPrefixFilterGRAI96(fields []string) (string, error) {
	nFields := len(fields) // filter, companyPrefix, assetType, serial
	if nFields == 0 {
		return "", fmt.Errorf("wrong fields: %q", fields)
	}

	// filter
	filter := getFilter(fields[0])
	if nFields == 1 {
		return "00110011" + string(filter), nil
	}

	// companyPrefix
	companyPrefix := getCompanyPrefix(fields[1], GRAI96PartitionTable)
	partition := []rune(fmt.Sprintf("%.3b", GRAI96PartitionTable[len(fields[1])][PValue]))
	if nFields == 2 {
		return "00110011" + string(filter) + string(partition) + string(companyPrefix), nil
	}

	// assetType
	assetType := getAssetType(fields[2], GRAI96PartitionTable[len(fields[1])])
	if nFields == 3 {
		return "00110011" + string(filter) + string(partition) + string(companyPrefix) + string(assetType), nil
	}

	// serial
	serial := getSerial(fields[3], 38)
	if nFields == 4 {
		return "00110011" + string(filter) + string(partition) + string(companyPrefix) + string(assetType) + string(serial), nil
	}

	return "", fmt.Errorf("unknown fields provided %q", fields)
}

// NewPrefixFilterSGTIN96 takes field values in a slice and return a prefix filter string
func NewPrefixFilterSGTIN96(fields []string) (string, error) {
	nFields := len(fields) // filter, companyPrefix, itemReference, serial
	if nFields == 0 {
		return "", fmt.Errorf("wrong fields: %q", fields)
	}

	// filter
	filter := getFilter(fields[0])
	if nFields == 1 {
		return "00110000" + string(filter), nil
	}

	// companyPrefix
	companyPrefix := getCompanyPrefix(fields[1], SGTIN96PartitionTable)
	partition := []rune(fmt.Sprintf("%.3b", SGTIN96PartitionTable[len(fields[1])][PValue]))
	if nFields == 2 {
		return "00110000" + string(filter) + string(partition) + string(companyPrefix), nil
	}

	// itemReference
	itemReference := getItemReference(fields[2], SGTIN96PartitionTable[len(fields[1])])
	if nFields == 3 {
		return "00110000" + string(filter) + string(partition) + string(companyPrefix) + string(itemReference), nil
	}

	// serial
	serial := getSerial(fields[3], 38)
	if nFields == 4 {
		return "00110000" + string(filter) + string(partition) + string(companyPrefix) + string(itemReference) + string(serial), nil
	}

	return "", fmt.Errorf("unknown fields provided %q", fields)
}

// NewPrefixFilterSSCC96 takes field values in a slice and return a prefix filter string
func NewPrefixFilterSSCC96(fields []string) (string, error) {
	nFields := len(fields) // filter, companyPrefix, extension
	if nFields == 0 {
		return "", fmt.Errorf("wrong fields: %q", fields)
	}

	// filter
	filter := getFilter(fields[0])
	if nFields == 1 {
		return "00110001" + string(filter), nil
	}

	// companyPrefix
	companyPrefix := getCompanyPrefix(fields[1], SSCC96PartitionTable)
	partition := []rune(fmt.Sprintf("%.3b", SSCC96PartitionTable[len(fields[1])][PValue]))
	if nFields == 2 {
		return "00110001" + string(filter) + string(partition) + string(companyPrefix), nil
	}

	// extension
	extension := getExtension(fields[2], SSCC96PartitionTable[len(fields[1])])
	if nFields == 3 {
		return "00110001" + string(filter) + string(partition) + string(companyPrefix) + string(extension), nil
	}

	return "", fmt.Errorf("unknown fields provided %q", fields)
}
