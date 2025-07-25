package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n int
	s string
	t string
}

func solve(n int, s, t string) string {
	type pair struct {
		i    int
		s, t byte
	}
	ms := make([]pair, 0)
	for i := 0; i < n; i++ {
		if s[i] != t[i] {
			ms = append(ms, pair{i, s[i], t[i]})
		}
	}
	cnt := len(ms)
	var pos [26][26]int
	for k, p := range ms {
		pos[p.s-'a'][p.t-'a'] = k + 1
	}
	for _, p := range ms {
		if opp := pos[p.t-'a'][p.s-'a']; opp != 0 {
			return fmt.Sprintf("%d\n%d %d\n", cnt-2, p.i+1, ms[opp-1].i+1)
		}
	}
	for _, p := range ms {
		row := p.t - 'a'
		for c := 0; c < 26; c++ {
			if pos[row][c] != 0 {
				j := pos[row][c] - 1
				return fmt.Sprintf("%d\n%d %d\n", cnt-1, p.i+1, ms[j].i+1)
			}
		}
	}
	return fmt.Sprintf("%d\n-1 -1\n", cnt)
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(20) + 1
		b1 := make([]byte, n)
		b2 := make([]byte, n)
		for i := 0; i < n; i++ {
			b1[i] = byte('a' + rnd.Intn(26))
			b2[i] = byte('a' + rnd.Intn(26))
		}
		tests = append(tests, testCase{n, string(b1), string(b2)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.s, tc.t)
		expected := solve(tc.n, tc.s, tc.t)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		want := strings.TrimSpace(expected)
		if got != want {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\n got: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
