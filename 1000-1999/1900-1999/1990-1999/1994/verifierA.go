package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
4 4 12 10 1 14 3 16 11 4 6 13 15 7 8 9 5 2
5 5 17 2 22 15 7 13 1 25 20 8 9 19 21 24 6 14 12 23 18 16 11 3 4 10 5
1 5 5 2 1 3 4
1 2 2 1
2 2 2 3 1 4
3 5 4 1 14 3 11 10 7 6 13 12 9 5 2 8 15
5 4 15 13 7 20 9 16 12 2 18 5 17 1 3 14 4 6 10 8 19 11
5 4 14 4 2 1 3 12 6 19 18 5 13 16 10 11 15 7 8 20 9 17
4 5 5 13 20 16 10 15 18 17 19 14 3 6 4 12 2 9 1 8 7 11
1 1 1
2 5 9 4 6 5 8 3 1 7 2 10
2 2 3 4 2 1
1 1 1
5 4 10 15 12 19 14 13 18 16 1 17 7 6 5 11 2 8 3 9 4 20
3 4 12 6 9 5 8 2 7 11 1 4 3 10
2 1 1 2
4 5 11 14 4 7 19 15 5 20 6 2 8 1 18 16 3 9 13 12 10 17
5 3 2 8 12 9 1 4 14 7 5 3 13 15 10 6 11
3 2 1 3 5 4 6 2
1 4 2 1 3 4
2 4 1 5 7 4 8 6 3 2
5 5 5 4 8 6 13 20 17 19 22 24 3 12 18 21 9 7 14 23 15 10 11 16 2 1 25
5 5 17 16 18 12 3 19 13 15 5 1 2 11 22 14 8 23 6 9 10 24 21 20 25 7 4
1 2 1 2
1 1 1
1 2 1 2
5 3 13 9 4 15 2 5 3 7 11 8 14 12 1 10 6
1 2 2 1
3 3 5 2 4 7 3 1 8 9 6
2 2 1 4 2 3
5 2 7 4 10 3 1 9 6 8 2 5
3 4 10 1 8 3 12 4 11 9 2 7 5 6
4 1 3 1 2 4
3 3 2 8 7 1 4 5 6 9 3
2 2 3 4 2 1
1 4 1 4 2 3
5 2 8 1 4 10 5 3 9 6 2 7
3 1 2 3 1
5 4 5 6 17 18 14 20 8 12 2 3 10 15 1 11 19 9 13 16 4 7
2 5 2 1 9 5 4 3 6 8 10 7
5 5 3 4 12 15 6 19 21 24 13 20 25 16 10 7 14 18 9 5 17 2 22 23 11 1 8
4 3 2 4 11 9 1 10 5 7 6 8 12 3
4 2 3 1 2 6 4 7 8 5
1 1 1
2 3 4 1 2 3 6 5
4 3 8 3 4 12 9 6 2 1 10 11 7 5
1 1 1
4 4 13 15 11 9 4 16 1 7 6 3 2 12 10 14 8 5
1 4 4 3 1 2
5 5 17 23 14 3 25 19 9 2 6 24 8 20 10 1 15 22 16 13 7 12 18 4 11 21 5
3 1 2 1 3
1 4 3 1 2 4
4 3 4 7 6 12 5 9 3 8 10 1 2 11
5 5 21 11 13 9 23 16 6 12 8 1 24 15 18 2 7 22 3 20 19 25 5 4 10 17 14
1 3 2 3 1
3 1 1 3 2
2 1 2 1
3 1 2 1 3
2 1 2 1
1 1 1
3 3 3 5 9 8 6 2 7 4 1
3 3 8 5 7 4 6 9 2 1 3
5 3 8 1 7 11 6 12 4 15 5 14 9 13 2 10 3
2 2 1 4 2 3
1 2 2 1
4 1 3 2 4 1
3 5 9 14 11 1 13 4 8 2 6 5 15 10 3 12 7
1 4 4 1 2 3
3 3 7 4 5 6 3 1 9 2 8
4 4 15 2 5 9 1 11 10 12 16 8 4 3 14 7 13 6
1 4 2 3 1 4
3 1 1 2 3
5 5 21 20 17 7 6 19 4 12 10 24 15 25 2 1 11 9 13 18 8 22 3 16 23 14 5
3 3 8 3 2 7 6 1 9 4 5
5 3 13 1 12 6 10 11 8 14 5 7 15 3 4 9 2
4 5 1 6 8 16 17 4 5 19 7 14 2 20 3 11 18 12 15 10 9 13
4 2 7 1 2 3 5 6 4 8
4 3 2 5 7 11 6 9 12 10 4 1 8 3
1 4 3 4 2 1
1 5 5 2 4 3 1
2 2 2 4 1 3
4 4 10 4 1 9 5 2 13 8 14 12 16 15 6 3 7 11
2 4 3 1 2 5 8 4 7 6
2 4 2 3 7 8 1 6 5 4
1 2 1 2
1 2 2 1
5 3 13 14 3 4 11 15 1 8 7 6 5 9 12 2 10
4 3 3 11 10 2 9 7 1 5 6 4 8 12
3 1 2 3 1
4 2 3 1 8 6 4 2 5 7
5 1 1 2 4 5 3
3 3 5 4 6 7 8 3 2 9 1
2 2 1 2 4 3
3 5 15 3 2 10 11 9 12 13 1 4 6 8 5 7 14
1 4 1 3 2 4
1 1 1
1 2 2 1
2 4 6 2 3 7 1 5 8 4
1 2 2 1
4 2 4 1 2 5 6 8 7 3
`

type testCase struct {
	n    int
	m    int
	vals []int
}

func parseTests(raw string) ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("bad line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n in line: %q", line)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad m in line: %q", line)
		}
		need := n*m + 2
		if len(fields) != need {
			return nil, fmt.Errorf("expected %d values got %d in line: %q", need, len(fields), line)
		}
		vals := make([]int, n*m)
		for i := 0; i < n*m; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("bad value in line: %q", line)
			}
			vals[i] = v
		}
		tests = append(tests, testCase{n: n, m: m, vals: vals})
	}
	return tests, nil
}

func solve(tc testCase) string {
	n, m := tc.n, tc.m
	if n == 1 && m == 1 {
		return "-1"
	}
	total := n * m
	var sb strings.Builder
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := (tc.vals[idx] % total) + 1
			idx++
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(x))
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	idx := 0
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.vals[idx]))
			idx++
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
