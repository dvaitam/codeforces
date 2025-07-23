package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func countMultiples(l, r, a int64) int64 {
	if l > r {
		return 0
	}
	start := ((l + a - 1) / a) * a
	if start > r {
		return 0
	}
	return (r-start)/a + 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	m := n / k
	a := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}
	pow10 := int64(1)
	for i := 0; i < k; i++ {
		pow10 *= 10
	}
	base := pow10 / 10
	pow10Minus1 := pow10 - 1
	ans := int64(1)
	for i := 0; i < m; i++ {
		total := (pow10Minus1 / a[i]) + 1
		l := b[i] * base
		r := l + base - 1
		bad := countMultiples(l, r, a[i])
		good := (total - bad) % MOD
		if good < 0 {
			good += MOD
		}
		ans = (ans * good) % MOD
	}
	fmt.Println(ans)
}
