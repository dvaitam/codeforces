package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		// Prepare sorted copies for counting >= price.
		as := append([]int64(nil), a...)
		bs := append([]int64(nil), b...)
		sort.Slice(as, func(i, j int) bool { return as[i] < as[j] })
		sort.Slice(bs, func(i, j int) bool { return bs[i] < bs[j] })

		// Collect candidate prices (all a_i and b_i).
		cand := make([]int64, 0, 2*n)
		cand = append(cand, a...)
		cand = append(cand, b...)
		sort.Slice(cand, func(i, j int) bool { return cand[i] < cand[j] })
		uniq := cand[:0]
		for _, v := range cand {
			if len(uniq) == 0 || uniq[len(uniq)-1] != v {
				uniq = append(uniq, v)
			}
		}

		// Iterate candidates from largest to smallest, maintaining counts of values >= price.
		idxA, idxB := len(as)-1, len(bs)-1
		var cntA, cntB int // number of a_i / b_i that are >= current price
		var best int64

		for i := len(uniq) - 1; i >= 0; i-- {
			p := uniq[i]
			for idxA >= 0 && as[idxA] >= p {
				cntA++
				idxA--
			}
			for idxB >= 0 && bs[idxB] >= p {
				cntB++
				idxB--
			}

			total := cntB
			positive := cntA
			negative := total - positive
			if negative <= k {
				revenue := p * int64(total)
				if revenue > best {
					best = revenue
				}
			}
		}

		fmt.Fprintln(out, best)
	}
}
