package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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

func expected(n, m, rb, cb, rd, cd, p int64) int64 {
	q := p * modInv(100) % MOD
	r := (100 - p) * modInv(100) % MOD
	cycleR := int64(2 * (n - 1))
	cycleC := int64(2 * (m - 1))
	cycle := lcm(cycleR, cycleC)
	row := rb
	col := cb
	dr := int64(1)
	dc := int64(1)
	prefixPow := int64(1)
	B := int64(0)
	for t := int64(0); t < cycle; t++ {
		good := row == rd || col == cd
		if good {
			term := (t % MOD) * q % MOD * prefixPow % MOD
			B = (B + term) % MOD
			prefixPow = prefixPow * r % MOD
		}
		if row+dr > n || row+dr < 1 {
			dr = -dr
		}
		row += dr
		if col+dc > m || col+dc < 1 {
			dc = -dc
		}
		col += dc
	}
	F := prefixPow
	Lmod := cycle % MOD
	numerator := (B + F*Lmod%MOD) % MOD
	denom := (1 - F + MOD) % MOD
	return numerator * modInv(denom) % MOD
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := int64(rng.Intn(5) + 2)
		m := int64(rng.Intn(5) + 2)
		rb := int64(rng.Intn(int(n)) + 1)
		cb := int64(rng.Intn(int(m)) + 1)
		rd := int64(rng.Intn(int(n)) + 1)
		cd := int64(rng.Intn(int(m)) + 1)
		p := int64(rng.Intn(99) + 1)
		input := fmt.Sprintf("1\n%d %d %d %d %d %d %d\n", n, m, rb, cb, rd, cd, p)
		exp := fmt.Sprintf("%d", expected(n, m, rb, cb, rd, cd, p))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
