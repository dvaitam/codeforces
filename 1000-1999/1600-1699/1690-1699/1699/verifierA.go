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

func checkCase(n int64, out string) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if n%2 == 0 {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}
	parts := strings.Fields(out)
	if len(parts) != 3 {
		return fmt.Errorf("expected 3 integers, got %d", len(parts))
	}
	vals := make([]int64, 3)
	for i, p := range parts {
		v, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", p)
		}
		if v < 0 || v > 1_000_000_000 {
			return fmt.Errorf("value out of range: %d", v)
		}
		vals[i] = v
	}
	a, b, c := vals[0], vals[1], vals[2]
	sum := (a ^ b) + (b ^ c) + (a ^ c)
	if sum != n {
		return fmt.Errorf("expected sum %d got %d", n, sum)
	}
	if n%2 == 1 {
		// when n is odd there is no solution
		return fmt.Errorf("n is odd, should output -1")
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Int63n(1_000_000_000) + 1
	input := fmt.Sprintf("1\n%d\n", n)
	return input, n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkCase(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
