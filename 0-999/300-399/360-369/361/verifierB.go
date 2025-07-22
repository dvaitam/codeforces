package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func runTest(binary string, n, k int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %v", err)
	}
	fmt.Fprintf(stdin, "%d %d\n", n, k)
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("program error: %v", err)
	}

	out := strings.TrimSpace(stdout.String())
	tokens := strings.Fields(out)
	if n == k {
		if len(tokens) != 1 || tokens[0] != "-1" {
			return fmt.Errorf("expected -1 for n=k=%d", n)
		}
		return nil
	}

	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	seen := make([]bool, n+1)
	perm := make([]int, n)
	for i, t := range tokens {
		v, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("token %d not integer", i+1)
		}
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range 1..%d", v, n)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
		perm[i] = v
	}
	good := 0
	for i, v := range perm {
		if gcd(i+1, v) > 1 {
			good++
		}
	}
	if good != k {
		return fmt.Errorf("good elements %d != %d", good, k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB <binary>")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := [][2]int{}
	for n := 1; n <= 14; n++ {
		for k := 0; k <= n; k++ {
			tests = append(tests, [2]int{n, k})
		}
	}
	for idx, t := range tests {
		if err := runTest(binary, t[0], t[1]); err != nil {
			fmt.Printf("test %d (n=%d k=%d) failed: %v\n", idx+1, t[0], t[1], err)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
