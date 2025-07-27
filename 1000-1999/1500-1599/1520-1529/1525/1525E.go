package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &dist[i][j])
		}
	}

	fact := int64(1)
	for i := 2; i <= n; i++ {
		fact *= int64(i)
	}
	factMod := fact % mod
	invFact := modPow(factMod, mod-2)

	sumNot := int64(0)
	r := make([]int, n)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			d := dist[i][j]
			rVal := n - d + 2
			if rVal < 1 {
				rVal = 1
			}
			r[i] = rVal
		}
		sort.Ints(r)
		idx := 0
		prod := int64(1)
		for pos := 1; pos <= n; pos++ {
			for idx < n && r[idx] <= pos {
				idx++
			}
			choices := idx - (pos - 1)
			if choices <= 0 {
				prod = 0
				break
			}
			prod *= int64(choices)
		}
		sumNot = (sumNot + prod%mod) % mod
	}

	total := (int64(m) % mod * factMod) % mod
	ans := (total - sumNot + mod) % mod
	ans = ans * invFact % mod
	fmt.Fprintln(out, ans)
}
