package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1000000007

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

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		res = res * ((n - k + i) % MOD) % MOD
		res = res * modInv(i%MOD) % MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int64
	var n, m int
	if _, err := fmt.Fscan(in, &k, &n, &m); err != nil {
		return
	}
	type Interval struct{ l, r int64 }
	A := make([]Interval, n)
	B := make([]Interval, m)
	points := []int64{0, k}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i].l, &A[i].r)
		points = append(points, A[i].l-1, A[i].r)
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &B[i].l, &B[i].r)
		points = append(points, B[i].l-1, B[i].r)
	}
	sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })
	uniq := []int64{}
	for _, p := range points {
		if len(uniq) == 0 || uniq[len(uniq)-1] != p {
			uniq = append(uniq, p)
		}
	}
	id := make(map[int64]int)
	for i, v := range uniq {
		id[v] = i
	}
	s := len(uniq)
	LA := make([][]int, s)
	LB := make([][]int, s)
	for _, it := range A {
		r := id[it.r]
		l := id[it.l-1]
		LA[r] = append(LA[r], l)
	}
	for _, it := range B {
		r := id[it.r]
		l := id[it.l-1]
		LB[r] = append(LB[r], l)
	}
	LBnd := make([]int64, s)
	UBnd := make([]int64, s)
	UBnd[0] = 0
	for i := 1; i < s; i++ {
		UBnd[i] = UBnd[i-1] + (uniq[i] - uniq[i-1])
	}
	for i := 1; i < s; i++ {
		if LBnd[i] < LBnd[i-1] {
			LBnd[i] = LBnd[i-1]
		}
		for _, l := range LA[i] {
			if LBnd[l]+1 > LBnd[i] {
				LBnd[i] = LBnd[l] + 1
			}
		}
		if UBnd[i] > UBnd[i-1]+(uniq[i]-uniq[i-1]) {
			UBnd[i] = UBnd[i-1] + (uniq[i] - uniq[i-1])
		}
		for _, l := range LB[i] {
			v := UBnd[l] + (uniq[i] - uniq[l] - 1)
			if v < UBnd[i] {
				UBnd[i] = v
			}
		}
		if LBnd[i] > UBnd[i] {
			fmt.Println(0)
			return
		}
	}
	dp := map[int64]int64{0: 1}
	for i := 1; i < s; i++ {
		length := uniq[i] - uniq[i-1]
		newDP := map[int64]int64{}
		for sum, cnt := range dp {
			for add := int64(0); add <= length; add++ {
				ns := sum + add
				if ns < LBnd[i] || ns > UBnd[i] {
					continue
				}
				ways := comb(length, add)
				newDP[ns] = (newDP[ns] + cnt*ways) % MOD
			}
		}
		dp = newDP
	}
	ans := int64(0)
	for _, v := range dp {
		ans = (ans + v) % MOD
	}
	fmt.Println(ans)
}
