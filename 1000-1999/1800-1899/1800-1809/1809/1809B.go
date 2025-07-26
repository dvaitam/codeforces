package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isqrt(x int64) int64 {
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		if n == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		r := isqrt(n - 1)
		fmt.Fprintln(writer, r)
	}
}
