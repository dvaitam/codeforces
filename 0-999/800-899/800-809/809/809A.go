package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	xs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i])
	}
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}

	var ans int64
	for j := 1; j < n; j++ {
		ans = (ans + xs[j]*(pow2[j]-1)) % MOD
	}
	for i := 0; i < n-1; i++ {
		ans = (ans - xs[i]*(pow2[n-1-i]-1)) % MOD
	}

	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(writer, ans)
}
