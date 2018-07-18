// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

// Package tdt contains Tag Data Translation module from binary to Pure Identity
package tdt

import (
	"fmt"
	"strings"

	"github.com/iomz/go-llrp/binutil"
)

// getISO6346CD returns check digit for container serial number
func getISO6346CD(cn string) (int, error) {
	if len(cn) != 10 {
		return 0, fmt.Errorf("Invalid ISO6346 code provided: %v", cn)
	}
	n := 0.0
	d := 0.5
	for i := 0; i < 10; i++ {
		d *= 2
		n += d * float64(strings.Index("0123456789A?BCDEFGHIJK?LMNOPQRSTU?VWXYZ", string(cn[i])))
	}
	return (int(n) - int(n/11)*11) % 10, nil
}

// pad6BitEncodingRuneSlice returns a new length
// and 16-bit (word-length) padded binary string in rune slice
// @ISO15962
func pad6BitEncodingRuneSlice(bs []rune) ([]rune, int) {
	length := len(bs)
	remainder := length % 16
	var padding []rune
	if remainder != 0 {
		padRuneSlice := binutil.ParseDecimalStringToBinRuneSlice("32") // pad string "100000"
		for i := 0; i < 16-remainder; i++ {
			padding = append(padding, padRuneSlice[i%6])
		}
		bs = append(bs, padding...)
		length += 16 - remainder
	}
	return bs, length
}

// NewPrefixFilterISO17363 takes fields and return the prefix filter in string
func NewPrefixFilterISO17363(fields []string) (string, error) {
	nFields := len(fields) // ownerCode, equipmentIdentifier, containerSerialNumber

	if nFields == 0 {
		return "", fmt.Errorf("wrong fields: %q", fields)
	}

	// dataIdentifier
	dataIdentifier := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[0]))
	if nFields == 1 {
		return string(dataIdentifier), nil
	}

	// ownerCode
	ownerCode := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[1]))
	if nFields == 2 {
		return string(dataIdentifier) + string(ownerCode), nil
	}

	// equipmentIdentifier
	equipmentIdentifier := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[2]))
	if nFields == 3 {
		return string(dataIdentifier) + string(ownerCode) + string(equipmentIdentifier), nil
	}

	// containerSerialNumber
	csn := fields[3]
	if 6 > len(csn) {
		leftPadding := binutil.GenerateNLengthZeroPaddingRuneSlice(6 - len(csn))
		csn = string(leftPadding) + csn
	} else if 6 < len(csn) {
		return "", fmt.Errorf("Invalid csn: %v", csn)
	}
	cd, err := getISO6346CD(fields[0] + fields[1] + csn)
	if err != nil {
		return "", err
	}
	containerSerialNumber := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(csn + fmt.Sprintf("%v", cd)))
	if nFields == 4 {
		return string(dataIdentifier) + string(ownerCode) + string(equipmentIdentifier) + string(containerSerialNumber), nil
	}

	return "", fmt.Errorf("unknown fields provided: %q", fields)
}

// NewPrefixFilterISO17365 takes fields and return the prefix filter in string
func NewPrefixFilterISO17365(fields []string) (string, error) {
	nFields := len(fields) // dataIdentifier, issuingAgencyCode, companyIdentification, serialNumber

	if nFields == 0 {
		return "", fmt.Errorf("wrong fields: %q", fields)
	}

	// dataIdentifier
	dataIdentifier := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[0]))
	if nFields == 1 {
		return string(dataIdentifier), nil
	}

	// issuingAgencyCode
	issuingAgencyCode := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[1]))
	if nFields == 2 {
		return string(dataIdentifier) + string(issuingAgencyCode), nil
	}

	// companyIdentification
	companyIdentification := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[2]))
	if nFields == 3 {
		return string(dataIdentifier) + string(issuingAgencyCode) + string(companyIdentification), nil
	}

	// serialNumber
	serialNumber := binutil.ParseRuneSliceTo6BinRuneSlice([]rune(fields[3]))
	if nFields == 4 {
		return string(dataIdentifier) + string(issuingAgencyCode) + string(companyIdentification) + string(serialNumber), nil
	}

	return "", fmt.Errorf("unknown fields provided: %q", fields)
}
