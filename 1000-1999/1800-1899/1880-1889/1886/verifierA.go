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
50
98
54
6
34
66
63
52
39
62
46
75
28
65
18
37
18
97
13
80
33
69
91
78
19
40
13
94
10
88
43
61
72
13
46
56
41
79
82
27
71
62
57
67
34
8
71
2
12
93
52
91
86
81
1
79
64
43
32
94
42
91
9
25
73
29
31
19
70
58
12
11
41
66
63
14
39
71
38
91
16
71
43
70
27
78
71
76
37
57
12
77
50
41
74
31
38
24
25
24
`

type testCase struct {
	n      int
	input  string
	output string
}

func solveCase(n int) string {
	if n < 7 || n == 9 {
		return "NO"
	}
	if n%3 == 0 {
		return fmt.Sprintf("YES\n1 4 %d", n-5)
	}
	return fmt.Sprintf("YES\n1 2 %d", n-3)
}

func loadCases() ([]testCase, string, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, "", fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, "", fmt.Errorf("bad test count: %w", err)
	}
	if len(fields) != t+1 {
		return nil, "", fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
	}
	var inputBuilder strings.Builder
	inputBuilder.WriteString(fields[0])
	inputBuilder.WriteByte('\n')

	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, "", fmt.Errorf("bad n at case %d: %w", i+1, err)
		}
		inputBuilder.WriteString(fields[i+1])
		inputBuilder.WriteByte('\n')
		cases = append(cases, testCase{
			n:      n,
			output: solveCase(n),
		})
	}
	return cases, inputBuilder.String(), nil
}

func buildExpected(cases []testCase) string {
	var sb strings.Builder
	for idx, tc := range cases {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(tc.output)
	}
	return sb.String()
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
	cases, input, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	expected := buildExpected(cases)
	got, err := run(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != expected {
		fmt.Println("expected:\n" + expected)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
