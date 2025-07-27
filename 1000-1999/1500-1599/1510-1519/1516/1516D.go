package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXA = 100000

func sieve() []int {
	spf := make([]int, MAXA+1)
	for i := 2; i <= MAXA; i++ {
		if spf[i] == 0 {
			for j := i; j <= MAXA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	return spf
}

func factorize(x int, spf []int) []int {
	res := []int{}
	for x > 1 {
		p := spf[x]
		res = append(res, p)
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(in, &n, &q)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	spf := sieve()
	fac := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		fac[i] = factorize(a[i], spf)
	}

	used := make([]int, MAXA+1)
	nxt := make([][]int, 17)
	for i := range nxt {
		nxt[i] = make([]int, n+2)
	}

	r := 0
	for l := 1; l <= n; l++ {
		if r < l-1 {
			r = l - 1
		}
		for r+1 <= n {
			ok := true
			for _, p := range fac[r+1] {
				if used[p] > 0 {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
			for _, p := range fac[r+1] {
				used[p]++
			}
			r++
		}
		nxt[0][l] = r + 1
		for _, p := range fac[l] {
			used[p]--
		}
	}
	nxt[0][n+1] = n + 1

	for k := 1; k < 17; k++ {
		for i := 1; i <= n+1; i++ {
			nxt[k][i] = nxt[k-1][nxt[k-1][i]]
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		ans := 0
		pos := l
		for k := 16; k >= 0; k-- {
			for nxt[k][pos] <= r {
				pos = nxt[k][pos]
				ans += 1 << k
			}
		}
		if pos <= r {
			ans++
		}
		fmt.Fprintln(out, ans)
	}
	out.Flush()
}
