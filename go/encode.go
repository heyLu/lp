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

func decode(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b - 'A'
	} else if 'a' <= b && b <= 'z' {
		return 26 + b - 'a'
	} else if '0' <= b && b <= '9' {
		return 52 + b - '0'
	} else if b == '+' {
		return 62
	} else if b == '/' {
		return 63
	} else {
		fmt.Println("oops")
		return ' '
	}
}

func base64decode(b []byte) []byte {
	res := make([]byte, len(b) / 4 * 3)
	//fmt.Println(len(b), "->", len(res))
	for i := 0; i < len(b); i += 4 {
		c1, c2, c3, c4 := decode(b[i]), decode(b[i+1]), decode(b[i+2]), decode(b[i+3])
		r1 := (c1 << 2) + (c2 >> 4)
		r2 := ((c2 ^ ((c2 >> 4) << 4)) << 4) + (c3 >> 2)
		r3 := ((c3 ^ ((c3 >> 2) << 2)) << 6) + c4
		base := i - i / 4
		//fmt.Print(i, base, c1, c2, c3, c4, "\t", string(b[i]), string(b[i+1]), string(b[i+2]), string(b[i+3]))
		//fmt.Printf("\t -> '%s' '%s' '%s'\n", string(r1), string(r2), string(r3))
		res[base + 0] = r1
		res[base + 1] = r2
		res[base + 2] = r3
	}
	return res
}

func main() {
	base64 := base64encode([]byte("Hello, World!  "))
	//                                          ^^     (padded to a multiple of three)
	fmt.Println(string(base64))
	plain := base64decode(base64)
	fmt.Println(string(plain))

	secret := []byte("SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t")
	fmt.Println(string(secret))
	decoded := base64decode(secret)
	fmt.Println(string(decoded))
}
