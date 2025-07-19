package main

import "fmt"

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	const MOD int64 = 998244353
	var ans int64
	switch {
	case n >= 3:
		ans = (n - 1 + n - 2) % MOD
		for k := int64(3); k < n; k++ {
			ans = (k*ans - 1) % MOD
			if ans < 0 {
				ans += MOD
			}
		}
		ans = (n * ans) % MOD
	case n == 1:
		ans = 1
	case n == 2:
		ans = 2
	}
	fmt.Println(ans)
}
