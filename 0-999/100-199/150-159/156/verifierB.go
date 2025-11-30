package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `7 6 +3 -4 -4 -5 +5 +3 +7
2 2 -1 -1
2 2 -2 +2
7 5 +5 -4 -1 +1 -6 +5 -7
6 1 -6 +2 +2 +5 -1 +3
9 7 +5 -2 -9 +9 -8 +7 -4 -3 +3
1 1 -1
2 2 +1 +1
9 6 -9 +4 -5 -8 -2 -2 -6 +4 +5
2 2 +2 +2
7 0 +7 +7 +1 +1 +6 +5 +4
2 1 +1 +1
3 0 -1 +3 +3
7 1 -1 +1 -3 -2 +5 -1 +6
7 3 -3 -7 +6 +7 +7 +7 +3
9 4 +8 +1 -7 -6 -5 +9 +8 +6 +9
5 1 +4 -5 -3 +3 -4
2 0 +2 +1
4 3 -4 +4 -1 +4
2 1 +2 -1
1 1 -1
8 0 -4 +3 +7 -6 +4 +1 +4 +4
5 2 +1 -4 +1 -4 +3
3 2 +1 -1 +1
4 2 -3 +4 -4 -2
4 3 -1 +2 -3 -3
2 1 +1 -1
3 2 -2 +2 +2
4 0 -2 +3 -3 -4
2 0 -2 -2
2 1 +2 -1
5 2 +2 -1 +1 +2 +4
1 0 -1
8 7 +7 +6 +5 +7 +6 +2 +8 +2
8 6 -4 +4 +2 +8 -6 +2 -3 -7
9 7 -8 -4 +1 -6 -1 +5 +7 -8 +2
9 0 +4 +1 -1 -6 +3 -6 -9 +2 +7
4 2 -4 -2 +1 +3
9 9 -6 +8 -5 -7 -1 +3 +5 -1 +8
7 2 -7 +6 +7 -4 +5 -4 -1
2 1 +1 +1
6 0 -2 -4 +1 -5 -6 +6
10 4 +7 -2 +1 -8 -4 -6 +1 +8 -1 +4
2 2 -2 +1
7 6 +1 -1 +4 +6 +7 -4 -1
8 5 -2 -2 +6 +6 -3 +4 -2 +4
1 0 +1
5 2 +5 +2 +4 +4 -3
3 0 +2 -2 -2
9 5 +2 +9 -3 -3 +4 -9 +4 +5 -7
1 0 +1
2 2 +1 -2
9 6 +7 -6 +4 -6 +7 -1 -6 -4 -5
8 1 +2 -2 +8 -3 -7 +4 -6 +6
8 1 -4 -1 -8 +4 -2 -3 -8 +8
4 4 -3 +1 -1 +3
7 6 -1 -3 -7 +5 +1 +7 -3
9 5 +9 +3 +7 -5 +9 +9 +7 -5 -8
6 4 +2 +6 +4 -5 -2 -3
8 6 +8 -6 -1 -5 +8 +4 +6 -7
1 0 +1
1 1 +1
10 0 -5 +3 -4 +7 -6 -3 -7 -3 -3 -3
4 1 -3 -4 -4 +2
8 3 +7 +4 +3 -1 +4 -2 -6 -8
1 1 -1
5 5 -2 -3 +1 +2 -1
5 5 +2 +4 +5 -5 +4
4 2 +1 -3 -3 -1
9 0 +6 +4 -4 +5 -5 -5 -6 +1 -9
2 0 -2 -1
7 7 -4 +1 +2 +7 +4 +4 +7
7 6 +2 +4 +2 -3 -4 +5 -5
9 3 -8 -9 -5 +1 +2 +7 +4 -1 -1
2 1 -1 -1
9 4 +4 +9 -6 +6 -5 +3 +9 -6 +3
7 2 +5 +2 +7 -5 +2 +7 +7
7 5 +6 -1 +5 -4 -3 +2 +6
8 7 -2 -3 -4 -5 -1 +1 -4 -8
4 2 -2 -1 -1 +3
1 1 +1
4 0 +2 -1 +3 -2
4 3 -2 -3 +3 -4
1 1 -1
1 0 -1
3 0 +3 +1 -1
3 1 +1 +3 +2
10 1 -3 +5 -7 +1 -2 -5 +8 +9 -3 -6
5 3 -1 -5 -5 -1 -1
8 6 +7 -8 +1 -1 -5 +6 -5 +2
10 3 +9 -5 +9 +2 +10 +9 -7 -5 -1 -5
5 4 -3 +4 +1 +5 -3
8 0 -4 +6 +4 -8 +4 +5 +7 +8
4 3 +2 +4 +4 -3
5 3 -5 -1 -1 -1 -2
7 6 -7 -6 +5 -3 +5 -3 +4
8 8 +1 +6 -1 +4 +8 +6 +6 -3
8 1 -2 +3 +8 -5 -4 -2 +4 -5
8 8 +3 +8 +8 +2 +5 -5 -6 -2
8 8 +2 +5 +3 +4 +1 -8 -5 -4
8 1 +6 -5 -6 -6 +6 +5 +4 +4
8 2 -7 -3 +1 +8 +3 -5 +6 +8
8 8 +7 -7 +3 +1 +8 -7 -7 -7
8 0 -7 -4 +1 +7 +8 +8 +4 +3
8 8 +3 +3 +8 +4 +7 -3 +2 +3
8 7 -3 -2 +1 -1 +8 +4 +8 -8
8 4 -7 +1 -8 +2 +5 +1 -5 +2
8 2 +6 -7 +1 +3 +5 -5 +3 -2
8 8 +3 +1 -3 +2 +6 +8 -5 -5
8 5 +1 +7 +4 +8 -2 +2 +1 -7
8 3 +5 -3 -1 -7 +6 +7 -3 +1
8 6 +8 +5 -8 +2 +5 -1 -6 +3
8 3 -7 -5 -8 -2 +1 +7 +5 +5
8 2 +4 +3 -1 -5 -4 +2 +2 -2
8 8 +1 +4 +2 +5 +5 -4 -2 -5
8 3 +2 +1 -8 +7 +5 +5 -5 -8
8 6 -3 -4 +2 +6 -3 -3 -3 -2
8 6 -5 -7 -7 +1 +5 +5 +5 +5
8 3 -3 +2 -3 +1 -5 +8 +2 +6`

