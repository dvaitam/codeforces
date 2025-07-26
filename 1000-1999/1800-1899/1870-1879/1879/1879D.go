package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] ^ a[i]
	}

	pow2 := make([]int64, 31)
	pow2[0] = 1
	for i := 1; i < 31; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}

	var ans int64
	for bit := 0; bit < 31; bit++ {
		var cnt0, sum0 int64 = 1, 0 // prefix 0 has parity 0 at index 0
		var cnt1, sum1 int64
		for j := 1; j <= n; j++ {
			if ((prefix[j] >> bit) & 1) == 0 {
				tmp := int64(j)*cnt1 - sum1
				ans = (ans + (tmp%MOD)*pow2[bit]) % MOD
				cnt0++
				sum0 += int64(j)
			} else {
				tmp := int64(j)*cnt0 - sum0
				ans = (ans + (tmp%MOD)*pow2[bit]) % MOD
				cnt1++
				sum1 += int64(j)
			}
		}
	}

	fmt.Fprintln(out, ans%MOD)
}
