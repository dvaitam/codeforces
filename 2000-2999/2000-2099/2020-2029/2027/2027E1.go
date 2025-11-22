package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func grundy(a, x int) int {
	if a == 0 || x == 0 {
		return 0
	}
	k := bits.Len(uint(a)) - 1 // highest bit present in a
	pow := 1 << k              // 2^k
	b := a - pow               // lower part of a
	x &= (pow << 1) - 1        // bits above k never change
	low := x & (pow - 1)       // bits strictly below k
	highBit := (x >> k) & 1
	if highBit == 0 {
		return bits.OnesCount(uint(low))
	}
	if low == 0 {
		return 1
	}
	lsb := low & -low // value of the lowest set bit
	pop := bits.OnesCount(uint(low))

	if lsb > b {
		return 1 ^ pop // behaves like b = 0, only bit k removable with u = 0
	}
	if low > b && low-lsb <= b {
		return 0 // removing bit k and that lowest bit empties the movable part
	}
	return pop + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		x := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}

		xorAll := 0
		for i := 0; i < n; i++ {
			xorAll ^= grundy(a[i], x[i])
		}

		if xorAll != 0 {
			fmt.Fprintln(out, "Alice")
		} else {
			fmt.Fprintln(out, "Bob")
		}
	}
}
