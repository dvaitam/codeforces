package main

import (
	"bufio"
	"fmt"
	"os"
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

func modInv(a int64) int64 {
	return modPow(a, mod-2)
}

func prepareFactorial(n int) ([]int64, []int64) {
	fac := make([]int64, n+1)
	inv := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	inv[n] = modInv(fac[n])
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fac, inv
}

func comb(n, k int, fac, inv []int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * inv[k] % mod * inv[n-k] % mod
}

func countPairs(free []bool) []int64 {
	n := len(free) - 1
	maxPairs := n / 2
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, maxPairs+1)
	}
	dp[0][0] = 1
	for i := 1; i <= n; i++ {
		for k := 0; k <= maxPairs; k++ {
			dp[i][k] = dp[i-1][k]
			if k > 0 && i >= 2 && free[i] && free[i-1] {
				dp[i][k] = (dp[i][k] + dp[i-2][k-1]) % mod
			}
		}
	}
	res := make([]int64, maxPairs+1)
	for k := 0; k <= maxPairs; k++ {
		res[k] = dp[n][k]
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var h, w, n int
	if _, err := fmt.Fscan(reader, &h, &w, &n); err != nil {
		return
	}

	rowFree := make([]bool, h+1)
	colFree := make([]bool, w+1)
	for i := 1; i <= h; i++ {
		rowFree[i] = true
	}
	for j := 1; j <= w; j++ {
		colFree[j] = true
	}

	for i := 0; i < n; i++ {
		var r1, c1, r2, c2 int
		fmt.Fscan(reader, &r1, &c1, &r2, &c2)
		rowFree[r1] = false
		rowFree[r2] = false
		colFree[c1] = false
		colFree[c2] = false
	}

	freeRows := 0
	for i := 1; i <= h; i++ {
		if rowFree[i] {
			freeRows++
		}
	}
	freeCols := 0
	for j := 1; j <= w; j++ {
		if colFree[j] {
			freeCols++
		}
	}

	rowPairs := countPairs(rowFree)
	colPairs := countPairs(colFree)

	maxN := h
	if w > maxN {
		maxN = w
	}
	fac, inv := prepareFactorial(maxN)

	ans := int64(0)
	for a := 0; a < len(rowPairs); a++ {
		for b := 0; b < len(colPairs); b++ {
			remRows := freeRows - 2*a
			remCols := freeCols - 2*b
			if remRows < 0 || remCols < 0 {
				continue
			}
			if remRows < b || remCols < a {
				continue
			}
			ways := rowPairs[a] * colPairs[b] % mod
			ways = ways * comb(remRows, b, fac, inv) % mod
			ways = ways * comb(remCols, a, fac, inv) % mod
			ways = ways * fac[a] % mod
			ways = ways * fac[b] % mod
			ans = (ans + ways) % mod
		}
	}

	fmt.Fprintln(writer, ans)
}
