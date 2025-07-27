package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		adj[i] = make([]bool, n)
		for j, c := range []byte(s) {
			if c == '1' {
				adj[i][j] = true
			}
		}
	}
	total := 1 << n
	dp := make([][][]uint64, total)
	for mask := 1; mask < total; mask++ {
		k := bits.OnesCount(uint(mask))
		dp[mask] = make([][]uint64, n)
		if k > 0 {
			size := 1 << (k - 1)
			for v := 0; v < n; v++ {
				if mask&(1<<v) != 0 {
					dp[mask][v] = make([]uint64, size)
				}
			}
		}
	}
	for v := 0; v < n; v++ {
		dp[1<<v][v][0] = 1
	}
	for mask := 1; mask < total; mask++ {
		k := bits.OnesCount(uint(mask))
		if k == n {
			continue
		}
		for last := 0; last < n; last++ {
			if mask&(1<<last) == 0 {
				continue
			}
			arr := dp[mask][last]
			if arr == nil {
				continue
			}
			for next := 0; next < n; next++ {
				if mask&(1<<next) != 0 {
					continue
				}
				bit := 0
				if adj[last][next] {
					bit = 1
				}
				newmask := mask | (1 << next)
				newarr := dp[newmask][next]
				for p := 0; p < len(arr); p++ {
					newarr[(p<<1)|bit] += arr[p]
				}
			}
		}
	}
	full := total - 1
	res := make([]uint64, 1<<(n-1))
	for last := 0; last < n; last++ {
		arr := dp[full][last]
		if arr == nil {
			continue
		}
		for i, val := range arr {
			res[i] += val
		}
	}
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
