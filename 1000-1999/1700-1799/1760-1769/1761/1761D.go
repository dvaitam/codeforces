package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1_000_000_007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if k > n {
		fmt.Fprintln(writer, 0)
		return
	}
	dp0 := make([]int, k+2)
	dp1 := make([]int, k+2)
	dp0[0] = 1
	for i := 0; i < n; i++ {
		ndp0 := make([]int, k+2)
		ndp1 := make([]int, k+2)
		for j := 0; j <= k && j <= i; j++ {
			if dp0[j] != 0 {
				ndp0[j] = (ndp0[j] + 3*dp0[j]) % mod
				ndp1[j+1] = (ndp1[j+1] + dp0[j]) % mod
			}
			if dp1[j] != 0 {
				ndp0[j] = (ndp0[j] + dp1[j]) % mod
				ndp1[j+1] = (ndp1[j+1] + 3*dp1[j]) % mod
			}
		}
		dp0, dp1 = ndp0, ndp1
	}
	ans := (dp0[k] + dp1[k]) % mod
	fmt.Fprintln(writer, ans)
}
