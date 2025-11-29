package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 3 6 10 0 5 0 9
2 3 10 6 3 10 0 4
3 3 4 9 8 3 10 1 8
2 1 9 0 4 5
2 1 8 10 2 8
3 1 10 10 9 9 6
2 2 5 9 8 10 3
2 1 4 3 7 3
1 2 6 3 9 5
1 1 7 3 3
3 3 4 6 2 7 4 9 1
3 1 5 3 0 0 5
1 2 5 4 7 6
2 1 9 5 8 9
1 3 6 6 5 1 8
1 3 8 9 9 6 8
2 2 9 2 10 2 1
1 1 5 0 10
1 3 9 5 6 8 5
3 1 9 9 7 0 3
3 2 1 6 10 6 9 8
3 1 2 1 5 6 8
2 3 4 7 7 2 10 10
3 3 10 4 8 8 7 5 1
3 1 6 5 7 2 8
2 1 2 7 2 7
3 2 9 2 2 9 1 10
1 3 3 7 0 1 8
2 2 7 5 3 1 6
3 2 10 9 8 1 9 9
3 3 8 4 0 5 8 7 1
3 3 9 4 1 9 5 0 9
1 2 3 5 7 3
2 2 9 0 5 5 9
2 1 8 6 9 9
2 2 7 7 3 2 7
1 1 1 7 0
1 1 3 8 3
3 3 2 10 1 6 3 8 0
1 2 3 1 6 5
3 1 8 0 1 1 3
3 1 2 8 1 9 1
1 1 5 7 9
3 1 5 10 4 9 9
3 1 6 4 2 10 7
1 2 4 3 2 9
2 2 7 0 9 2 9
1 2 2 0 6 0
1 1 5 1 10
1 1 10 7 6
2 3 8 3 4 10 1 5
2 2 6 9 0 10 0
1 3 10 7 7 3 7
3 1 5 3 0 3 5
3 3 5 0 6 5 3 3 3
2 1 10 1 0 1
1 2 10 1 8 6
2 3 4 2 7 7 7 3
3 1 9 8 2 1 5
3 1 2 4 6 4 3
1 1 4 10 3
2 2 1 1 2 8 7
2 3 3 1 1 3 9 8
1 1 3 7 3
2 3 8 1 5 6 8 2
1 2 10 9 2 9
3 1 7 5 8 8 3
3 1 6 3 9 7 10
1 2 10 10 7 0
3 2 4 5 2 0 5 7
1 3 10 6 4 9 8
1 1 1 7 1
3 3 4 7 4 9 4 4 0
3 3 6 4 10 1 0 9 0
2 1 3 1 0 4
2 3 2 1 1 3 9 9
1 3 6 5 9 9 7
1 3 10 7 6 4 9
1 1 6 5 5
1 2 4 6 0 0
2 3 9 3 4 2 9 0
3 2 5 3 0 0 9 2
3 2 10 2 3 9 5 4
3 2 9 1 8 7 3 4
2 1 6 9 4 9
2 1 5 5 8 4
2 2 8 1 1 6 3
3 1 4 4 5 10 8
2 3 8 9 6 3 8 7
2 1 4 9 7 4
3 1 1 0 1 10 3
3 2 6 8 1 1 7 4
3 2 8 0 4 8 10 8
2 2 10 7 2 10 8
1 3 1 7 9 3 3
3 1 1 4 1 6 5
2 2 3 9 6 7 2
2 1 4 0 5 7
3 2 2 10 5 4 9 2
3 3 5 0 3 8 4 4 3`

type testCase struct {
	n, m, c int
	a, b    []int
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: not enough numbers", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		m, err2 := strconv.Atoi(fields[1])
		c, err3 := strconv.Atoi(fields[2])
		if err := firstErr(err1, err2, err3); err != nil {
			return nil, fmt.Errorf("line %d: bad header: %w", idx+1, err)
		}
		expected := 3 + n + m
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, expected, len(fields))
		}
		a := make([]int, n)
		b := make([]int, m)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[3+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid A value: %w", idx+1, err)
			}
			a[i] = v
		}
		for i := 0; i < m; i++ {
			v, err := strconv.Atoi(fields[3+n+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid B value: %w", idx+1, err)
			}
			b[i] = v
		}
		cases = append(cases, testCase{n: n, m: m, c: c, a: a, b: b})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	return cases, nil
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func runCase(bin string, tc testCase) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d %d", tc.n, tc.m, tc.c)
	for _, v := range tc.a {
		fmt.Fprintf(&input, " %d", v)
	}
	for _, v := range tc.b {
		fmt.Fprintf(&input, " %d", v)
	}
	input.WriteByte('\n')

	expected := "0"

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())

	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
