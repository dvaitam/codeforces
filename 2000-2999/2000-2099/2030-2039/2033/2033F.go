package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	cache := make(map[int]int64)
	getRank := func(k int) int64 {
		if val, ok := cache[k]; ok {
			return val
		}
		if k == 1 {
			cache[k] = 1
			return 1
		}
		a, b := 0, 1%k
		for i := 1; i <= 6*k; i++ {
			if b == 0 {
				cache[k] = int64(i)
				return int64(i)
			}
			a, b = b, (a+b)%k
		}
		cache[k] = 0
		return 0
	}

	for ; t > 0; t-- {
		var n int64
		var k int
		fmt.Fscan(in, &n, &k)
		r := getRank(k)
		ans := (n % mod) * (r % mod) % mod
		fmt.Fprintln(out, ans)
	}
}
