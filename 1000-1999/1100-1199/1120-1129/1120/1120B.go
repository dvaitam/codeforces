package main

import (
	"bufio"
	"fmt"
	"os"
)

type operation struct {
	pos int
	val int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var sa, sb string
	fmt.Fscan(in, &sa)
	fmt.Fscan(in, &sb)

	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int(sa[i] - '0')
		b[i] = int(sb[i] - '0')
	}

	ops := make([]operation, 0)
	total := 0
	possible := true

	for i := 0; i < n-1 && possible; i++ {
		diff := b[i] - a[i]
		if diff == 0 {
			continue
		}
		if diff > 0 {
			if a[i+1]+diff > 9 {
				possible = false
				break
			}
			a[i] += diff
			a[i+1] += diff
			total += diff
			cap := 100000 - len(ops)
			if cap > 0 {
				limit := diff
				if limit > cap {
					limit = cap
				}
				for j := 0; j < limit; j++ {
					ops = append(ops, operation{i + 1, 1})
				}
			}
			// If diff > remaining capacity, we skip storing but count remains.
		} else { // diff < 0
			if a[i+1]+diff < 0 {
				possible = false
				break
			}
			a[i] += diff
			a[i+1] += diff
			total -= diff
			need := -diff
			cap := 100000 - len(ops)
			if cap > 0 {
				limit := need
				if limit > cap {
					limit = cap
				}
				for j := 0; j < limit; j++ {
					ops = append(ops, operation{i + 1, -1})
				}
			}
		}
	}

	if !possible || a[n-1] != b[n-1] {
		fmt.Fprintln(out, -1)
		return
	}

	fmt.Fprintln(out, total)
	for _, op := range ops {
		fmt.Fprintf(out, "%d %d\n", op.pos, op.val)
	}
}
