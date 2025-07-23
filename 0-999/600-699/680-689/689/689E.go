package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func prepareFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	inv := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fact, inv
}

func C(n, r int64, fact, inv []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	events := make(map[int64]int)
	coords := make([]int64, 0, 2*n)
	for i := 0; i < n; i++ {
		var l, r int64
		fmt.Fscan(reader, &l, &r)
		events[l]++
		events[r+1]--
	}
	for x := range events {
		coords = append(coords, x)
	}
	sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })

	fact, inv := prepareFact(n)
	var ans int64
	active := 0
	var prev int64
	for i, x := range coords {
		if i > 0 {
			length := x - prev
			if length > 0 {
				ans = (ans + C(int64(active), int64(k), fact, inv)*length) % mod
			}
		}
		active += events[x]
		prev = x
	}
	fmt.Fprintln(writer, ans%mod)
}
