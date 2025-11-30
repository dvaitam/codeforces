package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `4 0 -1 -5 3 -5 -2 0 -4
2 3 0 -2 -2
3 5 -1 -1 3 1 -1
4 0 -2 -5 -1 3 -4 -5 2
4 2 -5 1 2 2 2 -4 -4
1 -2 -4
2 1 -2 2 4
1 1 3
4 -5 -3 -2 2 -2 -3 -1 0
3 1 -4 3 -1 4 3
2 -1 2 3 4
4 3 5 -1 -1 -2 -5 -4 4
1 -3 1
2 -2 -1 5 -5
4 -5 -4 1 5 -1 -4 4 0
2 5 3 5 -1
2 -2 -4 3 -1
3 -2 0 5 2 -1 4
2 -3 -5 3 3
3 0 4 5 -5 -3 1
2 -3 3 -4 -3
2 2 4 -2 -2
2 -2 1 0 4
2 5 2 -4 4
1 3 4
3 2 2 -1 -5 -2 3
2 5 2 2 3
3 -4 -1 -3 4 1 -2
3 -1 1 -5 -2 -5 0
2 0 2 5 5
2 -1 0 5 -3
3 -5 0 4 3 -5 5
2 0 -5 2 5
1 -5 -2
1 -5 -2
3 -4 -5 0 5 1 -3
2 2 1 -3 0
3 -3 5 0 1 1 -5
4 -1 3 3 5 2 -5 4 -4
4 1 -3 -5 3 -3 4 5 3
2 -4 0 -2 -3
2 -5 -3 5 3
2 -4 1 4 -4
4 -3 4 4 -5 -1 0 1 -5
1 2 -4
3 -1 5 -3 2 -2 3
3 -3 1 0 -1 2 1
1 -1 3
3 3 2 -5 3 4 3
3 5 -5 2 1 -4 1
3 2 -5 -5 -1 -5 -1
3 5 -2 3 3 0 1
3 -2 -4 4 0 -2 4
3 -3 -3 0 -5 4 -5
2 0 0 -1 5
3 0 2 1 4 1 -3
1 -3 4
1 2 -3
3 -5 2 5 5 -1 4
2 -4 3 1 -1
2 3 -3 -4 5
2 4 -4 3 5
4 1 -1 -1 -1 -5 1 -1 -1
3 0 -2 1 -3 -5 3
2 5 4 1 0
4 -5 3 1 5 4 -2 -5 0
2 5 -2 5 0
4 -5 -2 4 -2 -1 -3 1 -4
4 -2 2 3 -4 -2 -3 2 -4
4 5 1 -1 -1 1 0 4 0
1 -1 -5
4 -5 -1 -2 1 1 1 5 5
4 -5 4 2 0 4 -3 4 -1
3 -5 1 2 3 -3 -5
1 4 0
3 -5 -4 -2 -4 5 3
4 -5 0 -5 0 1 -3 5 -1
4 5 -3 4 -3 1 -1 3 -5
2 -3 -3 2 5
1 3 -5
4 -3 0 4 -4 -4 3 -3 -1
2 -1 0 -1 -1
4 -3 2 3 -3 -5 5 4 -3
1 0 -4
2 5 2 4 -2
4 3 -3 0 5 -3 2 3 -5
1 3 0
1 -4 -4
4 4 0 4 2 0 1 3 0
1 -3 0
1 -3 -3
1 0 4
2 -5 1 5 -5
3 1 -5 4 -3 0 -4
4 -5 2 0 4 4 -1 5 -1
4 1 -3 -5 2 -1 -2 1 -4
3 -4 -4 -5 0 -5 -3
4 4 5 -5 0 2 3 2 2
1 -5 3
4 -1 -5 5 3 -4 -4 0 0`

type testCase struct {
	input    string
	expected string
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solve(arr []int64) int64 {
	n := len(arr) / 2
	sortVals := append([]int64{}, arr...)
	// sort
	for i := 0; i < len(sortVals); i++ {
		for j := i + 1; j < len(sortVals); j++ {
			if sortVals[j] < sortVals[i] {
				sortVals[i], sortVals[j] = sortVals[j], sortVals[i]
			}
		}
	}
	if n == 1 {
		return absInt64(sortVals[0] - sortVals[1])
	}
	const inf int64 = 1 << 60
	ans := inf
	if n == 2 {
		var tmp int64
		for i := 0; i < 4; i++ {
			tmp += absInt64(sortVals[i] - 2)
		}
		if tmp < ans {
			ans = tmp
		}
	}
	var sum int64
	for _, v := range sortVals {
		sum += absInt64(v)
	}
	if sum < ans {
		ans = sum
	}
	if n%2 == 0 {
		var tmp int64
		for i := 0; i < 2*n-1; i++ {
			tmp += absInt64(sortVals[i] + 1)
		}
		tmp += absInt64(int64(n) - sortVals[2*n-1])
		if tmp < ans {
			ans = tmp
		}
	}
	return ans
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
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
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d: expected %d values got %d", idx+1, 1+2*n, len(fields))
		}
		arr := make([]int64, 2*n)
		for i := 0; i < 2*n; i++ {
			v, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value: %w", idx+1, err)
			}
			arr[i] = v
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(arr), 10),
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
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
