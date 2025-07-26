package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Fenwick tree for prefix sums
type fenwick struct {
	n   int
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int, n+2)}
}

func (f *fenwick) add(i, v int) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *fenwick) sum(i int) int {
	res := 0
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

func (f *fenwick) search(val int) int {
	// largest idx such that prefix sum <= val
	idx := 0
	cur := 0
	bit := 1
	for bit <= f.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= f.n && cur+f.bit[next] <= val {
			cur += f.bit[next]
			idx = next
		}
		bit >>= 1
	}
	return idx
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		freq := make(map[int]int)
		for _, x := range arr {
			freq[x]++
		}
		type pair struct{ val, cnt int }
		vals := make([]pair, 0, len(freq))
		for v, c := range freq {
			vals = append(vals, pair{v, c})
		}
		sort.Slice(vals, func(i, j int) bool { return vals[i].val < vals[j].val })

		type fpair struct{ freq, val int }
		fs := make([]fpair, len(vals))
		for i, p := range vals {
			fs[i] = fpair{p.cnt, p.val}
		}
		sort.Slice(fs, func(i, j int) bool { return fs[i].freq < fs[j].freq })

		idxMap := make(map[int]int)
		for i, p := range fs {
			idxMap[p.val] = i + 1 // 1-based
		}

		nUnique := len(vals)
		bitCnt := newFenwick(nUnique)
		bitSum := newFenwick(nUnique)
		for i, p := range fs {
			bitCnt.add(i+1, 1)
			bitSum.add(i+1, p.freq)
		}

		missing := 0
		active := nUnique
		ans := active // worst case
		for mex := 0; mex <= n; mex++ {
			if missing <= k {
				pos := bitSum.search(k)
				removed := bitCnt.sum(pos)
				cost := active - removed
				if cost < ans {
					ans = cost
				}
			}
			if c, ok := freq[mex]; ok {
				// remove this value from active set
				idx := idxMap[mex]
				bitCnt.add(idx, -1)
				bitSum.add(idx, -c)
				active--
			} else {
				missing++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
