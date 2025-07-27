package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(out *bufio.Writer, in *bufio.Reader, idx []int) uint64 {
	fmt.Fprintf(out, "? %d", len(idx))
	for _, v := range idx {
		fmt.Fprintf(out, " %d", v)
	}
	fmt.Fprintln(out)
	out.Flush()
	var res uint64
	fmt.Fscan(in, &res)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	bits := 0
	for (1<<bits) < n && bits < 13 {
		bits++
	}
	if bits == 0 {
		bits = 1
	}

	results := make([]uint64, bits)
	for b := 0; b < bits; b++ {
		subset := make([]int, 0)
		for i := 1; i <= n; i++ {
			if (i>>b)&1 == 1 {
				subset = append(subset, i)
			}
		}
		if len(subset) > 0 {
			results[b] = query(out, in, subset)
		}
	}

	ans := make([]uint64, n)
	for i := 1; i <= n; i++ {
		var val uint64
		for b := 0; b < bits; b++ {
			if (i>>b)&1 == 0 {
				val |= results[b]
			}
		}
		ans[i-1] = val
	}

	fmt.Fprint(out, "!")
	for _, v := range ans {
		fmt.Fprintf(out, " %d", v)
	}
	fmt.Fprintln(out)
	out.Flush()
}
