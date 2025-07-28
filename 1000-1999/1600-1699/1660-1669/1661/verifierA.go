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

func solveCase(a, b []int64) int64 {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] > b[i] {
			a[i], b[i] = b[i], a[i]
		}
	}
	var ans int64
	abs := func(x int64) int64 {
		if x < 0 {
			return -x
		}
		return x
	}
	for i := 0; i+1 < n; i++ {
		ans += abs(a[i]-a[i+1]) + abs(b[i]-b[i+1])
	}
	return ans
}

type testCase struct {
	input    string
	expected string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// deterministic edge cases
	cases = append(cases, func() testCase {
		a := []int64{1, 1}
		b := []int64{1, 1}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString("2\n1 1\n1 1\n")
		exp := fmt.Sprintf("%d", solveCase(append([]int64(nil), a...), append([]int64(nil), b...)))
		return testCase{input: sb.String(), expected: exp}
	}())
	for len(cases) < 100 {
		n := rng.Intn(24) + 2
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Int63n(1_000_000_000) + 1
			b[i] = rng.Int63n(1_000_000_000) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[i]))
		}
		sb.WriteByte('\n')
		exp := fmt.Sprintf("%d", solveCase(append([]int64(nil), a...), append([]int64(nil), b...)))
		cases = append(cases, testCase{input: sb.String(), expected: exp})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
