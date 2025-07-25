package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN = 200000 + 5

var fac [MAXN]int64
var ifac [MAXN]int64
var invNum [MAXN]int64

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

func initComb() {
	fac[0] = 1
	for i := 1; i < MAXN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[MAXN-1] = modPow(fac[MAXN-1], MOD-2)
	for i := MAXN - 1; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
	for i := 1; i < MAXN; i++ {
		invNum[i] = modPow(int64(i), MOD-2)
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
}

func prefixCount(L int, P []int) int64 {
	if L == 0 {
		if len(P) == 0 {
			return 1
		}
		return 0
	}
	if len(P) == 0 || P[0] != 1 || P[len(P)-1] > L {
		return 0
	}
	for i := 1; i < len(P); i++ {
		if P[i] <= P[i-1] {
			return 0
		}
	}
	ans := fac[L-1]
	for i := 1; i < len(P); i++ {
		ans = ans * invNum[P[i]-1] % MOD
	}
	return ans
}

func suffixCount(L int, S []int) int64 {
	if L == 0 {
		if len(S) == 0 {
			return 1
		}
		return 0
	}
	if len(S) == 0 || S[len(S)-1] != L {
		return 0
	}
	for i := 1; i < len(S); i++ {
		if S[i] <= S[i-1] {
			return 0
		}
	}
	ans := fac[L-1]
	for i := 0; i < len(S)-1; i++ {
		ans = ans * invNum[L-S[i]] % MOD
	}
	return ans
}

func solveCase(n int, P, S []int) int64 {
	if len(P) == 0 || len(S) == 0 {
		return 0
	}
	if P[0] != 1 || S[len(S)-1] != n {
		return 0
	}
	// check intersection
	mp := make(map[int]struct{}, len(P))
	for _, v := range P {
		mp[v] = struct{}{}
	}
	inter := 0
	for _, v := range S {
		if _, ok := mp[v]; ok {
			inter++
		}
	}
	if inter != 1 {
		return 0
	}
	x := P[len(P)-1]
	if S[0] != x {
		return 0
	}
	for _, v := range P {
		if v > x {
			return 0
		}
	}
	for _, v := range S {
		if v < x {
			return 0
		}
	}
	L := x - 1
	R := n - x
	Pleft := []int{}
	for _, v := range P {
		if v < x {
			Pleft = append(Pleft, v)
		}
	}
	Sright := []int{}
	for _, v := range S {
		if v > x {
			Sright = append(Sright, v-x)
		}
	}
	A := prefixCount(L, Pleft)
	B := suffixCount(R, Sright)
	if A == 0 || B == 0 {
		return 0
	}
	ans := comb(n-1, L) * A % MOD
	ans = ans * B % MOD
	return ans
}

func main() {
	initComb()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m1, m2 int
		fmt.Fscan(in, &n, &m1, &m2)
		P := make([]int, m1)
		for i := 0; i < m1; i++ {
			fmt.Fscan(in, &P[i])
		}
		S := make([]int, m2)
		for i := 0; i < m2; i++ {
			fmt.Fscan(in, &S[i])
		}
		ans := solveCase(n, P, S)
		fmt.Fprintln(out, ans)
	}
}
