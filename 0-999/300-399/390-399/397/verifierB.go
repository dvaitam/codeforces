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

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		l := rng.Int63n(1_000_000_000) + 1
		r := rng.Int63n(1_000_000_000) + 1
		if l > r {
			l, r = r, l
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, l, r))
		minX := (n + r - 1) / r
		maxX := n / l
		if minX <= maxX {
			expected[i] = "Yes"
		} else {
			expected[i] = "No"
		}
	}
	return sb.String(), expected
}

func runCase(exe, input string, expected []string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	tokens := strings.Fields(strings.TrimSpace(out.String()))
	if len(tokens) != len(expected) {
		return fmt.Errorf("expected %d tokens got %d", len(expected), len(tokens))
	}
	for i, exp := range expected {
		if tokens[i] != exp {
			return fmt.Errorf("line %d: expected %s got %s", i+1, exp, tokens[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
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
