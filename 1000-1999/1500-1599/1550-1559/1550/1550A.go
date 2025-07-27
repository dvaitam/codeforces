package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func ceilSqrt(x int) int {
	if x <= 0 {
		return 0
	}
	r := int(math.Sqrt(float64(x)))
	if r*r < x {
		r++
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s int
		fmt.Fscan(reader, &s)
		ans := ceilSqrt(s)
		fmt.Fprintln(writer, ans)
	}
}
