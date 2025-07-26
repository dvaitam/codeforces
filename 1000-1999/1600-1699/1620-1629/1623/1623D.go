package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(a, e int64) int64 {
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

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 { return a / gcd(a, b) * b }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	inv100 := modInv(100)
	for ; T > 0; T-- {
		var n, m int64
		var rb, cb, rd, cd int64
		var p int64
		fmt.Fscan(in, &n, &m, &rb, &cb, &rd, &cd, &p)

		q := p * inv100 % MOD         // success probability
		r := (100 - p) * inv100 % MOD // failure probability

		cycleR := int64(2 * (n - 1))
		cycleC := int64(2 * (m - 1))
		cycle := lcm(cycleR, cycleC)

		// simulate positions for one cycle
		row := rb
		col := cb
		dr := int64(1)
		dc := int64(1)

		prefixPow := int64(1) // r^prefix_count
		prefixCount := int64(0)
		B := int64(0)

		for t := int64(0); t < cycle; t++ {
			good := row == rd || col == cd
			if good {
				term := (t % MOD) * q % MOD * prefixPow % MOD
				B = (B + term) % MOD
				prefixPow = prefixPow * r % MOD
				prefixCount++
			}

			// move
			if row+dr > n || row+dr < 1 {
				dr = -dr
			}
			row += dr
			if col+dc > m || col+dc < 1 {
				dc = -dc
			}
			col += dc
		}

		F := prefixPow // r^{count_good}
		Lmod := cycle % MOD
		numerator := (B + F*Lmod%MOD) % MOD
		denom := (1 - F + MOD) % MOD
		ans := numerator * modInv(denom) % MOD

		fmt.Fprintln(out, ans)
	}
}
