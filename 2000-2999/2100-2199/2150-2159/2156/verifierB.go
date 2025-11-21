package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	q int
	s string
	a []int64
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, q: 2, s: "BA", a: []int64{3, 4}},
		{n: 3, q: 1, s: "AAA", a: []int64{5}},
		{n: 4, q: 1, s: "BABA", a: []int64{20}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('A')
		} else {
			sb.WriteByte('B')
		}
	}
	arr := make([]int64, q)
	for i := 0; i < q; i++ {
		arr[i] = int64(rng.Intn(50) + 1)
	}
	return testCase{n: n, q: q, s: sb.String(), a: arr}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		for i, val := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([][]int64, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(tests) {
		return nil, fmt.Errorf("expected %d output lines, got %d", len(tests), len(lines))
	}
	results := make([][]int64, len(tests))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != tests[i].q {
			return nil, fmt.Errorf("test %d: expected %d values, got %d", i+1, tests[i].q, len(fields))
		}
		vals := make([]int64, tests[i].q)
		for j, f := range fields {
			v, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, f)
			}
			vals[j] = v
		}
		results[i] = vals
	}
	return results, nil
}

func simulate(tc testCase) []int64 {
	results := make([]int64, tc.q)
	hasB := strings.Contains(tc.s, "B")
	for idx, start := range tc.a {
		if !hasB {
			results[idx] = start
			continue
		}
		pos := 0
		cur := start
		steps := int64(0)
		for cur > 0 {
			if tc.s[pos] == 'A' {
				cur--
			} else {
				cur /= 2
			}
			steps++
			pos++
			if pos == tc.n {
				pos = 0
			}
		}
		results[idx] = steps
	}
	return results
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	input := buildInput(tests)
	out, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(out, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		exp := simulate(tc)
		for j := 0; j < tc.q; j++ {
			if got[i][j] != exp[j] {
				fmt.Fprintf(os.Stderr, "Mismatch on test %d query %d: expected %d got %d\nn=%d q=%d s=%s a=%v\n",
					i+1, j+1, exp[j], got[i][j], tc.n, tc.q, tc.s, tc.a)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
