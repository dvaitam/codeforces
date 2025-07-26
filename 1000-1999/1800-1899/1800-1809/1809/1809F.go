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
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		res := solveCase(n, k, a, b)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}

func solveCase(n int, k int64, a []int64, b []int) []int64 {
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}
	L := pref[n]
	cost1 := []int{}
	for i := 0; i < n; i++ {
		if b[i] == 1 {
			cost1 = append(cost1, i)
		}
	}
	if len(cost1) == 0 {
		ans := make([]int64, n)
		for i := range ans {
			ans[i] = 2 * L
		}
		return ans
	}

	dist := func(i, j int) int64 {
		if j >= i {
			return pref[j] - pref[i]
		}
		return L - (pref[i] - pref[j])
	}

	base := int64(0)
	m := len(cost1)
	for idx := 0; idx < m; idx++ {
		i := cost1[idx]
		j := cost1[(idx+1)%m]
		var D int64
		if m == 1 {
			D = L
		} else {
			D = dist(i, j)
		}
		if D > k {
			base += D - k
		}
	}
	cost1Cost := L + base

	ans := make([]int64, n)
	for s := 0; s < n; s++ {
		if b[s] == 1 {
			ans[s] = cost1Cost
			continue
		}
		pos := sort.SearchInts(cost1, s)
		prev := cost1[(pos-1+len(cost1))%len(cost1)]
		nxt := cost1[0]
		if pos < len(cost1) {
			nxt = cost1[pos]
		}
		var D int64
		if prev == nxt {
			D = L
		} else {
			D = dist(prev, nxt)
		}
		x := dist(prev, s)
		var add int64
		if D <= k {
			add = D - x
		} else if k > x {
			add = k - x
		}
		ans[s] = cost1Cost + add
	}
	return ans
}
