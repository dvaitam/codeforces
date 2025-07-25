package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	even, odd := int64(1), int64(0)
	pos, neg := int64(0), int64(0)
	parity := 0

	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		if x < 0 {
			parity ^= 1
		}
		if parity == 0 {
			pos += even
			neg += odd
			even++
		} else {
			pos += odd
			neg += even
			odd++
		}
	}

	fmt.Fprintf(out, "%d %d\n", neg, pos)
}
