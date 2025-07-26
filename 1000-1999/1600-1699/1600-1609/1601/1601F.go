package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod = int64(1000000007)
	m   = int64(998244353)
)

func nextNum(curr, n int64) int64 {
	if curr*10 <= n {
		return curr * 10
	}
	for curr%10 == 9 || curr+1 > n {
		curr /= 10
	}
	return curr + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	curr := int64(1)
	var ans int64
	for i := int64(1); i <= n; i++ {
		diff := (i - curr) % m
		if diff < 0 {
			diff += m
		}
		ans += diff
		ans %= mod
		if i == n {
			break
		}
		curr = nextNum(curr, n)
	}
	fmt.Fprintln(out, ans%mod)
}
