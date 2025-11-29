package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesData = `
2 4 0 16
1 4 5
2 2 3 3
1 4 13
4 2 0 4 2 0
6 3 5 5 7 2 8 1
1 3 1
1 3 2
4 2 1 4 2 1
5 4 16 3 3 15 3
4 4 5 14 10 4
3 4 12 2 16
2 3 7 3
4 3 6 0 7 0
4 4 7 13 7 8
4 4 4 7 14 9
6 3 7 2 8 1 3 4
1 1 0
1 3 5
4 2 3 1 3 3
2 2 3 0
6 3 8 3 5 2 8 8
2 4 15 6
1 4 7
5 4 10 6 0 1 1
2 2 3 4
6 2 0 4 3 2 2 3
2 2 2 4
6 3 7 6 7 5 7 1
2 2 0 1
4 4 6 11 1 0
2 4 2 13
5 4 6 0 4 16 15
4 3 8 8 1 1
1 4 6
5 4 11 9 2 0 7
5 1 2 1 0 1 1
5 1 1 1 0 1 1
2 3 4 1
2 1 1 2
1 4 16
4 4 10 12 0 3
5 4 4 4 0 3 13
1 2 2
4 3 1 5 4 5
6 2 4 4 0 2 4 4
2 1 1 0
3 3 2 5 1
1 4 8
5 1 0 2 0 1 0
3 2 0 4 1
6 2 1 3 1 4 3 2
3 2 2 4 1
4 2 2 1 2 4
1 1 2
3 3 5 0 1
3 1 0 0 2
2 3 8 8
4 3 1 6 8 8
1 4 2
3 4 5 9 7
6 4 9 1 0 0 8 2
3 3 8 4 0
6 2 4 1 2 4 0 3
3 4 7 8 10
4 4 3 13 11 5
2 1 2 1
1 4 7
3 1 0 0 0
1 4 11
1 4 16
4 2 2 2 3 3
6 4 1 5 16 11 6 7
2 2 2 1
5 2 2 4 0 1 1
6 4 5 7 4 5 7 16
6 3 0 7 3 0 6 2
5 4 12 6 12 16 15
2 1 0 2
1 1 0
6 2 3 3 2 3 2 3
6 3 1 2 3 3 0 0
6 1 1 0 1 2 2 0
1 3 2
6 3 5 5 4 8 4 2
4 1 2 0 2 1
4 3 4 7 2 7
6 3 5 5 6 6 6 8
6 3 3 3 6 3 4 8
5 2 1 4 4 4 4
3 4 8 13 9
6 4 16 1 6 8 14 2
2 1 2 0
1 3 3
5 2 1 2 4 2 0
1 1 2
3 2 1 3 3
3 3 1 6 1
2 1 0 2
3 3 3 6 3
`

type testCase struct {
	n   int
	k   int
	arr []int
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	var cases []testCase
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		if len(fields)-2 != n {
			return nil, fmt.Errorf("testcase mismatch n=%d but %d numbers", n, len(fields)-2)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, err
			}
			arr[i] = val
		}
		cases = append(cases, testCase{n: n, k: k, arr: arr})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

// solve mirrors 1054D.go so the verifier is self-contained.
func solve(tc testCase) int64 {
	n := tc.n
	k := tc.k
	a := make([]int, n+1)
	m := (1 << uint(k)) - 1
	for i := 1; i <= n; i++ {
		xi := a[i-1] ^ tc.arr[i-1]
		if t := xi ^ m; t < xi {
			a[i] = t
		} else {
			a[i] = xi
		}
	}
	sort.Ints(a)
	ln := int64(n + 1)
	ans := ln * (ln - 1) / 2
	for i := 0; i < len(a); {
		j := i + 1
		for j < len(a) && a[j] == a[i] {
			j++
		}
		cnt := int64(j - i)
		c1 := cnt / 2
		c2 := cnt - c1
		ans -= c1 * (c1 - 1) / 2
		ans -= c2 * (c2 - 1) / 2
		i = j
	}
	return ans
}

func expected(tc testCase) string {
	return strconv.FormatInt(solve(tc), 10)
}

func formatInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d", tc.n, tc.k)
	for _, v := range tc.arr {
		fmt.Fprintf(&b, " %d", v)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		input := formatInput(tc)
		exp := expected(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
