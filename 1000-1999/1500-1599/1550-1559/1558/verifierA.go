package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	a int
	b int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
100
49 97
53 5
33 65
62 51
100 38
61 45
74 27
64 17
36 17
96 12
79 32
68 90
77 18
39 12
93 9
87 42
60 71
12 45
55 40
78 81
26 70
61 56
66 33
7 70
1 11
92 51
90 100
85 80
0 78
63 42
31 93
41 90
8 24
72 28
30 18
69 57
11 10
40 65
62 13
38 70
37 90
15 70
42 69
26 77
70 75
36 56
11 76
49 40
73 30
37 23
24 23
4 78
84 33
60 8
11 86
96 16
19 4
10 89
69 87
50 90
67 35
66 30
27 86
75 53
74 35
57 63
84 82
89 45
10 41
78 14
62 75
80 42
24 31
2 93
34 14
90 28
47 21
42 54
7 12
100 18
89 28
5 73
81 68
77 87
9 3
15 81
24 77
73 15
50 11
47 14
4 77
2 24
23 91
15 61
26 93
7 86
2 69
54 79
12 33
8 28
`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// possible mirrors 1558A.go.
func possible(a, b int) []int {
	n := a + b
	set := make(map[int]struct{})
	for start := 0; start < 2; start++ {
		var serveA, serveB int
		if start == 0 {
			serveA = (n + 1) / 2
			serveB = n / 2
		} else {
			serveA = n / 2
			serveB = (n + 1) / 2
		}
		diff := a - serveA
		L := max(0, -diff)
		U := min(serveA, serveB-diff)
		if L > U {
			continue
		}
		for w := L; w <= U; w++ {
			k := 2*w + diff
			set[k] = struct{}{}
		}
	}
	res := make([]int, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	sort.Ints(res)
	return res
}

func solve(tc testCase) string {
	ans := possible(tc.a, tc.b)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(ans)))
	if len(ans) > 0 {
		sb.WriteByte('\n')
		for i, v := range ans {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	}
	return sb.String()
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	if len(lines)-1 != t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(lines)-1)
	}
	res := make([]testCase, 0, t)
	for i := 1; i < len(lines); i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d: expected 2 ints", i)
		}
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad a: %v", i, err)
		}
		b, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad b: %v", i, err)
		}
		res = append(res, testCase{a: a, b: b})
	}
	return res, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("1\n%d %d\n", tc.a, tc.b)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
