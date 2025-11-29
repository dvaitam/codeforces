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
8 5 2
100 53 35
14 14 14
6 4 5
72 43 44
8 3 3
44 36 22
90 32 21
2 2 1
78 49 9
2 2 2
14 1 2
86 19 17
4 3 4
74 62 34
62 54 3
100 40 44
68 62 27
78 68 73
42 1 26
100 83 66
56 44 35
82 77 63
68 54 48
66 5 24
94 11 63
86 34 22
44 22 36
52 5 47
60 25 37
100 44 3
26 19 3
82 47 2
46 15 26
70 65 8
82 37 64
98 81 93
30 18 16
24 5 18
58 25 49
60 3 1
18 15 7
4 3 3
88 35 56
82 54 1
14 10 2
66 46 12
12 7 1
70 49 10
70 52 4
42 17 31
6 5 3
54 23 2
82 82 21
54 5 46
88 19 48
86 31 11
82 8 67
50 35 20
20 10 13
48 9 35
44 34 8
82 22 76
50 28 4
96 18 59
86 3 45
58 56 48
96 93 92
86 58 39
56 51 50
4 2 3
68 66 18
44 33 7
70 48 29
74 44 12
36 33 16
28 11 9
52 3 33
12 4 9
26 25 1
26 9 19
30 3 20
42 6 22
72 69 18
70 18 62
64 51 4
34 32 22
52 15 33
18 1 16
44 38 8
90 7 47
60 24 43
54 21 23
28 25 15
44 13 12
50 4 27
50 33 17
24 18 13
48 5 19
16 9 11
`

type testCase struct {
	n int
	x int
	y int
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
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		x, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{n: n, x: x, y: y})
	}
	return cases, nil
}

func isYes(tc testCase) bool {
	n := tc.n
	x := tc.x
	y := tc.y
	if n == 2 {
		return false
	}
	mid := n / 2
	if (x == mid && y == mid) ||
		(x == mid && y == mid+1) ||
		(x == mid+1 && y == mid) ||
		(x == mid+1 && y == mid+1) {
		return false
	}
	return true
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.x, tc.y)
		expected := "NO"
		if isYes(tc) {
			expected = "YES"
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
