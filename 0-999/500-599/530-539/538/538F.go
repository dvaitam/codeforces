package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	if n == 1 {
		return
	}

	maxK := n - 1
	diff := make([]int, n+3)

	for v := 2; v <= n; v++ {
		t := v - 2
		limit := maxK
		if t < limit {
			limit = t
		}

		k := 1
		for k <= limit {
			q := t / k
			if q == 0 {
				break
			}
			kMax := t / q
			if kMax > limit {
				kMax = limit
			}
			u := q + 1
			if a[v] < a[u] {
				diff[k]++
				diff[kMax+1]--
			}
			k = kMax + 1
		}

		if t < maxK && a[v] < a[1] {
			L := t + 1
			if L < 1 {
				L = 1
			}
			diff[L]++
			diff[maxK+1]--
		}
	}

	out := bufio.NewWriter(os.Stdout)
	curr := 0
	for k := 1; k <= maxK; k++ {
		curr += diff[k]
		if k > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, curr)
	}
	fmt.Fprintln(out)
	out.Flush()
}
