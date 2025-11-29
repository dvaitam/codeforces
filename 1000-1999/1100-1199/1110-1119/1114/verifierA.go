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
12 13 1 8 16 15
12 9 15 11 18 6
16 4 9 4 24 3
19 8 17 22 25 19
4 9 3 23 2 21
10 15 17 3 11 13
10 19 20 6 17 15
14 16 8 1 25 17
0 2 12 22 25 21
20 0 19 15 10 7
10 2 6 18 7 7
4 17 14 2 2 10
16 15 3 9 17 9
3 17 10 17 6 25
19 17 18 9 14 2
19 12 10 18 7 9
5 6 5 1 19 21
8 15 2 2 21 24
4 4 1 2 22 17
12 16 8 16 25 7
6 18 13 18 8 14
15 20 11 2 10 19
3 15 18 20 10 6
7 0 8 3 22 7
11 5 10 13 1 3
4 7 1 18 20 17
19 2 0 3 20 6
19 18 3 12 2 11
3 1 19 0 6 5
3 15 6 23 25 1
0 17 13 19 3 8
2 7 2 20 9 11
13 5 1 16 14 1
19 3 12 6 8 11
15 18 5 22 21 6
1 5 5 10 16 8
3 19 14 21 5 0
15 13 18 16 9 20
11 12 8 4 17 22
0 14 2 10 23 1
17 8 4 7 24 15
11 19 9 21 11 18
20 19 4 22 9 12
13 20 2 0 19 6
10 5 7 7 20 14
12 18 13 1 12 22
18 13 1 5 14 2
8 5 14 16 15 17
19 0 1 15 10 9
14 1 13 6 17 20
2 4 0 12 21 13
10 0 6 0 22 24
0 16 19 3 6 3
19 20 6 9 8 22
5 3 15 12 20 2
0 8 14 25 25 3
8 4 20 16 20 20
11 3 4 8 0 1
1 6 8 17 10 11
18 1 19 20 15 22
20 14 20 13 11 17
5 6 12 18 9 0
4 4 8 10 10 25
11 2 10 24 19 1
1 8 5 4 18 9
11 12 17 4 9 3
15 7 1 9 5 16
2 9 12 10 9 13
3 3 17 15 15 10
10 3 15 3 22 15
13 1 9 10 23 21
4 5 20 18 12 25
20 2 2 25 2 6
7 1 12 0 3 12
17 16 9 14 15 25
18 6 13 2 11 7
8 18 5 13 6 11
3 2 0 16 14 24
6 3 15 12 8 6
20 1 6 19 4 3
6 14 12 11 17 4
3 19 15 4 18 12
20 13 16 15 21 10
15 15 20 21 6 17
19 7 0 10 22 23
10 10 1 16 4 8
19 4 12 18 9 22
15 2 2 16 1 2
7 4 1 9 0 24
14 10 5 25 4 20
14 11 16 12 16 16
1 18 2 21 25 25
16 19 2 23 13 24
6 9 17 19 13 15
12 19 18 7 25 0
0 5 9 16 18 8
10 2 15 8 9 24
13 12 12 1 5 20
4 7 9 23 10 1
1 15 13 4 15 19
`

type testCase struct {
	x int64
	y int64
	z int64
	a int64
	b int64
	c int64
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
		if len(parts) != 6 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		vals := make([]int64, 6)
		for i, p := range parts {
			v, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				return nil, err
			}
			vals[i] = v
		}
		cases = append(cases, testCase{x: vals[0], y: vals[1], z: vals[2], a: vals[3], b: vals[4], c: vals[5]})
	}
	return cases, nil
}

func expected(tc testCase) string {
	if tc.a < tc.x {
		return "NO"
	}
	if tc.a+tc.b < tc.x+tc.y {
		return "NO"
	}
	if tc.a+tc.b+tc.c < tc.x+tc.y+tc.z {
		return "NO"
	}
	return "YES"
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := fmt.Sprintf("%d %d %d %d %d %d\n", tc.x, tc.y, tc.z, tc.a, tc.b, tc.c)
		expect := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
