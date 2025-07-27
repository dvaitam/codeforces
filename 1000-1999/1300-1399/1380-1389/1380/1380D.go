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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var x, k, y int64
	fmt.Fscan(in, &x, &k, &y)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
	}

	ai, bi := 0, 0
	last := -1
	var total int64

	for bi < m {
		// find b[bi] in a starting from ai
		start := ai
		for ai < n && a[ai] != b[bi] {
			ai++
		}
		if ai == n {
			fmt.Fprintln(out, -1)
			return
		}
		// process segment [start, ai)
		cost, ok := processSegment(a[start:ai], getBound(a, last), a[ai], x, k, y)
		if !ok {
			fmt.Fprintln(out, -1)
			return
		}
		total += cost
		last = ai
		ai++
		bi++
	}
	// process suffix
	cost, ok := processSegment(a[ai:], getBound(a, last), 0, x, k, y)
	if !ok {
		fmt.Fprintln(out, -1)
		return
	}
	total += cost
	fmt.Fprintln(out, total)
}

func getBound(a []int64, idx int) int64 {
	if idx >= 0 && idx < len(a) {
		return a[idx]
	}
	return 0
}

func processSegment(seg []int64, left, right int64, x, k, y int64) (int64, bool) {
	l := len(seg)
	if l == 0 {
		return 0, true
	}
	maxV := left
	if right > maxV {
		maxV = right
	}
	maxInside := seg[0]
	for _, v := range seg {
		if v > maxInside {
			maxInside = v
		}
	}
	needFire := maxInside > maxV
	if needFire && int64(l) < k {
		return 0, false
	}
	var res int64
	if needFire {
		// remove one fireball covering the strongest
		res += x
		l -= int(k)
	}
	if k*y < x {
		// Berserk is cheaper than fireball for remaining
		res += int64(l) * y
	} else {
		res += (int64(l)/k)*x + (int64(l)%k)*y
	}
	// also compare pure berserk without using extra fireballs
	if !needFire {
		berserkOnly := int64(len(seg)) * y
		if res > berserkOnly {
			res = berserkOnly
		}
	}
	return res, true
}
