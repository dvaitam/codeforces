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

func expected(n, k int, s string) string {
	cnt := 0
	for i := 0; i < n; i++ {
		if s[i] == '#' {
			cnt++
			if cnt >= k {
				return "NO"
			}
		} else {
			cnt = 0
		}
	}
	return "YES"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 2 // at least 2
	k := rng.Intn(n) + 1
	b := make([]byte, n)
	b[0] = '.'
	b[n-1] = '.'
	for i := 1; i < n-1; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '.'
		} else {
			b[i] = '#'
		}
	}
	s := string(b)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%s\n", n, k, s)
	return sb.String(), expected(n, k, s)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	got = strings.ToUpper(got)
	expected = strings.ToUpper(expected)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
