package main

import (
	"bufio"
	"fmt"
	"os"
)

func applyCeil(x, m int64) int64 {
	for m > 0 && x > 1 {
		x = (x + 1) >> 1
		m--
	}
	return x
}

func applyFloor(x, n int64) int64 {
	for n > 0 && x > 0 {
		x >>= 1
		n--
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, n, m int64
		fmt.Fscan(in, &x, &n, &m)
		minVal := applyFloor(applyCeil(x, m), n)
		maxVal := applyCeil(applyFloor(x, n), m)
		fmt.Fprintln(out, minVal, maxVal)
	}
}
