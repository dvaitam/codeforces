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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	blocked := make([]bool, n)
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x >= 0 && x < n {
			blocked[x] = true
		}
	}

	costs := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		var c int
		fmt.Fscan(in, &c)
		costs[i] = int64(c)
	}

	if blocked[0] {
		fmt.Fprintln(out, -1)
		return
	}

	pre := make([]int, n)
	prev := -1
	for i := 0; i < n; i++ {
		if !blocked[i] {
			prev = i
		}
		pre[i] = prev
	}

	prevAllowed := 0
	Lmin := 1
	for i := 1; i < n; i++ {
		if !blocked[i] {
			diff := i - prevAllowed
			if diff > Lmin {
				Lmin = diff
			}
			prevAllowed = i
		}
	}
	diff := n - prevAllowed
	if diff > Lmin {
		Lmin = diff
	}

	if k < Lmin {
		fmt.Fprintln(out, -1)
		return
	}

	best := int64(-1)
	for L := Lmin; L <= k; L++ {
		pos := 0
		last := -1
		cnt := 0
		for pos < n {
			start := pre[pos]
			if start <= last {
				cnt = -1
				break
			}
			last = start
			cnt++
			pos = start + L
		}
		if cnt != -1 {
			cost := int64(cnt) * costs[L]
			if best == -1 || cost < best {
				best = cost
			}
		}
	}

	fmt.Fprintln(out, best)
}