type testCase struct {
	n    int
	m    int
	sign []byte
	a    []int
}

func solveCase(tc testCase) []string {
	n, m := tc.n, tc.m
	sign := tc.sign
	a := tc.a
	countPlus := make([]int, n+1)
	countMinus := make([]int, n+1)
	totalMinus := 0
	for i := 0; i < n; i++ {
		if sign[i] == '+' {
			countPlus[a[i]]++
		} else {
			countMinus[a[i]]++
			totalMinus++
		}
	}
	inS := make([]bool, n+1)
	var sSize int
	for c := 1; c <= n; c++ {
		t := countPlus[c] + (totalMinus - countMinus[c])
		if t == m {
			inS[c] = true
			sSize++
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		x := a[i]
		if sign[i] == '+' {
			if inS[x] && sSize == 1 {
				res[i] = "Truth"
			} else if !inS[x] {
				res[i] = "Lie"
			} else {
				res[i] = "Not defined"
			}
		} else {
			if inS[x] && sSize == 1 {
				res[i] = "Lie"
			} else if !inS[x] {
				res[i] = "Truth"
			} else {
				res[i] = "Not defined"
			}
		}
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d statements got %d", idx+1, n, len(fields)-2)
		}
		sign := make([]byte, n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			s := fields[2+i]
			if len(s) < 2 || (s[0] != '+' && s[0] != '-') {
				return nil, fmt.Errorf("line %d: bad statement %q", idx+1, s)
			}
			sign[i] = s[0]
			val, err := strconv.Atoi(s[1:])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value: %v", idx+1, err)
			}
			a[i] = val
		}
		cases = append(cases, testCase{n: n, m: m, sign: sign, a: a})
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
	return strings.TrimSpace(out.String()), nil
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

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for idx, s := range tc.sign {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte(s)
			sb.WriteString(strconv.Itoa(tc.a[idx]))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != tc.n {
			fmt.Printf("case %d: expected %d lines, got %d\n", i+1, tc.n, len(gotLines))
			os.Exit(1)
		}
		for j, line := range gotLines {
			if strings.TrimSpace(line) != expected[j] {
				fmt.Printf("case %d line %d failed\nexpected: %s\ngot: %s\n", i+1, j+1, expected[j], strings.TrimSpace(line))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
