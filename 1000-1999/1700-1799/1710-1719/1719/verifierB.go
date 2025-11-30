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
	n int
	k int
}

func runProg(bin string, input string) (string, error) {
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

// Expected result for a single testcase using the logic from 1719B.go.
func expected(n, k int) string {
	k %= 4
	if k == 0 {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	if k != 2 {
		for i := 1; i <= n; i += 2 {
			fmt.Fprintf(&sb, "%d %d\n", i, i+1)
		}
	} else {
		flip := true
		for i := 1; i <= n; i += 2 {
			if flip {
				fmt.Fprintf(&sb, "%d %d\n", i+1, i)
			} else {
				fmt.Fprintf(&sb, "%d %d\n", i, i+1)
			}
			flip = !flip
		}
	}
	return strings.TrimSpace(sb.String())
}

// Embedded copy of testcasesB.txt.
const testcaseData = `
2 1
4 2
6 3
8 4
10 5
12 6
14 7
16 8
18 9
20 10
22 11
24 12
26 13
28 14
30 15
32 16
34 17
36 18
38 19
40 20
42 21
44 22
46 23
48 24
50 25
52 26
54 27
56 28
58 29
60 30
62 31
64 32
66 33
68 34
70 35
72 36
74 37
76 38
78 39
80 40
82 41
84 42
86 43
88 44
90 45
92 46
94 47
96 48
98 49
100 50
102 51
104 52
106 53
108 54
110 55
112 56
114 57
116 58
118 59
120 60
122 61
124 62
126 63
128 64
130 65
132 66
134 67
136 68
138 69
140 70
142 71
144 72
146 73
148 74
150 75
152 76
154 77
156 78
158 79
160 80
162 81
164 82
166 83
168 84
170 85
172 86
174 87
176 88
178 89
180 90
182 91
184 92
186 93
188 94
190 95
192 96
194 97
196 98
198 99
200 100
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
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", i+1, len(fields))
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := expected(tc.n, tc.k)
		if strings.Fields(got) == nil {
			fmt.Print("")
		}
		if strings.Join(strings.Fields(expect), " ") != strings.Join(strings.Fields(got), " ") {
			fmt.Printf("case %d failed:\nexpected:\n%s\n\ngot:\n%s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
