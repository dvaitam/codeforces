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

func expected(a, b, c int) string {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	cnt := 0
	if a == m {
		cnt++
	}
	if b == m {
		cnt++
	}
	if c == m {
		cnt++
	}
	ansA := m + 1 - a
	ansB := m + 1 - b
	ansC := m + 1 - c
	if cnt == 1 {
		if a == m {
			ansA = 0
		}
		if b == m {
			ansB = 0
		}
		if c == m {
			ansC = 0
		}
	}
	return fmt.Sprintf("%d %d %d", ansA, ansB, ansC)
}

func generateCase(rng *rand.Rand) (int, int, int) {
	a := rng.Intn(1_000_000_001)
	b := rng.Intn(1_000_000_001)
	c := rng.Intn(1_000_000_001)
	return a, b, c
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a, b, c := generateCase(rng)
		input := fmt.Sprintf("1\n%d %d %d\n", a, b, c)
		expectedOutput := expected(a, b, c)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expectedOutput {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expectedOutput, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
