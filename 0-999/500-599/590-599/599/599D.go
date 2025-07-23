package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	n int64
	m int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var x int64
	if _, err := fmt.Fscan(in, &x); err != nil {
		return
	}

	var res []pair
	for n := int64(1); ; n++ {
		minSq := n * (n + 1) * (2*n + 1) / 6
		if minSq > x {
			break
		}
		denom := n * (n + 1)
		sixx := 6 * x
		if sixx%denom != 0 {
			continue
		}
		t := sixx / denom
		if (t+n-1)%3 != 0 {
			continue
		}
		m := (t + n - 1) / 3
		if m < n {
			continue
		}
		if n*(n+1)*(3*m-n+1)/6 == x {
			res = append(res, pair{n, m})
			if m != n {
				res = append(res, pair{m, n})
			}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].n == res[j].n {
			return res[i].m < res[j].m
		}
		return res[i].n < res[j].n
	})

	fmt.Fprintln(out, len(res))
	for _, p := range res {
		fmt.Fprintf(out, "%d %d\n", p.n, p.m)
	}
}
