package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	a int64
	b int64
	x int64
	y int64
}

// solve is the embedded logic from 16C.go.
func solve(a, b, x, y int64) (int64, int64) {
	g := gcd(x, y)
	x /= g
	y /= g
	k1 := a / x
	k2 := b / y
	k := k1
	if k2 < k {
		k = k2
	}
	if k > 0 {
		return k * x, k * y
	}
	return 0, 0
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Embedded copy of testcasesC.txt (each line: a b x y).
const testcaseData = `
12 8 12 45
24 2 15 3
1 25 35 1
9 8 39 39
13 6 48 30
13 1 34 40
27 5 35 12
15 15 27 25
31 1 28 14
25 49 3 40
18 2 38 23
46 24 22 44
30 42 9 39
34 6 17 7
46 7 18 2
45 10 40 50
43 9 25 14
37 43 21 13
27 33 33 8
36 7 45 31
21 6 48 5
23 20 23 22
33 15 24 3
46 2 26 4
16 6 5 10
33 1 49 33
47 43 40 6
49 6 31 21
40 19 26 41
27 45 43 36
22 4 31 5
25 10 20 8
25 13 8 47
43 19 1 22
44 42 5 10
18 27 37 5
35 34 12 46
48 37 12 5
30 38 34 11
11 32 4 28
6 30 23 37
45 50 40 34
24 37 30 13
16 11 37 1
4 23 12 14
48 36 26 26
40 38 34 32
17 12 41 50
13 24 14 4
48 28 25 32
28 35 17 16
7 26 2 32
11 50 12 13
50 17 19 24
44 44 9 16
15 46 33 5
42 35 22 45
23 28 39 20
12 47 7 26
8 35 42 44
40 4 14 34
23 37 18 46
45 30 6 37
15 36 16 46
45 48 49 39
44 26 17 48
39 7 46 4
19 19 38 29
28 14 18 15
15 3 29 37
30 26 33 29
1 8 8 6
50 12 44 2
2 15 5 31
19 4 28 6
42 2 39 34
41 16 16 43
35 44 43 4
9 31 3 16
31 14 19 17
41 17 10 21
9 43 19 36
13 22 23 17
20 45 13 28
19 17 33 17
29 18 25 27
10 10 26 32
30 38 20 13
2 24 3 19
30 15 36 50
12 26 25 4
15 23 42 7
21 15 28 40
22 43 20 3
29 21 43 29
35 41 18 26
16 48 35 44
44 47 32 13
35 46 11 40
27 21 35 45
41 3 48 32
25 9 9 13
9 20 32 46
14 16 24 7
26 34 38 35
43 8 30 39
33 30 22 26
20 40 33 43
50 14 40 44
25 35 27 30
20 14 39 31
40 39 9 36
21 37 6 33
7 3 50 2
19 21 4 9
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 fields got %d", i+1, len(fields))
		}
		a, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad a: %v", i+1, err)
		}
		b, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad b: %v", i+1, err)
		}
		x, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad x: %v", i+1, err)
		}
		y, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad y: %v", i+1, err)
		}
		tests = append(tests, testCase{a: a, b: b, x: x, y: y})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d %d %d\n", tc.a, tc.b, tc.x, tc.y)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		xRes, yRes := solve(tc.a, tc.b, tc.x, tc.y)
		exp := fmt.Sprintf("%d %d", xRes, yRes)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
