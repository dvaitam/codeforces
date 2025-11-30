package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `1 2 2 3 2 3
5 10 4 1 3 7 7 6 8 5 1 1 6 8 6 7 7 3
9 3 4 4 1 3 6 3 3 6 3 8 7 6 6
6 8 3 4 6 6 4 6 5 2 4 3 4 5 5 3 6
8 8 6 10 12 9 12 8 8 11 4 6 12 3 10 5 8 5 5
9 9 9 17 14 10 7 16 17 12 3 11 1 7 4 2 2 9 8 4 17
3 5 4 4 1 7 1 1 6 6 3
4 1 2 1 1 1 1 1
6 5 3 2 6 2 5 6 1 4 5 1 2 2
1 1 6 10 11
2 5 6 8 1 5 8 9 10 12
1 5 7 14 10 12 3 8 4
2 6 2 1 4 2 4 4 3 2 3
5 5 10 14 1 18 5 2 9 2 5 6 6
2 8 4 1 4 4 8 2 5 2 4 6 5
7 5 9 1 5 2 13 14 6 4 17 3 8 4 4
1 3 4 2 4 1 8
8 5 9 13 7 7 14 14 17 1 2 14 17 6 4 16
6 1 9 4 12 10 12 10 1 14
2 2 5 4 1 8 1
7 8 8 7 3 1 10 1 12 10 3 8 16 7 4 12 13 15
3 6 7 2 5 2 2 2 10 14 6 11
7 4 2 1 4 1 4 3 3 4 2 3 3 4
9 8 7 8 14 11 5 7 4 3 8 10 5 9 7 12 11 12 2 10
10 2 2 3 2 2 4 1 1 1 2 3 4 2 3
8 3 9 10 4 5 18 14 4 11 17 8 17 9
3 3 8 8 13 12 5 15 15
1 10 7 12 3 7 9 1 8 5 7 5 12 12
7 8 6 9 6 12 12 11 2 12 4 9 10 4 7 11 7 11
1 6 8 15 6 4 1 13 7 13
4 2 7 14 9 13 13 4 5
10 10 4 8 3 1 7 8 5 3 8 4 2 6 1 8 2 8 6 8 5 8 1
2 10 6 3 7 5 11 11 12 3 1 3 8 7 8
5 3 1 2 2 1 2 1 2 2 1
5 8 3 4 6 5 6 3 1 3 3 3 3 6 3 6
7 9 2 2 4 2 1 3 1 2 4 2 3 1 1 1 4 3 2
6 6 2 3 4 3 2 4 4 3 4 2 4 2 3
6 3 2 2 4 2 3 2 3 2 2 2
5 9 7 7 13 14 12 6 5 12 10 9 10 12 12 6 12
7 5 9 3 12 10 13 16 6 9 12 15 16 3 6
6 7 3 1 1 3 2 3 1 6 6 4 1 5 3 2
10 7 9 10 16 5 12 11 7 16 4 5 7 11 9 5 14 12 9 3
6 4 4 4 1 6 6 1 3 3 2 7 8
5 3 6 9 10 2 6 11 12 10 7
4 1 7 13 8 8 10 14
6 9 10 20 3 19 17 18 16 13 15 6 14 13 17 15 2 4
8 10 3 1 6 5 2 1 4 3 4 6 1 3 1 6 3 2 2 1 2
7 2 6 11 8 1 8 4 2 8 3 9
1 3 9 18 2 2 7
9 1 9 11 17 8 5 12 16 1 5 18 4
4 2 8 7 2 7 13 11 13
9 9 3 5 1 2 6 2 2 4 2 3 3 4 2 4 2 4 3 3 1
9 2 8 9 10 16 9 8 14 5 4 1 7 7
4 7 10 2 5 1 9 16 18 2 8 5 20 11
1 4 2 2 2 1 4 3
4 3 6 12 5 9 10 2 7 7
1 8 5 2 5 1 4 7 6 5 9 7
10 9 4 7 3 3 8 8 6 7 8 5 4 8 8 4 8 6 5 2 3 6
10 8 4 3 5 4 5 2 1 1 4 6 1 6 5 6 8 2 7 8 1
5 10 10 5 7 5 6 20 13 3 19 15 9 3 16 16 8 5
10 5 4 4 6 7 7 4 4 1 5 4 3 7 7 2 8 7
7 8 7 5 4 4 4 1 9 9 14 2 10 9 11 1 1 7
7 7 4 5 2 6 6 8 2 8 4 5 1 1 8 1 3
3 4 6 4 9 1 10 3 11 5
2 9 9 3 5 14 5 2 10 17 9 16 2 18
6 6 2 3 1 3 3 3 4 3 2 1 1 3 4
1 6 9 2 3 18 17 14 14 14
4 3 3 5 1 1 5 6 3 6
3 5 1 1 1 1 2 1 2 1 1
4 4 9 7 4 1 13 3 17 9 8
1 9 9 17 13 14 5 5 14 5 15 12 2
10 3 9 15 14 15 6 16 5 12 5 1 9 6 5 14
10 5 8 16 15 7 14 14 9 8 12 2 13 1 14 10 1 16
10 5 5 4 8 8 6 9 10 8 4 9 9 3 8 5 6 7
2 9 4 7 2 7 8 8 2 7 8 6 6 3
3 4 3 4 4 4 6 6 2 4
5 6 10 13 10 9 11 1 13 19 2 7 15 4
2 1 6 6 10 6
7 3 6 2 9 10 8 7 10 4 8 2 10
1 6 1 2 2 1 1 1 2 2
8 10 1 2 1 1 2 1 2 2 2 2 2 2 1 1 2 1 2 2 1
4 9 10 8 7 19 15 8 13 9 19 7 17 6 1 13
8 6 3 6 2 6 6 5 2 4 5 1 6 2 2 6 2
1 10 2 4 2 2 3 4 1 1 3 1 3 1
10 6 8 13 4 3 15 6 5 15 3 9 5 16 3 10 10 1 7
2 7 3 2 5 3 4 2 1 6 3 4
1 3 4 6 5 2 4
3 4 1 2 1 2 1 2 1 1
6 3 7 12 12 4 5 9 13 1 13 1
6 3 2 1 2 4 4 3 4 3 3 1
3 2 1 1 1 2 1 2
10 9 3 1 2 5 6 6 2 3 1 2 6 4 2 4 6 6 1 3 4 4
9 6 3 5 1 6 1 2 6 4 2 4 6 2 6 3 5 5
2 2 9 6 7 3 18
6 3 9 8 8 15 13 15 4 11 8 2
6 2 7 1 10 2 12 13 5 6 9
4 6 7 4 1 3 3 3 12 14 3 7 3
2 5 7 8 1 2 7 1 12 2
1 4 7 1 4 13 14 7
6 3 3 3 3 3 5 2 5 3 3 5`

