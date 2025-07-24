package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

// beauty computes the beauty value for all subsegments of array a using
// a naive O(n^2) approach. It returns the sum modulo MOD.
func beautySum(a []int64) int64 {
	n := len(a)
	var result int64
	for l := 0; l < n; l++ {
		prefix := int64(0)
		prefixMax := int64(0)
		maxSub := int64(0)
		curr := int64(0)
		for r := l; r < n; r++ {
			val := a[r]
			prefix += val
			if prefix > prefixMax {
				prefixMax = prefix
			}
			curr += val
			if curr < 0 {
				curr = 0
			}
			if curr > maxSub {
				maxSub = curr
			}
			beauty := prefixMax + maxSub
			result += beauty
			if result >= MOD {
				result %= MOD
			}
		}
	}
	return result % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		ans := beautySum(arr)
		fmt.Fprintln(writer, ans)
	}
}
