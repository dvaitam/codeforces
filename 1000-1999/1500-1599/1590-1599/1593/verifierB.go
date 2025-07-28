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

func expected(s string) string {
	ans := len(s)
	targets := []string{"00", "25", "50", "75"}
	for _, p := range targets {
		p0 := p[0]
		p1 := p[1]
		pos1 := -1
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == p1 {
				pos1 = i
				break
			}
		}
		if pos1 == -1 {
			continue
		}
		pos0 := -1
		for i := pos1 - 1; i >= 0; i-- {
			if s[i] == p0 {
				pos0 = i
				break
			}
		}
		if pos0 == -1 {
			continue
		}
		moves := (len(s) - pos1 - 1) + (pos1 - pos0 - 1)
		if moves < ans {
			ans = moves
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(15) + 2
	digits := make([]byte, n)
	for i := range digits {
		digits[i] = byte(rng.Intn(10)) + '0'
	}
	targets := []string{"00", "25", "50", "75"}
	pair := targets[rng.Intn(len(targets))]
	i := rng.Intn(n - 1)
	j := rng.Intn(n-i-1) + i + 1
	digits[i] = pair[0]
	digits[j] = pair[1]
	if digits[0] == '0' {
		digits[0] = '1'
	}
	return string(digits)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		input := fmt.Sprintf("1\n%s\n", s)
		expectedOutput := expected(s)
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
