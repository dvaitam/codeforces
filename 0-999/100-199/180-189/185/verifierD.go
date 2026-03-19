package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Correct brute-force solver for small cases:
// Compute LCM(k^(2^l)+1, k^(2^(l+1))+1, ..., k^(2^r)+1) mod p
// For small k and small 2^r, we can compute this directly using big.Int.
func solveBrute(k, l, r, p int64) int64 {
	lcm := big.NewInt(1)
	bk := big.NewInt(k)
	bp := big.NewInt(p)

	for i := l; i <= r; i++ {
		// a_i = k^(2^i) + 1
		exp := new(big.Int).Exp(big.NewInt(2), big.NewInt(i), nil) // 2^i as big int
		ai := new(big.Int).Exp(bk, exp, nil)                       // k^(2^i)
		ai.Add(ai, big.NewInt(1))                                   // k^(2^i) + 1

		g := new(big.Int).GCD(nil, nil, new(big.Int).Set(lcm), new(big.Int).Set(ai))
		lcm.Mul(lcm, ai)
		lcm.Div(lcm, g)
	}
	return new(big.Int).Mod(lcm, bp).Int64()
}

func genCase(rng *rand.Rand) (string, string) {
	primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}
	// Keep k small (<=50) and exponents small (l<=4, r-l<=3) to keep big.Int computation feasible
	k := int64(rng.Intn(50) + 1)
	l := rng.Int63n(5)
	r := l + rng.Int63n(4)
	p := primes[rng.Intn(len(primes))]
	input := fmt.Sprintf("1\n%d %d %d %d\n", k, l, r, p)
	expected := fmt.Sprintf("%d", solveBrute(k, l, r, p))
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genCase(rng)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: runtime error: %v\n%s\ninput:\n%s", i+1, err, stderr.String(), input)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
