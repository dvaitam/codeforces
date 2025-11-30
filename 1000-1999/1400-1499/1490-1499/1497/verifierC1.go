package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	g := gcd(a, b)
	if g == 0 {
		return 0
	}
	res := a / g * b
	return res
}

func checkSolution(n int, parts []int64) bool {
	if len(parts) != 3 {
		return false
	}
	sum := int64(0)
	for _, v := range parts {
		if v <= 0 {
			return false
		}
		sum += v
	}
	if sum != int64(n) {
		return false
	}
	cur := parts[0]
	for i := 1; i < 3; i++ {
		cur = lcm(cur, parts[i])
		if cur > int64(n)/2 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Println("failed to load embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(outFields) != 3 {
			fmt.Printf("case %d: expected 3 numbers got %d\n", i+1, len(outFields))
			os.Exit(1)
		}
		parts := make([]int64, 3)
		for j, s := range outFields {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Printf("case %d: invalid integer %q\n", i+1, s)
				os.Exit(1)
			}
			parts[j] = v
		}
		if !checkSolution(tc.n, parts) {
			fmt.Printf("case %d failed validation\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

type testCase struct {
	n int
	k int
}

// Embedded copy of testcasesC1.txt (one "n k" pair per line).
const testcaseData = `
17 3
26 3
24 3
95 3
46 3
191 3
174 3
81 3
67 3
158 3
57 3
158 3
12 3
151 3
177 3
43 3
113 3
166 3
103 3
188 3
133 3
98 3
142 3
116 3
131 3
71 3
12 3
10 3
96 3
122 3
84 3
100 3
111 3
137 3
45 3
146 3
48 3
63 3
62 3
9 3
48 3
86 3
47 3
37 3
133 3
133 3
95 3
134 3
175 3
146 3
49 3
117 3
109 3
191 3
137 3
198 3
96 3
154 3
93 3
95 3
117 3
44 3
196 3
105 3
186 3
192 3
121 3
170 3
138 3
66 3
128 3
74 3
130 3
131 3
134 3
93 3
172 3
119 3
121 3
92 3
148 3
188 3
145 3
188 3
119 3
127 3
171 3
59 3
86 3
182 3
45 3
160 3
71 3
200 3
125 3
82 3
80 3
183 3
132 3
146 3
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields got %d", i+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %v", i+1, err)
		}
		tests = append(tests, testCase{n: n, k: k})
	}
	return tests, nil
}
