package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pst struct {
	left  []int
	right []int
	sum   []int
	roots []int
	m     int
}

func newPST(n, m int) *pst {
	// estimate nodes: (n+1)*log2(m)+1
	est := (n + 1) * 20
	return &pst{
		left:  make([]int, 1, est),
		right: make([]int, 1, est),
		sum:   make([]int, 1, est),
		roots: make([]int, n+1),
		m:     m,
	}
}

func (p *pst) newNode(from int) int {
	idx := len(p.sum)
	p.left = append(p.left, p.left[from])
	p.right = append(p.right, p.right[from])
	p.sum = append(p.sum, p.sum[from])
	return idx
}

func (p *pst) update(prevRoot, pos, l, r int) int {
	cur := p.newNode(prevRoot)
	p.sum[cur]++
	if l != r {
		mid := (l + r) >> 1
		if pos <= mid {
			p.left[cur] = p.update(p.left[cur], pos, l, mid)
		} else {
			p.right[cur] = p.update(p.right[cur], pos, mid+1, r)
		}
	}
	return cur
}

func (p *pst) query(rootL, rootR, ql, qr, l, r int) int {
	if ql > r || qr < l || rootR == rootL {
		return 0
	}
	if ql <= l && r <= qr {
		return p.sum[rootR] - p.sum[rootL]
	}
	mid := (l + r) >> 1
	return p.query(p.left[rootL], p.left[rootR], ql, qr, l, mid) +
		p.query(p.right[rootL], p.right[rootR], ql, qr, mid+1, r)
}

func upperBound(arr []int, x int) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i] > x })
}

func lowerBound(arr []int, x int) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i] >= x })
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			vals[i] = a[i]
		}

		sort.Ints(vals)
		vals = unique(vals)
		m := len(vals)

		// compress values
		comp := make([]int, n)
		for i, v := range a {
			comp[i] = lowerBound(vals, v)
		}

		// suffix maximum value and its rightmost position
		sufMaxVal := make([]int, n+1)
		sufMaxPos := make([]int, n+1)
		sufMaxVal[n] = -1 // sentinel
		maxVal := -1
		maxPos := n
		for i := n - 1; i >= 0; i-- {
			if a[i] > maxVal {
				maxVal = a[i]
				maxPos = i
			}
			// keep rightmost occurrence for equal values: do nothing when a[i]==maxVal
			sufMaxVal[i] = maxVal
			sufMaxPos[i] = maxPos
		}

		pst := newPST(n, m)
		for i := 0; i < n; i++ {
			pst.roots[i+1] = pst.update(pst.roots[i], comp[i], 0, m-1)
		}

		ans := 1
		for i := 0; i < n; i++ {
			if sufMaxVal[i+1] > a[i] { // need larger element after i
				r := sufMaxPos[i+1]
				minVal := a[i]
				maxVal := sufMaxVal[i+1]
				lIdx := upperBound(vals, minVal)
				rIdx := lowerBound(vals, maxVal) - 1
				midCnt := 0
				if lIdx <= rIdx && r > i+1 {
					midCnt = pst.query(pst.roots[i+1], pst.roots[r], lIdx, rIdx, 0, m-1)
				}
				cur := midCnt + 2
				if cur > ans {
					ans = cur
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}

func unique(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	j := 1
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[i-1] {
			arr[j] = arr[i]
			j++
		}
	}
	return arr[:j]
}
