package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
1 NO
2 YES
3 YES
4 YES
5 YES
6 NO
7 NO
8 NO
9 NO
10 NO
11 NO
12 YES
13 NO
14 NO
15 NO
16 NO
17 NO
18 NO
19 NO
20 NO
21 NO
22 NO
23 NO
24 NO
25 NO
26 NO
27 NO
28 NO
29 NO
30 YES
31 YES
32 NO
33 NO
34 NO
35 YES
36 NO
37 NO
38 NO
39 NO
40 NO
41 NO
42 NO
43 YES
44 NO
45 NO
46 YES
47 NO
48 NO
49 NO
50 NO
51 NO
52 YES
53 NO
54 NO
55 NO
56 NO
57 NO
58 NO
59 NO
60 NO
61 NO
62 NO
63 NO
64 YES
65 NO
66 NO
67 NO
68 NO
69 NO
70 NO
71 NO
72 NO
73 NO
74 NO
75 NO
76 NO
77 NO
78 NO
79 NO
80 NO
81 NO
82 NO
83 NO
84 NO
85 NO
86 YES
87 NO
88 NO
89 NO
90 NO
91 NO
92 NO
93 NO
94 NO
95 NO
96 NO
97 NO
98 NO
99 NO
100 NO
`

type testCase struct {
	a    int
	want string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesData, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{a: n, want: parts[1]})
	}
	return cases, nil
}

func isWinning(a int) bool {
	switch a {
	case 2, 3, 4, 5, 12, 30, 31, 35, 43, 46, 52, 64, 86:
		return true
	default:
		return false
	}
}

func run(bin string, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d\n", tc.a)
		expect := "NO"
		if isWinning(tc.a) {
			expect = "YES"
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expect {
			fmt.Printf("Case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(testcases))
}
