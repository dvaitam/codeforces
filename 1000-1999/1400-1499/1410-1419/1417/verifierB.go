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
	n int
	t int64
	a []int64
}

// solve mirrors 1417B reference logic.
func solve(tc testCase) []int {
	n, t := tc.n, tc.t
	a := tc.a
	res := make([]int, n)
	var c int64
	for i := 0; i < n; i++ {
		x := a[i]
		if t%2 == 0 && x == t/2 {
			res[i] = int(c % 2)
			c++
		} else if 2*x < t {
			res[i] = 0
		} else {
			res[i] = 1
		}
	}
	return res
}

// Embedded testcases from testcasesB.txt (each line: n t a1 ... an).
const testcaseData = `
3 37 48 4 16
2 32 48 28
8 42 24 50 13 6 31 1 24 27
10 49 49 0 44 28 17 46 14 37 6 20
1 2 1
9 1 24 43 13 27 46 1 33 14 48
8 32 35 14 22 14 43 14 48 29
5 2 26 35 41 6 11
5 8 47 21 46 45 32
7 33 42 12 19 18 37 31 32
7 38 2 30 15 47 25 26 42
3 24 35 44 49
6 6 28 42 32 6 49 10
9 26 23 31 46 1 30 2 19 45 39
10 38 25 41 10 10 32 14 0 49 12 34
9 15 25 32 22 36 22 29 17 42 35
10 47 0 24 50 47 32 8 33 49 35 13
7 4 30 23 36 35 12 32 26
8 23 26 22 0 34 34 39 50 39
6 30 38 1 14 40 11 35
3 49 6 48 1
7 50 10 30 23 13 12 20 11
1 26 31
2 33 49 15
1 25 36
7 26 25 15 12 35 13 42 20
2 32 46 14
9 47 20 35 6 9 45 18 16 39 7
1 33 27
9 40 0 7 9 6 9 32 20 13 12
2 29 10 29
6 17 24 34 40 16 38 12
1 33 1
9 3 36 3 44 9 22 38 20 17 24
3 6 32 24 36
9 47 7 14 18 31 17 17 15 23 28
3 38 5 14 5
5 9 17 15 9 14 0
7 7 29 36 26 37 16 2 19
9 12 18 30 21 44 37 15 24 20 32
9 36 19 35 25 18 1 12 34 15 42
5 32 4 19 41 0 1
3 36 18 12 33
10 41 25 38 40 44 41 13 16 30 20 22
1 7 6
4 11 43 39 2 26
1 31 13
4 42 5 27 2 2
7 36 32 25 23 16 18 31 42
1 41 39
9 14 25 34 34 46 22 23 10 42 39
10 9 42 11 25 1 12 34 17 39 29 5
8 16 24 3 8 14 9 23 28 5
8 25 5 35 36 4 46 21 18 10
1 45 28
1 30 31
4 41 19 10 13 34
9 26 20 37 28 27 17 22 9 6 8
9 39 11 25 39 10 45 16 9 26 5
4 10 11 14 19 27
6 3 34 3 1 34 19 41
9 29 5 7 7 42 11 49 34 17 31
3 50 47 23 7
3 13 4 3 15
4 29 2 34 24 37
8 25 10 12 32 19 33 12 42 9
8 2 32 2 15 15 44 46 50 9
6 50 50 21 28 12 42 0
7 22 28 33 35 23 44 21 10
4 23 17 5 18 32
6 17 26 31 0 41 22 31
10 37 21 17 32 31 28 42 27 36 7 49
10 15 22 35 41 35 14 18 32 48 8 11
6 22 41 11 18 8 20 38
2 44 0 0
1 14 14
6 3 33 11 17 42 2 8
4 40 33 1 41 38
3 3 44 23 15
6 26 1 48 23 50 36 15
7 18 28 21 8 44 25 8 49
9 25 17 24 0 50 18 40 31 9 10
8 5 46 3 17 38 16 33 20 6
3 48 16 21 16
3 24 43 14 48
3 35 18 11 19
2 18 16 23
2 20 17 7
3 2 13 36 12
4 40 3 40 25 33
5 22 44 45 43 33 19
3 12 7 22 32
7 35 15 34 41 40 19 6 16
9 38 40 16 11 25 31 13 7 31 15
5 2 12 44 10 21 6
2 28 14 16
4 7 29 1 45 25
5 41 23 38 41 28 21
6 2 39 7 41 4 44 17
4 49 32 12 9 21
4 25 40 0 6 18
4 44 1 13 0 44
4 27 24 7 41 10
10 33 24 25 41 20 39 31 46 22 25 15
10 8 38 46 8 48 32 22 29 11 48 21
9 3 16 41 1 43 18 38 11 14 19
7 47 48 30 22 6 31 5 44
10 5 35 4 26 29 12 15 6 20 11 1
6 1 8 19 10 7 44 23
10 40 41 40 19 14 12 40 39 35 11 22
2 34 25 49
3 2 44 6 30
4 22 6 32 31 33
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: invalid format", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		t, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad t: %v", i+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, 2+n, len(fields))
		}
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			v, err := strconv.ParseInt(fields[2+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad a[%d]: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		cases = append(cases, testCase{n: n, t: t, a: arr})
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.t)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')

		expectedLabels := solve(tc)
		expectedStrs := make([]string, len(expectedLabels))
		for i, v := range expectedLabels {
			expectedStrs[i] = strconv.Itoa(v)
		}
		expected := strings.Join(expectedStrs, " ")

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
