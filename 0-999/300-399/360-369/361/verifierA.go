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

	outLines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(outLines) != n {
		return fmt.Errorf("expected %d lines, got %d", n, len(outLines))
	}
	table := make([][]int, n)
	for i := 0; i < n; i++ {
		fields := strings.Fields(outLines[i])
		if len(fields) != n {
			return fmt.Errorf("line %d: expected %d numbers, got %d", i+1, n, len(fields))
		}
		row := make([]int, n)
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("line %d col %d: not integer", i+1, j+1)
			}
			if v < -1000 || v > 1000 {
				return fmt.Errorf("line %d col %d: value %d out of range", i+1, j+1, v)
			}
			row[j] = v
		}
		table[i] = row
	}

	// check sums
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			sum += table[i][j]
		}
		if sum != k {
			return fmt.Errorf("row %d sum %d != %d", i+1, sum, k)
		}
	}
	for j := 0; j < n; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			sum += table[i][j]
		}
		if sum != k {
			return fmt.Errorf("column %d sum %d != %d", j+1, sum, k)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA <binary>")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := [][2]int{}
	for n := 1; n <= 10; n++ {
		for k := 1; k <= 10; k++ {
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
