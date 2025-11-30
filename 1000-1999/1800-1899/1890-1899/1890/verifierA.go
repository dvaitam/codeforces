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
100
5 2 1 2 3 2
5 2 2 2 3 1
6 1 2 1 1 3 2
6 3 3 1 2 1 3
2 3 2
5 3 1 2 2 2
6 3 1 3 2 2 3
4 1 3 1 1
5 3 3 3 1 3
5 2 1 3 2 3
2 1 3
3 1 1 3
5 1 1 2 3 2
2 2 3
4 3 1 3 2
6 1 3 3 3 2 2
2 3 2
4 3 1 2 1
3 1 1 3
4 2 1 1 3
3 1 1 1
6 3 2 3 3 2 3
3 1 3 3
5 3 2 2 2 3
4 1 2 3 1
5 3 3 2 1 1
2 3 2
2 3 1
4 1 2 2 1
2 1 3
3 1 3 3
6 3 3 1 1 1 3
3 3 3 1
5 1 2 1 1 3
2 1 1
2 2 1
2 3 1
6 2 3 1 2 1 1
2 3 2
4 2 1 1 3
5 1 3 1 3 2
3 2 2 3
5 3 1 3 3 1
2 3 1
3 2 3 2
2 3 2
3 1 2 3
5 3 3 2 3 2
5 3 2 1 3 3
2 2 3
2 2 3
2 3 2
3 1 2 2
6 2 3 2 3 3 3
3 3 2 2
5 3 1 1 3 1
4 1 1 1 3
5 2 3 3 3 2
2 2 3
6 2 3 3 1 1 2
2 2 3
3 2 3 2
6 3 1 1 2 2 2
5 1 2 1 3 3
2 3 1
2 2 3
5 2 1 1 1 3
2 3 3
6 1 1 1 3 3 1
4 2 3 1 1
5 2 3 1 1 2
5 1 2 1 3 3
4 1 1 2 1
2 1 1
4 3 2 2 3
2 3 3
6 3 2 3 3 2 3
5 2 3 1 1 2
6 2 1 1 1 2 2
4 2 3 1 2
6 1 1 2 1 1 3
4 2 2 3 1
4 1 2 3 1
2 2 1
6 3 1 2 2 2 2
5 1 1 3 2 2
4 2 1 2 1
5 2 1 2 2 3
3 1 3 3
5 3 1 1 1 1
3 1 2 1
2 2 3
6 2 2 2 3 3 3
3 2 1 2
3 2 3 1
5 1 2 1 1 3
2 3 2
3 1 2 2
4 1 3 1 1
6 1 1 1 2 2 2
`

type testCase struct {
	input    string
	expected string
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(nums []int) string {
	freq := make(map[int]int)
	for _, x := range nums {
		freq[x]++
	}
	if len(freq) == 1 {
		return "Yes"
	}
	if len(freq) != 2 {
		return "No"
	}
	counts := make([]int, 0, 2)
	for _, c := range freq {
		counts = append(counts, c)
	}
	if absInt(counts[0]-counts[1]) <= 1 {
		return "Yes"
	}
	return "No"
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d: not enough numbers", caseIdx+1)
		}
		nums := make([]int, n)
		numStr := make([]string, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseIdx+1, err)
			}
			nums[i] = val
			numStr[i] = fields[pos+i]
		}
		pos += n
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(strings.Join(numStr, " "))
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solve(nums),
		})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra data after parsing")
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
		fmt.Println("usage: verifierA /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
