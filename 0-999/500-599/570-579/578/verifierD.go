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

func lcs(a, b string) int {
	n := len(a)
	m := len(b)
	dp := make([]int, m+1)
	for i := 1; i <= n; i++ {
		prev := 0
		for j := 1; j <= m; j++ {
			cur := dp[j]
			if a[i-1] == b[j-1] {
				dp[j] = prev + 1
			} else if dp[j] < dp[j-1] {
				dp[j] = dp[j-1]
			}
			prev = cur
		}
	}
	return dp[m]
}

func countBrute(n, m int, s string) int {
	letters := []byte("abcdefghijklmnopqrstuvwxyz"[:m])
	var dfs func(pos int, t []byte) int
	dfs = func(pos int, t []byte) int {
		if pos == n {
			if lcs(s, string(t)) == n-1 {
				return 1
			}
			return 0
		}
		count := 0
		for _, c := range letters {
			t[pos] = c
			count += dfs(pos+1, t)
		}
		return count
	}
	return dfs(0, make([]byte, n))
}

func run(bin string, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		header := strings.TrimSpace(scanner.Text())
		if header == "" {
			continue
		}
		idx++
		parts := strings.Fields(header)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad header on case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing string\n", idx)
			os.Exit(1)
		}
		s := strings.TrimSpace(scanner.Text())
		expect := countBrute(n, m, s)
		input := fmt.Sprintf("%d %d\n%s\n", n, m, s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d non-integer output %s\n", idx, got)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expect, val)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
