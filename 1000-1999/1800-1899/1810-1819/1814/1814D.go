package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func feasible(L, k int64, fi []int64, maxF int64) bool {
	if L+k < maxF {
		return false
	}
	for _, f := range fi {
		if f <= k {
			continue
		}
		r := L % f
		if r != 0 && r < f-k {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int64
		fmt.Fscan(in, &n, &k)
		f := make([]int64, n)
		for i := range f {
			fmt.Fscan(in, &f[i])
		}
		d := make([]int64, n)
		p := make([]int64, n)
		maxF := int64(0)
		for i := range d {
			fmt.Fscan(in, &d[i])
			p[i] = f[i] * d[i]
			if f[i] > maxF {
				maxF = f[i]
			}
		}
		// sort firepower values
		sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })

		// collect unique fire rates
		uniqFMap := make(map[int64]struct{})
		uniqF := make([]int64, 0)
		for _, x := range f {
			if _, ok := uniqFMap[x]; !ok {
				uniqFMap[x] = struct{}{}
				uniqF = append(uniqF, x)
			}
		}

		lowerBound := maxF - k
		if lowerBound < 0 {
			lowerBound = 0
		}

		// candidate L values
		cands := make([]int64, 0, 2*n+1)
		cands = append(cands, lowerBound)
		for _, val := range p {
			cands = append(cands, val)
			cands = append(cands, val-k)
		}
		sort.Slice(cands, func(i, j int) bool { return cands[i] < cands[j] })
		// unique
		uniq := cands[:0]
		var prev int64 = -1 << 63
		for _, v := range cands {
			if v != prev {
				uniq = append(uniq, v)
				prev = v
			}
		}
		cands = uniq

		best := 0
		l, r := 0, -1
		nInt := int(n)
		for _, L := range cands {
			if L < lowerBound {
				continue
			}
			for r+1 < nInt && p[r+1] <= L+k {
				r++
			}
			for l <= r && p[l] < L {
				l++
			}
			keep := r - l + 1
			if keep <= best {
				if keep < 0 {
					keep = 0
				}
				// even if keep <= best, we might continue to check feasibility to skip loops? but we skip to save time
			}
			if feasible(L, k, uniqF, maxF) {
				if keep > best {
					best = keep
				}
			}
		}
		ans := int(n) - best
		fmt.Fprintln(out, ans)
	}
}
