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

func tryOps(orig []byte, target byte) ([]int, bool) {
	n := len(orig)
	a := make([]byte, n)
	copy(a, orig)
	ops := make([]int, 0, n)
	for i := 0; i < n-1; i++ {
		if a[i] != target {
			ops = append(ops, i+1)
			a[i] = target
			if a[i+1] == 'W' {
				a[i+1] = 'B'
			} else {
				a[i+1] = 'W'
			}
		}
	}
	for i := 0; i < n; i++ {
		if a[i] != target {
			return nil, false
		}
	}
	return ops, true
}

func expected(n int, s string) string {
	orig := []byte(s)
	if ops, ok := tryOps(orig, 'W'); ok {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(ops))
		for i, v := range ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		if len(ops) > 0 {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
		return strings.TrimRight(sb.String(), "\n")
	}
	if ops, ok := tryOps(orig, 'B'); ok {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(ops))
		for i, v := range ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		if len(ops) > 0 {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
		return strings.TrimRight(sb.String(), "\n")
	}
	return "-1"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(199) + 2 // 2..200
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'W'
		} else {
			b[i] = 'B'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, expected(n, s)
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
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
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
		input, expect := generateCase(rng)
		if err := runCase(exe, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
