package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(n, m int) int {
	count := 0
	for a := 0; a*a <= n; a++ {
		b := n - a*a
		if b >= 0 && a+b*b == m {
			count++
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, m int
		if _, err := fmt.Sscan(line, &n, &m); err != nil {
			fmt.Fprintf(os.Stderr, "invalid line %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect := solve(n, m)
		input := fmt.Sprintf("%d %d\n", n, m)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output %q\n", idx, outStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d (n=%d m=%d)\n", idx, expect, got, n, m)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
