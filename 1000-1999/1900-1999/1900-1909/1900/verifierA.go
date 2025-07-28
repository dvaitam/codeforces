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

type testCase struct {
	input    string
	expected string
}

func solveCase(n int, s string) int {
	empty := 0
	hasTriple := false
	for i := 0; i < n; i++ {
		if s[i] == '.' {
			empty++
		}
		if i+2 < n && s[i] == '.' && s[i+1] == '.' && s[i+2] == '.' {
			hasTriple = true
		}
	}
	if empty == 0 {
		return 0
	}
	if hasTriple {
		return 2
	}
	return empty
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		bytesStr := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				bytesStr[j] = '.'
			} else {
				bytesStr[j] = '#'
			}
		}
		s := string(bytesStr)
		in.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
		out.WriteString(fmt.Sprintf("%d\n", solveCase(n, s)))
	}
	return testCase{input: in.String(), expected: out.String()}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "1\n1\n.\n", expected: "1\n"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
