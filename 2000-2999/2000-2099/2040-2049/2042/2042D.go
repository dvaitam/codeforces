package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type bitMax struct {
	n    int
	data []int
}

func newBitMax(n int) *bitMax {
	b := &bitMax{n: n, data: make([]int, n+1)}
	for i := range b.data {
		b.data[i] = -1
	}
	return b
}

func (b *bitMax) update(pos, val int) {
	for pos <= b.n {
		if val > b.data[pos] {
			b.data[pos] = val
		}
		pos += pos & -pos
	}
}

func (b *bitMax) query(pos int) int {
	res := -1
	for pos > 0 {
		if b.data[pos] > res {
			res = b.data[pos]
		}
		pos -= pos & -pos
	}
	return res
}

type bitMin struct {
	n    int
	data []int
}

func newBitMin(n int) *bitMin {
	b := &bitMin{n: n, data: make([]int, n+1)}
	const inf = int(1e9 + 7)
	for i := range b.data {
		b.data[i] = inf
	}
	return b
}

func (b *bitMin) update(pos, val int) {
	for pos <= b.n {
		if val < b.data[pos] {
			b.data[pos] = val
		}
		pos += pos & -pos
	}
}

func (b *bitMin) query(pos int) int {
	const inf = int(1e9 + 7)
	res := inf
	for pos > 0 {
		if b.data[pos] < res {
			res = b.data[pos]
		}
		pos -= pos & -pos
	}
	return res
}

type interval struct {
	l, r, idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		l := make([]int, n)
		r := make([]int, n)
		allR := make([]int, 0, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
			allR = append(allR, r[i])
		}

		sort.Ints(allR)
		allR = unique(allR)

		getIdx := func(x int) int {
			// 1-based index
			pos := sort.SearchInts(allR, x)
			return pos + 1
		}

		m := len(allR)
		revIdx := func(x int) int {
			return m - getIdx(x) + 1
		}

		ints := make([]interval, n)
		for i := 0; i < n; i++ {
			ints[i] = interval{l: l[i], r: r[i], idx: i}
		}

	sort.Slice(ints, func(i, j int) bool {
		return ints[i].l < ints[j].l
	})

	const infVal = int(1e9 + 7)
	bitMaxL := newBitMax(m)
	bitMinR := newBitMin(m)

	Lmax := make([]int, n)
	Rmin := make([]int, n)

		for i := 0; i < n; {
			j := i
			curL := ints[i].l
			for j < n && ints[j].l == curL {
				j++
			}

			// group statistics for r values
			cnt := make(map[int]int)
			rVals := make([]int, 0, j-i)
			for k := i; k < j; k++ {
				cnt[ints[k].r]++
			}
			for v := range cnt {
				rVals = append(rVals, v)
			}
			sort.Ints(rVals)

			maxR := rVals[len(rVals)-1]
			countMax := cnt[maxR]
			posMap := make(map[int]int, len(rVals))
			for idx, v := range rVals {
				posMap[v] = idx
			}

			// answer queries for this group using only earlier intervals + other same-l intervals
			for k := i; k < j; k++ {
				idxR := revIdx(ints[k].r)
				oldL := bitMaxL.query(idxR)
				oldR := bitMinR.query(idxR)

				groupL := -1
				if j-i >= 2 {
					if ints[k].r < maxR || (ints[k].r == maxR && countMax > 1) {
						groupL = curL
					}
				}

				groupR := infVal
				if cnt[ints[k].r] >= 2 {
					groupR = ints[k].r
				} else {
					pos := posMap[ints[k].r]
					if pos+1 < len(rVals) {
						groupR = rVals[pos+1]
					}
				}

				if oldL > groupL {
					groupL = oldL
				}
				if oldR < groupR {
					groupR = oldR
				}

				Lmax[ints[k].idx] = groupL
				Rmin[ints[k].idx] = groupR
			}

			// insert all intervals with this l
			for k := i; k < j; k++ {
				idxR := revIdx(ints[k].r)
				bitMaxL.update(idxR, ints[k].l)
				bitMinR.update(idxR, ints[k].r)
			}
			i = j
		}

		for i := 0; i < n; i++ {
			if Lmax[i] == -1 || Rmin[i] == infVal {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, 0)
				continue
			}
			ans := int64(l[i]-Lmax[i]) + int64(Rmin[i]-r[i])
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	k := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[k] = a[i]
			k++
		}
	}
	return a[:k]
}
