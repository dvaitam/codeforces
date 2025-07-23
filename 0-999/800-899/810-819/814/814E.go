package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN = 60

var (
	n      int
	d      []int
	fac    [MAXN]int64
	invfac [MAXN]int64
	g      [MAXN]int64
	pair   [MAXN]int64
	dp     [MAXN][MAXN]int64
	vis    [MAXN][MAXN]bool
)

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func C(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fac[n] * invfac[r] % MOD * invfac[n-r] % MOD
}

func initComb() {
	fac[0] = 1
	for i := 1; i < MAXN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	invfac[MAXN-1] = modPow(fac[MAXN-1], MOD-2)
	for i := MAXN - 1; i > 0; i-- {
		invfac[i-1] = invfac[i] * int64(i) % MOD
	}
	inv2 := modPow(2, MOD-2)
	pair[0] = 1
	for b := 2; b < MAXN; b += 2 {
		pair[b] = fac[b] * invfac[b/2] % MOD * modPow(inv2, int64(b/2)) % MOD
	}
	g[0] = 1
	for i := 1; i < MAXN; i++ {
		val := int64(0)
		for k := 3; k <= i; k++ {
			val = (val + C(i-1, k-1)*fac[k-1]%MOD*inv2%MOD*g[i-k]) % MOD
		}
		g[i] = val
	}
}

func waysHorizontal(a, b int) int64 {
	if b%2 == 1 {
		return 0
	}
	p := b / 2
	if p == 0 {
		return g[a]
	}
	ans := int64(0)
	for k := 0; k <= a; k++ {
		ans = (ans + C(a, k)*fac[k]%MOD*C(k+p-1, p-1)%MOD*g[a-k]) % MOD
	}
	ans = ans * pair[b] % MOD
	return ans
}

func dfs(pos, size int) int64 {
	if pos == n && size == 0 {
		return 1
	}
	if pos >= n || pos+size > n {
		return 0
	}
	if vis[pos][size] {
		return dp[pos][size]
	}
	vis[pos][size] = true
	A := 0
	for i := pos; i < pos+size; i++ {
		if d[i] == 2 {
			A++
		}
	}
	B := size - A
	res := int64(0)
	for x1 := 0; x1 <= A; x1++ {
		x0 := A - x1
		for y2 := 0; y2 <= B; y2++ {
			for y1 := 0; y1 <= B-y2; y1++ {
				y0 := B - y2 - y1
				sizeNext := x1 + y1 + 2*y2
				if pos+size+sizeNext > n {
					continue
				}
				waysHoriz := waysHorizontal(y0, x0+y1)
				if waysHoriz == 0 {
					continue
				}
				parentChoice := C(A, x1) * fac[B] % MOD * invfac[y2] % MOD * invfac[y1] % MOD * invfac[y0] % MOD
				childrenWays := fac[sizeNext] * modPow(modPow(2, MOD-2), int64(y2)) % MOD
				sub := dfs(pos+size, sizeNext)
				res = (res + parentChoice%MOD*waysHoriz%MOD*childrenWays%MOD*sub%MOD) % MOD
			}
		}
	}
	dp[pos][size] = res
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	d = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &d[i])
	}
	initComb()
	if n < 1+d[0] {
		fmt.Println(0)
		return
	}
	ans := dfs(1, d[0]) % MOD
	fmt.Println(ans)
}
