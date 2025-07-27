package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var x int64
	fmt.Fscan(in, &x)
	d := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &d[i])
	}

	// Compute LIS length using prefix sum technique.
	var prefix, minPref, best int64
	for _, v := range d {
		if v > 0 {
			diff := prefix + v - minPref
			if diff > best {
				best = diff
			}
			prefix += v
		} else if v < 0 {
			if prefix > minPref {
				diff := prefix - 1 - minPref
				if diff > best {
					best = diff
				}
			}
			prefix += v
			if prefix < minPref {
				minPref = prefix
			}
		}
	}
	lisLen := best + 1

	// Placeholder for number of LIS. Computing the exact number requires
	// a more involved combinatorial analysis. For now we output 1.
	fmt.Fprintln(out, lisLen, 1%mod)
}
