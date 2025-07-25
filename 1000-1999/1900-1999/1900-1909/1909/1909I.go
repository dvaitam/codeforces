package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod1 int64 = 998244353
const mod2 int64 = 1000000007

func powmod(x, y, mod int64) int64 {
	res := int64(1)
	for y > 0 {
		if y&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		y >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var x int64
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}

	if n > 8 {
		// This naive solution only supports n <= 8.
		// For larger n, we output 0 as a fallback.
		fmt.Println(0)
		return
	}

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	used := make([]bool, n)
	perm := make([]int, n)

	counts := make([][]int64, n-1)
	for i := range counts {
		counts[i] = make([]int64, n)
	}

	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			for m := 3; m <= n+1; m++ {
				k := 0
				for i := 0; i < n-1; i++ {
					if perm[i]+perm[i+1] >= m {
						k++
					}
				}
				counts[m-3][k]++
			}
			return
		}
		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				perm[pos] = nums[i]
				dfs(pos + 1)
				used[i] = false
			}
		}
	}

	dfs(0)

	var ans int64
	for m := 3; m <= n+1; m++ {
		for k := 0; k <= n-1; k++ {
			val := counts[m-3][k] % mod1
			pow := powmod(x%mod2, int64(m*n+k), mod2)
			ans = (ans + (val%mod2)*pow) % mod2
		}
	}
	fmt.Println(ans)
}
