package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	N   = 710
	mod = 998244353
)

var (
	a   [N]int
	b   [N]int
	c   [N]int
	C   [N][N]int
	f   [2][N][N]int64
	pre [N][N]int64
	suf [N][N]int64
	g   [N][N]int64
	rdr = bufio.NewReader(os.Stdin)
	wtr = bufio.NewWriter(os.Stdout)
)

func upd(x *int64, y int64) {
	*x = (*x + y) % mod
}

func main() {
	defer wtr.Flush()
	var n int
	fmt.Fscan(rdr, &n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(rdr, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(rdr, &b[i])
	}
	for i := 0; i <= n-1; i++ {
		fmt.Fscan(rdr, &c[i])
	}
	C[0][0] = 1
	for i := 1; i <= n; i++ {
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			C[i][j] = (C[i-1][j] + C[i-1][j-1]) % mod
		}
	}
	f[1][1][0] = 1
	pre[0][0] = int64(a[1])
	suf[0][0] = int64(b[1])
	for i := 1; i <= n; i++ {
		nw := i & 1
		nx := nw ^ 1
		for j := range f[nx] {
			for k := range f[nx][j] {
				f[nx][j][k] = 0
			}
		}
		for j := 1; j <= i; j++ {
			for k := 0; k <= i-1; k++ {
				if f[nw][j][k] != 0 {
					val := f[nw][j][k]
					upd(&pre[i][k], val*int64(a[j+1]))
					upd(&suf[i][i-1-k], val*int64(b[j+1]))
					upd(&f[nx][j+1][k+1], val)
					upd(&f[nx][j][k], val*int64(k+1))
					upd(&f[nx][j][k+1], val*int64(i-k-1))
				}
			}
		}
	}
	for i := 0; i <= n; i++ {
		for x := 0; x <= i; x++ {
			for y := 0; y <= n-x-1; y++ {
				upd(&g[i][y], pre[i][x]*int64(c[x+y+func() int {
					if i > 0 {
						return 1
					} else {
						return 0
					}
				}()]))
			}
		}
	}
	for i := 1; i <= n; i++ {
		ans := int64(0)
		for p := 1; p <= i; p++ {
			for y := 0; y <= i-p; y++ {
				upd(&ans, g[p-1][y]*suf[i-p][y]%mod*int64(C[i-1][p-1])%mod)
			}
		}
		fmt.Fprint(wtr, ans, " ")
	}
}
