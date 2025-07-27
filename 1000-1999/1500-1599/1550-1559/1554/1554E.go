package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 998244353

func powMod(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func divisors(x int) []int {
	ds := []int{}
	for i := 1; i*i <= x; i++ {
		if x%i == 0 {
			ds = append(ds, i)
			if i*i != x {
				ds = append(ds, x/i)
			}
		}
	}
	return ds
}

func check(g [][]int, k int) bool {
	n := len(g)
	if (n-1)%k != 0 {
		return false
	}
	ok := true
	var dfs func(int, int) int
	dfs = func(u, p int) int {
		cnt := 0
		for _, v := range g[u] {
			if v == p {
				continue
			}
			val := dfs(v, u)
			if val < 0 {
				ok = false
				continue
			}
			cnt += val
		}
		cnt %= k
		if p == -1 {
			if cnt != 0 {
				ok = false
			}
			return 0
		}
		if cnt == 0 {
			return 1
		}
		if cnt == k-1 {
			return 0
		}
		ok = false
		return -1
	}
	dfs(0, -1)
	return ok
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
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		F := make([]int, n+1)
		F[1] = powMod(2, n-1)
		divs := divisors(n - 1)
		for _, d := range divs {
			if d == 1 {
				continue
			}
			if check(g, d) {
				F[d] = 1
			}
		}

		ans := make([]int, n+1)
		for k := n; k >= 1; k-- {
			val := F[k]
			for m := k * 2; m <= n; m += k {
				val -= ans[m]
				if val < 0 {
					val += mod
				}
			}
			ans[k] = val % mod
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
