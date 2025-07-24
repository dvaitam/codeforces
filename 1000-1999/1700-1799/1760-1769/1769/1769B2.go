package main

import (
	"bufio"
	"fmt"
	"os"
)

func ceilDiv(a, b int64) int64 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b // not used for negative in this problem
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		total += a[i]
	}

	found := make([]bool, 101)
	var prefix int64
	for i := 0; i < n; i++ {
		ai := a[i]
		for p := 0; p <= 100; p++ {
			// range of x giving first progress = p
			l1 := ceilDiv(int64(p)*ai, 100)
			r1 := (int64(p+1)*ai - 1) / 100
			if l1 > r1 {
				continue
			}
			// range of x giving overall progress = p
			l2 := ceilDiv(int64(p)*total, 100) - prefix
			r2 := (int64(p+1)*total-1)/100 - prefix
			if l2 > r2 {
				continue
			}
			// intersect ranges with [0, ai]
			l := l1
			if l2 > l {
				l = l2
			}
			if l < 0 {
				l = 0
			}
			r := r1
			if r2 < r {
				r = r2
			}
			if r > ai {
				r = ai
			}
			if l <= r {
				found[p] = true
			}
		}
		prefix += ai
	}

	first := true
	for p := 0; p <= 100; p++ {
		if found[p] {
			if !first {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, p)
			first = false
		}
	}
	fmt.Fprintln(out)
}
