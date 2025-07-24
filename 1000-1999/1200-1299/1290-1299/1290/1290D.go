package main

import (
	"bufio"
	"fmt"
	"os"
)

// Interactive solution for Codeforces problem 1290D - Coffee Varieties (hard version).
// The approach divides indices into blocks of size k. We first check duplicates
// inside each block, then for each pair of blocks. At most k queries are issued
// without a reset, so the memory constraint is respected. Total queries are n^2/k,
// which is within the allowed 3n^2/(2k) bound.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	distinct := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		distinct[i] = true
	}

	blocks := n / k
	// check duplicates within each block
	for b := 0; b < blocks; b++ {
		fmt.Fprintln(out, "R")
		out.Flush()
		for i := b*k + 1; i <= (b+1)*k; i++ {
			fmt.Fprintf(out, "? %d\n", i)
			out.Flush()
			var resp string
			if _, err := fmt.Fscan(in, &resp); err != nil {
				return
			}
			if resp == "Y" {
				distinct[i] = false
			}
		}
	}

	// check duplicates across pairs of blocks
	for b1 := 0; b1 < blocks; b1++ {
		for b2 := b1 + 1; b2 < blocks; b2++ {
			fmt.Fprintln(out, "R")
			out.Flush()
			// query entire first block
			for i := b1*k + 1; i <= (b1+1)*k; i++ {
				fmt.Fprintf(out, "? %d\n", i)
				out.Flush()
				var resp string
				if _, err := fmt.Fscan(in, &resp); err != nil {
					return
				}
			}
			// query second block, mark duplicates
			for j := b2*k + 1; j <= (b2+1)*k; j++ {
				fmt.Fprintf(out, "? %d\n", j)
				out.Flush()
				var resp string
				if _, err := fmt.Fscan(in, &resp); err != nil {
					return
				}
				if resp == "Y" {
					distinct[j] = false
				}
			}
		}
	}

	count := 0
	for i := 1; i <= n; i++ {
		if distinct[i] {
			count++
		}
	}

	fmt.Fprintf(out, "! %d\n", count)
	out.Flush()
}
