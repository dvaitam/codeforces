package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesC = `1 -1 -1
3 -1 1 1 0 0 1
2 1 -1 1 1
2 0 1 0 1
5 0 1 0 1 0 -1 -1 0 0 0
4 0 1 -1 1 -1 -1 -1 -1
2 0 -1 -1 1
5 0 1 1 1 -1 0 0 1 1 0
5 0 0 0 -1 0 1 1 0 1 1
2 0 0 0 1
5 0 1 0 0 0 1 1 1 1 0
4 1 -1 0 1 -1 1 0 0
3 0 1 1 1 1 1
6 1 1 0 0 1 -1 0 1 0 1 1 -1
3 1 -1 -1 1 -1 -1
5 1 -1 0 1 -1 1 -1 1 -1 0
2 -1 -1 0 1
1 -1 0
3 -1 -1 1 -1 -1 -1
1 -1 -1
6 -1 0 0 -1 -1 1 -1 1 1 -1 0 1
1 -1 -1
1 -1 0
5 1 1 1 -1 0 0 0 -1 0 0
5 1 1 -1 0 0 1 1 -1 0 -1
1 1 1
3 -1 -1 0 -1 1 1
4 0 1 0 -1 0 0 0 1
4 1 -1 1 1 -1 1 -1 0
1 -1 -1
2 -1 0 1 -1
5 1 -1 -1 -1 1 0 -1 0 -1 1
2 1 1 1 0
3 1 0 0 1 -1 -1
1 0 0
2 -1 1 1 -1
2 -1 -1 -1 -1
2 -1 -1 -1 1
6 0 0 0 1 1 0 -1 1 -1 1 0 0
5 -1 1 1 -1 0 1 1 -1 -1 1
4 0 -1 1 -1 1 0 0 1
3 0 -1 1 0 -1 -1
3 -1 1 -1 0 -1 0
6 0 0 -1 1 1 -1 -1 0 -1 0 0 1
1 -1 0
2 -1 1 0 0
6 0 -1 0 0 -1 0 -1 -1 -1 1 0 1
4 -1 1 -1 -1 1 1 0 -1
6 1 0 0 0 0 -1 0 0 0 1 0 1
6 0 0 1 0 0 -1 -1 0 1 0 1 0
6 1 1 -1 1 1 1 -1 -1 0 -1 1 -1
4 -1 -1 1 1 -1 -1 0 0
2 1 1 1 0
4 -1 1 0 -1 -1 1 0 -1
3 1 -1 1 1 0 -1
2 0 1 -1 0
3 1 1 -1 0 0 1
1 1 0
6 -1 0 1 -1 0 0 0 0 1 1 0 1
6 0 0 1 0 1 1 1 -1 1 -1 1 1
2 0 1 0 1
1 0 0
5 1 0 1 -1 -1 -1 0 -1 1 1
5 0 -1 -1 0 1 -1 0 1 1 1
2 0 1 -1 -1
5 1 0 0 0 1 1 -1 0 1 -1
1 0 -1
4 1 1 1 -1 1 0 1 0
4 0 1 0 -1 -1 1 0 -1
4 0 1 1 1 -1 -1 -1 0
4 0 1 0 -1 -1 0 1 0
1 0 -1
5 0 1 0 -1 1 0 0 1 -1 0
6 1 1 0 -1 0 0 0 0 1 0 1 1
4 1 -1 1 1 -1 0 1 1
2 1 1 -1 0
1 -1 0
5 -1 1 0 -1 -1 -1 1 0 0 -1
3 0 -1 0 0 0 -1
4 0 0 0 -1 1 0 1 -1
3 0 -1 -1 -1 0 -1
6 0 -1 0 -1 -1 -1 0 1 1 0 0 1
3 0 1 1 1 1 1
6 0 1 0 1 1 1 1 0 1 1 1 1
1 0 0
4 0 -1 0 0 0 0 -1 -1
3 0 -1 -1 -1 0 -1
3 -1 1 1 0 -1 1
3 -1 1 0 1 0 0
6 -1 0 0 -1 0 -1 -1 -1 0 0 -1 0
3 0 -1 0 -1 -1 1
2 1 1 -1 0
3 1 1 -1 -1 -1 -1
4 0 0 -1 0 1 1 -1 0
6 1 1 0 -1 -1 0 0 0 1 0 1 1
5 -1 1 1 1 1 0 0 1 0 -1
4 0 1 0 -1 -1 0 1 -1
1 1 1
2 -1 0 0 0
6 -1 0 -1 1 0 -1 -1 -1 -1 0 1 -1
`

func solveCase(a, b []int64) int64 {
	var c, d, e int64
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] > b[i] {
			c += a[i]
		} else if a[i] < b[i] {
			d += b[i]
		} else if a[i] == 1 {
			e += a[i]
		} else if a[i] == -1 {
			e++
			c--
			d--
		}
	}
	x1 := c + e
	x2 := d + e
	x3 := (c + e + d) >> 1
	ans := x1
	if x2 < ans {
		ans = x2
	}
	if x3 < ans {
		ans = x3
	}
	return ans
}

type testCase struct {
	a []int64
	b []int64
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesC), "\n")
	tests := make([]testCase, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("bad line %d", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d length mismatch", i+1)
		}
		a := make([]int64, n)
		b := make([]int64, n)
		for j := 0; j < n; j++ {
			val, err := strconv.ParseInt(fields[1+j], 10, 64)
			if err != nil {
				return nil, err
			}
			a[j] = val
		}
		for j := 0; j < n; j++ {
			val, err := strconv.ParseInt(fields[1+n+j], 10, 64)
			if err != nil {
				return nil, err
			}
			b[j] = val
		}
		tests[i] = testCase{a: a, b: b}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.a)))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	fields := strings.Fields(output)
	if len(fields) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(fields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := strconv.FormatInt(solveCase(tc.a, tc.b), 10)
		if fields[i] != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, fields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
