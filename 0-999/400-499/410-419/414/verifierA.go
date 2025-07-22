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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func checkCase(n, k int, out string) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if !(k < n/2 || (n == 1 && k > 0)) {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}
	tokens := strings.Fields(out)
	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	if k < n/2 || (n == 1 && k > 0) {
		return fmt.Errorf("output provided but no solution should exist")
	}
	seen := map[int]bool{}
	arr := make([]int, n)
	for i, t := range tokens {
		val, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("invalid integer %q", t)
		}
		if val < 1 || val > 1_000_000_000 {
			return fmt.Errorf("value out of range: %d", val)
		}
		if seen[val] {
			return fmt.Errorf("duplicate value %d", val)
		}
		seen[val] = true
		arr[i] = val
	}
	sum := 0
	for i := 0; i+1 < n; i += 2 {
		sum += gcd(arr[i], arr[i+1])
	}
	if sum != k {
		return fmt.Errorf("expected gcd sum %d got %d", k, sum)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(20) + 1
	k := rng.Intn(40)
	if rng.Intn(5) == 0 {
		n = 1
		k = rng.Intn(3)
	}
	return fmt.Sprintf("%d %d\n", n, k), n, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, k := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkCase(n, k, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
