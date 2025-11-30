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
83
39
102
167
13
19
138
25
94
150
15
130
55
10
23
112
108
18
62
24
142
109
16
145
32
58
162
161
150
16
148
150
102
13
57
12
143
35
75
108
143
161
169
63
80
92
160
168
147
10
10
52
123
46
32
152
133
170
121
58
92
26
88
165
161
38
20
84
3
11
53
11
77
87
56
160
100
29
89
120
68
19
70
73
102
163
78
80
59
131
92
88
154
146
33
43
15
100
92
16
`

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func solve(n int) string {
	if isPrime(n) {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	cases := make([]testCase, 0, len(fields))
	for _, f := range fields {
		input := f + "\n"
		expected := solve(mustAtoi(f))
		cases = append(cases, testCase{input: input, expected: expected})
	}
	return cases, nil
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
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
		fmt.Println("usage: verifierH /path/to/binary")
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
