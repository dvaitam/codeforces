package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	limit  int64 = 1_000_000_000
	infLCM int64 = limit + 1
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lcmCap(a, b int64) int64 {
	if a == infLCM || b == infLCM {
		return infLCM
	}
	g := gcd(a, b)
	a /= g
	res := a * b
	if res > infLCM {
		return infLCM
	}
	return res
}

func solve(arr []int64) int {
	present := make(map[int64]bool)
	for _, v := range arr {
		present[v] = true
	}

	dp := make(map[int64]int)
	for _, x := range arr {
		next := make(map[int64]int)
		if next[x] < 1 {
			next[x] = 1
		}
		for l, lLen := range dp {
			nl := lcmCap(l, x)
			if lLen+1 > next[nl] {
				next[nl] = lLen + 1
			}
		}
		for l, lLen := range next {
			if lLen > dp[l] {
				dp[l] = lLen
			}
		}
	}

	ans := 0
	for l, lLen := range dp {
		if l > limit || !present[l] {
			if lLen > ans {
				ans = lLen
			}
		}
	}
	return ans
}

const testcasesC = `1 3
1 12
2 10 9
5 7 20 2 19 6
4 13 17 12 18
4 17 9 2 1
3 15 11 13
4 17 6 18 6
2 8 1
2 11 6
2 17 17
3 17 18 6
4 14 17 12 19
3 12 15 6
4 15 17 8 16
3 16 17 17
3 15 15 12
5 18 15 16 8 11
6 6 20 9 16 10 10
6 17 18 17 17 20 19
4 10 7 16 17
3 20 3 11
6 1 7 4 2 19 2
3 19 8 4
5 5 9 8 7 2
4 2 2 12 12
2 8 1
1 4
1 1
1 1
3 9 5 6
6 6 17 1 13 19 2
2 5 2
1 12
5 4 10 11 16 1
3 15 18 20
6 2 9 13 20 5 16
2 3 11
1 1
4 5 17 19 13
4 17 11 5 11
3 9 20 14
6 1 18 5 2 9 2
2 6 6
1 15
6 8 17 2 8 8 15
1 9
1 19
2 20 20
6 12 9 14 9 17 1
2 2 13
4 6 4 17 3
2 4 4
1 6
2 4 7
1 17
6 15 15 10 18 13 7
6 7 14 14 17 1 19
5 2 14 17 19 6
1 16
3 1 17 4
5 12 10 12 10 1
6 14 4 4 10 7 1
4 2 14 16 15
2 19 20
1 1
3 1 12 10
6 3 8 16 7 4 19
3 13 15 5
3 13 4 9
1 4
1 20
3 13 7 4
1 20
6 16 2 16 10 12 15
2 12 9
4 17 16 14 16
6 10 13 8 6 16 20
3 18 14 3
5 19 4 3 12 6
5 5 14 3 3 2
2 10 13
2 11 15
2 17 10
1 5
5 14 4 11 17 8
6 17 9 6 6 15 8
4 12 19 5 15
4 1 20 13 6
4 17 2 16 9
4 9 14 16 12
5 11 3 8 18 20
2 13 13
6 1 11 15 17 15 6
1 1
4 7 19 20 13
2 4 13
5 7 9 19 19 7
4 20 5 1 20
6 14 16 9 17 19 6`

type testCase struct {
	arr []int64
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(testcasesC, "\n")
	var tests []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n+1, len(fields))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", idx+1, err)
			}
			arr[i] = v
		}
		tests = append(tests, testCase{arr: arr})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(tc.arr)))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solve(tc.arr)
		input := buildInput(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(want) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %d\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
