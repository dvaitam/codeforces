package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

func solve(a, b, c int) int {
	total := a + b + c
	base := ceilDiv(total, 3)
	cand1 := max(base, max(a, c))
	cand2 := max(base, max(c, ceilDiv(a, 2)))
	cand3 := max(base, max(a, ceilDiv(c, 2)))
	ans := cand1
	if cand2 < ans {
		ans = cand2
	}
	if cand3 < ans {
		ans = cand3
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
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		fmt.Fprintln(writer, solve(a, b, c))
	}
}
