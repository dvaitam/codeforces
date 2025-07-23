package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

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

// comb computes C(n,k) modulo mod using precomputed factorials.
func comb(n, k int, fact, invFact []int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var h, w, n int
	if _, err := fmt.Fscan(in, &h, &w, &n); err != nil {
		return
	}
	type Point struct{ x, y int }
	pts := make([]Point, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}
	pts[n] = Point{h, w}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})

	maxN := h + w
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}

	dp := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		p := pts[i]
		dp[i] = comb(p.x+p.y-2, p.x-1, fact, invFact)
		for j := 0; j < i; j++ {
			q := pts[j]
			if q.x <= p.x && q.y <= p.y {
				ways := dp[j] * comb(p.x-q.x+p.y-q.y, p.x-q.x, fact, invFact) % mod
				dp[i] = (dp[i] - ways + mod) % mod
			}
		}
	}
	fmt.Println(dp[n])
}
