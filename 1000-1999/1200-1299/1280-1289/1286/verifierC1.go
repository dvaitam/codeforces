package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseC1 struct {
	n   int
	s   string
	exp string
}

func solveC1(n int, s string) string {
	return s
}

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

func generateTests() []testCaseC1 {
	rng := rand.New(rand.NewSource(3))
	cases := make([]testCaseC1, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := range cases {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteRune(letters[rng.Intn(len(letters))])
		}
		s := sb.String()
		cases[i] = testCaseC1{n: n, s: s, exp: solveC1(n, s)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
