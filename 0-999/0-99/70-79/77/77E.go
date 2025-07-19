package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for i := 0; i < T; i++ {
		var R, r, K float64
		fmt.Fscan(reader, &R, &r, &K)
		delta := R - r
		x1 := R + r
		y1 := 2 * delta * K
		d1 := math.Sqrt(x1*x1 + y1*y1)
		ans := 2 * R * r / (d1 - delta)
		ans -= 2 * R * r / (d1 + delta)
		fmt.Fprintf(writer, "%.9f\n", ans)
	}
}
