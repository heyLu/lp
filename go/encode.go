package main

import "fmt"

// encodes a six-bit integer (max = 2^6 - 1 = 63) as a character in A-Za-z0-9+/
func encode(b byte) byte {
	if b < 26 {
		return 'A' + b
	} else if b < 52 {
		return 'a' + (b - 26)
	} else if b < 62 {
		return '0' + (b - 52)
	} else if b == 62 {
		return '+'
	} else if b == 63 {
		return '/'
	} else {
		fmt.Println("oops")
		return ' '
	}
}

func base64encode(b []byte) []byte {
	res := make([]byte, len(b) / 3 * 4)
	//fmt.Println(len(b), "->", len(res))
	for i := 0; i < len(b); i += 3 {
		c1, c2, c3 := b[i], b[i+1], b[i+2]
		r1 := c1 >> 2
		r2 := ((c1 ^ (r1 << 2)) << 4) + (c2 >> 4)
		r3 := ((c2 ^ ((c2 >> 4) << 4)) << 2) + (c3 >> 6)
		r4 := c3 ^ ((c3 >> 6) << 6)
		//fmt.Println(i, c1, c2, c3)
		base := i + i / 3
		res[base + 0] = encode(r1)
		res[base + 1] = encode(r2)
		res[base + 2] = encode(r3)
		res[base + 3] = encode(r4)
	}
	return res
}

func main() {
	base64 := base64encode([]byte("Hello, World!  "))
	//                                          ^^     (padded to a multiple of three)
	fmt.Println(string(base64))
}
