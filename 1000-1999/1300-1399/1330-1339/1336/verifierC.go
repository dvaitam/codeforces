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

const MOD int = 998244353

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(S, T string) int {
	n := len(S)
	m := len(T)
	s := []byte(S)
	t := []byte(T)
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}
	if m <= n {
		for i := 0; i < n; i++ {
			if i >= m || t[i] == s[0] {
				dp[i][i] = 2
			}
		}
		for length := 1; length < n; length++ {
			c := s[length]
			for l := 0; l+length-1 < n; l++ {
				r := l + length - 1
				val := dp[l][r]
				if val == 0 {
					continue
				}
				if l > 0 {
					if l-1 >= m || t[l-1] == c {
						dp[l-1][r] = (dp[l-1][r] + val) % MOD
					}
				}
				if r+1 < n {
					if r+1 >= m || t[r+1] == c {
						dp[l][r+1] = (dp[l][r+1] + val) % MOD
					}
				}
			}
		}
		ans := 0
		for i := m - 1; i < n; i++ {
			ans = (ans + dp[0][i]) % MOD
		}
		return ans
	}
	return 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
			fmt.Fprintf(os.Stderr, "case %d invalid line\n", idx)
			os.Exit(1)
		}
		S := parts[0]
		T := parts[1]
		expectedAns := expected(S, T)
		input := fmt.Sprintf("%s\n%s\n", S, T)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var ans int
		if _, err := fmt.Sscan(got, &ans); err != nil || ans != expectedAns {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx, expectedAns, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
