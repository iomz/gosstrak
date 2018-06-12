// Copyright (c) 2018 Iori Mizutani
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tdt

func translate(afi byte, uii []byte) (string, error) {
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
