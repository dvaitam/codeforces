package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesRaw = `4 35 9 9 4 5 1 2 5
5 11 5 10 5 8 1 3 8 2 4
3 8 7 1 3 2 2
5 36 8 1 1 5 1 9 8 4 6
6 29 2 6 6 7 7 5 2 4 8 9 6
5 28 7 10 5 3 3 1 6 6 7
2 41 10 6 10
3 10 2 9 4 8 4
4 40 9 3 4 5 3 3 7
5 32 6 1 9 2 1 6 4 3 4
5 29 9 10 5 7 10 6 8 6 2
6 40 1 3 9 8 3 2 1 2 1 3 5
3 47 8 7 9 9 5
4 36 7 2 7 8 4 2 6
3 44 10 1 7 1 5
4 50 1 10 8 6 10 1 9
4 47 7 1 10 8 2 7 7
2 37 1 1 9
6 27 6 3 7 1 3 5 9 10 7 3 10
5 47 5 10 10 5 1 7 9 10 7
3 21 3 8 7 10 9
3 33 2 10 10 10 7
4 26 8 1 5 3 5 7 5
2 17 1 2 2
5 10 8 4 4 1 4 2 2 2 1
6 43 2 1 5 7 3 6 2 1 7 10 10
3 11 9 10 8 3 6
6 26 9 10 3 6 9 2 1 1 10 5 2
5 6 1 1 5 9 5 10 10 5 8
5 8 4 5 3 9 9 1 6 8 2
5 44 3 5 2 6 5 4 6 10 3
6 15 10 1 4 8 6 3 7 6 7 10 8
2 17 1 9 5
6 21 4 4 4 4 7 6 5 1 8 9 3
5 49 8 2 9 5 2 4 2 7 7
3 8 8 9 4 3 4
4 24 6 6 5 10 3 1 4
4 31 10 9 1 6 1 3 4
4 42 4 2 7 6 6 4 2
2 26 6 10 6
5 22 10 5 7 10 5 6 10 2 7
3 40 8 6 5 1 2
6 34 1 3 10 4 9 8 5 7 7 10 1
2 26 3 10 4
5 43 7 8 2 7 3 8 4 5 9
2 20 5 3 5
6 20 8 3 7 6 9 6 4 5 1 5 9
6 19 8 5 5 3 5 5 6 3 5 7 8
5 47 3 7 1 2 10 4 6 1 9
4 3 7 2 10 6 3 1 6
3 40 6 9 7 4 9
2 3 6 1 8
2 11 5 10 4
5 19 3 1 1 8 7 9 2 7 5
5 4 4 6 7 10 10 8 10 4 10
6 44 2 6 7 3 4 9 8 2 7 7 4
4 50 1 5 1 5 2 3 10
4 49 8 7 5 2 5 1 8
3 17 9 4 3 1 7
6 1 10 9 5 1 7 6 2 5 3 10 4
2 12 10 7 9
6 47 1 4 7 1 1 9 7 10 3 1 7
5 13 3 4 2 10 8 9 9 6 5
3 33 10 5 7 4 5
6 17 3 5 6 10 5 9 4 4 9 1 2
3 18 3 6 4 3 1
6 15 7 5 5 4 5 7 1 1 3 8 7
4 24 7 6 10 4 5 5 5
5 40 3 10 6 3 7 1 2 5 2
5 14 8 5 1 5 6 1 10 8 7
5 28 6 10 8 4 7 7 5 2 2
3 46 6 6 10 7 7
2 25 1 7 10
3 8 4 8 7 3 3
3 40 2 6 6 9 8
3 25 8 10 3 1 9
3 44 4 3 2 5 9
2 1 6 5 4
2 20 3 3 2
3 38 7 5 3 2 4
3 39 7 4 7 9 8
3 39 2 8 4 4 2
3 16 4 10 10 3 9
4 5 10 6 2 6 9 5 3
6 31 7 9 10 9 4 3 9 10 2 7 7
4 15 8 9 7 6 10 3 1
2 24 10 8 3
5 37 6 6 3 8 4 9 8 9 5
6 14 3 10 4 5 2 2 4 7 9 4 5
5 46 1 7 4 1 5 5 10 4 7
2 49 8 6 7
3 42 2 3 9 1 7
3 2 9 8 8 6 9
5 2 10 7 7 4 9 1 9 1 4
4 3 3 8 3 8 3 3 5
6 25 2 9 2 2 7 8 6 4 1 9 7
3 31 4 3 4 4 9
4 42 5 8 10 10 9 4 4
4 8 10 1 1 10 6 2 9
3 30 9 2 7 3 9`

type testCase struct {
	n int
	c int64
	a []int64
	b []int64
}

func solveCase(tc testCase) int64 {
	curDay := int64(0)
	curMoney := int64(0)
	ans := int64(1 << 60)
	for i := 0; i < tc.n; i++ {
		if curMoney >= tc.c {
			if curDay < ans {
				ans = curDay
			}
		} else {
			need := tc.c - curMoney
			days := curDay + (need+tc.a[i]-1)/tc.a[i]
			if days < ans {
				ans = days
			}
		}
		if i == tc.n-1 {
			break
		}
		if curMoney < tc.b[i] {
			need := tc.b[i] - curMoney
			d := (need + tc.a[i] - 1) / tc.a[i]
			curDay += d
			curMoney += d * tc.a[i]
		}
		curMoney -= tc.b[i]
		curDay++
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("case %d: malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		c, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse c: %v", idx+1, err)
		}
		expectedFields := 2*n + 1
		if len(fields) != expectedFields {
			return nil, fmt.Errorf("case %d: expected %d values got %d", idx+1, expectedFields, len(fields))
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[2+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a%d: %v", idx+1, i+1, err)
			}
			a[i] = val
		}
		b := make([]int64, n-1)
		for i := 0; i < n-1; i++ {
			val, err := strconv.ParseInt(fields[2+n+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse b%d: %v", idx+1, i+1, err)
			}
			b[i] = val
		}
		res = append(res, testCase{n: n, c: c, a: a, b: b})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		expected := strconv.FormatInt(solveCase(tc), 10)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.c))
		for idx, v := range tc.a {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.b {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
