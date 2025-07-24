package main

import (
	"bufio"
	"fmt"
	"os"
)

func edgesCounts(n int) []int64 {
	e := make([]int64, n+1)
	for i := n; i >= 1; i-- {
		c := n / i
		e[i] = int64(c) * int64(c-1) / 2
		for j := i * 2; j <= n; j += i {
			e[i] -= e[j]
		}
	}
	return e
}

func minCost(n int, m int64) int64 {
	edges := edgesCounts(n)
	var cost int64
	for d := n; d >= 2 && m > 0; d-- {
		avail := edges[d] / int64(d-1)
		if avail == 0 {
			continue
		}
		need := m / int64(d-1)
		if need > avail {
			need = avail
		}
		cost += need * int64(d)
		m -= need * int64(d-1)
	}
	if m > 0 {
		return -1
	}
	return cost
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(in, &n, &m)
		res := minCost(n, m)
		fmt.Fprintln(out, res)
	}
}
