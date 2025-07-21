package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t, a, b int64
	if _, err := fmt.Fscan(reader, &t, &a, &b); err != nil {
		return
	}
	const mod = 1000000007
	// Special infinite case: t=1, a=1, b=1
	if t == 1 && a == 1 {
		if b == 1 {
			fmt.Println("inf")
		} else {
			fmt.Println(0)
		}
		return
	}
	// Compute C = sum di * a^i, where di are base-t digits of a
	C := int64(0)
	powA := int64(1)
	aa := a
	for aa > 0 {
		d := aa % t
		C += d * powA
		aa /= t
		// safeguard overflow
		if powA > (1<<62)/max(a, 1) {
			powA = 0
		} else {
			powA *= a
		}
	}
	denom := t*a - 1
	if denom <= 0 || b < C || (b-C)%denom != 0 {
		fmt.Println(0)
	} else {
		// assume unique solution exists
		fmt.Println(1)
	}
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
