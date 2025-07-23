package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const PHI int64 = MOD - 1

func modPow(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	nMod := int64(1)
	parity := int64(1)
	for i := 0; i < k; i++ {
		var a int64
		fmt.Fscan(reader, &a)
		nMod = nMod * (a % PHI) % PHI
		parity = parity * (a % 2) % 2
	}
	exp := (nMod - 1 + PHI) % PHI
	pow2 := modPow(2, exp)
	inv3 := modPow(3, MOD-2)

	var numer int64
	if parity == 0 { // n is even
		numer = (pow2 + 1) % MOD
	} else { // n is odd
		numer = (pow2 - 1 + MOD) % MOD
	}
	numer = numer * inv3 % MOD
	denom := pow2

	fmt.Fprintf(writer, "%d/%d\n", numer, denom)
}
