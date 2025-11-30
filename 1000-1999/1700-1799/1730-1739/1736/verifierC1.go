package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
1 1
2 2 1
5 3 5 2 5 1
10 3 7 7 9 6 9 8 9 5 1
1 1
8 6 7 7 3 3 4 4 1
3 2 1 1
9 9 6 9 9 3 8 7 9 6
10 6 6 8 3 7 8 9 4 8 5
8 6 8 8 6 8 8 4 6
3 3 2 2
5 3 5 5 5 5
10 10 7 5 4 8 9 6 10 2 6
1 1
2 1 1
5 5 2 1 5 2
5 2 2 1 4 1
1 1
6 2 2 6 1 1 1
2 1 1
1 1
5 2 2 2 5 1
7 5 1 7 2 2 1 1
6 5 6 6 6 1 3
6 4 1 3 4 5 5
1 1
7 7 5 6 2 4 2 1
6 1 1 4 2 5 5
7 4 5 3 2 7 3 3
5 5 4 1 5 2
1 1
1 1
3 1 1 2
4 1 2 2 4
2 2 1
10 4 10 10 6 5 7 5 9 1 3
1 1
7 2 1 5 6 1 2 1
2 1 1
4 1 2 1 4
8 5 7 4 4 7 7 1 1
7 5 5 2 1 6 7 4
6 1 5 1 5 3 3
6 3 1 6 4 1 1
5 2 1 4 1 4
8 8 4 2 1 5 1 6 5
2 1 2
4 1 3 4 4
3 2 2 1
5 1 1 1 5 3
7 2 6 1 1 5 6 4
1 1
5 3 4 2 3 3
8 8 7 8 5 7 4 3 8
10 5 9 7 2 10 10 2 2 6 3
9 3 7 2 2 1 3 5 7 4
6 4 2 5 3 1 2
9 7 2 6 9 4 9 5 3 3
8 4 7 6 3 8 8 1 7
3 2 3 1
8 5 7 5 7 8 6 6 2
4 2 4 4 1
6 4 5 6 4 6 2
2 1 2
4 4 2 1 4
9 4 5 4 8 3 1 7 8 5
9 3 8 4 2 6 1 8 9 2
10 8 6 8 5 9 8 1 2 10 6
3 2 2 3
3 1 1 2
7 4 6 3 2 1 3 5
8 1 6 1 7 8 4 5 8
3 2 3 3
5 1 3 3 3 3
5 4 5 1 5 2
7 5 5 7 7 2 7 5
2 2 1
4 4 2 3 1
2 1 2
6 2 3 3 1 3 4
6 2 4 4 3 4 2
8 4 5 6 3 2 4 8 4
6 2 3 2 2 2 3
9 7 7 6 5 9 6 7 5 9
10 2 6 5 7 8 3 5 6 8 8
2 1 2
7 2 1 1 3 2 3 1
7 1 5 3 2 7 7 5
7 5 3 4 6 2 3 3
4 4 1 2 2
6 3 2 4 3 3 1
6 2 2 6 2 6 5
1 1
6 6 5 1 2 2 1
7 4 7 3 2 3 5 5
2 2 2
4 1 4 4 4
10 6 9 10 10 2 10 9 9 8 7
8 3 7 7 8 1 2 8 3
2 1 1
`

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

type testCase struct {
	n int
	a []int
}

func parseTestcases(raw string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Buffer(make([]byte, 1024), 1<<20)
	tests := make([]testCase, 0)
	lineNo := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		lineNo++
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n on line %d: %w", lineNo, err)
		}
		if len(parts) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", lineNo, 1+n, len(parts))
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(parts[i+1])
			if err != nil {
				return nil, fmt.Errorf("invalid value on line %d", lineNo)
			}
			a[i] = val
		}
		tests = append(tests, testCase{n: n, a: a})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("no testcases parsed")
	}
	return tests, nil
}

func solve(tc testCase) string {
	n := tc.n
	a := make([]int, n+1)
	for i := 0; i < n; i++ {
		a[i+1] = tc.a[i]
	}
	pref := 0
	var res int64
	for r := 1; r <= n; r++ {
		val := r - a[r] + 1
		if val > pref {
			pref = val
		}
		L := pref
		if L < 1 {
			L = 1
		}
		if L <= r {
			res += int64(r - L + 1)
		}
	}
	return fmt.Sprintf("%d", res)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		expected := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
