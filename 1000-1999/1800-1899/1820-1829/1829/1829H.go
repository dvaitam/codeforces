package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		freq := make([]int, 64)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			freq[x]++
		}

		pow2 := make([]int64, n+1)
		pow2[0] = 1
		for i := 1; i <= n; i++ {
			pow2[i] = pow2[i-1] * 2 % MOD
		}

		cntSuper := make([]int, 64)
		for mask := 0; mask < 64; mask++ {
			cnt := 0
			for val := 0; val < 64; val++ {
				if val&mask == mask {
					cnt += freq[val]
				}
			}
			cntSuper[mask] = cnt
		}

		dp := make([]int64, 64)
		for mask := 63; mask >= 0; mask-- {
			total := (pow2[cntSuper[mask]] - 1 + MOD) % MOD
			for sup := mask + 1; sup < 64; sup++ {
				if sup&mask == mask {
					total -= dp[sup]
					if total < 0 {
						total += MOD
					}
				}
			}
			dp[mask] = total
		}

		var result int64
		for mask := 0; mask < 64; mask++ {
			if bits.OnesCount(uint(mask)) == k {
				result += dp[mask]
				if result >= MOD {
					result -= MOD
				}
			}
		}

		fmt.Fprintln(writer, result)
	}
}
