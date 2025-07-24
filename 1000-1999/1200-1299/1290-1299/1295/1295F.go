package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func combInterval(length int64, t int, invFact []int64) int64 {
	if t == 0 {
		return 1
	}
	res := int64(1)
	for i := 0; i < t; i++ {
		res = res * ((length + int64(i)) % MOD) % MOD
	}
	res = res * invFact[t] % MOD
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	l := make([]int64, n)
	r := make([]int64, n)
	values := make([]int64, 0, 2*n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &l[i], &r[i])
		values = append(values, l[i])
		values = append(values, r[i]+1)
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	uniq := make([]int64, 0, len(values))
	for _, v := range values {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	arr := uniq
	m := len(arr)

	index := make(map[int64]int, m)
	for i, v := range arr {
		index[v] = i
	}
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		L[i] = index[l[i]]
		R[i] = index[r[i]+1]
	}

	// precompute factorials
	fact := make([]int64, n+50)
	invFact := make([]int64, n+50)
	fact[0] = 1
	for i := 1; i < len(fact); i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[len(fact)-1] = modPow(fact[len(fact)-1], MOD-2)
	for i := len(fact) - 1; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	dp := make([]int64, n+1)
	dp[0] = 1
	for j := m - 2; j >= 0; j-- {
		length := arr[j+1] - arr[j]
		dp2 := make([]int64, n+1)
		for pos := 0; pos <= n; pos++ {
			val := dp[pos]
			if val == 0 {
				continue
			}
			k := 0
			for pos+k < n && L[pos+k] <= j && j < R[pos+k] {
				k++
			}
			for t := 0; t <= k; t++ {
				add := val * combInterval(length, t, invFact) % MOD
				dp2[pos+t] = (dp2[pos+t] + add) % MOD
			}
		}
		dp = dp2
	}

	total := int64(1)
	for i := 0; i < n; i++ {
		total = total * ((r[i] - l[i] + 1) % MOD) % MOD
	}
	ans := dp[n] * modPow(total, MOD-2) % MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(out, ans)
}
