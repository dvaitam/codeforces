package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func digits(x int) int {
	cnt := 0
	for x > 0 {
		cnt++
		x /= 10
	}
	if cnt == 0 {
		return 1
	}
	return cnt
}

func trailingZeros(x int) int {
	cnt := 0
	for x%10 == 0 {
		cnt++
		x /= 10
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		lens := make([]int, n)
		tz := make([]int, n)
		base := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			d := digits(x)
			z := trailingZeros(x)
			lens[i] = d
			tz[i] = z
			base += d - z
		}
		sort.Slice(tz, func(i, j int) bool { return tz[i] > tz[j] })
		leftover := 0
		for i := 1; i < n; i += 2 {
			leftover += tz[i]
		}
		finalLen := base + leftover
		if finalLen >= m+1 {
			fmt.Fprintln(out, "Sasha")
		} else {
			fmt.Fprintln(out, "Anna")
		}
	}
}
