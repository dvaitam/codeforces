package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	const maxN = 300000
	f := make([]int64, maxN+1)
	f[0] = 1
	if maxN >= 1 {
		f[1] = 1
	}
	for i := 2; i <= maxN; i++ {
		f[i] = (f[i-1] + 2*int64(i-1)%MOD*f[i-2]) % MOD
	}

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		used := make([]bool, n+1)
		usedCount := 0
		for i := 0; i < k; i++ {
			var r, c int
			fmt.Fscan(reader, &r, &c)
			if !used[r] {
				used[r] = true
				usedCount++
			}
			if !used[c] {
				used[c] = true
				usedCount++
			}
		}
		remaining := n - usedCount
		fmt.Fprintln(writer, f[remaining])
	}
}
