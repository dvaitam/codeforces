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

func buildPermutation(n int) []int {
	p := []int{}
	if n%2 == 1 {
		for i := n - 2; i >= 5; i -= 2 {
			p = append(p, i, i-1)
		}
		p = append(p, 1, 2, 3)
	} else {
		for i := n - 2; i >= 2; i -= 2 {
			p = append(p, i, i-1)
		}
	}
	p = append(p, n-1, n)
	return p
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	tests := []testCase{}
	tests = append(tests, testCase{input: "1\n4\n", expected: "2 1 3 4"})
	tests = append(tests, testCase{input: "1\n5\n", expected: "3 2 1 4 5"})
	for len(tests) < 100 {
		n := rng.Intn(97) + 4
		perm := buildPermutation(n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d\n", n))
		expSb := strings.Builder{}
		for i, v := range perm {
			if i > 0 {
				expSb.WriteByte(' ')
			}
			expSb.WriteString(fmt.Sprintf("%d", v))
		}
		tests = append(tests, testCase{input: sb.String(), expected: expSb.String()})
	}
	return tests
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
