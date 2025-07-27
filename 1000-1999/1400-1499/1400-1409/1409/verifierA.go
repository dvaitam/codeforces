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

func solve(a, b int64) int64 {
	d := a - b
	if d < 0 {
		d = -d
	}
	moves := d / 10
	if d%10 != 0 {
		moves++
	}
	return moves
}

func generateCase(rng *rand.Rand) (string, int64) {
	a := rng.Int63n(1_000_000_000) + 1
	b := rng.Int63n(1_000_000_000) + 1
	input := fmt.Sprintf("1\n%d %d\n", a, b)
	expected := solve(a, b)
	return input, expected
}

func runCase(exe string, input string, expected int64) error {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
