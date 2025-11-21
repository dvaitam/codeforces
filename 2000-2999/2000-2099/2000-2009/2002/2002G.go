package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type suffixDS struct {
	bitsets [][]uint64
	nz      [][]int
	counts  []int
	total   int
	words   int
	maxVal  int
	tmpIdx  []int
	tmpVals []uint64
}

func newSuffixDS(masks []uint64, maxVal int) *suffixDS {
	total := len(masks)
	words := (total + 63) >> 6
	bitsets := make([][]uint64, maxVal+1)
	for i := 0; i <= maxVal; i++ {
		bitsets[i] = make([]uint64, words)
	}
	counts := make([]int, maxVal+1)
	for idx, mask := range masks {
		word := idx >> 6
		bit := uint(idx & 63)
		flag := uint64(1) << bit
		mm := mask
		for mm != 0 {
			v := bits.TrailingZeros64(mm)
			bitsets[v][word] |= flag
			counts[v]++
			mm &= mm - 1
		}
	}
	nz := make([][]int, maxVal+1)
	for v := 0; v <= maxVal; v++ {
		arr := bitsets[v]
		if len(arr) == 0 {
			continue
		}
		list := make([]int, 0)
		for w, val := range arr {
			if val != 0 {
				list = append(list, w)
			}
		}
		nz[v] = list
	}
	return &suffixDS{
		bitsets: bitsets,
		nz:      nz,
		counts:  counts,
		total:   total,
		words:   words,
		maxVal:  maxVal,
		tmpIdx:  make([]int, 0, words),
		tmpVals: make([]uint64, 0, words),
	}
}

func (ds *suffixDS) maxMex(pmask uint64) int {
	candIdx := ds.tmpIdx[:0]
	candVals := ds.tmpVals[:0]
	hasCand := false
	for val := 0; val <= ds.maxVal; val++ {
		if (pmask>>uint(val))&1 != 0 {
			continue
		}
		c := ds.counts[val]
		if c == 0 {
			ds.tmpIdx = candIdx[:0]
			ds.tmpVals = candVals[:0]
			return val
		}
		if c == ds.total {
			continue
		}
		if !hasCand {
			nz := ds.nz[val]
			candIdx = candIdx[:0]
			candVals = candVals[:0]
			for _, w := range nz {
				candIdx = append(candIdx, w)
				candVals = append(candVals, ds.bitsets[val][w])
			}
			if len(candIdx) == 0 {
				ds.tmpIdx = candIdx[:0]
				ds.tmpVals = candVals[:0]
				return val
			}
			hasCand = true
		} else {
			bitsArr := ds.bitsets[val]
			for i := 0; i < len(candIdx); {
				w := candIdx[i]
				newVal := candVals[i] & bitsArr[w]
				if newVal == 0 {
					last := len(candIdx) - 1
					candIdx[i] = candIdx[last]
					candVals[i] = candVals[last]
					candIdx = candIdx[:last]
					candVals = candVals[:last]
				} else {
					candVals[i] = newVal
					i++
				}
			}
			if len(candIdx) == 0 {
				ds.tmpIdx = candIdx[:0]
				ds.tmpVals = candVals[:0]
				return val
			}
		}
	}
	ds.tmpIdx = candIdx[:0]
	ds.tmpVals = candVals[:0]
	return ds.maxVal + 1
}

func combRow(m int) []int {
	res := make([]int, m+1)
	res[0] = 1
	for k := 1; k <= m; k++ {
		res[k] = res[k-1] * (m - (k - 1)) / k
	}
	return res
}

func solveCase(n int, d [][]int, r [][]int) int {
	diag := n
	pref := make([][]uint64, diag)
	suf := make([][]uint64, diag)
	comb := combRow(n - 1)
	for i := 0; i < diag; i++ {
		cap := comb[i]
		pref[i] = make([]uint64, 0, cap)
		suf[i] = make([]uint64, 0, cap)
	}
	var dfsPref func(x, y, steps int, mask uint64)
	dfsPref = func(x, y, steps int, mask uint64) {
		if steps == n-1 {
			pref[x] = append(pref[x], mask)
			return
		}
		if x+1 < n {
			val := d[x][y]
			dfsPref(x+1, y, steps+1, mask|(uint64(1)<<uint(val)))
		}
		if y+1 < n {
			val := r[x][y]
			dfsPref(x, y+1, steps+1, mask|(uint64(1)<<uint(val)))
		}
	}
	dfsPref(0, 0, 0, 0)

	var dfsSuf func(x, y, steps int, mask uint64)
	dfsSuf = func(x, y, steps int, mask uint64) {
		if steps == n-1 {
			suf[x] = append(suf[x], mask)
			return
		}
		if x > 0 {
			val := d[x-1][y]
			dfsSuf(x-1, y, steps+1, mask|(uint64(1)<<uint(val)))
		}
		if y > 0 {
			val := r[x][y-1]
			dfsSuf(x, y-1, steps+1, mask|(uint64(1)<<uint(val)))
		}
	}
	dfsSuf(n-1, n-1, 0, 0)

	maxVal := 2*n - 2
	best := 0
	for idx := 0; idx < n; idx++ {
		if len(suf[idx]) == 0 {
			continue
		}
		ds := newSuffixDS(suf[idx], maxVal)
		for _, pmask := range pref[idx] {
			mex := ds.maxMex(pmask)
			if mex > best {
				best = mex
			}
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		d := make([][]int, n-1)
		for i := 0; i < n-1; i++ {
			d[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &d[i][j])
			}
		}
		r := make([][]int, n)
		for i := 0; i < n; i++ {
			r[i] = make([]int, n-1)
			for j := 0; j < n-1; j++ {
				fmt.Fscan(in, &r[i][j])
			}
		}
		ans := solveCase(n, d, r)
		fmt.Fprintln(out, ans)
	}
}
