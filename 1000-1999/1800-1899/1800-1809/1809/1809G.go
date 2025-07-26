package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	if n == 1 {
		fmt.Fprintln(writer, 1)
		return
	}

	groups := []int{}
	start := 0
	for i := 1; i < n; i++ {
		if a[i]-a[i-1] > k {
			groups = append(groups, i-start)
			start = i
		}
	}
	groups = append(groups, n-start)

	if len(groups) == 1 {
		fmt.Fprintln(writer, 0)
		return
	}

	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}

	sumBad := int64(0)
	for _, s := range groups {
		sumBad = (sumBad + int64(s)*(int64(s)-1)) % MOD
	}

	ans := (fact[n] - fact[n-2]*sumBad) % MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(writer, ans)
}
