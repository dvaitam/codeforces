package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modAdd(a, b int64) int64 {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}
func modSub(a, b int64) int64 {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}
func modMul(a, b int64) int64 { return (a % MOD) * (b % MOD) % MOD }
func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = modMul(res, a)
		}
		a = modMul(a, a)
		e >>= 1
	}
	return res
}

func leftContribution(pos int, preA, preB []int64) int64 {
	sA := preA[pos-1]
	sB := preB[pos-1]
	return modSub(modMul(int64(pos), sA), sB)
}

func rightContribution(pos, n int, preA, preB []int64) int64 {
	sA := modSub(preA[n], preA[pos])
	sB := modSub(preB[n], preB[pos])
	return modSub(sB, modMul(int64(pos), sA))
}

func get(pre []int64, l, r int) int64 {
	if l > r {
		return 0
	}
	return modSub(pre[r], pre[l-1])
}

func g(L, R int, preA, preB []int64) int64 {
	if L+1 >= R {
		return 0
	}
	mid := (L + R) / 2
	sumA1 := get(preA, L+1, mid)
	sumB1 := get(preB, L+1, mid)
	part1 := modSub(sumB1, modMul(int64(L), sumA1))

	sumA2 := get(preA, mid+1, R-1)
	sumB2 := get(preB, mid+1, R-1)
	part2 := modSub(modMul(int64(R), sumA2), sumB2)

	return modAdd(part1, part2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	pos := make([]int, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &pos[i])
	}

	preA := make([]int64, n+1)
	preB := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		preA[i] = modAdd(preA[i-1], a[i])
		preB[i] = modAdd(preB[i-1], modMul(a[i], int64(i)))
	}

	inv2 := (MOD + 1) / 2
	invPow2 := make([]int64, m+2)
	invPow2[0] = 1
	for i := 1; i <= m+1; i++ {
		invPow2[i] = modMul(invPow2[i-1], inv2)
	}

	norm := modSub(1, invPow2[m])
	normInv := modPow(norm, MOD-2)

	ans := int64(0)
	for i := 1; i <= m; i++ {
		lc := leftContribution(pos[i], preA, preB)
		ans = modAdd(ans, modMul(invPow2[i], lc))
	}

	for i := 1; i <= m; i++ {
		rc := rightContribution(pos[i], n, preA, preB)
		ans = modAdd(ans, modMul(invPow2[m-i+1], rc))
	}

	for i := 1; i < m; i++ {
		for j := i + 1; j <= m; j++ {
			contrib := g(pos[i], pos[j], preA, preB)
			w := invPow2[j-i+1]
			ans = modAdd(ans, modMul(w, contrib))
		}
	}

	ans = modMul(ans, normInv)
	fmt.Fprintln(writer, ans)
}
