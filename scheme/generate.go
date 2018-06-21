// Package scheme contains scheme related utilities
package scheme

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/iomz/go-llrp/binutil"
)

// CheckIfStringInSlice checks if string exists in a string slice
// TODO: fix the way it is, it should be smarter
func CheckIfStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// MakeEPC generates EPC in binary string and PC in hex string or prefix and elements
func MakeEPC(pf bool, cs string, fv string, cp string, ir string, ext string, ser string, iar string, at string) (string, string) {
	var uii []byte
	var f string
	var elem string

	switch strings.ToUpper(cs) {
	case "GIAI-96":
		uii, f, elem, _ = MakeGIAI96(pf, fv, cp, iar)
	case "GRAI-96":
		uii, f, elem, _ = MakeGRAI96(pf, fv, cp, at, ser)
	case "SGTIN-96":
		uii, f, elem, _ = MakeSGTIN96(pf, fv, cp, ir, ser)
	case "SSCC-96":
		uii, f, elem, _ = MakeSSCC96(pf, fv, cp, ext)
	}

	// If only prefix flag is on, return prefix as epc
	if pf {
		return f, elem
	}

	// TODO: update pc when length changed (for non-96-bit codes)
	pc := binutil.Pack([]interface{}{
		uint8(48), // L4-0=11000(6words=96bits), UMI=0, XI=0
		uint8(0),  // RFU=0
	})

	uiibs, _ := binutil.ParseHexStringToBinString(hex.EncodeToString(uii))

	return uiibs, hex.EncodeToString(pc)
	/*
		length := uint16(18)
		epclen := uint16(96)
	*/
}

// MakeISO returns ISO UII in binary string and PCbits in hex string or prefix and elements
func MakeISO(pf bool, std string, oc string, ei string, csn string, di string, iac string, cin string, sn string) (string, string) {
	var uii []byte
	var pc []byte
	var length int
	var f string
	var elem string

	switch std {
	case "17363":
		afi := "A9" // 0xA9 ISO 17363 freight containers
		uii, length, f, elem, _ = MakeISO17363(pf, oc, ei, csn)
		pc = MakeISOPC(length, afi)
	case "17365":
		afi := "A2" // 0xA2 ISO 17365 transport uit
		uii, length, f, elem, _ = MakeISO17365(pf, di, iac, cin, sn)
		pc = MakeISOPC(length, afi)
	}

	// If only prefix flag is on, return prefix as iso uii
	if pf {
		return f, elem
	}

	uiibs, _ := binutil.ParseHexStringToBinString(hex.EncodeToString(uii))

	return uiibs, hex.EncodeToString(pc)
	/*
		return hex.EncodeToString(pc) + "," +
			strconv.FormatUint(uint64(length/16), 10) + "," +
			strconv.FormatUint(uint64(length), 10) + "," +
			hex.EncodeToString(uii) + "\n" +
			uiibs
	*/
}

// MakeISOPC returns PC bits in []byte
func MakeISOPC(length int, afi string) []byte {
	l := []rune(fmt.Sprintf("%.5b", length/16))
	pc1, err := binutil.ParseBinRuneSliceToUint8Slice(append(l, rune('0'), rune('0'), rune('1'))) // L, UMI, XI, T
	if err != nil {
		panic(err)
	}
	c, _ := strconv.ParseUint(afi, 16, 8)
	return binutil.Pack([]interface{}{
		pc1[0],
		uint8(c), // AFI
	})
}

// Sixenc2bin 6bit encoded string to binary
func Sixenc2bin(sixenc []rune) []rune {
	var bs []rune
	for i := 0; i < len(sixenc); i++ {
		r := binutil.ParseRuneTo6BinRuneSlice(sixenc[i])
		bs = append(bs, r...)
	}
	return bs
}

// PrintGoBytes print the ID in go-test ready []byte format
func PrintGoBytes(bs string, opt string) string {
	if len(bs) == 0 {
		return ""
	}
	ds, _ := binutil.ParseBinStringToDecArrayString(bs)
	pc, _ := binutil.ParseHexStringToDecArrayString(opt)
	return "[]byte{" + pc + "},[]byte{" + ds + "}"
}

// PrintID print the ID for text files
func PrintID(bs string, opt string) string {
	if len(bs) == 0 {
		return ""
	}
	return opt + "," + bs
}

/*
	bs, opt = MakeEPC(*prefixFilter, *epcScheme, *epcFilter, *epcCompanyPrefix, *epcItemReference, *epcExtension, *epcSerial, *epcIndivisualAssetReference, *epcAssetType)
	bs, opt = MakeISO(*prefixFilter, "17363", *isoOwnerCode, *isoEquipmentCategoryIdentifier, *isoContainerSerialNumber, *isoDataIdeintifier, *isoIssuingAgencyCode, *isoCompanyIdentification, *isoSerialNumber)
	bs, opt = MakeISO(*prefixFilter, "17365", *isoOwnerCode, *isoEquipmentCategoryIdentifier, *isoContainerSerialNumber, *isoDataIdeintifier, *isoIssuingAgencyCode, *isoCompanyIdentification, *isoSerialNumber)
*/
