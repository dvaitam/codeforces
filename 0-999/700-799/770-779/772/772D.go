package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAX int = 1000000

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	cnt := make([]int, MAX)
	sum := make([]int64, MAX)
	sq := make([]int64, MAX)

	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		cnt[x]++
		sum[x] = (sum[x] + int64(x)) % MOD
		sq[x] = (sq[x] + int64(x)*int64(x)%MOD) % MOD
	}

	pow10 := []int{1, 10, 100, 1000, 10000, 100000}
	// accumulate values for all numbers having digits >= current mask
	for pos := 0; pos < 6; pos++ {
		step := pow10[pos]
		for i := MAX - 1 - step; i >= 0; i-- {
			if (i/step)%10 < 9 {
				j := i + step
				cnt[i] += cnt[j]
				sum[i] += sum[j]
				if sum[i] >= MOD {
					sum[i] -= MOD
				}
				sq[i] += sq[j]
				if sq[i] >= MOD {
					sq[i] -= MOD
				}
			}
		}
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	inv2 := (MOD + 1) / 2

	f := make([]int64, MAX)
	for i := 0; i < MAX; i++ {
		c := cnt[i]
		if c == 0 {
			continue
		}
		var mul int64
		if c >= 2 {
			mul = pow2[c-2]
		} else { // c==1
			mul = inv2
		}
		s := sum[i]
		s2 := sq[i]
		val := (s * s) % MOD
		val = (val + s2) % MOD
		f[i] = val * mul % MOD
	}

	// Möbius transform to get exact contributions
	for pos := 0; pos < 6; pos++ {
		step := pow10[pos]
		for i := 0; i < MAX-step; i++ {
			if (i/step)%10 < 9 {
				j := i + step
				f[i] -= f[j]
				if f[i] < 0 {
					f[i] += MOD
				}
			}
		}
	}

	var ans uint64
	for i := 0; i < MAX; i++ {
		if f[i] != 0 {
			ans ^= uint64(i) * uint64(f[i])
		}
	}
	fmt.Fprintln(out, ans)
}
