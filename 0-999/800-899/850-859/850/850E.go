package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 1000000007

func fwht(a []int64, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				u := a[i+j]
				v := a[i+j+step]
				a[i+j] = u + v
				a[i+j+step] = u - v
			}
		}
	}
	if invert {
		for i := range a {
			a[i] /= int64(n)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	N := 1 << n

	F := make([]int64, N)
	G := make([]int64, N)
	for i := 0; i < N; i++ {
		if s[i] == '1' {
			F[i] = 1
		} else {
			G[i] = 1
		}
	}

	fwht(F, false)
	fwht(G, false)
	H := make([]int64, N)
	for i := 0; i < N; i++ {
		H[i] = F[i] * G[i]
	}
	fwht(H, true)

	var count int64
	for mask := 0; mask < N; mask++ {
		weight := int64(1) << bits.OnesCount(uint(mask))
		count += weight * H[mask]
	}

	result := (3 * (count % MOD)) % MOD
	if result < 0 {
		result += MOD
	}
	fmt.Fprintln(out, result)
}
