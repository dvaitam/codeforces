package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases previously stored in testcasesA.txt.
const testcasesAData = `100
4 18 2 9 4
15 14 15 13 7
3 15 0 13 14
19 0 14 9 8
18 3 10 1 1
0 20 17 1 13
6 13 0 17 8
14 15 17 8 12
7 7 14 10 1
13 17 20 4 6
20 9 3 11 17
13 16 6 10 10
18 15 16 13 19
1 15 7 13 14
5 11 17 12 3
14 16 3 6 17
12 11 15 1 16
1 9 19 19 19
12 20 5 6 17
7 0 6 18 18
7 12 16 12 19
11 14 8 18 20
0 12 16 5 17
17 6 13 2 16
11 18 17 7 17
13 15 11 14 12
0 17 17 20 20
10 14 19 1 8
20 5 17 19 6
2 17 8 2 3
2 0 14 1 9
7 8 3 20 6
11 9 2 6 6
8 16 5 9 10
14 10 15 16 4
0 9 12 11 14
6 8 3 9 17
6 19 13 1 8
0 12 4 2 6
14 16 13 18 8
20 16 14 8 17
20 0 12 19 11
20 13 1 10 5
6 1 9 3 3
9 9 5 14 19
8 4 0 18 2
18 6 18 15 6
19 16 1 13 7
11 3 6 19 14
18 6 15 4 13
9 16 15 1 11
19 12 9 1 6
6 10 18 5 11
13 6 8 4 13
17 11 17 16 18
7 2 1 3 5
5 5 17 7 9
10 19 16 9 12
10 10 3 10 8
19 15 4 19 18
3 10 1 14 3
12 4 4 11 4
19 18 12 3 19
17 7 18 3 9
11 9 18 18 4
14 8 3 2 10
0 19 0 3 14
3 1 6 8 19
13 5 3 15 6
7 5 3 14 13
17 9 17 9 16
10 3 6 11 2
0 0 9 20 11
14 12 10 13 3
2 10 19 15 4
8 6 19 18 16
11 8 5 18 7
9 6 7 12 3
8 2 14 3 19
20 10 7 13 10
1 10 5 11 19
9 7 10 4 18
19 18 19 3 8
7 0 7 13 3
8 17 2 3 1
20 0 9 12 16
15 4 3 17 11
2 16 5 6 5
4 10 9 4 17
19 9 4 7 5
17 1 10 20 18
6 5 9 14 18
5 1 7 9 3
14 13 17 9 18
14 17 14 1 13
10 5 8 16 1
20 13 18 1 2
11 18 4 19 5
4 8 8 13 19
12 5 19 3 8`

type testCase struct {
	b int
	p int
	f int
	h int
	c int
}

// solve mirrors the logic from 1207A.go for a single test case.
func solve(tc testCase) int {
	buns := tc.b / 2
	profit := 0
	if tc.h >= tc.c {
		x := min(tc.p, buns)
		profit += x * tc.h
		buns -= x
		y := min(tc.f, buns)
		profit += y * tc.c
	} else {
		x := min(tc.f, buns)
		profit += x * tc.c
		buns -= x
		y := min(tc.p, buns)
		profit += y * tc.h
	}
	return profit
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx+4 >= len(tokens) {
			return nil, fmt.Errorf("test %d missing numbers", i+1)
		}
		b, _ := strconv.Atoi(tokens[idx])
		p, _ := strconv.Atoi(tokens[idx+1])
		f, _ := strconv.Atoi(tokens[idx+2])
		h, _ := strconv.Atoi(tokens[idx+3])
		c, _ := strconv.Atoi(tokens[idx+4])
		idx += 5
		cases = append(cases, testCase{b: b, p: p, f: f, h: h, c: c})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("embedded data has %d extra tokens", len(tokens)-idx)
	}
	return cases, nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("1\n%d %d %d %d %d\n", tc.b, tc.p, tc.f, tc.h, tc.c)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strconv.Itoa(solve(tc))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestCases(testcasesAData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
