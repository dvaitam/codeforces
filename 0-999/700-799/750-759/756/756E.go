package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const mod int64 = 1000000007

type Group struct {
	counts []int
	sum    int
	A      int64
	r      int
}

func groupWays(cnt []int) []int64 {
	total := 0
	for _, c := range cnt {
		total += c
	}
	ways := make([]int64, total+1)
	ways[0] = 1
	cur := 0
	for _, c := range cnt {
		nw := make([]int64, cur+c+1)
		for i := 0; i <= cur; i++ {
			if ways[i] == 0 {
				continue
			}
			val := ways[i]
			for t := 0; t <= c; t++ {
				nw[i+t] = (nw[i+t] + val) % mod
			}
		}
		ways = nw
		cur += c
	}
	return ways
}

func firstValid(base, modv int) int {
	b := base % modv
	if b < 0 {
		b += modv
	}
	start := b
	if start < base {
		diff := base - start
		add := (diff + modv - 1) / modv
		start += add * modv
	}
	return start
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int, n)
	totalCoins := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
		totalCoins += b[i]
	}
	var mStr string
	fmt.Fscan(in, &mStr)
	m := new(big.Int)
	m.SetString(mStr, 10)

	// build groups
	var groups []Group
	curCnt := []int{}
	curSum := 0
	for i := 0; i < n; i++ {
		curCnt = append(curCnt, b[i])
		curSum += b[i]
		if i == n-1 || a[i] != 1 {
			ratio := int64(0)
			if i < n-1 {
				ratio = a[i]
			}
			groups = append(groups, Group{counts: curCnt, sum: curSum, A: ratio})
			curCnt = []int{}
			curSum = 0
		}
	}

	g := len(groups)
	remainders := make([]int, g)
	for i := 0; i < g-1; i++ {
		div := big.NewInt(groups[i].A)
		rem := new(big.Int)
		q := new(big.Int)
		q.QuoRem(m, div, rem)
		remainders[i] = int(rem.Int64())
		m = q
	}
	// last remainder
	if m.Cmp(big.NewInt(int64(totalCoins))) > 0 {
		fmt.Fprintln(out, 0)
		return
	}
	if m.Sign() < 0 {
		fmt.Fprintln(out, 0)
		return
	}
	remainders[g-1] = int(m.Int64())

	for i := 0; i < g; i++ {
		groups[i].r = remainders[i]
	}

	// precompute group ways
	gways := make([][]int64, g)
	for i, gr := range groups {
		gways[i] = groupWays(gr.counts)
	}

	dp := make([]int64, totalCoins+1)
	dp[0] = 1
	maxCarry := 0

	for idx := 0; idx < g-1; idx++ {
		gr := groups[idx]
		ways := gways[idx]
		next := make([]int64, totalCoins+1)
		newMax := 0
		A := int(gr.A)
		r := gr.r
		for carry := 0; carry <= maxCarry; carry++ {
			curVal := dp[carry]
			if curVal == 0 {
				continue
			}
			start := firstValid(r-carry, A)
			for k := start; k <= gr.sum; k += A {
				val := ways[k]
				if val == 0 {
					continue
				}
				t := carry + k - r
				if t < 0 {
					continue
				}
				cOut := t / A
				if cOut > totalCoins {
					continue
				}
				next[cOut] = (next[cOut] + curVal*val) % mod
				if cOut > newMax {
					newMax = cOut
				}
			}
		}
		dp = next
		maxCarry = newMax
	}

	// last group
	last := groups[g-1]
	ways := gways[g-1]
	r := last.r
	var ans int64
	for carry := 0; carry <= maxCarry; carry++ {
		curVal := dp[carry]
		if curVal == 0 {
			continue
		}
		k := r - carry
		if k < 0 || k > last.sum {
			continue
		}
		if k >= 0 && k < len(ways) {
			ans = (ans + curVal*ways[k]) % mod
		}
	}
	fmt.Fprintln(out, ans)
}
