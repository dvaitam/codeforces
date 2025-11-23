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

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func randPrime(rng *rand.Rand) int {
	for {
		v := rng.Intn(100000-2) + 2
		if isPrime(v) {
			return v
		}
	}
}

func genCase(rng *rand.Rand) (string, []int) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	nums := make([]int, t)
	for i := 0; i < t; i++ {
		n := randPrime(rng)
		nums[i] = n
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return sb.String(), nums
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func check(inputs []int, output string) error {
	reader := strings.NewReader(output)
	var m int
	for i, n := range inputs {
		if _, err := fmt.Fscan(reader, &m); err != nil {
			return fmt.Errorf("failed to read output for case %d (n=%d): %v", i+1, n, err)
		}

		// Check 1: m is prime
		if !isPrime(m) {
			return fmt.Errorf("case %d (n=%d): output m=%d is not prime", i+1, n, m)
		}

		// Check 2: n + m is not prime
		sum := n + m
		if isPrime(sum) {
			return fmt.Errorf("case %d (n=%d): n+m=%d is prime (should be composite)", i+1, n, sum)
		}
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
	for i := 1; i <= 100; i++ {
		input, inputs := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if err := check(inputs, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}