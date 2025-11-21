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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		var v int64
		fmt.Fscan(in, &n, &m, &v)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + a[i-1]
		}

		prefSeg := make([]int, n+1)
		acc := int64(0)
		cnt := 0
		for i := 1; i <= n; i++ {
			acc += a[i-1]
			if acc >= v {
				cnt++
				acc = 0
			}
			prefSeg[i] = cnt
		}

		if prefSeg[n] < m {
			fmt.Fprintln(out, -1)
			continue
		}

		suf := make([]int, n+2)
		acc = 0
		cnt = 0
		for i := n; i >= 1; i-- {
			acc += a[i-1]
			if acc >= v {
				cnt++
				acc = 0
			}
			suf[i] = cnt
		}

		maxSeg := prefSeg[n]
		firstIdx := make([]int, maxSeg+1)
		for i := range firstIdx {
			firstIdx[i] = n + 1
		}
		firstIdx[0] = 0
		ptr := 1
		for i := 1; i <= n && ptr <= maxSeg; i++ {
			for ptr <= prefSeg[i] {
				firstIdx[ptr] = i
				ptr++
				if ptr > maxSeg {
					break
				}
			}
		}

		best := int64(0)
		for r := 0; r <= n; r++ {
			need := m - suf[r+1]
			if need < 0 {
				need = 0
			}
			if need > maxSeg {
				continue
			}
			idx := firstIdx[need]
			if idx == n+1 || idx > r {
				continue
			}
			val := prefix[r] - prefix[idx]
			if val > best {
				best = val
			}
		}

		fmt.Fprintln(out, best)
	}
}
