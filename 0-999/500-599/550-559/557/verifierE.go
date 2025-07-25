package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveE(s string, k int) string {
	n := len(s)
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	for l := 1; l <= n; l++ {
		for i := 0; i+l-1 < n; i++ {
			j := i + l - 1
			if l == 1 {
				dp[i][j] = true
			} else if l == 2 {
				dp[i][j] = s[i] == s[j]
			} else if s[i] == s[j] && (l <= 4 || dp[i+2][j-2]) {
				dp[i][j] = true
			}
		}
	}
	substrs := make([]string, 0)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if dp[i][j] {
				substrs = append(substrs, s[i:j+1])
			}
		}
	}
	sort.Strings(substrs)
	if k-1 < len(substrs) {
		return substrs[k-1]
	}
	return ""
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(2))
	}
	return string(b)
}

func genCase() (string, string) {
	n := rand.Intn(6) + 1
	s := randString(n)
	maxSub := n * (n + 1) / 2
	k := rand.Intn(maxSub) + 1
	res := solveE(s, k)
	input := fmt.Sprintf("%s\n%d\n", s, k)
	return input, res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input, expected := genCase()
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
