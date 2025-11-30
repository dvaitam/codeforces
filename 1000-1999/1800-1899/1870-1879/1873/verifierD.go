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
100
8 3
WWBBWWBB
16 16
WBBBWBBBBWBWWWWW
19 15
BWBBBWBWWWWWWWBWBWB
11 9
BBWWBBWWBWB
14 3
BWWWBBBWWWBBWB
3 1
BBW
10 10
WBBWWWBWWW
17 13
BWWBWWWWWBWWBWBBW
15 6
WWWBBBWWWWWBWBW
12 10
WWWBBBWBWBWB
14 11
BBWWBWBBWBWWBB
2 1
WB
9 6
BWBWWWWWW
10 7
WBWBBBWWBB
15 14
WWBBWBBBBBWBBWB
6 5
WBBWBB
11 3
WWBWBBBBBBW
5 1
WWBWB
9 9
WBWWBBBBB
2 1
WB
3 3
WWB
11 2
WWWWWWBWWBB
18 1
WBBBWWWBWBWWWWWWBB
7 5
WBWBWBB
11 6
WBWBWBWBBWB
8 7
BBBWBBWB
1 1
W
10 10
WBBBBBBWBB
8 3
BWWWWWWW
3 2
BWB
14 12
WBWBBBWWWBBWBB
19 17
WBWWWBWBWWWBBWBWWBW
16 7
WBWBWBWBBBBWBBWB
6 5
BWBBWB
14 2
BBBWWBWWWBBWWW
13 8
BBWWBWBWBWBWW
4 3
WWWW
4 3
BWWB
10 10
BBBWWBBBWW
18 10
BWWWBBBBWWBWWWBBBW
17 9
BWWWWBWWBBWBBBBBW
13 12
BBWBBBBWWWBWW
13 6
WWWWBBBBBWBWB
7 7
BBWBWBB
15 13
BBBWBBWBWWWWBBB
3 3
WBB
7 4
BWWBBWB
7 6
BWWWBWB
16 10
BWWWWBWBWWWWWWBB
7 3
WBWBBBB
6 5
WBBBBB
15 8
BBBWWWWBBBWBWWB
8 4
BWWBBWBW
4 1
WWWW
16 8
WBWWWBBWBBBBBBBB
16 2
BWBWWWWBWWWBWBWB
17 13
WWBWWBWWWBWWWWBWW
2 1
BB
15 2
BWBBBWWWWBBBWBW
19 10
BWBBBBBBBBBWBBWWBWB
20 3
BBBBBWBBWWBBWWWWBWWW
7 1
BWWBWWB
16 3
WWBBBWBWBWWBWBWB
7 6
BWWWBWW
14 8
BWBBBBBWBWBBBW
16 14
BWWBWBBWWWWBBBWW
13 13
WBWWWWBWWWBWW
12 5
BBBBBBBBBBWW
12 9
BWBBBWWWBBWW
9 4
WWWWBBWBB
19 8
BWWBWBBBBWBWWWWBWBB
5 5
BWWWW
9 5
WWBWBBWWB
15 15
WBWBWBWWBWBWBWW
7 3
WBWWWBW
7 1
WWWWBBW
9 7
BWBBBBBBW
3 2
WWW
9 5
BBBWBBBBW
13 11
BBWBBBWBBBBBB
13 13
BWWBBWBBBBWBW
12 10
BBBWBWBBWWBW
16 1
WWWWWWBWWBWBWBWW
7 2
BBWBBBB
4 4
WWBW
18 15
WBBWWWWWWWWWWBBBWB
10 7
WWBWWBWBBW
14 9
WBBWWBBWBBWBWW
18 12
BWWWWWWBWBWBBWBBBB
7 3
WWBBBWW
12 12
BWWBWBWBWWBB
6 5
BWBWWB
16 6
WWBWBWWWWBBBWWBB
11 8
WBWWBWBWBBB
17 17
BBBBWWWWWBBWWBWWW
19 17
BWBWWBWBBBWBWBWBBWW
8 7
BBWWWBBW
9 6
WBWBBBWWB
19 6
BWBBWBBWWBWWWBBWWBB
9 9
WWBWWWWWB
`

func solve(n, k int, s string) int {
	ops := 0
	for i := 0; i < n; {
		if s[i] == 'B' {
			ops++
			i += k
		} else {
			i++
		}
	}
	return ops
}

type testCase struct {
	input    string
	expected string
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+2 > len(fields) {
			return nil, fmt.Errorf("case %d: missing n/k/s", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		k, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad k: %w", caseIdx+1, err)
		}
		s := fields[pos+2]
		pos += 3
		if len(s) != n {
			return nil, fmt.Errorf("case %d: expected string len %d got %d", caseIdx+1, n, len(s))
		}
		input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
		cases = append(cases, testCase{
			input:    input,
			expected: strconv.Itoa(solve(n, k, s)),
		})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra data after parsing testcases")
	}
	return cases, nil
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
		fmt.Println("usage: verifierD /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
