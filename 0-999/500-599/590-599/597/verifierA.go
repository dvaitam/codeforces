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

func floorDiv(a, b int64) int64 {
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func generateCase(rng *rand.Rand) (string, int64) {
	k := rng.Int63n(1_000_000_000) + 1             // 1..1e9
	a := rng.Int63n(2_000_000_001) - 1_000_000_000 // -1e9..1e9
	b := rng.Int63n(2_000_000_001) - 1_000_000_000 // -1e9..1e9
	if a > b {
		a, b = b, a
	}
	input := fmt.Sprintf("%d %d %d\n", k, a, b)
	expected := floorDiv(b, k) - floorDiv(a-1, k)
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
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
