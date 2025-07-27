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

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(a, b, x, y, n int64) int64 {
	decA := min64(n, a-x)
	a1 := a - decA
	rem1 := n - decA
	decB := min64(rem1, b-y)
	b1 := b - decB

	decB2 := min64(n, b-y)
	b2 := b - decB2
	rem2 := n - decB2
	decA2 := min64(rem2, a-x)
	a2 := a - decA2

	prod1 := a1 * b1
	prod2 := a2 * b2
	if prod1 < prod2 {
		return prod1
	}
	return prod2
}

func generateCase(rng *rand.Rand) (string, int64) {
	x := rng.Int63n(1_000_000_000) + 1
	y := rng.Int63n(1_000_000_000) + 1
	a := x + rng.Int63n(1_000_000_000-x+1)
	b := y + rng.Int63n(1_000_000_000-y+1)
	n := rng.Int63n(1_000_000_000) + 1
	input := fmt.Sprintf("1\n%d %d %d %d %d\n", a, b, x, y, n)
	expected := solve(a, b, x, y, n)
	return input, expected
}

func runCase(exe, input string, expected int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
