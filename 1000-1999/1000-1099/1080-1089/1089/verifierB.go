package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `19 15 285
-16 9 -144
14 6 84
18 -19 -342
7 -19 -133
0 19 0
10 -9 -90
-2 -16 32
8 15 120
-14 4 -56
-10 -1 10
20 -11 -220
1 -7 -7
-2 -6 12
-16 -17 272
-7 13 -91
-18 2 -36
-7 -5 35
6 -17 -102
-14 8 -112
11 -6 -66
7 4 28
-19 -3 57
-14 12 -168
7 12 84
14 -14 -196
0 2 0
11 15 165
8 7 56
19 -19 -361
10 -12 -120
9 16 144
-11 2 -22
-10 8 -80
-4 -6 24
2 16 32
16 4 64
-9 6 -54
1 -19 -19
18 -12 -216
8 4 32
17 13 221
-16 9 -144
5 9 45
-12 -14 168
-1 14 -14
18 7 126
-3 -14 42
14 -13 -182
14 -4 -56
2 -5 -10
-10 4 -40
-10 11 -110
-6 15 -90
-3 9 -27
-5 -3 15
11 14 154
14 -16 -224
-6 -8 48
-17 -13 221
-12 -17 204
-14 10 -140
12 18 216
1 14 14
6 -8 -48
-8 9 -72
4 -8 -32
11 6 66
8 -2 -16
-6 16 -96
7 6 42
0 17 0
4 -11 -44
-19 17 -323
12 9 108
-9 2 -18
18 -12 -216
-4 15 -60
-15 -9 135
19 2 38
-15 -6 90
-4 4 -16
13 6 78
15 16 240
14 -4 -56
14 -6 -84
6 -12 -72
-17 -20 340
20 5 100
15 1 15
-11 13 -143
-5 -6 30
19 -19 -361
1 14 14
-14 -17 238
-13 -19 247
-13 -15 195
17 4 68
5 1 5
14 -2 -28`

type testCase struct {
	a int64
	b int64
}

// solve embeds the simple multiplication logic from 1089B.go.
func solve(a, b int64) int64 {
	return a * b
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
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
	want := strconv.FormatInt(solve(tc.a, tc.b), 10)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields)%3 != 0 {
		return nil, fmt.Errorf("invalid testcase data: count not multiple of 3")
	}
	tests := make([]testCase, 0, len(fields)/3)
	for i := 0; i < len(fields); i += 3 {
		a, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return nil, err
		}
		exp, err := strconv.ParseInt(fields[i+2], 10, 64)
		if err != nil {
			return nil, err
		}
		if solve(a, b) != exp {
			return nil, fmt.Errorf("embedded expected mismatch for %d %d", a, b)
		}
		tests = append(tests, testCase{a: a, b: b})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(tests))
}
