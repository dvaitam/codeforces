package main

import (
	"bufio"
	"fmt"
	"os"
)

// The task is to make a and b equal using the following operations:
// 1) increment a, 2) increment b and 3) set a := a | b.
// We iterate over all candidates B >= b up to 2*b. We first increment b to B,
// then apply the OR operation once which changes a to a|B. If a|B is still
// larger than B, we additionally increment b to match it. The total number of
// operations for a given B is (B-b) + 1 + (a|B-B). The answer is the minimum
// of this value over all B, also considering simply incrementing a to b.
func solveCase(a, b int) int {
	if b == a {
		return 0
	}
	ans := b - a
	for B := b; B <= b*2; B++ {
		cost := B - b
		y := a | B
		cost += y - B
		cost++ // one OR operation
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		fmt.Fprintln(writer, solveCase(a, b))
	}
}
