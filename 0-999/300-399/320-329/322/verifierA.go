package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `100
50 98
54 6
34 66
63 52
39 62
46 75
28 65
18 37
18 97
13 80
33 69
91 78
19 40
13 94
10 88
43 61
72 13
46 56
41 79
82 27
71 62
57 67
34 8
71 2
12 93
52 91
86 81
1 79
64 43
32 94
42 91
9 25
73 29
31 19
70 58
12 11
41 66
63 14
39 71
38 91
16 71
43 70
27 78
71 76
37 57
12 77
50 41
74 31
38 24
25 24
5 79
85 34
61 9
12 87
97 17
20 5
11 90
70 88
51 91
68 36
67 31
28 87
76 54
75 36
58 64
85 83
90 46
11 42
79 15
63 76
81 43
25 32
3 94
35 15
91 29
48 22
43 55
8 13
19 90
29 6
74 82
69 78
88 10
4 16
82 25
78 74
16 51
12 48
15 5
78 3
25 24
92 16
62 27
94 8
87 3
70 55
80 13
34 9
29 10
83 39`

type testCase struct {
	m int
	n int
}

// solveCase mirrors 322A.go.
func solveCase(m, n int) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, m+n-1)
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", 1, i)
	}
	for i := 2; i <= m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", i, 1)
	}
	return strings.TrimSpace(sb.String())
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Fields(testcasesRaw)
	if len(lines) < 1 {
		return nil, fmt.Errorf("no test count")
	}
	idx := 0
	t, err := strconv.Atoi(lines[idx])
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	idx++
	if len(lines) != 1+2*t {
		return nil, fmt.Errorf("expected %d pairs got %d values", t, (len(lines)-1)/2)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx+1 >= len(lines) {
			return nil, fmt.Errorf("case %d incomplete", i+1)
		}
		m, err := strconv.Atoi(lines[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %w", i+1, err)
		}
		n, err := strconv.Atoi(lines[idx+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %w", i+1, err)
		}
		idx += 2
		cases = append(cases, testCase{m: m, n: n})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		input := fmt.Sprintf("%d %d\n", tc.m, tc.n)
		expect := solveCase(tc.m, tc.n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
