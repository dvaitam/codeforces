package main

import (
	"bufio"
	"fmt"
	"os"
)

func nextPermutation(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, s int
	if _, err := fmt.Fscan(in, &n, &s); err != nil {
		return
	}
	mod := 998244353
	if n > 8 {
		// naive solution only works for small n
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	ans := make([]int, n)
	for {
		cntS := 0
		for i := 0; i < n-1; i++ {
			if p[i]+1 == p[i+1] {
				cntS++
			}
		}
		if cntS == s {
			cntK := 0
			for i := 0; i < n-1; i++ {
				if p[i] < p[i+1] {
					cntK++
				}
			}
			ans[cntK] = (ans[cntK] + 1) % mod
		}
		if !nextPermutation(p) {
			break
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, ans[i]%mod)
	}
	out.WriteByte('\n')
}