func runProg(bin, input string) (string, error) {
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

func solve(n, m, k int, aVals, bVals []int) string {
	presentA := make([]bool, k+1)
	presentB := make([]bool, k+1)
	cntA, cntB := 0, 0
	for _, x := range aVals {
		if x <= k && !presentA[x] {
			presentA[x] = true
			cntA++
		}
	}
	for _, x := range bVals {
		if x <= k && !presentB[x] {
			presentB[x] = true
			cntB++
		}
	}
	need := (k + 1) / 2
	ok := cntA >= need && cntB >= need
	if ok {
		for i := 1; i <= k; i++ {
			if !presentA[i] && !presentB[i] {
				ok = false
				break
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	n, m, k int
	a, b    []int
}

func parseCases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d too short", idx+1)
		}
		ints := make([]int, len(fields))
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d invalid int: %v", idx+1, err)
			}
			ints[i] = v
		}
		n, m, k := ints[0], ints[1], ints[2]
		expectedLen := 3 + n + m
		if len(ints) != expectedLen {
			return nil, fmt.Errorf("line %d length mismatch expected %d got %d", idx+1, expectedLen, len(ints))
		}
		a := make([]int, n)
		copy(a, ints[3:3+n])
		b := make([]int, m)
		copy(b, ints[3+n:])
		res = append(res, testCase{n: n, m: m, k: k, a: a, b: b})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases(testcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		want := solve(tc.n, tc.m, tc.k, tc.a, tc.b)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.m, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runProg(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
