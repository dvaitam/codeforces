package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
8 -6 -4 2 -10 -1 -19 18 -17
3 2 -19 11
1 -19
4 -18 -20 -6 0
2 -17 2
7 -12 -7 8 7 -11 2 -1
3 1 6 4
1 6
5 14 9 -18 16 -13
7 4 -10 -20 12 -12 17 -11
2 1 -5
3 -5 -19 -10
3 -15 7 18
2 19 20
8 -11 19 18 -18 -4 1 4 -19
1 11
2 2 -2
3 9 -5 12
6 -10 5 1 -3 11 19
1 -1
5 15 10 -18 14 16
5 -18 9 5 -13 2
8 -17 -19 -3 -18 -4 -2 -7 13
6 4 -4 -7 -13 16 1
4 17 14 2 -10
3 1 -20 17
1 16
3 2 3 -2
5 0 11 5 18 7
3 -20 -11 16
1 8
3 1 -20 10
5 19 -8 -16 15 7
5 -9 13 -10 -16 20
3 17 -13 12
7 7 -3 -1 -2 -20 20 19
5 14 13 15 0 1
4 7 -11 -20 12
3 16 4 3
8 -18 15 6 -6 -19 3 13 -10
4 -18 -2 14 -17
2 -19 18
7 -17 13 -16 16 16 14 11
1 -20
9 15 7 15 -9 14 18 14 2 -10
7 1 9 1 -6 11 16 2
10 -15 -9 16 17 -2 18 1 -1 4 -9
3 4 7 3
10 -18 20 -2 6 11 -5 14 16 12 1
5 11 -10 1 -8 -12
10 -17 20 16 9 12 -8 -12 10 -17 16
7 16 -12 -2 15 19 -8 -6
8 12 -14 11 18 17 20 -8 -4
4 -9 19 -16 18
1 8
4 -20 0 13 0
4 11 10 2 6
3 -6 1 4
6 4 19 17 13 18 2
4 -7 -7 -10 20
2 15 8
1 6
9 20 -5 2 -6 15 -18 -13 0 15
9 18 -17 4 19 -20 15 -19 5 10
1 17
8 0 14 0 -8 -19 -6 5 12
5 -11 1 -4 -12 -12
1 -2
10 -7 3 -16 6 12 19 18 0 15 18
2 19 20
9 -3 10 12 4 11 14 7 20 17
1 -4
2 17 19
7 10 -17 -12 10 7 -13 6
8 5 20 20 13 13 20 19 3
7 17 0 17 -1 -10 6 19
3 -9 8 19
10 17 -14 2 3 -3 -2 -5 17 -7 6
10 4 -16 -14 -20 0 8 -4 15 -9 2
5 8 19 20 19 0
6 15 -4 13 -1 -8 -4
4 -10 -2 -5 19
7 18 -13 16 -5 4 20 12
9 7 16 6 13 10 4 17 8 12
6 -8 9 7 7 -15 -10
9 11 18 -20 -3 19 -6 10 16 7
7 19 13 -8 13 12 -18 9
9 0 5 11 -16 9 -4 1 3 -16
8 -6 2 11 3 3 4 -10 5
4 -7 3 15 -5
10 16 1 1 1 14 -15 -16 19 4 -15
10 -8 19 8 -10 -10 -16 18 10 -7 12
1 4
2 20 17
8 4 20 9 -13 -8 11 1 -9
8 -18 -8 -6 -2 14 -16 -1 -16
10 -5 20 -8 -5 -6 2 -10 20 -7 8
8 2 2 14 19 14 -12 17 4
5 -14 -8 5 17 3
2 7 19
10 5 16 12 15 8 1 20 20 -11 7
4 2 5 -18 -16
5 10 -11 -5 -10 -14
5 18 20 -13 10 15
`

type testCase struct {
	input    string
	expected string
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(nums []int64) string {
	n := len(nums)
	maxv := nums[0]
	for _, v := range nums {
		if v > maxv {
			maxv = v
		}
	}
	if n == 1 {
		return "1"
	}
	g := int64(0)
	for _, v := range nums {
		g = gcd(g, maxv-v)
	}
	sumOps := int64(0)
	set := make(map[int64]struct{}, n)
	for _, v := range nums {
		sumOps += (maxv - v) / g
		set[v] = struct{}{}
	}
	k := int64(1)
	for {
		candidate := maxv - k*g
		if _, ok := set[candidate]; !ok {
			break
		}
		k++
	}
	return strconv.FormatInt(sumOps+k, 10)
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n, len(fields)-1)
		}
		nums := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value: %w", idx+1, err)
			}
			nums[i] = v
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(nums),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
