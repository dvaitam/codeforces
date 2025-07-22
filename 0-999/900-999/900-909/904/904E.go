package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func apply(mask int, n int, full uint32, adjOrig []uint32) bool {
	var adj [22]uint32
	copy(adj[:], adjOrig)
	for v := 0; v < n; v++ {
		if mask&(1<<uint(v)) != 0 {
			neigh := adj[v]
			for sub := neigh; sub != 0; sub &= sub - 1 {
				i := bits.TrailingZeros32(sub)
				adj[i] |= neigh
			}
		}
	}
	for i := 0; i < n; i++ {
		if adj[i]|(1<<uint(i)) != full {
			return false
		}
	}
	return true
}

func nextComb(x int) int {
	u := x & -x
	v := u + x
	return v + (((v ^ x) / u) >> 2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	adjOrig := make([]uint32, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adjOrig[u] |= 1 << uint(v)
		adjOrig[v] |= 1 << uint(u)
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	full := uint32(1<<uint(n)) - 1

	for k := 0; k <= n; k++ {
		if k == 0 {
			if apply(0, n, full, adjOrig) {
				fmt.Fprintln(out, 0)
				fmt.Fprintln(out)
				return
			}
			continue
		}
		limit := 1 << uint(n)
		for mask := (1 << uint(k)) - 1; mask < limit; mask = nextComb(mask) {
			if bits.OnesCount(uint(mask)) != k {
				break
			}
			if apply(mask, n, full, adjOrig) {
				fmt.Fprintln(out, k)
				first := true
				for i := 0; i < n; i++ {
					if mask&(1<<uint(i)) != 0 {
						if !first {
							fmt.Fprint(out, " ")
						}
						fmt.Fprint(out, i+1)
						first = false
					}
				}
				fmt.Fprintln(out)
				return
			}
		}
	}
}
