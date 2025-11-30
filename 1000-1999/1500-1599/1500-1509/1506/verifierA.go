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
const testcasesRaw = `13 14 11
9 17 125
13 10 123
12 19 56
17 5 37
5 4 20
9 18 155
5 10 7
3 11 31
18 4 46
14 11 53
18 16 227
17 9 16
18 1 3
13 1 10
16 11 63
11 3 13
19 8 62
5 18 58
3 3 6
17 16 56
10 18 75
4 18 43
18 7 124
20 18 301
10 15 24
20 13 163
19 8 75
6 7 12
2 20 17
16 3 6
5 5 2
3 18 44
13 17 71
17 8 56
19 14 141
15 16 170
12 3 21
20 4 63
19 11 49
8 1 5
4 8 24
6 11 55
2 4 3
8 2 3
1 4 2
20 19 62
13 3 24
4 2 1
7 6 8
16 7 94
2 1 2
20 4 34
3 8 3
10 12 56
6 2 9
15 2 20
4 13 13
9 12 94
16 19 87
7 2 13
6 6 22
17 9 31
20 15 90
1 16 14
19 17 160
12 13 65
5 18 89
1 15 12
3 11 3
18 9 35
8 16 91
20 10 173
12 19 163
20 5 92
10 13 107
3 1 3
7 11 21
8 8 58
13 19 224
14 2 13
19 14 24
6 15 9
9 6 29
17 16 1
2 16 21
10 15 13
14 7 71
3 5 1
13 14 81
1 7 1
1 17 4
7 4 20
7 10 36
6 4 16
13 3 2
9 15 30
9 5 42
17 12 30
5 9 2`

type testCase struct {
	n int64
	m int64
	x int64
}

func solveCase(tc testCase) int64 {
	x := tc.x - 1
	row := x % tc.n
	col := x / tc.n
	return row*tc.m + col + 1
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("malformed embedded testcases")
	}
	res := make([]testCase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		n, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse n at case %d: %v", i/3+1, err)
		}
		m, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse m at case %d: %v", i/3+1, err)
		}
		x, err := strconv.ParseInt(fields[i+2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse x at case %d: %v", i/3+1, err)
		}
		res = append(res, testCase{n: n, m: m, x: x})
	}
	return res, nil
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		expected := strconv.FormatInt(solveCase(tc), 10)
		input := fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.m, tc.x)
		got, err := runCandidate(bin, input)
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
