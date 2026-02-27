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

func optimal(s string) string {
	k := strings.Count(s, "a")
	return strings.Repeat("a", k) + strings.Repeat("b", len(s)-k)
}

func simulate(s string, ops []int) string {
	b := []byte(s)
	for i, op := range ops {
		if op == 1 {
			for l, r := 0, i; l < r; l, r = l+1, r-1 {
				b[l], b[r] = b[r], b[l]
			}
		}
	}
	return string(b)
}

func runCase(bin, s string) error {
	input := s + "\n"
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	parts := strings.Fields(strings.TrimSpace(out.String()))
	n := len(s)
	if len(parts) != n {
		return fmt.Errorf("expected %d integers, got %d: %q", n, len(parts), out.String())
	}
	ops := make([]int, n)
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || (v != 0 && v != 1) {
			return fmt.Errorf("invalid value %q at position %d", p, i)
		}
		ops[i] = v
	}
	got := simulate(s, ops)
	want := optimal(s)
	if got != want {
		return fmt.Errorf("resulting string %q is not optimal %q", got, want)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []string{"abba", "aaaa"}
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			if rng.Intn(2) == 0 {
				b[j] = 'a'
			} else {
				b[j] = 'b'
			}
		}
		tests = append(tests, string(b))
	}

	for idx, s := range tests {
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s\n", idx+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
