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
	a []int
	b []int
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
100
5
3 4
15 10
16 19
1 11
7 12
7
15 11
12 13
11 10
6 17
6 1
16 5
19 4
2
5 12
9 16
7
17 1
2 14
8 6
19 8
17 14
8 10
7 17
6
14 7
1 7
17 19
10 7
6 20
17 19
3
9 8
18 8
9 1
6
6 15
11 17
15 5
12 17
8 6
3 17
3
5 9
12 2
9 6
7
12 10
15 4
19 15
13 1
10 13
7 11
16 4
7
3 19
19 1
10 7
10 10
14 1
12 1
7 15
1
13 17
4
20 14
7 5
16 6
4 18
4
15 12
5 4
5 19
17 18
4
11 18
9 19
15 10
15 7
5
7 5
20 1
9 7
12 14
5 9
6
3 6
4 18
1 3
19 16
12 14
7 20
6
20 9
16 16
3 14
10 16
18 4
13 1
6
19 5
6 19
9 10
7 15
1 18
4 19
7
10 20
13 14
18 11
4 11
12 17
12 11
15 10
5
19 2
20 14
11 12
20 15
1 19
4
4 2
1 14
6 12
20 10
2
16 4
20 18
1
19 7
1
9 12
3
7 6
13 5
10 10
3
9 16
10 12
11 9
2
12 1
17 15
2
15 10
9 14
1
19 8
7
16 10
8 14
20 20
19 16
18 6
6 11
8 5
1
19 7
4
7 10
11 3
16 8
6 4
6
11 18
16 12
12 10
1 7
9 14
9 3
7
7 3
12 6
9 5
12 15
4 10
2 2
18 17
1
8 3
1
3 14
3
6 19
19 19
20 2
5
13 13
5 10
11 6
3 2
15 1
3
3 14
14 14
18 17
4
19 6
18 8
11 1
14 5
6
16 1
9 20
16 15
4 20
16 13
15 11
7
6 11
14 8
3 3
7 2
14 18
15 2
10 13
3
1 9
15 2
2 12
1
9 5
6
10 13
7 10
12 8
15 14
14 11
8 8
4
14 13
15 5
4 19
12 17
6
7 9
15 1
19 6
4 2
11 9
1 14
3
8 5
5 8
6 13
4
18 12
16 12
14 19
7 8
4
20 15
9 11
17 19
3 7
5
11 10
14 11
19 8
16 19
2 10
2
3 13
17 9
5
18 7
2 10
20 7
7 19
19 4
3
1 15
11 8
1 11
7
9 14
18 12
5 19
11 13
1 3
8 19
20 2
4
15 15
15 2
20 7
6 4
4
3 3
14 18
19 7
5 7
5
1 14
19 5
4 15
4 3
9 3
7
9 6
2 9
13 3
10 4
4 17
1 10
11 5
6
18 16
8 4
16 5
20 8
2 10
17 15
6
9 19
1 6
4 7
20 3
10 12
15 9
1
1 20
2
15 10
7 9
6
16 6
17 10
19 5
9 17
2 1
15 12
2
8 1
6 17
1
1 2
2
7 9
10 8
1
8 7
7
1 20
14 7
4 17
11 1
9 14
6 20
17 9
6
15 19
8 17
18 1
6 10
11 20
16 14
5
11 9
4 17
4 4
16 8
4 13
6
14 1
1 10
14 16
20 6
2 19
3 10
7
13 1
8 18
12 1
10 18
7 15
19 15
20 19
6
1 15
17 13
14 7
16 9
20 18
5 20
7
20 15
11 4
9 1
10 10
20 16
15 6
20 6
5
15 11
5 8
10 3
1 13
16 20
4
16 9
2 2
17 9
13 8
6
11 10
16 9
17 15
9 14
1 13
19 6
2
14 10
3 5
6
13 3
9 9
5 3
1 14
17 1
4 9
6
2 7
16 13
20 7
3 3
4 10
19 4
3
9 4
19 17
13 11
4
17 2
7 15
3 6
3 15
5
19 11
8 14
1 3
6 3
13 16
6
11 8
14 4
19 9
4 7
19 19
7 13
2
4 2
16 9
6
16 12
10 8
9 17
6 10
20 19
6 9
3
3 4
17 16
15 5
3
9 7
16 12
6 4
1
12 18
6
15 15
7 1
6 10
15 10
19 9
1 20
6
11 13
6 15
4 8
4 3
1 17
3 8
4
19 10
4 8
9 15
10 10
6
12 13
14 3
20 6
8 13
3 10
18 8
5
20 17
6 2
8 16
5 3
5 11
6
20 10
16 7
14 10
18 16
13 12
6 14
6
5 7
8 17
1 6
12 8
9 7
9 18
1
1 2
6
8 12
8 10
11 6
5 16
1 20
7 4
3
11 16
3 13
5 14
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			ai, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad a[%d]: %v", i+1, j, err)
			}
			bi, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad b[%d]: %v", i+1, j, err)
			}
			a[j] = ai
			b[j] = bi
		}
		res = append(res, testCase{n: n, a: a, b: b})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of test data")
	}
	return res, nil
}

// solve mirrors 1503D.go.
func solve(tc testCase) int {
	n := tc.n
	a := tc.a
	b := tc.b
	if n > 16 {
		return -1
	}
	m := 1 << n
	minFlip := -1
	for mask := 0; mask < m; mask++ {
		front := make([]int, n)
		back := make([]int, n)
		flips := 0
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				front[i] = b[i]
				back[i] = a[i]
				flips++
			} else {
				front[i] = a[i]
				back[i] = b[i]
			}
		}
		idx := make([]int, n)
		for i := 0; i < n; i++ {
			idx[i] = i
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if front[idx[i]] > front[idx[j]] {
					idx[i], idx[j] = idx[j], idx[i]
				}
			}
		}
		ok := true
		for i := 1; i < n; i++ {
			if back[idx[i-1]] <= back[idx[i]] {
				ok = false
				break
			}
		}
		if ok {
			if minFlip == -1 || flips < minFlip {
				minFlip = flips
			}
		}
	}
	return minFlip
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i], tc.b[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		expect := solve(tc)
		out, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d: %v", idx+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || val != expect {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
