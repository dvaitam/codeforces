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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		sums := make([]uint64, n)
		for i := 0; i < n; i++ {
			var s uint64
			for j := 1; j <= m; j++ {
				var v uint64
				fmt.Fscan(in, &v)
				s += v * uint64(j)
			}
			sums[i] = s
		}
		base := majorityValue(sums)
		var idx int
		var diff uint64
		for i, v := range sums {
			if v != base {
				idx = i + 1
				if v > base {
					diff = v - base
				} else {
					diff = base - v
				}
				break
			}
		}
		fmt.Fprintf(out, "%d %d\n", idx, diff)
	}
}

func majorityValue(arr []uint64) uint64 {
	var cand uint64
	cnt := 0
	for _, v := range arr {
		if cnt == 0 {
			cand = v
			cnt = 1
		} else if v == cand {
			cnt++
		} else {
			cnt--
		}
	}
	return cand
}
