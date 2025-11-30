package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n   int
	arr []int64
}

// Embedded testcases from testcasesC.txt.
const testcasesRaw = `7 -4 -9 2 4 1 -4 4
5 10 -8 -9 -9 5
5 -10 6 8 8 -4
4 -8 10 6 6
6 6 -1 -7 -6 3 8
6 -8 -7 3 -8 -7 3
4 -10 4 3 3
3 5 0 -2
3 1 -8 -7
5 -10 1 1 -5 -10
4 1 -8 9 -6
4 -10 -4 -7 -10
5 1 -10 9 -3 -6
4 4 -7 5 1
8 -2 -6 -10 -4 1 0 5 -1
5 7 10 0 -5 8
3 -7 7 8
5 -5 2 -6 -6 -3
5 6 -3 -3 -5 -1
5 3 -9 -6 9 -10
6 -8 -8 -6 3 -1 7
6 -6 8 3 -1 10 1
3 -3 4 10
5 10 6 -9 2 3
3 3 0 4
4 1 -1 5 -8
4 -7 -2 -7 7
7 -6 4 2 -5 3 3 -5
4 4 0 6 -6
5 4 10 10 -8 5
4 -1 -10 4 9
6 -10 -4 -1 -7 10 -1
7 9 -6 3 5 -8 5 -3
7 2 -2 10 -10 -7 -2 -9
3 -2 2 6
7 2 4 -7 -2 1 -1 -4
7 -8 -9 -8 -2 -1 7 0
3 6 -3 -5
3 3 -1 -1
7 -6 8 6 10 -4 7 -7
6 10 7 2 -2 -1 4
5 8 10 -6 -5 -7
8 -7 2 2 8 4 -6 7 -1
5 10 5 3 -4 5
6 6 0 5 10 -9 4
5 -6 5 -9 9 -4
3 1 5 2
3 6 -8 -8
8 2 -10 1 -9 -7 9 -10 -2
8 -1 -3 -6 8 -1 -4 -7 3
6 0 2 -5 0 3 10
8 3 -6 4 -6 6 0 -6 -4
4 4 1 2 3
6 2 -3 -4 4 -4 8
8 -9 2 -9 -3 10 -8 -5 1
3 10 -5 -3
7 -1 9 -8 6 -1 1 3
6 -9 10 6 10 7 3
7 4 5 -2 5 -4 0 -2
3 -9 -9 -5
5 -10 -1 10 -10 -6
3 3 -3 9
6 7 -3 4 -4 0 9
3 9 -8 0
5 7 4 0 -2 -10
7 -9 -4 1 -8 -4 6 1
4 -4 -2 -1 -1
7 2 -2 5 1 -3 -9 -1
7 -8 -10 4 5 4 -9 3
6 4 4 -7 -8 -8 -3
3 -6 3 -4
6 9 -8 3 7 2 -9
4 -3 5 -3 -6
5 1 0 3 -7 7
5 9 7 -4 -1 4
7 9 4 7 10 -2 -2 -3
3 -7 9 -7
4 3 -3 -4 -1
8 -10 7 6 3 -9 -7 2 10
5 -7 8 1 -3 7
8 -1 -3 -3 -8 6 -1 0 -3
5 10 5 -1 8 -5
4 -10 7 6 0
5 8 10 -10 -6 2
4 -5 6 -8 -6
4 5 8 -4 -3
8 -6 -3 2 1 9 8 -6 10
6 -7 9 -10 6 9 1
6 4 -1 -10 -3 7 10
4 5 5 7 0
8 -8 -2 -6 9 2 -4 0 -1
6 -9 -4 -9 0 -3 0
6 -3 -2 1 -5 -1 -10
5 8 7 -9 10 -6
5 -10 5 10 -9 -10
4 -9 -10 -3 10
5 -8 -9 1 3 -6
4 4 3 -6 1
5 -5 10 0 3 2
3 3 -2 7`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testcase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n+1, len(fields))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i], err = strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value", idx+1)
			}
		}
		res = append(res, testcase{n: n, arr: arr})
	}
	return res, nil
}

// Embedded solver logic from 466C.go.
func countWays(a []int64) int64 {
	n := len(a)
	if n < 3 {
		return 0
	}
	var total int64
	for _, v := range a {
		total += v
	}
	if total%3 != 0 {
		return 0
	}
	target := total / 3
	var prefix int64
	var cntT int64
	var ans int64
	for i := 0; i < n-1; i++ {
		prefix += a[i]
		if i > 0 && prefix == 2*target {
			ans += cntT
		}
		if prefix == target {
			cntT++
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		exp := countWays(tc.arr)

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d\n", tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprint(exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
