package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		total += a[i]
	}
	sort.Ints(a)
	dp := make([]int64, total+1)
	dp[0] = 1
	var deltaSum int64
	idx := 0
	for idx < n {
		val := a[idx]
		j := idx
		for j < n && a[j] == val {
			j++
		}
		cnt := j - idx
		// compute contribution for this group using prefix dp
		for s := 0; s < val && s <= total; s++ {
			if dp[s] == 0 {
				continue
			}
			t := val + s
			ceilHalf := (t + 1) / 2
			delta := int64(val - ceilHalf)
			if delta <= 0 {
				continue
			}
			contrib := dp[s] * delta % mod
			deltaSum = (deltaSum + int64(cnt)*contrib) % mod
		}
		// update dp to include cnt copies of value val
		for c := 0; c < cnt; c++ {
			for s := total; s >= val; s-- {
				dp[s] = (dp[s] + dp[s-val]) % mod
			}
		}
		idx = j
	}
	var sum1 int64
	for s := 0; s <= total; s++ {
		if dp[s] == 0 {
			continue
		}
		ceilHalf := int64((s + 1) / 2)
		sum1 = (sum1 + dp[s]*ceilHalf) % mod
	}
	ans := (sum1 + deltaSum) % mod
	fmt.Println(ans)
}
