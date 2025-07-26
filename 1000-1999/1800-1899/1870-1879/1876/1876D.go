package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func powMod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	first := make(map[int]int)
	last := make(map[int]int)
	cnt := make(map[int]int)
	oddSeq := make([]int, 0, n)
	evenSeq := make([]int, 0, n)
	for i, v := range a {
		if _, ok := first[v]; !ok {
			first[v] = i
		}
		last[v] = i
		cnt[v]++
		if cnt[v]%2 == 1 {
			oddSeq = append(oddSeq, v)
		} else {
			evenSeq = append(evenSeq, v)
		}
	}

	d := len(cnt)
	total := powMod(2, int64(d))

	equalCount := int64(0)
	if len(oddSeq) == len(evenSeq) {
		same := true
		for i := range oddSeq {
			if oddSeq[i] != evenSeq[i] {
				same = false
				break
			}
		}
		if same {
			// DSU on values with overlapping intervals
			vals := make([]int, 0, d)
			idx := make(map[int]int)
			for v := range cnt {
				idx[v] = len(vals)
				vals = append(vals, v)
			}
			parent := make([]int, len(vals))
			for i := range parent {
				parent[i] = i
			}
			var find func(int) int
			find = func(x int) int {
				if parent[x] != x {
					parent[x] = find(parent[x])
				}
				return parent[x]
			}
			union := func(a, b int) {
				ra, rb := find(a), find(b)
				if ra != rb {
					parent[rb] = ra
				}
			}
			for i := 0; i < len(vals); i++ {
				vi := vals[i]
				for j := i + 1; j < len(vals); j++ {
					vj := vals[j]
					fi, li := first[vi], last[vi]
					fj, lj := first[vj], last[vj]
					if !(li < fj || lj < fi) {
						union(i, j)
					}
				}
			}
			comp := 0
			for i := range parent {
				if find(i) == i {
					comp++
				}
			}
			equalCount = powMod(2, int64(comp))
		}
	}

	inv2 := (mod + 1) / 2
	ans := ((total-equalCount)%mod + mod) % mod
	ans = ans * int64(inv2) % mod
	fmt.Fprintln(out, ans)
}
