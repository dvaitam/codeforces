package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func addCandidates(cand map[int64]struct{}, arr *[]int64, val, L, R int64) {
	if val < L || val > R || val <= 0 {
		return
	}
	if _, ok := cand[val]; ok {
		return
	}
	cand[val] = struct{}{}
	*arr = append(*arr, val)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y, z, k int64
		fmt.Fscan(in, &x, &y, &z, &k)
		ans := int64(1 << 62)
		for q := int64(0); ; q++ {
			forcedDamage := k * q * (q + 1) / 2
			if forcedDamage >= z {
				break
			}
			L := q*k + 1
			R := (q + 1) * k
			C := z - forcedDamage
			if C <= 0 {
				continue
			}
			candMap := make(map[int64]struct{})
			candidates := make([]int64, 0, 16)

			for delta := int64(0); delta < 4; delta++ {
				addCandidates(candMap, &candidates, L+delta, L, R)
				addCandidates(candMap, &candidates, R-delta, L, R)
			}

			approxFloat := math.Sqrt(float64(C) * float64(y) / float64(x))
			approx := int64(approxFloat)
			for delta := -5; delta <= 5; delta++ {
				addCandidates(candMap, &candidates, approx+int64(delta), L, R)
			}

			if len(candidates) == 0 {
				continue
			}
			for _, cur := range candidates {
				attacksFinal := (C + cur - 1) / cur
				cost := cur*x + q*y + attacksFinal*y
				if cost < ans {
					ans = cost
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
