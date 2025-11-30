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
8 1 2 1 4 1 8 9 2 7 6 2 6 6 2 3 7
6 8 2 5 7 10 10 1 7 3 4 0 7
8 0 3 4 1 10 4 8 9 2 6 7 1 10 7 3 8
7 4 10 0 1 4 10 0 0 4 6 8 9 6 7
2 4 5 4 10
4 9 1 0 1 4 4 8 5
2 8 3 2 1
7 4 4 8 2 9 8 10 3 8 1 6 10 8 6
5 4 7 5 9 10 2 2 1 1 6
7 9 7 2 8 10 4 5 10 7 6 3 7 7 8
6 7 10 0 7 4 2 7 0 9 3 0 5
8 6 0 8 1 10 1 10 10 6 0 5 0 1 9 0 4
5 3 2 9 4 3 1 6 7 5 6
3 5 6 10 10 6 2
8 2 8 5 2 3 2 7 5 6 6 7 6 3 3 7 3
1 6 0
4 10 1 2 5 0 10 10 2
4 9 4 9 1 8 4 5 6
8 0 10 8 10 10 8 6 9 7 7 4 7 3 5 4 0
1 0 2
6 0 4 10 0 2 1 6 10 3 9 6 8
4 7 3 5 9 1 9 1 5
6 8 7 5 4 0 8 0 3 5 1 3 8
6 3 3 4 10 4 4 8 6 4 7 5 3
1 4 8
2 0 7 7 7
1 6 7
8 7 1 1 1 3 1 2 6 3 7 9 1 6 8 6 0
3 3 7 3 2 4 5
6 6 1 8 4 9 8 3 4 7 8 9 7
5 4 3 0 1 9 1 2 6 3 3
5 10 0 8 8 6 0 1 6 10 4
2 9 5 3 10
5 3 3 1 8 4 10 5 3 5 10
8 4 9 2 2 0 8 8 5 5 9 10 0 2 6 2 2
2 2 3 7 9
4 3 2 3 6 5 9 9 2
8 1 9 0 8 9 5 7 7 4 0 3 8 10 2 10 7
8 8 5 1 4 2 9 6 3 5 4 6 0 3 0 5 3
6 7 10 10 10 3 4 5 10 2 4 0 5
1 10 2
6 0 7 10 0 0 3 0 0 3 10 5 1
1 5 10
7 2 3 7 6 2 5 4 2 10 5 6 6 0 6
5 8 8 10 7 0 9 1 6 6 2
1 8 2
3 1 5 3 2 3 0
3 10 8 2 1 6 9
2 9 10 7 2
1 4 5
7 0 10 0 7 1 5 4 10 2 7 3 8 5 2
7 5 4 7 6 0 4 8 4 8 7 0 8 9 8
5 10 0 7 6 1 6 5 7 0 0
5 0 4 10 10 9 4 10 3 8 8
6 6 4 3 1 9 5 3 9 10 8 10 5
3 2 5 0 9 0 9
3 5 5 4 10 4 5
8 6 9 6 2 0 2 9 0 7 2 5 0 7 10 10 4
4 1 8 6 4 2 8 2 1
3 9 1 8 10 8 9
7 6 4 4 4 0 6 4 4 8 8 8 5 5 3
7 2 0 8 2 10 9 6 5 7 0 8 6 10 9
4 0 5 8 2 10 3 10 5
8 0 3 9 3 4 2 6 1 9 7 3 7 8 1 3 2
8 1 6 10 6 4 4 6 5 9 5 1 4 0 7 0 4
4 6 6 6 10 10 10 6 0
8 5 9 2 9 4 5 0 6 7 8 2 0 1 9 5 5
1 1 3
2 10 8 7 0
6 0 5 6 2 10 4 6 10 2 9 2 6
5 8 0 2 2 2 7 10 0 8 0
7 2 5 9 1 1 8 2 4 3 4 5 4 4 8
8 2 7 8 2 0 10 9 2 10 8 0 5 1 3 10 7
4 7 8 2 5 10 2 7 8
1 8 1
6 0 1 1 6 9 5 9 7 5 6 8 5
2 2 5 0 2
3 0 5 9 3 0 6
1 4 6
1 9 2
6 1 6 0 7 5 9 9 4 10 4 9 7
7 2 0 7 4 3 6 1 5 1 1 0 5 0 2
7 9 10 0 5 7 8 7 7 1 0 8 6 4 0
2 1 5 5 1
8 0 2 8 10 4 0 0 6 5 2 8 2 2 2 2 10
4 9 5 0 7 10 10 6 0
4 3 10 4 5 2 3 5 3
3 6 7 5 9 2 6
1 2 9
1 10 6
3 2 0 0 5 8 0
1 0 1
3 2 10 6 0 6 6
6 3 2 5 8 3 8 6 1 2 6 9 10
6 1 6 6 3 7 6 3 6 3 10 7 6
2 4 4 8 5
1 9 9
8 3 4 0 9 5 6 10 1 8 0 2 6 0 6 6 6
2 7 9 7 2
3 5 7 6 2 9 4
`

type testCase struct {
	n int
	a []int
	b []int
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, a, b []int) float64 {
	var minDiff int64 = 1<<63 - 1
	var maxSum int64 = -1 << 63
	for i := 0; i < n; i++ {
		diff := int64(a[i] - b[i])
		sum := int64(a[i] + b[i])
		if diff < minDiff {
			minDiff = diff
		}
		if sum > maxSum {
			maxSum = sum
		}
	}
	return float64(minDiff+maxSum) * 0.5
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, 1+2*n, len(fields))
		}
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			if a[i], err = strconv.Atoi(fields[1+i]); err != nil {
				return nil, fmt.Errorf("line %d parse a[%d]: %v", idx+1, i, err)
			}
			if b[i], err = strconv.Atoi(fields[1+n+i]); err != nil {
				return nil, fmt.Errorf("line %d parse b[%d]: %v", idx+1, i, err)
			}
		}
		cases = append(cases, testCase{n: n, a: a, b: b})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		expect := expected(tc.n, tc.a, tc.b)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseFloat(strings.Fields(got)[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse output %q\n", idx+1, got)
			os.Exit(1)
		}
		if diff := gotVal - expect; diff > 1e-6 || diff < -1e-6 {
			fmt.Printf("case %d failed: expected %.6f got %.8f\n", idx+1, expect, gotVal)
			fmt.Printf("input:\n%s", input.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
