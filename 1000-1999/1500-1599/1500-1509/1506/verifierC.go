package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func longestCommonSubstr(a, b string) int {
	n := len(a)
	m := len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	maxLen := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
				}
			}
		}
	}
	return maxLen
}

func expected(a, b string) int {
	l := longestCommonSubstr(a, b)
	return len(a) + len(b) - 2*l
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesC.txt: %v\n", err)
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
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		a := parts[0]
		b := parts[1]
		expect := fmt.Sprintf("%d", expected(a, b))
		input := fmt.Sprintf("1\n%s\n%s\n", a, b)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
