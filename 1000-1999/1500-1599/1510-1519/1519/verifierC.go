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

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `100
1
1
3
6
2 6 6 3 3 5
7 20 2 19 6 14
7
7 6 7 5 3 5 4
17 9 2 1 12 15 11
7
4 5 2 5 2 2 2
1 6 11 6 5 17 17
6
5 6 5 2 4 4
17 12 19 12 12 15
3
2 3 3
15 17 8
8
5 8 6 8 8 6 8 8
8 11 6 20 9 16 10 10
7
3 6 2 4 5 3 6
20 3 11 1 7 4 2
1
1
19
4
1 2 3 2
7 2 14 2
1
1
12
3
1 3 1
3 4 3
1
1
1
6
3 2 2 6 2 5
1 13 19 2 8 5
1
1
12
2
2 2
16 1
5
4 5 5 1 3
13 20 5 16 8
2
2 1
1 15
3
3 3 2
16 17 11
3
2 2 2
20 14 1
3
3 1 2
2 5 6
3
1 2 3
8 17 2
4
2 4 1 3
3 19 8 20
6
3 6 4 3 5 1
5 2 13 14 6 4
2
1 1
4 1
3
1 1 1
1 17 15
8
5 7 4 4 7 7 1 1
14 17 19 6 4 16 12 1
2
2 2
12 10
1
1
4
2
2 1
1 15
1
1
16
8
4 2 1 5 1 6 5 2
8 16 7 4 19 12 13 15
3
2 2 1
9 4 4
2
2 2
7 4
1
1
2
8
5 6 8 3 6 5 8 8
14 16 10 13 8 6 16 20
5
5 4 1 5 5
4 3 12 6 18
3
2 1 1
2 5 10
7
2 6 6 6 3 4 2
17 10 4 5 18 14 4
6
5 2 6 5 3 2
6 15 8 13 12 19
3
2 2 3
1 20 13
3
2 3 1
16 9 13
5
4 4 3 5 3
3 8 18 20 7
7
7 6 4 6 1 3 4
17 15 6 4 1 13 7
7
2 1 4 7 5 7 7
7 9 19 19 7 16 20
3
1 3 3
14 16 9
3
2 3 1
3 12 1
8
2 8 6 8 5 8 1 2
20 12 6 13 9 5 2 6
8
7 8 5 3 1 5 8 1
12 2 18 13 19 15 7 10
8
3 8 5 2 5 6 5 6
10 13 17 3 17 7 13 20
3
3 3 1
10 2 8
8
4 5 1 2 2 7 6 4
11 12 3 11 15 12 6 16
8
5 8 3 8 4 5 6 3
4 8 16 7 12 6 12 5
3
1 2 3
13 13 11
5
5 5 5 3 4
10 18 20 3 12
5
4 4 2 3 3
15 16 3 6 11
7
2 1 1 3 2 3 1
14 1 18 11 8 20 13
5
4 2 3 3 2
16 4 5 7 11
5
2 4 3 3 1
11 7 8 8 20
1
1
12
1
1
6
2
2 2
9 5
6
5 5 1 3 6 6
20 13 8 2 13 16
8
6 2 8 7 8 3 7 7
17 15 2 4 15 19 5 4
3
1 2 2
15 1 9
2
2 1
6 1
3
2 3 1
11 15 2
8
4 2 8 3 1 3 1 1
7 18 1 17 11 17 8 5
6
4 1 2 5 1 2
4 15 7 2 20 7
7
3 5 6 7 4 6 5
17 6 17 4 5 7 6
7
2 3 3 4 2 4 2
13 11 10 4 18 4 16
5
3 5 4 3 2
14 5 18 4 1
4
2 2 4 1
5 1 9 16
1
1
5
6
1 6 2 1 2 6
18 6 3 15 10 7
3
2 3 2
17 19 3
7
4 6 6 1 4 3 6
4 9 1 7 14 11 9
7
5 5 6 2 4 7 2
6 15 15 12 13 16 20
5
5 2 5 4 4
7 16 19 11 10
2
1 2
20 16
4
2 3 2 3
4 1 1 7
6
1 3 5 3 6 6
11 15 3 14 16 1
5
5 5 2 2 2
6 20 13 3 19
8
5 2 8 8 4 3 5 4
7 20 11 19 20 13 17 14
4
2 1 3 2
5 20 13 14
2
2 2
13 16
7
3 2 2 2 1 5 5
3 20 18 1 2 13 14
7
2 5 3 1 3 5 3
17 16 19 3 15 8 9
1
1
16
1
1
5
4
3 2 1 2
10 4 18 18
2
1 2
5 2
5
5 3 4 1 5
12 11 4 20 12
2
2 2
9 16
5
5 5 2 1 1
11 14 1 12 18
1
1
18
7
4 4 2 7 2 2 5
2 1 19 12 6 10 1
1
1
19
4
4 1 3 1
20 3 8 8
4
1 1 4 1
17 9 19 8`

type testCase struct {
	n int
	u []int
	s []int64
}

func solveCase(tc testCase) []int64 {
	n := tc.n
	groups := make(map[int][]int64)
	for i := 0; i < n; i++ {
		groups[tc.u[i]] = append(groups[tc.u[i]], tc.s[i])
	}

	prefix := make(map[int][]int64)
	for _, arr := range groups {
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		for i := 1; i < len(arr); i++ {
			arr[i] += arr[i-1]
		}
		size := len(arr)
		ps := prefix[size]
		if ps == nil {
			ps = make([]int64, size)
			prefix[size] = ps
		}
		for i := 0; i < size; i++ {
			ps[i] += arr[i]
		}
	}

	ans := make([]int64, n)
	for size, ps := range prefix {
		for k := 1; k <= size; k++ {
			teams := size / k
			if teams == 0 {
				break
			}
			idx := teams*k - 1
			ans[k-1] += ps[idx]
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	rd := strings.NewReader(testcasesRaw)
	in := bufio.NewReader(rd)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	res := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		u := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &u[i]); err != nil {
				return nil, fmt.Errorf("case %d: parse u%d: %v", caseIdx+1, i+1, err)
			}
		}
		s := make([]int64, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(in, &s[i]); err != nil {
				return nil, fmt.Errorf("case %d: parse s%d: %v", caseIdx+1, i+1, err)
			}
		}
		res = append(res, testCase{n: n, u: u, s: s})
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expect := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for idx, v := range tc.u {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.s {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(strings.TrimSpace(got))
		if len(gotFields) != tc.n {
			fmt.Printf("case %d failed\nexpected %d numbers\ngot %d numbers\n", i+1, tc.n, len(gotFields))
			os.Exit(1)
		}
		for idx, f := range gotFields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil || val != expect[idx] {
				fmt.Printf("case %d failed at k=%d\nexpected: %v\ngot: %s\n", i+1, idx+1, expect, strings.TrimSpace(got))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
