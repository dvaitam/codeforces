package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func modInv(x int64) int64 {
	return modPow(x, MOD-2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var d int
	if _, err := fmt.Fscan(reader, &d); err != nil {
		return
	}

	// Precompute factorials and inverse factorials up to d
	fact := make([]int64, d+1)
	invFact := make([]int64, d+1)
	fact[0] = 1
	for i := 1; i <= d; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[d] = modInv(fact[d])
	for i := d; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	n := d - 1
	bin := make([]int64, d) // coefficients for (d-1)-th forward difference
	if n >= 0 {
		for k := 0; k <= n; k++ {
			val := fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
			if (n-k)&1 == 1 {
				val = (MOD - val) % MOD
			}
			bin[k] = val
		}
	}

	var FA0, FA1, FB0 int64
	for i := 0; i <= d; i++ {
		var v int64
		fmt.Fscan(reader, &v)
		if i <= n {
			FA0 = (FA0 + v*bin[i]) % MOD
		}
		if i >= 1 && i-1 <= n {
			FA1 = (FA1 + v*bin[i-1]) % MOD
		}
	}
	for i := 0; i <= d; i++ {
		var v int64
		fmt.Fscan(reader, &v)
		if i <= n {
			FB0 = (FB0 + v*bin[i]) % MOD
		}
	}

	c := (FA1 - FA0) % MOD
	if c < 0 {
		c += MOD
	}
	diff := (FB0 - FA0) % MOD
	if diff < 0 {
		diff += MOD
	}
	s := diff * modInv(c) % MOD
	fmt.Fprintln(writer, s)
}
