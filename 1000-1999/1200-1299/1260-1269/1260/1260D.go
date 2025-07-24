package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Trap struct {
	l, r, d int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m, n, k, t int
	if _, err := fmt.Fscan(in, &m, &n, &k, &t); err != nil {
		return
	}
	agility := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &agility[i])
	}
	traps := make([]Trap, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &traps[i].l, &traps[i].r, &traps[i].d)
	}

	sort.Slice(agility, func(i, j int) bool { return agility[i] > agility[j] })
	sort.Slice(traps, func(i, j int) bool {
		if traps[i].l == traps[j].l {
			return traps[i].r < traps[j].r
		}
		return traps[i].l < traps[j].l
	})

	can := func(threshold int) bool {
		extra := 0
		currL, currR := -1, -1
		for _, tr := range traps {
			if tr.d <= threshold {
				continue
			}
			if currL == -1 {
				currL, currR = tr.l, tr.r
			} else if tr.l <= currR {
				if tr.r > currR {
					currR = tr.r
				}
			} else {
				extra += 2 * (currR - currL + 1)
				currL, currR = tr.l, tr.r
			}
		}
		if currL != -1 {
			extra += 2 * (currR - currL + 1)
		}
		return n+1+extra <= t
	}

	lo, hi := 0, m
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if can(agility[mid-1]) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}

	fmt.Fprintln(out, lo)
}
