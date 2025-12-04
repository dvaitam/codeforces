package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b)
}

func solveRef(s, t string) int {
	n := len(s)
	m := len(t)
	if m > n {
		return 0
	}

	// KMP prefix function
	pi := make([]int, m)
	for i, j := 1, 0; i < m; i++ {
		for j > 0 && t[i] != t[j] {
			j = pi[j-1]
		}
		if t[i] == t[j] {
			j++
		}
		pi[i] = j
	}

	// Limit array: limit[i] is the max start index (1-based) of a match ending <= i
	limit := make([]int, n+1)
	for i, j, lastStart := 0, 0, 0; i < n; i++ {
		for j > 0 && s[i] != t[j] {
			j = pi[j-1]
		}
		if s[i] == t[j] {
			j++
		}
		if j == m {
			lastStart = i - m + 2 // 1-based start index
			j = pi[j-1]
		}
		limit[i+1] = lastStart
	}

	mod := 1000000007
	dp := make([]int, n+1)
	sumDp := make([]int, n+1)
	dp[0] = 1
	sumDp[0] = 1

	for i := 1; i <= n; i++ {
		dp[i] = dp[i-1]
		if limit[i] > 0 {
			dp[i] = (dp[i] + sumDp[limit[i]-1]) % mod
		}
		sumDp[i] = (sumDp[i-1] + dp[i]) % mod
	}
	return (dp[n] - 1 + mod) % mod
}

func genTests() []string {
	rand.Seed(1)
	var tests []string
	for i := 0; i < 100; i++ {
		sl := rand.Intn(20) + 1
		tl := rand.Intn(10) + 1
		s := randString(sl)
		t := randString(tl)
		tests = append(tests, fmt.Sprintf("%s\n%s\n", s, t))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candB")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	tests := genTests()
	for i, input := range tests {
		lines := strings.Split(strings.TrimSpace(input), "\n")
		s := lines[0]
		t := lines[1]
		exp := solveRef(s, t)
		
		gotStr, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\nOutput: %s", i+1, err, gotStr)
			os.Exit(1)
		}

		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%d\nGot:\n%d\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}