package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	a1, a2, a3 int
}

func generateTestsA() []testCaseA {
	r := rand.New(rand.NewSource(1))
	tests := []testCaseA{{3, 4, 5}, {2, 2, 6}, {1, 6, 5}}
	for len(tests) < 120 {
		tests = append(tests, testCaseA{r.Intn(1e8) + 1, r.Intn(1e8) + 1, r.Intn(1e8) + 1})
	}
	return tests
}

func expectedA(tc testCaseA) string {
	sum := tc.a1 + tc.a2 + tc.a3
	if sum%3 == 0 {
		return "0"
	}
	return "1"
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsA()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d %d\n", tc.a1, tc.a2, tc.a3)
		want := expectedA(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
