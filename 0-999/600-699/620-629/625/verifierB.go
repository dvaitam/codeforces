package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(s, t string) int {
	m := len(t)
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && t[i] != t[j] {
			j = pi[j-1]
		}
		if t[i] == t[j] {
			j++
		}
		pi[i] = j
	}
	goTo := make([][26]int, m+1)
	for k := 0; k <= m; k++ {
		for c := 0; c < 26; c++ {
			if k < m && byte('a'+c) == t[k] {
				goTo[k][c] = k + 1
			} else if k == 0 {
				goTo[k][c] = 0
			} else {
				goTo[k][c] = goTo[pi[k-1]][c]
			}
		}
	}
	const inf = int(1e9)
	dp := make([]int, m)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		ndp := make([]int, m)
		for j := range ndp {
			ndp[j] = inf
		}
		for k := 0; k < m; k++ {
			if dp[k] == inf {
				continue
			}
			if dp[k]+1 < ndp[0] {
				ndp[0] = dp[k] + 1
			}
			next := goTo[k][c]
			if next < m && dp[k] < ndp[next] {
				ndp[next] = dp[k]
			}
		}
		dp = ndp
	}
	ans := inf
	for _, v := range dp {
		if v < ans {
			ans = v
		}
	}
	return ans
}

type testB struct {
	s, t string
}

func genTests() []testB {
	rand.Seed(2)
	tests := make([]testB, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := range tests {
		n := rand.Intn(20) + 1
		m := rand.Intn(6) + 1
		sb := make([]rune, n)
		for j := range sb {
			sb[j] = letters[rand.Intn(len(letters))]
		}
		tb := make([]rune, m)
		for j := range tb {
			tb[j] = letters[rand.Intn(len(letters))]
		}
		tests[i] = testB{s: string(sb), t: string(tb)}
	}
	// add edge case
	tests = append(tests, testB{s: "aaa", t: "a"})
	tests = append(tests, testB{s: "abcdef", t: "gh"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n%s\n", t.s, t.t)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := fmt.Sprintf("%d", solveB(t.s, t.t))
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
