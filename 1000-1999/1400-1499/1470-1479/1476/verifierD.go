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
	input  string
	output string
}

func solveCase(n int, s string) string {
	s = " " + s
	left := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if s[i] == 'L' {
			left[i] = 1
			if i >= 2 && s[i-1] == 'R' {
				left[i] = 2 + left[i-2]
			}
		}
	}
	right := make([]int, n+2)
	for i := n; i >= 1; i-- {
		if s[i] == 'R' {
			right[i] = 1
			if i+1 <= n && s[i+1] == 'L' {
				right[i] = 2 + right[i+2]
			}
		}
	}
	ans := make([]int, n+1)
	ans[0] = 1 + right[1]
	for i := 1; i < n; i++ {
		ans[i] = 1 + left[i] + right[i+1]
	}
	ans[n] = 1 + left[n]
	var sb strings.Builder
	for i := 0; i <= n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildCase(n int, s string) testCase {
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
	output := solveCase(n, s)
	return testCase{input: input, output: output}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	bytes := make([]byte, n)
	for i := range bytes {
		if rng.Intn(2) == 0 {
			bytes[i] = 'L'
		} else {
			bytes[i] = 'R'
		}
	}
	return buildCase(n, string(bytes))
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
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase(1, "L"))
	cases = append(cases, buildCase(2, "LR"))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
