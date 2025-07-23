package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

type pair struct {
	val int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	var l int64
	if _, err := fmt.Fscan(in, &n, &l, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	full := l / int64(n)
	rem := int(l % int64(n))

	// prepare sorted order and mapping
	pairs := make([]pair, n)
	for i := 0; i < n; i++ {
		pairs[i] = pair{val: a[i], idx: i}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].val == pairs[j].val {
			return pairs[i].idx < pairs[j].idx
		}
		return pairs[i].val < pairs[j].val
	})
	pos := make([]int, n)
	isPart := make([]bool, n)
	for r, p := range pairs {
		pos[p.idx] = r
		if p.idx < rem {
			isPart[r] = true
		}
	}

	// we only need up to min(k, full+1) lengths
	kmax := k
	if f := int(full) + 1; kmax > f {
		kmax = f
	}

	dp := make([][]int64, kmax+1)
	for i := 0; i <= kmax; i++ {
		dp[i] = make([]int64, n)
	}
	// len=1
	for i := 0; i < n; i++ {
		dp[1][i] = 1
	}
	countFull := make([]int64, kmax+1)
	countFull[1] = int64(n) % mod

	for length := 2; length <= kmax; length++ {
		prefix := int64(0)
		for i := 0; i < n; i++ {
			prefix += dp[length-1][i]
			if prefix >= mod {
				prefix -= mod
			}
			dp[length][i] = prefix
		}
		sum := int64(0)
		for i := 0; i < n; i++ {
			sum += dp[length][i]
			if sum >= mod {
				sum -= mod
			}
		}
		countFull[length] = sum
	}

	ans := int64(0)
	for length := 1; length <= kmax; length++ {
		if int64(length) <= full {
			cnt := (full - int64(length) + 1) % mod
			ans += cnt * countFull[length] % mod
			if ans >= mod {
				ans %= mod
			}
		}
	}

	if rem > 0 {
		ans += int64(rem)
		ans %= mod
		for length := 2; length <= kmax; length++ {
			if int64(length-1) > full {
				break
			}
			prefix := int64(0)
			cnt := int64(0)
			for i := 0; i < n; i++ {
				prefix += dp[length-1][i]
				if prefix >= mod {
					prefix -= mod
				}
				if isPart[i] {
					cnt += prefix
					if cnt >= mod {
						cnt %= mod
					}
				}
			}
			ans += cnt
			if ans >= mod {
				ans %= mod
			}
		}
	}

	fmt.Println(ans % mod)
}
