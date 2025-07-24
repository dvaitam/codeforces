package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func modPow(a, e int64) int64 {
	a %= MOD
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow((a%MOD+MOD)%MOD, MOD-2)
}

func solve(n, k int64) int64 {
	g := gcd(n, k)
	n /= g
	k /= g

	pow2k := modPow(2, k)
	invPow2k := modInv(pow2k)

	Lterm := n % MOD * modInv((pow2k-1+MOD)%MOD) % MOD

	invK := modInv(k)
	invNK := modInv(n % MOD * k % MOD)
	B := (n%MOD*invK%MOD - invNK + MOD) % MOD

	inv2 := modInv(2)
	A := ((n+1)%MOD*inv2%MOD - ((k-1)%MOD)*B%MOD*inv2%MOD + MOD) % MOD

	Dnumer := (1 - ((k+1)%MOD)*invPow2k%MOD + MOD) % MOD
	Ddenom := (1 - invPow2k + MOD) % MOD
	D := Dnumer * modInv(Ddenom) % MOD

	return (Lterm + A + B*D%MOD) % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		fmt.Fprintln(writer, solve(n, k))
	}
}
