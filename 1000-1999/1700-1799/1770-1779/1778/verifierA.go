package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

const testcaseData = `100
8
1 -1 1 1 1 1 1 1
5
-1 1 -1 -1 1
10
-1 1 -1 -1 1 1 -1 1 1 1
5
1 1 1 -1 -1
3
1 -1 1
7
-1 1 -1 -1 -1 -1 -1
10
1 -1 -1 1 1 -1 1 1 -1 1
10
-1 1 1 -1 1 1 -1 1 -1 -1
4
-1 1 1 -1
3
-1 -1 -1
3
1 1 -1
5
1 1 1 1 1
3
1 -1 1
7
-1 -1 -1 1 -1 -1 1
4
1 1 -1 -1
4
-1 -1 -1 -1
3
-1 -1 1
3
1 -1 -1
2
-1 -1
3
1 -1 -1
2
1 -1
6
-1 -1 -1 1 1 1
4
-1 1 -1 -1
8
-1 1 1 1 -1 -1 -1 -1
4
1 1 -1 1
4
-1 1 1 1
7
1 1 -1 -1 1 -1 1
2
1 -1
5
1 1 1 1 -1
6
1 1 -1 -1 -1 1
4
-1 -1 1 1
8
-1 1 1 -1 -1 1 -1 1
4
1 1 -1 -1
9
1 1 1 -1 1 -1 -1 -1 -1
8
1 1 -1 -1 -1 -1 -1 -1
3
-1 1 1
4
-1 1 1 -1
2
1 1
3
1 -1 1
3
-1 1 -1
2
-1 -1
6
1 1 -1 1 1 1
7
-1 -1 1 1 -1 -1 -1
6
1 1 1 -1 1 -1
2
1 -1
4
1 1 1 -1
6
-1 1 -1 -1 1 -1
10
-1 1 1 1 1 1 -1 -1 1 1
7
1 -1 1 -1 1 1 -1
6
1 -1 -1 1 -1 -1
3
-1 -1 -1
8
-1 -1 1 1 1 1 -1 1
3
1 -1 1
4
1 -1 1 -1
3
-1 1 -1
3
1 1 1
5
-1 -1 -1 -1 -1
9
1 1 -1 -1 1 -1 1 1 1
7
1 1 -1 -1 -1 1 1
7
-1 -1 1 -1 1 1 1
3
-1 -1 -1
5
-1 -1 1 -1 1
7
-1 -1 1 1 1 -1 -1
10
-1 1 -1 1 1 1 1 -1 -1 -1
4
1 1 1 -1
9
1 1 1 1 1 -1 -1 -1 -1
6
1 -1 -1 1 1 -1
9
-1 -1 1 1 -1 1 1 1 -1
3
1 -1 -1
2
-1 1
3
1 -1 1
9
-1 -1 1 1 -1 1 -1 1 1
3
-1 -1 1
9
1 -1 1 1 -1 -1 -1 1 1
2
-1 -1
3
1 1 -1
4
1 1 -1 -1
8
-1 -1 1 -1 -1 1 1 1
2
1 1
6
-1 1 -1 -1 1 -1
7
1 -1 -1 -1 1 -1 -1
5
-1 -1 -1 -1 1
7
-1 -1 -1 -1 1 -1 1
7
1 -1 -1 -1 1 1 1
6
1 1 -1 -1 -1 1
4
1 -1 -1 -1
7
-1 -1 -1 1 1 1 -1
4
-1 1 -1 -1
4
1 1 1 -1
8
1 1 -1 -1 1 1 -1 1
8
-1 1 1 1 -1 1 1 1
3
-1 -1 1
3
-1 1 1
4
1 1 -1 -1
9
1 -1 1 1 -1 1 -1 1 -1
9
1 -1 -1 1 -1 1 -1 1 1
3
1 -1 1
6
-1 -1 1 -1 -1 1
8
1 1 -1 1 1 1 -1 -1
2
-1 1`

func solve(tc []int) int {
	sum := 0
	for _, v := range tc {
		sum += v
	}
	hasNegPair := false
	hasOppPair := false
	for i := 0; i+1 < len(tc); i++ {
		if tc[i] == -1 && tc[i+1] == -1 {
			hasNegPair = true
		}
		if tc[i] != tc[i+1] {
			hasOppPair = true
		}
	}
	if hasNegPair {
		return sum + 4
	}
	if hasOppPair {
		return sum
	}
	return sum - 4
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d: not enough array values", caseNum+1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad array value: %w", caseNum+1, err)
			}
			arr[i] = v
		}
		pos += n

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.Itoa(solve(arr)),
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
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
