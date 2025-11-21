package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Operation struct {
	i, j int
	x    int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := append([]int64(nil), a...)
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

		isSorted := true
		for i := 0; i+1 < n; i++ {
			if a[i] > a[i+1] {
				isSorted = false
				break
			}
		}
		if isSorted {
			fmt.Fprintln(out, 0)
			continue
		}

		u := -1
		v := -1
		for i := 0; i < n && u == -1; i++ {
			for j := i + 1; j < n; j++ {
				if a[i]+a[j] == k {
					u = i
					v = j
					break
				}
			}
		}
		if u == -1 {
			fmt.Fprintln(out, -1)
			continue
		}

		ops := make([]Operation, 0, 2*n)
		order := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if i != u {
				order = append(order, i)
			}
		}
		order = append(order, u)

		for _, idx := range order {
			if idx == u {
				// should not happen for non-last iteration
				if len(order) == 1 && idx == u {
					// single element case
				}
			}
			cur := a[idx]
			target := b[idx]

			// Step1: adjust pair (u,v) so that a[u] = k - cur, a[v] = cur
			if a[u] != k-cur {
				x := a[u] - (k - cur)
				ops = append(ops, Operation{i: u + 1, j: v + 1, x: x})
				a[u] -= x
				a[v] += x
			}
			// Step2: adjust idx using pair (idx,u)
			x2 := cur - target
			ops = append(ops, Operation{i: idx + 1, j: u + 1, x: x2})
			a[idx] -= x2
			a[u] += x2

			// Update pair: new pair is (idx, previous u)
			if idx != u {
				v = u
				u = idx
			}
		}

		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			fmt.Fprintf(out, "%d %d %d\n", op.i, op.j, op.x)
		}
	}
}
