package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func modPow(base, exp, mod int64) int64 {
	result := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return result
}

func expected(n, m int64) int64 {
	if m == 1 {
		return 0
	}
	ans := modPow(3, n, m) - 1
	if ans < 0 {
		ans += m
	}
	return ans
}

func runCase(bin string, n, m int64) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(n, m)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.ParseInt(parts[0], 10, 64)
		m, _ := strconv.ParseInt(parts[1], 10, 64)
		if err := runCase(bin, n, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d m=%d)\n", idx, err, n, m)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
