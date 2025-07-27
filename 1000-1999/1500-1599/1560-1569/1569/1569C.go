package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	const maxN = 200000
	fact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Ints(a)
		maxVal := a[n-1]
		second := a[n-2]
		if maxVal-second > 1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		if maxVal == second {
			fmt.Fprintln(writer, fact[n])
			continue
		}
		// maxVal - second == 1 and maxVal is unique
		cnt := 0
		for i := n - 1; i >= 0 && a[i] == second; i-- {
			cnt++
		}
		ans := fact[n] * int64(cnt) % mod
		ans = ans * modPow(int64(cnt+1), mod-2) % mod
		fmt.Fprintln(writer, ans)
	}
}
