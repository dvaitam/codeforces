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
100
50 48 26
6 2 4
63 25 58
39 30 22
75 27 64
18 9 4
97 12 79
33 9 19
13 11 1
88 42 60
72 12 45
56 20 39
82 26 70
62 28 55
67 33 7
71 1 11
93 51 90
86 80 0
79 63 42
32 20 4
25 18 7
31 25 30
19 17 14
12 1 5
66 62 13
39 35 18
91 15 70
43 34 13
78 70 75
37 28 5
77 49 40
74 30 37
24 6 5
5 4 2
61 4 5
87 16 19
5 0 4
88 50 67
36 33 15
28 21 18
54 37 17
58 31 42
83 45 10
42 39 7
63 37 40
43 12 15
3 2 1
15 11 3
48 10 21
55 52 3
13 12 2
90 28 5
74 68 9
4 0 1
78 73 15
51 5 23
15 0 9
3 0 0
92 15 61
27 23 25
8 0 6
80 12 33
9 3 1
83 38 44
56 11 3
65 59 5
77 12 50
26 8 11
94 60 72
22 21 6
99 7 86
21 5 10
68 32 15
77 56 22
2 1 1
73 65 39
84 45 49
85 32 19
72 1 58
95 10 42
95 5 69
36 8 15
98 61 45
79 36 45
76 16 39
50 47 26
84 10 0
77 24 42
21 7 7
82 57 48
91 86 72
54 2 25
90 72 53
99 84 90
6 1 3
9 4 2
58 33 56
63 58 35
78 0 4
64 41 39
`

type testCase struct {
	n int
	a int
	b int
}

func solveCase(tc testCase) int {
	count := 0
	for i := 1; i <= tc.n; i++ {
		if i-1 >= tc.a && tc.n-i <= tc.b {
			count++
		}
	}
	return count
}

func parseTestcases() ([]testCase, string, error) {
	fields := strings.Fields(testcasesData)
	if len(fields) == 0 {
		return nil, "", fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, "", err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+2 >= len(fields) {
			return nil, "", fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, _ := strconv.Atoi(fields[pos])
		a, _ := strconv.Atoi(fields[pos+1])
		b, _ := strconv.Atoi(fields[pos+2])
		pos += 3
		cases = append(cases, testCase{n: n, a: a, b: b})
	}
	return cases, strings.TrimSpace(testcasesData) + "\n", nil
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, _, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for i, tc := range testcases {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.a, tc.b)
		want := strconv.Itoa(solveCase(tc))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
