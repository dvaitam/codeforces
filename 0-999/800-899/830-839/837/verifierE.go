package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func primeFactors(n int64) []int64 {
	factors := []int64{}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

func expected(x, y int64) string {
	primes := primeFactors(x)
	var ans int64
	for y > 0 {
		g := gcd(x, y)
		x1 := x / g
		if x1 == 1 {
			ans += y / g
			break
		}
		yg := y / g
		rmin := yg
		for _, p := range primes {
			if x1%p == 0 {
				r := yg % p
				if r < rmin {
					rmin = r
				}
			}
		}
		ans += rmin
		y -= rmin * g
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	x := rng.Int63n(100000) + 1
	y := rng.Int63n(100000) + 1
	input := fmt.Sprintf("%d %d\n", x, y)
	return input, expected(x, y)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
