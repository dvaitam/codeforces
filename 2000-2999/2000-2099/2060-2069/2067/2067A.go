package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumDigits(a int) int {
	s := 0
	for a > 0 {
		s += a % 10
		a /= 10
	}
	return s
}

func buildPairs(limit int) map[[2]int]struct{} {
	pairs := make(map[[2]int]struct{})
	for n := 0; n <= limit; n++ {
		x := sumDigits(n)
		y := sumDigits(n + 1)
		pairs[[2]int{x, y}] = struct{}{}
	}
	return pairs
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	pairs := buildPairs(2_000_000)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if _, ok := pairs[[2]int{x, y}]; ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
