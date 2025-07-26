package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	// count total question marks
	Q := 0
	bytes := []byte(s)
	for _, ch := range bytes {
		if ch == '?' {
			Q++
		}
	}

	// precompute powers k^e for k=1..17, e=0..Q
	powTab := make([][]int64, 18)
	for k := 1; k <= 17; k++ {
		powTab[k] = make([]int64, Q+1)
		powTab[k][0] = 1
		for e := 1; e <= Q; e++ {
			powTab[k][e] = powTab[k][e-1] * int64(k) % mod
		}
	}

	maxMask := 1 << 17
	// h[k][mask]
	h := make([][]int64, 18)
	for k := 1; k <= 17; k++ {
		h[k] = make([]int64, maxMask)
	}

	// helper to add value to h for current mask and reduction
	add := func(mask int, reduction int) {
		exp := Q - reduction
		for k := 1; k <= 17; k++ {
			h[k][mask] = (h[k][mask] + powTab[k][exp]) % mod
		}
	}

	nBytes := len(bytes)
	// odd length palindromes
	for center := 0; center < nBytes; center++ {
		mask := 0
		reduction := 0
		for l, r := center, center; l >= 0 && r < nBytes; l, r = l-1, r+1 {
			if l == r { // single char
				if bytes[l] == '?' {
					// nothing to do
				}
			} else {
				a := bytes[l]
				b := bytes[r]
				if a != '?' && b != '?' {
					if a != b {
						break
					}
				} else if a == '?' && b == '?' {
					reduction++
				} else if a == '?' { // b is letter
					reduction++
					mask |= 1 << (b - 'a')
				} else { // a is letter, b is '?'
					reduction++
					mask |= 1 << (a - 'a')
				}
			}
			add(mask, reduction)
		}
	}

	// even length palindromes
	for center := 0; center+1 < nBytes; center++ {
		mask := 0
		reduction := 0
		for l, r := center, center+1; l >= 0 && r < nBytes; l, r = l-1, r+1 {
			a := bytes[l]
			b := bytes[r]
			if a != '?' && b != '?' {
				if a != b {
					break
				}
			} else if a == '?' && b == '?' {
				reduction++
			} else if a == '?' {
				reduction++
				mask |= 1 << (b - 'a')
			} else {
				reduction++
				mask |= 1 << (a - 'a')
			}
			add(mask, reduction)
		}
	}

	// subset sums
	for k := 1; k <= 17; k++ {
		f := h[k]
		for bit := 0; bit < 17; bit++ {
			for mask := 0; mask < maxMask; mask++ {
				if mask&(1<<bit) != 0 {
					f[mask] = (f[mask] + f[mask^(1<<bit)]) % mod
				}
			}
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var t string
		fmt.Fscan(reader, &t)
		mask := 0
		for i := 0; i < len(t); i++ {
			mask |= 1 << (t[i] - 'a')
		}
		k := len(t)
		fmt.Fprintln(writer, h[k][mask]%mod)
	}
}
