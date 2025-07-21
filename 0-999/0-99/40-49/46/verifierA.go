package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n int) []int {
	pos := 1
	res := make([]int, 0, n-1)
	for i := 1; i < n; i++ {
		pos = (pos+i-1)%n + 1
		res = append(res, pos)
	}
	return res
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != n-1 {
		return fmt.Errorf("expected %d numbers got %d", n-1, len(fields))
	}
	exp := expected(n)
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad output %q", f)
		}
		if v != exp[i] {
			return fmt.Errorf("at position %d expected %d got %d", i+1, exp[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for n := 2; n <= 100; n++ {
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	if err := runCase(exe, 50); err != nil {
		fmt.Fprintf(os.Stderr, "case n=50 failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
