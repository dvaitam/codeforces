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

type rule struct {
	p1 int
	t1 byte
	p2 int
	t2 byte
}

type testCase struct {
	typeStr  string
	n        int
	rules    []rule
	s        string
	expected string
}

func letterTypeMap(typeStr string) []byte {
	m := make([]byte, len(typeStr))
	for i, ch := range typeStr {
		if ch == 'V' {
			m[i] = 'V'
		} else {
			m[i] = 'C'
		}
	}
	return m
}

func checkWord(word []byte, ltypes []byte, rules []rule) bool {
	for _, r := range rules {
		if (ltypes[int(word[r.p1]-'a')] == r.t1) && (ltypes[int(word[r.p2]-'a')] != r.t2) {
			return false
		}
	}
	return true
}

func bruteForce(typeStr string, n int, rules []rule, s string) string {
	l := len(typeStr)
	letters := make([]byte, l)
	for i := 0; i < l; i++ {
		letters[i] = byte('a' + i)
	}
	ltypes := letterTypeMap(typeStr)
	best := ""
	word := make([]byte, n)
	var dfs func(pos int, greater bool) bool
	dfs = func(pos int, greater bool) bool {
		if pos == n {
			if checkWord(word, ltypes, rules) {
				best = string(word)
				return true
			}
			return false
		}
		start := byte('a')
		if !greater {
			start = s[pos]
		}
		for ch := start; ch < byte('a')+byte(l); ch++ {
			word[pos] = ch
			if dfs(pos+1, greater || ch > s[pos]) {
				return true
			}
		}
		return false
	}
	if dfs(0, false) {
		return best + "\n"
	}
	return "-1\n"
}

func generateRandomCase(rng *rand.Rand) testCase {
	l := rng.Intn(3) + 1
	typeStr := make([]byte, l)
	for i := 0; i < l; i++ {
		if rng.Intn(2) == 0 {
			typeStr[i] = 'V'
		} else {
			typeStr[i] = 'C'
		}
	}
	n := rng.Intn(4) + 1
	maxRules := n
	rcount := rng.Intn(maxRules + 1)
	rules := make([]rule, rcount)
	for i := 0; i < rcount; i++ {
		p1 := rng.Intn(n)
		p2 := rng.Intn(n)
		t1 := byte('V')
		if rng.Intn(2) == 0 {
			t1 = 'C'
		}
		t2 := byte('V')
		if rng.Intn(2) == 0 {
			t2 = 'C'
		}
		rules[i] = rule{p1: p1, t1: t1, p2: p2, t2: t2}
	}
	letters := make([]byte, l)
	for i := 0; i < l; i++ {
		letters[i] = byte('a' + i)
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(l)])
	}
	s := sb.String()
	expected := bruteForce(string(typeStr), n, rules, s)
	return testCase{typeStr: string(typeStr), n: n, rules: rules, s: s, expected: expected}
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(tc.typeStr)
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.rules)))
	for _, r := range tc.rules {
		sb.WriteString(fmt.Sprintf("%d %c %d %c\n", r.p1+1, r.t1, r.p2+1, r.t2))
	}
	sb.WriteString(tc.s)
	sb.WriteByte('\n')
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
