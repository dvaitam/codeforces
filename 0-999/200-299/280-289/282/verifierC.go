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

func expected(a, b string) string {
	if len(a) != len(b) {
		return "NO"
	}
	hasA1 := strings.Contains(a, "1")
	hasB1 := strings.Contains(b, "1")
	if hasA1 == hasB1 {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, a, b string, idx int) error {
	input := a + "\n" + b + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	ans := strings.TrimSpace(strings.ToUpper(out.String()))
	exp := expected(a, b)
	if ans != exp {
		return fmt.Errorf("expected %s got %s", exp, ans)
	}
	return nil
}

func generateString(rng *rand.Rand, n int, forceZero bool) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if forceZero {
			b[i] = '0'
		} else {
			if rng.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 100; i++ {
		n := rng.Intn(50) + 1
		m := n
		// 10% chance of different lengths
		if rng.Intn(10) == 0 {
			m = rng.Intn(50) + 1
			if m == n {
				m++
			}
		}

		// Generate A and B with some probability of being all zeros
		isZeroA := rng.Intn(5) == 0
		a := generateString(rng, n, isZeroA)

		isZeroB := rng.Intn(5) == 0
		b := generateString(rng, m, isZeroB)

		if err := runCase(bin, a, b, i); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nInput:\n%s\n%s\n", i, err, a, b)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}