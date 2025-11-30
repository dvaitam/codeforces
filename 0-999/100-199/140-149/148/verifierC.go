package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 148C.go.
func solve(n, m, k int) (string, error) {
	a := make([]int, n)
	a[0] = 1
	sum, mx, cnt := 1, 1, 1
	if m == 0 && k == 0 {
		out := make([]string, n)
		for i := 0; i < n; i++ {
			out[i] = "1"
		}
		return strings.Join(out, " "), nil
	}
	if m+k+1 >= n && k == 0 {
		return "-1", nil
	}
	if m+k+1 < n {
		a[1] = 1
		sum = 2
		cnt = 2
	}
	for k > 0 {
		a[cnt] = sum + 1
		if a[cnt] > mx {
			mx = a[cnt]
		}
		sum += a[cnt]
		cnt++
		k--
	}
	for m > 0 {
		a[cnt] = mx + 1
		mx = a[cnt]
		cnt++
		m--
	}
	for i := 0; i < n; i++ {
		if a[i] > 50000 {
			return "-1", nil
		}
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		v := a[i]
		if v == 0 {
			v = 1
		}
		out[i] = strconv.Itoa(v)
	}
	return strings.Join(out, " "), nil
}

type testCase struct {
	n, m, k int
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
8 2 2
47 5 9
33 6 1
75 5 13
82 12 11
70 14 8
71 0 11
41 12 13
68 5 5
31 7 0
23 10 5
18 11 5
58 13 11
76 11 11
58 5 12
92 14 7
63 8 15
65 11 14
60 11 14
63 7 10
90 5 8
99 15 9
39 13 9
94 6 15
66 11 2
44 0 6
96 3 1
74 1 8
76 7 3
97 4 8
32 6 1
55 1 1
47 11 5
32 0 2
15 2 0
44 0 11
17 5 5
67 0 12
76 1 7
20 1 0
45 3 9
44 15 0
40 14 1
34 12 4
61 7 2
85 10 3
85 14 4
75 12 15
66 10 4
44 8 8
78 13 0
90 4 1
33 1 4
21 5 3
59 7 1
32 7 14
86 8 2
30 11 8
88 13 8
68 0 4
46 12 13
15 2 7
14 3 0
24 7 3
28 0 14
59 9 12
28 6 13
55 0 1
54 5 3
85 15 11
52 3 11
89 11 9
30 13 3
40 6 0
58 1 13
82 15 14
27 2 0
37 0 11
40 2 7
97 15 6
83 11 12
39 11 12
33 3 3
50 10 12
89 3 0
80 15 1
93 15 9
46 14 4
48 8 15
68 15 13
63 9 12
30 5 15
77 8 13
90 2 3
86 11 5
19 13 2
12 1 4
38 12 7
91 10 14
23 9 3
20 13 3
43 7 8
22 5 14
91 7 12
46 4 14
57 0 12
95 5 12
66 1 15
36 12 8
91 13 15
47 10 2
98 7 6
52 12 0
41 14 14
84 5 3
91 12 6
78 12 6
54 12 6
96 6 15
79 4 0
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("malformed test data")
	}
	res := make([]testCase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("bad n at case %d: %v", i/3+1, err)
		}
		m, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("bad m at case %d: %v", i/3+1, err)
		}
		k, err := strconv.Atoi(fields[i+2])
		if err != nil {
			return nil, fmt.Errorf("bad k at case %d: %v", i/3+1, err)
		}
		res = append(res, testCase{n: n, m: m, k: k})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k)
		expected, err := solve(tc.n, tc.m, tc.k)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
