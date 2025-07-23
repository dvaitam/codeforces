package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func cbrtInt(x int64) int64 {
	r := int64(math.Round(math.Cbrt(float64(x))))
	for (r+1)*(r+1)*(r+1) <= x {
		r++
	}
	for r*r*r > x {
		r--
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	for i := 0; i < n; i++ {
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		prod := a * b
		c := cbrtInt(prod)
		if c*c*c != prod || a%c != 0 || b%c != 0 {
			fmt.Fprintln(writer, "No")
			continue
		}
		x := a / c
		y := b / c
		if x*y == c {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
