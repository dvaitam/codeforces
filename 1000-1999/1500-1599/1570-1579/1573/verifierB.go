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

// Embedded testcases from testcasesB.txt.
const testcaseData = `
4 5 3 1 7 2 6 4 8
4 3 1 7 5 6 2 8 4
1 1 2
5 9 1 7 3 5 8 10 4 6 2
3 1 5 3 4 2 6
5 1 3 5 9 7 6 4 8 10 2
6 5 9 3 7 1 11 4 10 8 2 6 12
2 3 1 2 4
1 1 2
1 1 2
5 3 1 7 9 5 10 6 8 2 4
6 3 1 7 2 10 4 6 5 12 11 8 9
4 3 1 7 5 2 6 4 8
1 1 2
4 3 1 7 5 2 6 4 8
3 3 1 5 2 6 4
3 5 1 3 4 2 6
2 1 3 2 4
1 1 2
5 3 1 7 9 5 10 6 8 2 4
5 9 1 7 3 5 8 10 4 6 2
6 11 1 9 3 7 5 13 4 12 2 10 6
3 3 1 5 2 6 4
6 3 1 7 2 12 4 10 6 11 8 9 5
3 5 1 3 4 2 6
2 1 3 2 4
5 7 1 3 9 5 10 6 8 2 4
5 9 1 7 3 5 8 10 4 6 2
3 5 1 3 4 2 6
6 3 1 5 2 11 4 9 6 12 10 8 7
6 11 1 9 3 7 5 13 4 12 2 10 6
5 9 1 7 3 5 8 10 4 6 2
5 9 1 7 3 5 8 10 4 6 2
2 1 3 2 4
4 5 3 1 7 2 6 4 8
2 1 3 2 4
6 3 1 5 2 11 4 9 6 12 10 8 7
5 3 1 7 9 5 10 6 8 2 4
2 1 3 2 4
1 1 2
3 1 3 5 2 4 6
4 3 1 7 5 2 6 4 8
3 3 1 5 2 6 4
3 5 1 3 4 2 6
2 1 3 2 4
1 1 2
3 3 1 5 2 6 4
2 1 3 2 4
3 3 1 5 2 6 4
1 1 2
1 1 2
3 5 1 3 4 2 6
3 5 1 3 4 2 6
2 1 3 2 4
1 1 2
4 3 1 7 5 2 6 4 8
3 3 1 5 2 6 4
5 9 1 7 3 5 8 10 4 6 2
6 3 1 7 2 12 4 10 6 11 8 9 5
2 1 3 2 4
5 9 1 7 3 5 8 10 4 6 2
3 3 1 5 2 6 4
3 3 1 5 2 6 4
4 3 1 7 5 2 6 4 8
3 3 1 5 2 6 4
3 3 1 5 2 6 4
1 1 2
3 5 1 3 4 2 6
6 5 9 3 7 1 11 4 10 8 2 6 12
6 11 1 9 3 7 5 13 4 12 2 10 6
4 5 3 1 7 2 6 4 8
4 3 1 7 5 2 6 4 8
1 1 2
4 3 1 7 5 2 6 4 8
3 3 1 5 2 6 4
4 3 1 7 5 2 6 4 8
5 9 1 7 3 5 8 10 4 6 2
4 3 1 7 5 2 6 4 8
6 3 1 7 2 12 4 10 6 11 8 9 5
6 3 1 5 2 11 4 9 6 12 10 8 7
3 5 1 3 4 2 6
4 3 1 7 5 2 6 4 8
1 1 2
5 3 1 7 9 5 10 6 8 2 4
2 1 3 2 4
6 11 1 9 3 7 5 13 4 12 2 10 6
4 3 1 7 5 2 6 4 8
5 9 1 7 3 5 8 10 4 6 2
6 11 1 9 3 7 5 13 4 12 2 10 6
1 1 2
6 3 1 5 2 11 4 9 6 12 10 8 7
2 1 3 2 4
6 11 1 9 3 7 5 13 4 12 2 10 6
4 3 1 7 5 2 6 4 8
5 9 1 7 3 5 8 10 4 6 2
3 3 1 5 2 6 4
4 3 1 7 5 2 6 4 8
3 5 1 3 4 2 6
`

// solve mirrors 1573B.go.
func solve(tc testCase) int {
	n := tc.n
	posA := make([]int, 2*n+3)
	posB := make([]int, 2*n+3)
	for i, v := range tc.a {
		posA[v] = i + 1
	}
	for i, v := range tc.b {
		posB[v] = i + 1
	}
	const INF = int(1e9)
	bestEven := make([]int, 2*n+5)
	for e := 2 * n; e >= 2; e -= 2 {
		cur := posB[e]
		if cur == 0 {
			cur = INF
		}
		next := bestEven[e+2]
		if next == 0 {
			next = INF
		}
		if cur < next {
			bestEven[e] = cur
		} else {
			bestEven[e] = next
		}
	}
	ans := INF
	for o := 1; o < 2*n; o += 2 {
		pa := posA[o]
		if pa == 0 {
			pa = INF
		}
		pb := bestEven[o+1]
		if pa+pb < ans {
			ans = pa + pb
		}
	}
	return ans - 2
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("case %d too short", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", idx+1, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("case %d expected %d numbers, got %d", idx+1, 1+2*n, len(fields))
		}
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("case %d bad a%d: %v", idx+1, i, err)
			}
			a[i] = v
		}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+n+i])
			if err != nil {
				return nil, fmt.Errorf("case %d bad b%d: %v", idx+1, i, err)
			}
			b[i] = v
		}
		cases = append(cases, testCase{n: n, a: a, b: b})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := strconv.Itoa(solve(tc))
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
