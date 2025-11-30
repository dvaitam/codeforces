package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxN = 600005

var (
	a    [maxN]int
	t    [maxN]int
	bArr [maxN]int
)

const testcasesF = `6 8 5 4 2 1 3 6
7 1 5 7 6 4 3 1 2
6 3 6 4 3 1 2 5
5 4 5 3 1 4 2
3 3 2 3 1
2 0 2 1
3 1 3 2 1
4 1 4 3 1 2
5 4 4 5 2 3 1
3 2 3 2 1
4 6 4 2 3 1
2 1 1 2
4 3 4 2 1 3
5 2 4 2 5 3 1
5 0 1 3 2 4 5
6 0 6 5 3 2 1 4
2 0 1 2
4 4 1 2 4 3
5 1 4 1 2 3 5
5 1 1 5 4 3 2
3 2 3 1 2
4 5 4 1 2 3
2 0 1 2
5 2 4 2 3 5 1
2 0 1 2
4 3 2 3 1 4
7 1 6 5 1 2 3 4 7
7 13 6 5 7 3 4 2 1
5 4 5 1 3 4 2
6 14 1 5 2 6 3 4
5 6 3 4 5 1 2
7 15 4 6 5 1 2 3 7
6 11 4 5 2 6 3 1
6 8 6 1 2 5 3 4
6 0 1 6 2 3 5 4
5 4 1 4 2 3 5
7 11 1 7 5 2 3 4 6
4 2 4 1 3 2
3 2 1 3 2
7 2 1 5 3 2 4 7 6
4 4 2 1 4 3
3 1 1 2 3
6 5 3 4 5 1 6 2
3 0 2 3 1
4 1 2 3 1 4
5 7 2 1 5 4 3
6 6 4 5 6 1 3 2
4 6 3 2 1 4
2 1 2 1
4 4 1 2 4 3
4 6 1 3 4 2
6 11 1 6 3 5 2 4
6 15 3 4 5 1 6 2
5 1 1 5 4 2 3
5 6 1 3 5 4 2
7 12 1 6 7 3 5 2 4
5 4 5 3 2 1 4
5 4 5 3 1 4 2
7 19 5 6 7 1 3 2 4
2 1 1 2
3 3 2 3 1
5 4 1 4 5 3 2
6 15 2 3 4 1 6 5
4 2 1 2 4 3
7 16 2 7 6 5 3 4 1
2 0 1 2
6 14 3 4 5 2 1 6
5 9 1 3 2 4 5
5 10 4 1 5 3 2
4 6 4 2 3 1
6 9 2 5 1 3 4 6
5 8 2 1 5 3 4
5 9 2 3 4 1 5
4 6 2 3 4 1
4 5 3 2 1 4
7 15 4 5 3 2 6 7 1
7 4 5 7 1 2 4 6 3
5 2 4 3 5 2 1
7 12 2 4 5 7 3 1 6
7 12 3 6 7 4 1 2 5
7 19 2 6 1 5 3 7 4
7 9 7 1 4 3 6 5 2
4 0 1 3 4 2
6 4 1 6 5 4 3 2
3 1 2 1 3
7 10 3 2 4 6 7 5 1
3 3 1 3 2
4 0 4 2 3 1
2 0 1 2
7 7 2 3 5 4 7 1 6
7 19 1 7 5 6 2 3 4
5 6 1 5 4 3 2
7 13 2 4 6 1 5 3 7
7 5 2 3 7 5 6 1 4
3 1 3 2 1
5 10 4 5 2 1 3
4 4 1 3 2 4
5 1 4 1 3 5 2
2 1 1 2
5 2 4 2 1 3 5
`

func solve(n int, k int64, p []int) (string, []int) {
	for i := 1; i <= n; i++ {
		a[p[i-1]] = i
		t[i] = 0
		bArr[i] = 0
	}
	add := func(x int) {
		for i := x; i <= n; i += i & -i {
			t[i]++
		}
	}
	ask := func(x int) int {
		s := 0
		for i := x; i > 0; i -= i & -i {
			s += t[i]
		}
		return s
	}
	tot := int64(n * (n - 1) / 2)
	sum := tot
	for val := 1; val <= n; val++ {
		sum -= int64(ask(a[val] - 1))
		add(a[val])
	}
	k -= sum
	if k%2 != 0 || k/2 < 0 || k/2 > (tot-sum) {
		return "NO", nil
	}
	k /= 2
	for i := 1; i <= n; i++ {
		t[i] = 0
	}
	for iVal := 1; iVal <= n; iVal++ {
		xCnt := ask(a[iVal] - 1)
		if k > int64(xCnt) {
			k -= int64(xCnt)
			add(a[iVal])
			continue
		}
		if k > 0 && k <= int64(xCnt) {
			id := 0
			for j := 1; j <= iVal; j++ {
				if a[j] < a[iVal] && k > 0 {
					id = j
					k--
				}
			}
			for j := 1; j <= id; j++ {
				bArr[j] = iVal + 1 - j
			}
			for j := id + 1; j <= iVal-1; j++ {
				bArr[j] = iVal - j
			}
			bArr[iVal] = iVal - id
			add(a[iVal])
			continue
		}
		bArr[iVal] = iVal
		add(a[iVal])
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = bArr[i]
	}
	return "YES", res
}

type testCase struct {
	input    string
	expected string
}

func parseCase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return testCase{}, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %w", err)
	}
	k, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return testCase{}, fmt.Errorf("bad k: %w", err)
	}
	if len(fields) != 2+n {
		return testCase{}, fmt.Errorf("expected %d permutation values, got %d", n, len(fields)-2)
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return testCase{}, fmt.Errorf("bad perm[%d]: %w", i, err)
		}
		perm[i] = v
	}
	status, arr := solve(n, k, perm)
	var exp string
	if status == "NO" {
		exp = "NO"
	} else {
		s := make([]string, len(arr))
		for i, v := range arr {
			s[i] = strconv.Itoa(v)
		}
		exp = "YES\n" + strings.Join(s, " ")
	}
	var input strings.Builder
	input.WriteString("1\n")
	fmt.Fprintf(&input, "%d %d\n", n, k)
	for i, v := range perm {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	return testCase{input: input.String(), expected: exp}, nil
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcasesF, "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseCase(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runCandidate(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Printf("case %d failed: expected %q got %q\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
