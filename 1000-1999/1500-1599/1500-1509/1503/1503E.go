package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		res = res * (n - i + 1) / i
	}
	return res % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int64
	fmt.Fscan(in, &n, &m)
	if n == 1 || m == 1 {
		fmt.Println(0)
		return
	}
	if n == 2 {
		ans := 2 * comb(m+2, 4) % mod
		fmt.Println(ans)
		return
	}
	if m == 2 {
		ans := 2 * comb(n+2, 4) % mod
		fmt.Println(ans)
		return
	}
	// TODO: implement general case
	fmt.Println(0)
}
