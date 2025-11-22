package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBit = 60

type xorBasis struct {
	b [maxBit]uint64
}

func (xb *xorBasis) reduce(x uint64) uint64 {
	for i := maxBit - 1; i >= 0; i-- {
		if xb.b[i] != 0 && (x&(1<<i)) != 0 {
			x ^= xb.b[i]
		}
	}
	return x
}

func (xb *xorBasis) insert(x uint64) {
	for i := maxBit - 1; i >= 0; i-- {
		if (x & (1 << i)) == 0 {
			continue
		}
		if xb.b[i] == 0 {
			xb.b[i] = x
			return
		}
		x ^= xb.b[i]
	}
}

func msb(x uint64) int {
	for i := maxBit - 1; i >= 0; i-- {
		if (x & (1 << i)) != 0 {
			return i
		}
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]uint64, n)
		b := make([]uint64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		var c string
		fmt.Fscan(in, &c)

		var base uint64
		for i := 0; i < n; i++ {
			base ^= a[i]
		}
		ans := base
		var basis xorBasis

		for i := n - 1; i >= 0; i-- {
			d := a[i] ^ b[i]
			d = basis.reduce(d)
			if d == 0 {
				continue
			}
			pivot := msb(d)
			if c[i] == '0' {
				// Gellyfish minimizes: make pivot bit 0 if possible.
				if (ans>>pivot)&1 == 1 {
					ans ^= d
				}
			} else {
				// Flower maximizes: make pivot bit 1 if possible.
				if (ans>>pivot)&1 == 0 {
					ans ^= d
				}
			}
			basis.insert(d)
		}

		fmt.Fprintln(out, ans)
	}
}
