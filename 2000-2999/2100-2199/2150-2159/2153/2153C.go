package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type fenwick struct {
	n    int
	tree []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, tree: make([]int, n+2)}
}

func (f *fenwick) add(idx, delta int) {
	for idx <= f.n {
		f.tree[idx] += delta
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func (f *fenwick) findByOrder(k int) int {
	if k <= 0 {
		return 0
	}
	idx := 0
	bit := 1
	for bit<<1 <= f.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= f.n && f.tree[next] < k {
			k -= f.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}

func uniqueSorted(arr []int) []int {
	if len(arr) == 0 {
		return nil
	}
	res := []int{arr[0]}
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[i-1] {
			res = append(res, arr[i])
		}
	}
	return res
}

func solveCase(a []int) int64 {
	n := len(a)
	sort.Ints(a)
	values := uniqueSorted(a)
	idx := make(map[int]int, len(values))
	for i, v := range values {
		idx[v] = i + 1
	}
	ft := newFenwick(len(values))
	parity := make([]bool, len(values))
	var prefSum int64
	var oddTotal int64
	oddCount := 0
	best := int64(0)

	for i, val := range a {
		prefSum += int64(val)
		pos := idx[val]
		if !parity[pos-1] {
			parity[pos-1] = true
			oddCount++
			oddTotal += int64(val)
			ft.add(pos, 1)
		} else {
			parity[pos-1] = false
			oddCount--
			oddTotal -= int64(val)
			ft.add(pos, -1)
		}

		prefixLen := i + 1
		if prefixLen < 3 {
			continue
		}
		maxVal := int64(val)
		if prefSum-maxVal <= maxVal {
			continue
		}

		dropCount := 0
		var dropSum int64
		if oddCount > 2 {
			dropCount = oddCount - 2
			k := oddCount
			idxLargest := ft.findByOrder(k)
			idxSecond := ft.findByOrder(k - 1)
			sumLargestTwo := int64(values[idxLargest-1] + values[idxSecond-1])
			dropSum = oddTotal - sumLargestTwo
		}
		finalCount := prefixLen - dropCount
		if finalCount < 3 {
			continue
		}
		total := prefSum - dropSum
		if total-maxVal <= maxVal {
			continue
		}
		if total > best {
			best = total
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n < 3 {
			fmt.Fprintln(out, 0)
			continue
		}
		fmt.Fprintln(out, solveCase(a))
	}
}
