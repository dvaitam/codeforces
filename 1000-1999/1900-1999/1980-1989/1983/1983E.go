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
		vals := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &vals[i])
			sum += vals[i]
		}
		N := n - k
		var pSpec, pNorm int64
		invN1 := modPow(int64(N+1), MOD-2)
		pSpec = (int64(N/2) + 1) % MOD * invN1 % MOD
		if N > 0 {
			invN := modPow(int64(N), MOD-2)
			pNorm = (int64((N+1)/2) % MOD) * invN % MOD
		}
		var alice int64
		for i, v := range vals {
			if i < k {
				alice = (alice + v%MOD*pSpec) % MOD
			} else {
				alice = (alice + v%MOD*pNorm) % MOD
			}
		}
		bob := (sum%MOD - alice) % MOD
		if bob < 0 {
			bob += MOD
		}
		fmt.Fprintf(writer, "%d %d\n", alice, bob)
	}
}
