package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `4 8 10
3 8 10
8 10 10
10 10 10
5 9 9
4 9 10
9 10 10
3 6 7
9 10 10
2 4 10
10 10 10
1 5 8
10 10 10
7 10 10
6 6 6
3 10 10
5 10 10
5 8 10
7 9 10
10 10 10
8 10 10
7 8 10
7 10 10
8 9 10
5 7 10
6 9 10
8 9 10
7 10 10
7 9 9
8 10 10
7 9 10
6 10 10
6 8 10
4 10 10
9 9 10
8 10 10
10 10 10
7 8 10
3 4 7
8 9 10
8 10 10
7 10 10
8 10 10
8 10 10
2 8 10
8 10 10
5 10 10
4 8 10
5 7 10
7 9 10
6 8 10
7 8 10
7 8 10
8 10 10
8 10 10
7 7 10
8 8 10
7 7 10
7 8 10
8 8 10
7 10 10
3 6 10
6 6 10
5 5 10
3 6 9
6 7 10
9 10 10
10 10 10
6 8 10
2 2 4
7 9 10
5 5 8
9 10 10
10 10 10
9 10 10
9 10 10
8 10 10
4 8 9
10 10 10
8 10 10
7 9 9
7 10 10
9 9 10
5 9 10
10 10 10
6 8 9
10 10 10
5 9 10
9 10 10
9 10 10
8 8 10
8 9 10
9 9 10
9 10 10
9 10 10
6 10 10
7 8 9
10 10 10
10 10 10
9 10 10`

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Embedded reference logic from 1311D.go.
func solveCase(a, b, c int) (int, int, int, int) {
	const INF = 1000000000
	ans := INF
	var ansA, ansB, ansC int

	processA := func(A int) {
		costA := abs(A - a)
		if costA > ans {
			return
		}
		k := b / A
		tryB := func(B int) {
			costB := abs(B - b)
			if costA+costB > ans {
				return
			}
			jc := c / B
			tryC := func(C int) {
				total := costA + costB + abs(C-c)
				if total < ans {
					ans = total
					ansA, ansB, ansC = A, B, C
				}
			}
			if jc >= 1 {
				tryC(B * jc)
			}
			tryC(B * (jc + 1))
		}
		if k >= 1 {
			tryB(A * k)
		}
		tryB(A * (k + 1))
	}

	for A := a; A >= 1; A-- {
		if a-A > ans {
			break
		}
		processA(A)
	}
	for A := a + 1; A <= 10000; A++ {
		if A-a > ans {
			break
		}
		processA(A)
	}
	return ans, ansA, ansB, ansC
}

type testCase struct {
	a, b, c int
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("expected triples of integers, got %d tokens", len(fields))
	}
	res := make([]testCase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		a, err1 := strconv.Atoi(fields[i])
		b, err2 := strconv.Atoi(fields[i+1])
		c, err3 := strconv.Atoi(fields[i+2])
		if err := firstErr(err1, err2, err3); err != nil {
			return nil, fmt.Errorf("invalid numbers at position %d: %w", i+1, err)
		}
		res = append(res, testCase{a: a, b: b, c: c})
	}
	return res, nil
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expectedCost, A, B, C := solveCase(tc.a, tc.b, tc.c)
		input := fmt.Sprintf("1\n%d %d %d\n", tc.a, tc.b, tc.c)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expStr := fmt.Sprintf("%d\n%d %d %d", expectedCost, A, B, C)
		if strings.TrimSpace(got) != expStr {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
