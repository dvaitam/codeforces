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
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	a := make([]int64, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int64, n)
	for i := range b {
		fmt.Fscan(in, &b[i])
	}

	limit := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		if b[i] == 0 {
			limit[i] = k // theoretically infinite, but k suffices
		} else {
			limit[i] = a[i] / b[i]
		}
		total += limit[i]
	}

	result := make([]int64, n)

	if total < k {
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, 0)
		}
		fmt.Fprintln(out)
		return
	}

	suffix := make([]int64, n+2)
	for i := n - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] + limit[i]
	}

	required := k
	for i := 0; i < n; i++ {
		remainingCap := suffix[i+1]
		need := required - remainingCap
		if need < 0 {
			need = 0
		}
		if need > limit[i] {
			need = limit[i]
		}
		result[i] = need
		required -= need
	}

	for i, v := range result {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
