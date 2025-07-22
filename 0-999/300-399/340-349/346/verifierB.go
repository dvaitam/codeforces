package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
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

func prefixFunction(s string) []int {
	pi := make([]int, len(s))
	for i := 1; i < len(s); i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func buildAutomaton(virus string) [][]int {
	n := len(virus)
	pi := prefixFunction(virus)
	next := make([][]int, n+1)
	for k := 0; k <= n; k++ {
		next[k] = make([]int, 26)
		for c := 0; c < 26; c++ {
			if k == n {
				next[k][c] = n
				continue
			}
			t := k
			for t > 0 && byte('A'+c) != virus[t] {
				t = pi[t-1]
			}
			if byte('A'+c) == virus[t] {
				t++
			}
			next[k][c] = t
		}
	}
	return next
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expectedLength(s1, s2, virus string) int {
	n1, n2, n3 := len(s1), len(s2), len(virus)
	next := buildAutomaton(virus)
	dp := make([][][]int, n1+1)
	for i := range dp {
		dp[i] = make([][]int, n2+1)
		for j := range dp[i] {
			dp[i][j] = make([]int, n3+1)
			for k := range dp[i][j] {
				dp[i][j][k] = -1
			}
		}
	}
	var f func(i, j, k int) int
	f = func(i, j, k int) int {
		if i == n1 || j == n2 {
			return 0
		}
		if dp[i][j][k] != -1 {
			return dp[i][j][k]
		}
		res := max(f(i+1, j, k), f(i, j+1, k))
		if s1[i] == s2[j] {
			c := s1[i] - 'A'
			nk := next[k][c]
			if nk < n3 {
				res = max(res, 1+f(i+1, j+1, nk))
			}
		}
		dp[i][j][k] = res
		return res
	}
	return f(0, 0, 0)
}

func isSubsequence(s, sub string) bool {
	j := 0
	for i := 0; i < len(s) && j < len(sub); i++ {
		if s[i] == sub[j] {
			j++
		}
	}
	return j == len(sub)
}

func generateString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('A' + rng.Intn(4))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n1 := rng.Intn(5) + 1
	n2 := rng.Intn(5) + 1
	n3 := rng.Intn(3) + 1
	s1 := generateString(rng, n1)
	s2 := generateString(rng, n2)
	virus := generateString(rng, n3)
	return fmt.Sprintf("%s\n%s\n%s\n", s1, s2, virus)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{
		"ABC\nABCD\nBC\n",
		"AAAA\nAAAA\nAA\n",
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for idx, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		s1, s2, virus := lines[0], lines[1], lines[2]
		expectLen := expectedLength(s1, s2, virus)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expectLen == 0 {
			if got != "0" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected 0 got %s\ninput:\n%s", idx+1, got, tc)
				os.Exit(1)
			}
			continue
		}
		if len(got) != expectLen || !isSubsequence(s1, got) || !isSubsequence(s2, got) || strings.Contains(got, virus) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output %s\ninput:\n%s", idx+1, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
