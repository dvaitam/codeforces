package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseC struct {
	s string
}

func generateTestsC() []testCaseC {
	r := rand.New(rand.NewSource(1))
	tests := []testCaseC{{"aa"}, {"abc"}, {"abca"}}
	letters := []byte{'a', 'b', 'c'}
	for len(tests) < 120 {
		n := r.Intn(20) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = letters[r.Intn(3)]
		}
		tests = append(tests, testCaseC{string(b)})
	}
	return tests
}

func solveC(s string) int {
	if strings.Contains(s, "aa") {
		return 2
	}
	if strings.Contains(s, "aba") || strings.Contains(s, "aca") {
		return 3
	}
	if strings.Contains(s, "abca") || strings.Contains(s, "acba") {
		return 4
	}
	if strings.Contains(s, "abbacca") || strings.Contains(s, "accabba") {
		return 7
	}
	return -1
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsC()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n", len(tc.s), tc.s)
		exp := fmt.Sprintf("%d", solveC(tc.s))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
