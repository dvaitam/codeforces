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

type testCase struct {
	input    string
	expected string
}

func countWays(n, a, b, c int) int64 {
	var total int64
	zmax := n / 2
	if c < zmax {
		zmax = c
	}
	for z := 0; z <= zmax; z++ {
		S := 2*n - 4*z
		t := S - a
		L := (t + 1) / 2
		if L < 0 {
			L = 0
		}
		U := S / 2
		if U > b {
			U = b
		}
		if L <= U {
			total += int64(U - L + 1)
		}
	}
	return total
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	a := rng.Intn(50)
	b := rng.Intn(50)
	c := rng.Intn(50)
	input := fmt.Sprintf("%d %d %d %d\n", n, a, b, c)
	expected := fmt.Sprintf("%d", countWays(n, a, b, c))
	return testCase{input: input, expected: expected}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	var cases []testCase
	// simple deterministic
	cases = append(cases, testCase{input: "1 0 0 0\n", expected: "1"})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
