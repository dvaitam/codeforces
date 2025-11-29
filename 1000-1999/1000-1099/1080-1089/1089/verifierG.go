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
195 NO
154 NO
179 YES
107 YES
154 NO
28 NO
132 NO
78 NO
165 NO
17 YES
62 NO
182 NO
65 NO
103 YES
131 YES
178 NO
100 NO
35 NO
92 NO
54 NO
152 NO
25 NO
10 NO
117 NO
177 NO
7 YES
181 YES
10 NO
45 NO
171 NO
71 YES
146 NO
142 NO
98 NO
106 NO
52 NO
181 YES
60 NO
195 NO
102 NO
8 NO
19 YES
7 YES
103 YES
135 NO
92 NO
21 NO
186 NO
42 NO
155 NO
74 NO
192 NO
186 NO
122 NO
133 NO
184 NO
68 NO
102 NO
12 NO
42 NO
12 NO
178 NO
173 YES
59 YES
171 NO
173 YES
149 YES
99 NO
73 YES
47 YES
191 YES
116 NO
200 NO
15 NO
122 NO
123 NO
1 NO
89 YES
23 YES
12 NO
166 NO
55 NO
179 YES
90 NO
199 YES
179 YES
45 NO
192 NO
181 YES
44 NO
119 NO
160 NO
35 NO
142 NO
179 YES
98 NO
32 NO
87 NO
102 NO
11 YES
`

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := int64(3); i <= n/i; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

type testCase struct {
	n    int64
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
		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{n: n, want: parts[1]})
	}
	return cases, nil
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
		fmt.Fprintln(os.Stderr, "Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d\n", tc.n)
		expected := "NO"
		if isPrime(tc.n) {
			expected = "YES"
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected || got != tc.want {
			fmt.Printf("Case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(testcases))
}
