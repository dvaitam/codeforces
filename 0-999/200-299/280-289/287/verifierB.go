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

func expected(n, k int64) int64 {
	if n == 1 {
		return 0
	}
	target := n - 1
	maxTotal := k * (k - 1) / 2
	if maxTotal < target {
		return -1
	}
	l, r := int64(1), k-1
	ans := k
	for l <= r {
		m := (l + r) / 2
		sum := m * (2*k - m - 1) / 2
		if sum >= target {
			ans = m
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
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
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("bad test format on line %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.ParseInt(parts[0], 10, 64)
		k, _ := strconv.ParseInt(parts[1], 10, 64)
		want := fmt.Sprintf("%d", expected(n, k))
		input := fmt.Sprintf("%d %d\n", n, k)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed: expected %s got %s (n=%d k=%d)\n", idx, want, got, n, k)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
