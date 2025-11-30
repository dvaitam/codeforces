package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `
100
0 1 1
2
1
5 3 3
4 1 4
0 4 5
1 4 4
5 4 2 4
3 4 2 0
0 3 4
2 3 3
4 1 4 1
1 2 1
1 2
1
1 5 5
2 4 5 4 1
3 3 5 4 2
4 3 3
3 1 3
5 5 3
5 5 2
3 2 3 4 4
2 5
3 4 3
4 5 4 5
3 3 5
1 3 2
4 2 3
2 2
5 5 5
4 4 5 4 4
3 2 5 1 3
4 3 5
0 2 5
0 1 5 0 0
4 1 3
4
1 5 0
4 2 3
1 1
0 3 5
0 1 3
2
1 1 5
0 1 1
0
0
0 1 3
2
1 1 5
1 5 1
3 4 0 1 1
0
0 3 5
5 5 5
0 2 2 3 0
2 4 5
4 5 0 2
3 4 5 1 3
1 1 3
0
0 3 1
4 5 4
3 4 2 1 2
2 2 4 3
5 1 5
1
5 0 2 0 1
1 2 1
3 5
1
4 1 2
1
5 3
0 3 1
4 1 4
4
5 3 3
5 3 2
4 0 1
0 4 4
1 0 4 5
0 1 0 0
0 2 2
0 1
0 4
5 4 4
2 4 5 3
1 5 1 5
3 4 5
0 4 4 0
3 4 4 1 0
5 4 3
0 4 0 4
2 2 5
2 3 1
5 3 0
0
2 2 1
3 0
3
5 4 4
1 4 4 0
0 2 0 2
2 1 2
3
1 0
4 3 4
5 3 1
2 3 0 2
0 1 1
4
2
5 4 2
5 0 0 4
5 3
0 4 3
2 3 1 2
2 3 4
3 4 4
5 2 3 1
1 3 4 2
4 4 1
4 5 4 0
0
2 2 5
1 3
0 0 5 5 0
1 3 4
1 5 5
5 2 3 1
4 3 1
1 4 3
0
2 5 2
5 4 2 1 1
3 5
1 4 3
4 5 1 3
3 5 0
4 4 2
3 4 0 3
2 3
2 4 4
2 4 2 5
5 5 0 5
1 5 5
1 3 5 3 5
0 2 3 4 5
3 2 1
0 3
1
5 5 5
3 1 0 3 4
1 2 5 4 4
1 4 5
1 0 4 5
3 3 2 4 4
1 4 2
0 2 0 3
4 5
5 1 5
3
5 2 3 2 4
3 1 1
4
2
1 4 3
5 5 5 1
0 1 3
3 4 3
1 0 2 4
3 0 2
0 5 4
4 3 1 5 2
3 5 1 3
5 5 3
0 2 2 2 2
5 2 5
5 4 5
0 4 5 1
3 4 4 1 4
5 1 3
0
1 3 4
1 5 3
0 0 0 5 3
2 1 2
2 1 3
3
2 1 3
3 3 4
1 5 3
5 1 2 2
1 1 2
3
1 5
2 2 3
1 1
1 2 4
5 4 4
5 2 2 5
4 4 4 5
5 3 4
5 5 5
5 2 4 4
5 1 3
2
3 3 1
2 3 4
3 0 1
2 3 1 0
0 3 2
2 0 5
5 3
0 5 3
1 4 3 4 2
3 5 1
2 3 2
3 0 1
1 2
2 2 4
2 2
0 2 1 1
5 2 5
0 2
2 5 4 0 1
1 1 4
3
2 1 2 4
4 1 3
5
5 4 3
1 1 4
3
3 4 2 4
4 5 1
4 4 4 5 3
3
5 4 2
3 3 4 3
0 0
3 5 2
0 5 4 1 0
3 2
3 1 3
0
5 2 1
1 1 2
3
5 0
2 4 1
3 1 0 3
1
4 1 2
5
4 4
0 1 2
4
0 4
2 5 2
1 2 3 0 1
4 0
1 1 4
1
0 4 1 5
3 3 5
5 3 5
4 4 5 1 4
0 2 2
1 3
1 2
2 4 2
3 1 3 2
2 0
4 1 4
2
2 4 3 2
1 4 2
5 4 5 0
0 4
4 2 2
1 3
4 0
5 2 1
5 2
5
5 4 5
0 5 1 1
4 2 0 5 1
0 2 5
1 0
5 3 5 2 1
1 3 3
4 4 0
3 3 5
5 1 4
2
5 0 5 5

`

type testCase struct {
	k int
	a []int
	b []int
}

func solveCase(tc testCase) ([]int, bool) {
	res := make([]int, 0, len(tc.a)+len(tc.b))
	i, j := 0, 0
	lines := tc.k
	for i < len(tc.a) || j < len(tc.b) {
		if i < len(tc.a) && tc.a[i] == 0 {
			res = append(res, 0)
			lines++
			i++
			continue
		}
		if j < len(tc.b) && tc.b[j] == 0 {
			res = append(res, 0)
			lines++
			j++
			continue
		}
		if i < len(tc.a) && tc.a[i] <= lines {
			res = append(res, tc.a[i])
			i++
			continue
		}
		if j < len(tc.b) && tc.b[j] <= lines {
			res = append(res, tc.b[j])
			j++
			continue
		}
		return nil, false
	}
	return res, true
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+2 > len(fields) {
			return nil, fmt.Errorf("case %d: unexpected end of data", i+1)
		}
		k, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse k: %v", i+1, err)
		}
		n, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[pos+2])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %v", i+1, err)
		}
		pos += 3
		if pos+n+m > len(fields) {
			return nil, fmt.Errorf("case %d: not enough elements", i+1)
		}
		tc := testCase{k: k, a: make([]int, n), b: make([]int, m)}
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a[%d]: %v", i+1, j, err)
			}
			tc.a[j] = val
		}
		pos += n
		for j := 0; j < m; j++ {
			val, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse b[%d]: %v", i+1, j, err)
			}
			tc.b[j] = val
		}
		pos += m
		cases = append(cases, tc)
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("trailing data after parsing")
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

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

	var inputBuilder strings.Builder
	inputBuilder.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		inputBuilder.WriteString(fmt.Sprintf("%d %d %d\n", tc.k, len(tc.a), len(tc.b)))
		for i, v := range tc.a {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.Itoa(v))
		}
		inputBuilder.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.Itoa(v))
		}
		inputBuilder.WriteByte('\n')
	}

	gotOut, err := runCandidate(bin, inputBuilder.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimRight(gotOut, "\n"), "\n")
	if len(lines) != len(cases) {
		fmt.Printf("expected %d lines of output, got %d\n", len(cases), len(lines))
		os.Exit(1)
	}

	for i, tc := range cases {
		expectedAns, ok := solveCase(tc)
		var expected string
		if !ok {
			expected = "-1"
		} else {
			strs := make([]string, len(expectedAns))
			for idx, v := range expectedAns {
				strs[idx] = strconv.Itoa(v)
			}
			expected = strings.Join(strs, " ")
		}
		got := strings.TrimSpace(lines[i])
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %q\ngot: %q\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
