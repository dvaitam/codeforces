package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// F returns the sum of indices on the path from x to the root.
func F(x int64) int64 {
	return 2*x - int64(bits.OnesCount64(uint64(x)))
}

// pathSum returns the sum of indices along the simple path between u and v.
func pathSum(u, v int64) int64 {
	a, b := u, v
	for a != b {
		if a > b {
			a >>= 1
		} else {
			b >>= 1
		}
	}
	l := a
	return F(u) + F(v) - 2*F(l) + l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s int64
	fmt.Fscan(in, &s)

	var count int64
	for u := int64(1); u <= s; u++ {
		for v := u; v <= s; v++ {
			if pathSum(u, v) == s {
				count++
			}
		}
	}
	fmt.Println(count)
}
