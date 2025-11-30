package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `
26 50 6
39 50 3
16 24 1
17 63 24
34 77 34
6 45 31
13 39 23
40 77 13
24 48 21
14 23 7
43 54 9
22 27 5
34 77 5
31 51 12
3 20 17
1 28 6
24 58 5
30 63 12
39 68 1
23 31 3
22 33 9
1 11 9
19 24 1
1 59 26
22 32 3
15 52 0
34 74 10
7 26 13
16 39 6
35 38 3
31 38 0
15 24 8
1 38 35
15 23 0
23 36 10
36 40 4
19 31 4
13 29 16
39 55 11
11 35 17
36 37 0
17 74 12
36 45 6
3 16 9
16 36 3
8 19 2
24 39 2
1 31 7
18 40 19
17 23 6
7 54 12
17 24 7
39 58 7
15 69 31
8 31 7
5 70 35
35 52 13
19 39 1
17 28 7
18 31 9
36 39 3
22 37 0
5 11 5
32 33 0
1 74 4
8 24 8
39 75 17
7 12 3
9 19 3
35 55 17
30 56 13
19 53 8
22 47 2
21 74 15
33 60 13
17 54 14
29 40 3
19 74 24
2 32 16
30 51 13
10 75 16
26 35 5
11 40 23
4 39 8
24 26 0
7 16 5
29 34 5
33 41 8
3 18 2
9 69 24
31 45 2
11 34 19
28 38 9
9 15 5
21 25 1
13 50 10
15 16 1
12 31 15
10 50 17
14 43 6
`

type testCase struct {
	l int64
	r int64
	k int64
}

func solveCase(tc testCase) string {
	if tc.l == tc.r {
		if tc.l > 1 {
			return "YES"
		}
		return "NO"
	}
	odd := (tc.r+1)/2 - tc.l/2
	if tc.k >= odd {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 numbers got %d", idx+1, len(fields))
		}
		l, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse l: %v", idx+1, err)
		}
		r, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse r: %v", idx+1, err)
		}
		k, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
		}
		cases = append(cases, testCase{l: l, r: r, k: k})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.l, tc.r, tc.k))

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
