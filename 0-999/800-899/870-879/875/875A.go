package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumDigits(x int) int {
	s := 0
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	start := n - 100
	if start < 1 {
		start = 1
	}
	var ans []int
	for x := start; x <= n; x++ {
		if x+sumDigits(x) == n {
			ans = append(ans, x)
		}
	}

	fmt.Fprintln(out, len(ans))
	for _, v := range ans {
		fmt.Fprintln(out, v)
	}
}
