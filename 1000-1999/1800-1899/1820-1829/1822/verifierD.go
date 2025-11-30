package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
16
38
35
9
24
39
31
41
38
5
39
1
31
17
36
15
13
46
31
35
36
31
26
41
10
15
41
10
34
25
48
1
43
50
5
11
49
38
3
20
50
2
18
31
39
47
25
46
28
26
47
37
29
9
24
7
3
9
32
14
17
44
28
50
41
20
27
33
25
37
23
35
38
27
38
15
22
44
2
18
39
43
45
11
45
21
35
37
37
7
46
42
14
41
37
18
19
8
5
31`

type testCase struct {
	input    string
	expected string
}

func solveCase(n int) string {
	if n == 1 {
		return "1"
	}
	if n%2 == 1 {
		return "-1"
	}
	res := make([]int, 0, n)
	res = append(res, n)
	for i, j := n-1, 2; j <= n-2; i, j = i-2, j+2 {
		res = append(res, i)
		res = append(res, j)
	}
	res = append(res, 1)
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func loadCases() ([]testCase, error) {
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
		input := fmt.Sprintf("1\n%d\n", n)
		cases = append(cases, testCase{
			input:    input,
			expected: solveCase(n),
		})
	}
	return cases, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
