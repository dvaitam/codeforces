package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

type mat struct{ a, b, c, d int64 }

func mul(x, y mat) mat {
	return mat{
		(x.a*y.a + x.b*y.c) % mod,
		(x.a*y.b + x.b*y.d) % mod,
		(x.c*y.a + x.d*y.c) % mod,
		(x.c*y.b + x.d*y.d) % mod,
	}
}

func powJacob(n int64) mat {
	res := mat{1, 0, 0, 1}
	base := mat{1, 2, 1, 0}
	for n > 0 {
		if n&1 == 1 {
			res = mul(res, base)
		}
		base = mul(base, base)
		n >>= 1
	}
	return res
}

func jacobsthal(n int64) int64 {
	if n == 0 {
		return 0
	}
	p := powJacob(n - 1)
	return p.a % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	deg := make([]int, n+1)
	var total int64
	for i := 0; i < m; i++ {
		var a, b int
		var d int64
		fmt.Fscan(in, &a, &b, &d)
		deg[a]++
		deg[b]++
		total += d + 1
	}

	if m != n {
		fmt.Fprintln(out, 0)
		return
	}
	for i := 1; i <= n; i++ {
		if deg[i] != 2 {
			fmt.Fprintln(out, 0)
			return
		}
	}

	ans := jacobsthal(total - 1)
	ans = ans * 12 % mod
	fmt.Fprintln(out, ans)
}
