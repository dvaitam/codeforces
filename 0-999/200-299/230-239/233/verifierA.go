package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) ([]string, error) {
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
		return nil, fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.Fields(strings.TrimSpace(out.String())), nil
}

func verify(n int, tokens []string) error {
	if n%2 == 1 {
		if len(tokens) != 1 || tokens[0] != "-1" {
			return fmt.Errorf("expected -1 for odd n")
		}
		return nil
	}
	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	p := make([]int, n+1)
	used := make([]bool, n+1)
	for i := 0; i < n; i++ {
		var x int
		if _, err := fmt.Sscan(tokens[i], &x); err != nil {
			return fmt.Errorf("invalid integer at position %d", i+1)
		}
		if x < 1 || x > n {
			return fmt.Errorf("value %d out of range", x)
		}
		if used[x] {
			return fmt.Errorf("duplicate value %d", x)
		}
		used[x] = true
		p[i+1] = x
	}
	for i := 1; i <= n; i++ {
		if p[i] == i {
			return fmt.Errorf("p[%d] == %d", i, i)
		}
		if p[p[i]] != i {
			return fmt.Errorf("p[p[%d]] != %d", i, i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for n := 1; n <= 100; n++ {
		input := fmt.Sprintf("%d\n", n)
		tokens, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
		if err := verify(n, tokens); err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
