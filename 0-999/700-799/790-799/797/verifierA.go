package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// factorize returns prime factors of n in ascending order
func factorize(n int) []int {
	factors := []int{}
	d := 2
	for d*d <= n {
		for n%d == 0 {
			factors = append(factors, d)
			n /= d
		}
		d++
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

func generateCase(rng *rand.Rand) (string, int, int, bool) {
	n := rng.Intn(100000-2+1) + 2
	k := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	// Determine if solution exists
	factors := factorize(n)
	hasSolution := len(factors) >= k
	return sb.String(), n, k, hasSolution
}

func runCase(bin string, input string, n, k int, hasSolution bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	res := strings.TrimSpace(out.String())
	if res == "-1" {
		if hasSolution {
			return fmt.Errorf("expected solution but got -1")
		}
		return nil
	}
	fields := strings.Fields(res)
	if len(fields) != k {
		return fmt.Errorf("expected %d numbers got %d (%q)", k, len(fields), res)
	}
	prod := 1
	for _, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("non-integer output %q", f)
		}
		if v <= 1 {
			return fmt.Errorf("factor %d not >1", v)
		}
		prod *= v
	}
	if prod != n {
		return fmt.Errorf("product mismatch expected %d got %d", n, prod)
	}
	if !hasSolution {
		return fmt.Errorf("should be impossible but got factors")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, k, ok := generateCase(rng)
		if err := runCase(bin, in, n, k, ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
