package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(out *bufio.Writer, in *bufio.Reader, idx []int) int {
	fmt.Fprintf(out, "? %d", len(idx))
	for _, v := range idx {
		fmt.Fprintf(out, " %d", v)
	}
	fmt.Fprintln(out)
	out.Flush()
	var res int
	fmt.Fscan(in, &res)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, x, y int
	if _, err := fmt.Fscan(in, &n, &x, &y); err != nil {
		return
	}

	// First determine mask = p1 ^ p2
	mask := 0
	for b := 0; b < 10; b++ {
		var subset []int
		for i := 1; i <= n; i++ {
			if (i>>b)&1 == 1 {
				subset = append(subset, i)
			}
		}
		res := query(out, in, subset)
		m := len(subset)
		r1 := y
		if m%2 == 0 {
			r1 = y ^ x
		}
		if res == r1 {
			mask |= 1 << b
		}
	}

	// choose a bit set in mask
	bit := 0
	for b := 0; b < 10; b++ {
		if mask&(1<<b) != 0 {
			bit = b
			break
		}
	}

	// find index with this bit = 1
	candidate := make([]int, 0)
	for i := 1; i <= n; i++ {
		if (i>>bit)&1 == 1 {
			candidate = append(candidate, i)
		}
	}

	for len(candidate) > 1 {
		mid := len(candidate) / 2
		subset := candidate[:mid]
		res := query(out, in, subset)
		m := len(subset)
		r1 := y
		if m%2 == 0 {
			r1 = y ^ x
		}
		if res == r1 {
			candidate = subset
		} else {
			candidate = candidate[mid:]
		}
	}

	p1 := candidate[0]
	p2 := p1 ^ mask
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	fmt.Fprintf(out, "! %d %d\n", p1, p2)
	out.Flush()
}
