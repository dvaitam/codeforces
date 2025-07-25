package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

// Fenwick tree for prefix sums over integers (1-based indices).
type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+2)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.tree[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	if i > f.n {
		i = f.n
	}
	if i < 0 {
		return 0
	}
	s := 0
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

func invCount(arr []int) int64 {
	m := len(arr)
	fw := NewFenwick(m)
	inv := int64(0)
	for i, v := range arr {
		pos := v + 1
		inv += int64(i - fw.Sum(pos))
		fw.Add(pos, 1)
	}
	return inv
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		q := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &q[i])
		}

		invQ := invCount(q) % mod
		rowInv := invQ * int64(n) % mod

		maxDiff := 18
		if k-1 < maxDiff {
			maxDiff = k - 1
		}
		r := make([]int, 2*maxDiff+1)
		for d := -maxDiff; d <= maxDiff; d++ {
			r[d+maxDiff] = k - abs(d)
		}
		cNeg := 0
		if k-1 > maxDiff {
			for d := -k + 1; d <= -maxDiff-1; d++ {
				cNeg += k + d
			}
		}

		fw := NewFenwick(2*n + 2)
		prevCnt := 0
		cross := int64(0)
		for _, y := range p {
			if cNeg != 0 {
				cross = (cross + int64(prevCnt*cNeg%int(mod))) % mod
			} else {
				cross = (cross + 0) % mod
			}
			for d := -maxDiff; d <= maxDiff; d++ {
				cntPair := r[d+maxDiff]
				if cntPair == 0 {
					continue
				}
				var threshold int
				if d >= 0 {
					val := y << d
					if val > 2*n+1 {
						continue
					}
					threshold = val
				} else {
					threshold = y >> (-d)
				}
				greater := prevCnt - fw.Sum(threshold)
				if greater > 0 {
					cross = (cross + int64(cntPair*greater)%mod) % mod
				}
			}
			fw.Add(y, 1)
			prevCnt++
		}

		ans := (rowInv + cross) % mod
		fmt.Fprintln(out, ans)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
