package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxA = 1000000
const mod int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	freq := make([]int, maxA+1)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x <= maxA {
			freq[x]++
		}
	}

	// count of numbers divisible by g
	cnt := make([]int, maxA+1)
	for g := 1; g <= maxA; g++ {
		for m := g; m <= maxA; m += g {
			cnt[g] += freq[m]
		}
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] << 1) % mod
	}

	dp := make([]int64, maxA+1)
	for g := maxA; g >= 2; g-- {
		c := cnt[g]
		if c == 0 {
			continue
		}
		contrib := int64(c) * pow2[c-1] % mod
		contrib = contrib * int64(g) % mod
		var sub int64
		for m := g * 2; m <= maxA; m += g {
			sub += dp[m]
			if sub >= mod {
				sub -= mod
			}
		}
		val := contrib - sub
		val %= mod
		if val < 0 {
			val += mod
		}
		dp[g] = val
	}

	var ans int64
	for g := 2; g <= maxA; g++ {
		ans += dp[g]
		if ans >= mod {
			ans -= mod
		}
	}

	fmt.Fprintln(writer, ans)
}
