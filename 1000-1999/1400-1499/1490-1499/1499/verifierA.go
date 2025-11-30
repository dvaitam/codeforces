package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 1499A.go.
func solve(n, k1, k2, w, b int) string {
	white := k1 + k2
	maxWhite := white / 2
	black := 2*n - white
	maxBlack := black / 2
	if w <= maxWhite && b <= maxBlack {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	n, k1, k2 int
	w, b      int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
100
7 6 0
4 7
7 4 7
5 3
9 2 4
2 1
10 4 8
9 2
5 0 5
0 5
6 3 4
0 2
7 5 3
7 7
9 4 0
8 0
2 2 1
2 2
1 1 1
0 1
2 0 2
0 0
3 3 0
0 2
9 7 1
4 8
5 5 0
4 2
9 3 9
8 9
5 3 0
4 3
6 4 1
2 1
4 1 0
4 2
8 1 1
2 2
1 0 1
1 0
4 4 3
4 2
8 7 5
1 5
10 1 7
9 10
6 6 1
1 0
5 0 5
1 2
3 2 3
0 0
3 1 0
0 0
2 2 0
2 2
2 1 0
1 0
1 0 0
0 0
8 3 0
0 8
7 1 4
1 3
2 2 1
1 1
3 0 3
0 0
7 3 4
5 7
10 2 10
3 0
3 1 2
2 0
10 7 10
2 0
8 6 8
4 5
7 4 2
0 7
2 1 2
0 2
5 1 1
3 2
10 4 10
5 9
10 2 4
6 6
2 0 2
0 2
6 1 1
1 5
8 6 6
0 6
10 6 10
0 2
8 1 4
2 7
9 7 8
9 0
1 1 1
1 1
1 1 0
0 0
1 1 1
1 0
4 0 0
4 4
2 0 0
2 2
4 2 2
1 0
8 6 1
0 4
8 1 4
2 8
6 0 6
1 2
1 0 0
0 1
9 5 5
9 0
10 10 7
10 7
7 5 2
3 6
10 4 0
2 2
5 2 2
2 5
2 1 2
0 0
5 1 1
4 2
6 3 4
1 2
2 1 2
0 0
5 1 4
5 0
5 3 2
2 3
2 0 2
1 1
6 6 6
6 2
2 1 0
2 1
7 0 4
5 2
3 3 0
0 0
4 1 0
3 0
2 1 2
2 1
8 7 3
6 1
6 1 2
4 6
3 3 1
2 0
2 2 0
2 1
4 0 3
3 2
4 0 1
4 1
2 0 1
1 1
9 2 1
9 7
3 3 3
3 2
8 7 3
8 3
1 1 1
1 0
9 2 4
9 2
7 4 7
1 1
9 0 1
3 2
1 1 0
1 1
3 1 3
2 3
9 8 0
9 1
9 9 1
6 3
5 4 4
3 3
7 3 0
0 2
5 4 4
2 2
2 1 1
1 1
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	pos++
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+4 >= len(fields) {
			return nil, fmt.Errorf("case %d incomplete", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		k1, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k1: %v", i+1, err)
		}
		k2, err := strconv.Atoi(fields[pos+2])
		if err != nil {
			return nil, fmt.Errorf("case %d bad k2: %v", i+1, err)
		}
		w, err := strconv.Atoi(fields[pos+3])
		if err != nil {
			return nil, fmt.Errorf("case %d bad w: %v", i+1, err)
		}
		b, err := strconv.Atoi(fields[pos+4])
		if err != nil {
			return nil, fmt.Errorf("case %d bad b: %v", i+1, err)
		}
		pos += 5
		res = append(res, testCase{n: n, k1: k1, k2: k2, w: w, b: b})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end")
	}
	return res, nil
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("1\n%d %d %d\n%d %d\n", tc.n, tc.k1, tc.k2, tc.w, tc.b)
		expected := solve(tc.n, tc.k1, tc.k2, tc.w, tc.b)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
