package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	pow2 := make([]int64, 60)
	for b := 0; b < 60; b++ {
		pow2[b] = (1 << b) % mod
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		cnt := make([]int64, 60)
		for _, x := range arr {
			for b := 0; b < 60; b++ {
				if (x>>b)&1 == 1 {
					cnt[b]++
				}
			}
		}
		nMod := int64(n) % mod
		result := int64(0)
		for _, x := range arr {
			andVal := int64(0)
			orVal := int64(0)
			for b := 0; b < 60; b++ {
				if (x>>b)&1 == 1 {
					andVal = (andVal + cnt[b]*pow2[b]) % mod
					orVal = (orVal + nMod*pow2[b]) % mod
				} else {
					orVal = (orVal + cnt[b]*pow2[b]) % mod
				}
			}
			result = (result + andVal*orVal%mod) % mod
		}
		fmt.Fprintln(out, result)
	}
}
