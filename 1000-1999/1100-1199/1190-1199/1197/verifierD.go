package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
4 2 6
9 5 10 8
2 1 8
-2 7
4 2 8
7 7 5 2
3 2 3
6 2 -10
2 2 10
-9 -1
1 3 8
9
7 4 7
8 4 -6 1 -7 -9 -6
8 2 5
3 10 -1 3 6 2 8 1
7 2 6
-10 -2 9 -5 0 7 8
2 2 10
-2 -1
2 1 8
10 5
2 3 2
3 -6
1 3 7
3
2 1 10
9 -9
7 3 9
-2 6 -3 -9 -1 -10 -8
2 1 4
3 -1
5 2 1
0 0 1 -6 2
7 4 9
2 10 9 7 -7 9 6
5 4 4
-1 3 -2 6 -1
6 1 7
8 0 -10 2 9 8
3 1 6
4 1 1
5 4 1
8 -9 -10 1 -2
8 3 10
9 0 -5 1 -5 0 1 9
5 3 7
-7 -10 8 -6 -1
4 3 4
0 -5 3 10
2 1 10
0 0
4 4 3
-8 0 10 -4
8 3 4
-7 -9 6 -4 0 8 -5 -2
6 1 10
1 8 -6 3 -1 6
5 4 6
10 3 -1 3 8
7 1 7
-6 -4 -10 5 9 6 3
4 1 8
6 -1 7 0
4 1 10
-1 -7 -3 -9
1 2 7
8
1 1 8
-7
3 3 4
-10 6 7
7 1 10
-7 0 -6 -2 7 5 -9
6 2 4
-7 7 -7 -5 -3 -2
3 1 8
10 8 2
1 3 4
-2
7 1 8
0 -10 -9 -6 -9 -7 -9
2 4 1
-8 6
8 3 3
0 -8 1 2 10 2 8 -1
6 3 4
0 3 -7 -6 7 -10
7 1 10
-5 -9 1 4 9 10 7
7 1 10
3 -9 1 10 5 0 3
7 4 1
-3 -4 7 -2 8 -8 3
4 4 3
-10 0 1 7
5 1 8
-7 6 2 -7 0
2 1 8
-6 -3
7 1 9
-8 8 -7 2 -5 -10 0
2 1 2
5 -1
5 1 1
8 6 6 -3 -7
2 1 9
0 8
3 1 4
-5 10 -3
8 4 5
1 9 2 1 7 3 -8 2
4 4 3
3 8 8 6
8 2 7
-6 -5 -7 5 5 6 4 8
3 2 5
-4 -6 8
6 2 9
-1 3 9 8 8 -2
4 3 1
-2 5 2 -4
3 3 4
0 5 -6
7 4 10
-4 4 8 10 7 -10 5
2 4 1
4 -3
4 1 4
-2 -3 -4 -2
3 2 10
-9 -2 -5
1 3 3
3
2 1 2
-8 -2
5 1 6
4 8 0 -10 -10
6 3 7
2 5 -8 -4 10 8
8 4 3
7 0 -7 -2 -8 3 -7 4
5 1 9
1 1 4 -1 10
5 1 6
8 7 6 -7 5
6 1 5
8 -5 10 10 10 -6
3 3 8
-7 -7 7
3 3 10
3 7 -1
3 4 8
-1 -5 -8
2 2 9
7 8
7 3 2
-2 -2 2 -9 -6 -9 5
5 2 9
1 0 2 4 7
2 3 8
-7 -6
5 1 2
8 -7 -5 -4 8
7 4 3
8 9 -6 2 -4 7 6
3 2 4
-2 1 -1
1 4 7
2
6 3 8
6 -1 5 -10 9 -4
1 1 4
5
3 4 4
-4 6 -4
1 4 2
8
5 2 3
4 -8 9 -9 -10
6 2 9
-8 5 7 -10 0 0
6 3 3
-8 9 -9 -8 0 -4
2 2 7
-3 5
6 1 1
3 -8 -4 -5 2 5
8 1 9
3 -4 10 5 -1 -10 4 4
7 4 3
4 -9 -2 1 1 4 6
6 4 4
-10 -4 -2 1 -6 4
4 2 4
-10 -5 8 2
3 1 3
-7 9 -5
8 4 3
-9 -10 2 4 0 3 -9 -9
`

type testCase struct {
	n int
	m int
	k int64
	a []int64
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+2 >= len(fields) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, _ := strconv.Atoi(fields[pos])
		m, _ := strconv.Atoi(fields[pos+1])
		k, _ := strconv.ParseInt(fields[pos+2], 10, 64)
		pos += 3
		if pos+n > len(fields) {
			return nil, fmt.Errorf("not enough values for case %d", i+1)
		}
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			val, _ := strconv.ParseInt(fields[pos+j], 10, 64)
			arr[j] = val
		}
		pos += n
		cases = append(cases, testCase{n: n, m: m, k: k, a: arr})
	}
	return cases, nil
}

// solve mirrors 1197D.go.
func solve(tc testCase) int64 {
	n, m := tc.n, tc.m
	k := tc.k
	a := tc.a
	const INF int64 = 9e18
	mn := make([]int64, m)
	for i := 0; i < m; i++ {
		mn[i] = INF
	}
	mn[0] = 0
	var ans int64
	var prefix int64
	for r := 1; r <= n; r++ {
		prefix += a[r-1]
		cr := r % m
		qr := r / m
		var group2 int64 = INF
		for c := cr; c < m; c++ {
			if mn[c] < group2 {
				group2 = mn[c]
			}
		}
		var group1 int64 = INF
		for c := 0; c < cr; c++ {
			if mn[c] < group1 {
				group1 = mn[c]
			}
		}
		best := group2
		if group1+k < best {
			best = group1 + k
		}
		cur := prefix - k*int64(qr) - best
		if cur > ans {
			ans = cur
		}
		base := prefix - k*int64(qr)
		if base < mn[cr] {
			mn[cr] = base
		}
	}
	if ans < 0 {
		ans = 0
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		return
	}

	for idx, tc := range testcases {
		input := buildInput(tc)
		expected := strconv.FormatInt(solve(tc), 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
