package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

type dsu struct {
	p []int
}

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n)}
	for i := range d.p {
		d.p[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra != rb {
		d.p[rb] = ra
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = make([]int, n+1)
		for j := i; j <= n; j++ {
			fmt.Fscan(reader, &a[i][j])
			if i == j && a[i][j] == 2 {
				fmt.Fprintln(writer, 0)
				return
			}
		}
	}

	d := newDSU(n + 1)
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			if a[i][j] == 1 {
				for k := i; k < j; k++ {
					d.union(k, k+1)
				}
			}
		}
	}

	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			if a[i][j] == 2 && d.find(i) == d.find(j) {
				fmt.Fprintln(writer, 0)
				return
			}
		}
	}

	mustSame := make([]bool, n)
	for i := 1; i < n; i++ {
		if d.find(i) == d.find(i+1) {
			mustSame[i] = true
		}
	}

	need := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if a[i][j] == 2 {
				if i > need[j] {
					need[j] = i
				}
			}
		}
	}

	dp := make([][]int64, n)
	for i := range dp {
		dp[i] = make([]int64, n)
	}
	dp[0][0] = 1
	if need[1] > 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	for i := 1; i < n; i++ {
		for last := 0; last < i; last++ {
			val := dp[i-1][last]
			if val == 0 {
				continue
			}
			if !mustSame[i] {
				dp[i][i] = (dp[i][i] + val) % mod
			}
			dp[i][last] = (dp[i][last] + val) % mod
		}
		req := need[i+1]
		if req > 0 {
			for last := 0; last <= i; last++ {
				if last < req {
					dp[i][last] = 0
				}
			}
		}
	}

	var res int64
	for last := 0; last < n; last++ {
		res = (res + dp[n-1][last]) % mod
	}
	res = res * 2 % mod
	fmt.Fprintln(writer, res)
}
