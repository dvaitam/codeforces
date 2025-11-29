package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Helpers from 1250L reference solution.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

func solve(a, b, c int) int {
	total := a + b + c
	base := ceilDiv(total, 3)
	cand1 := max(base, max(a, c))
	cand2 := max(base, max(c, ceilDiv(a, 2)))
	cand3 := max(base, max(a, ceilDiv(c, 2)))
	ans := cand1
	if cand2 < ans {
		ans = cand2
	}
	if cand3 < ans {
		ans = cand3
	}
	return ans
}

type testcase struct {
	a, b, c int
}

const testcaseData = `
100
8 1 3
5 10 1
5 0 3
8 10 10
6 1 9
10 1 8
1 9 4
0 6 0
9 8 6
6 1 5
3 10 7
2 0 1
3 8 2
8 4 9
7 3 2
6 7 1
4 4 10
6 1 10
10 8 8
6 5 9
1 1 2
10 4 3
7 9 10
3 3 7
3 9 3
3 4 6
4 8 8
3 4 9
7 9 9
2 10 10
7 3 0
10 5 1
10 7 6
2 7 0
0 7 3
3 9 8
10 3 6
0 10 3
7 3 8
8 8 6
10 5 8
7 1 0
7 3 1
3 8 7
2 4 10
1 9 10
10 8 8
0 10 5
0 4 2
10 9 1
8 5 6
8 7 9
9 9 4
1 1 7
5 3 0
1 8 6
8 9 3
9 5 6
7 5 4
6 7 1
6 8 4
3 3 2
2 10 5
10 7 10
5 7 4
4 0 7
5 5 9
3 5 0
10 2 3
1 5 6
2 3 1
6 9 5
4 9 5
2 6 10
0 9 2
1 0 7
0 3 8
3 7 5
5 6 8
5 5 9
3 1 4
4 10 2
0 3 5
4 10 7
2 2 1
8 10 9
8 1 2
6 6 6
0 7 10
9 4 6
1 2 3
8 5 6
9 10 2
8 4 5
9 6 8
6 4 5
10 9 0
6 0 4
9 7 6
3 2 1
`

func loadTestcases() ([]testcase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %v", err)
	}
	expectedFields := 1 + 3*t
	if len(fields) != expectedFields {
		return nil, fmt.Errorf("testcase count mismatch: expected %d numbers, got %d", expectedFields, len(fields))
	}
	tests := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		a, _ := strconv.Atoi(fields[1+3*i])
		b, _ := strconv.Atoi(fields[1+3*i+1])
		c, _ := strconv.Atoi(fields[1+3*i+2])
		tests = append(tests, testcase{a: a, b: b, c: c})
	}
	return tests, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d %d\n", tc.a, tc.b, tc.c)
		expect := fmt.Sprintf("%d", solve(tc.a, tc.b, tc.c))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
