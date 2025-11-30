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
	n   int
	arr []int
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
100
11
11 3 10 8 2 4 7 1 9 6 5
7
6 5 3 1 7 2 4
8
1 4 5 6 3 2 7 8
8
4 2 3 7 5 6 1 8
7
3 6 5 4 2 1 7
7
1 3 6 5 7 4 2
10
3 9 6 2 10 5 7 1 8 4
1
1
6
4 6 3 2 5 1
8
6 3 1 2 7 8 4 5
3
3 1 2
7
7 4 6 2 1 3 5
2
1 2
11
6 11 10 7 9 4 5 8 2 3 1
2
2 1
9
3 6 9 7 4 5 8 2 1
2
2 1
10
4 8 9 6 7 2 10 5 3 1
1
1
6
2 1 5 4 6 3
10
9 8 10 7 6 5 4 3 1 2
5
5 2 1 3 4
9
3 1 4 5 6 7 8 9 2
9
6 9 5 7 8 1 3 2 4
9
2 5 6 7 9 8 1 3 4
9
1 5 6 7 8 4 2 3 9
1
1
1
1
3
1 2 3
7
3 6 2 7 4 1 5
3
1 3 2
9
8 1 5 3 6 9 7 2 4
8
7 6 8 1 4 3 5 2
2
2 1
8
3 4 1 5 6 2 8 7
11
6 4 1 7 10 9 5 3 2 8 11
5
3 5 1 4 2
6
4 1 6 3 2 5
8
2 4 8 5 6 1 7 3
3
3 2 1
5
4 3 1 2 5
10
10 4 5 2 6 1 3 9 8 7
10
1 4 9 5 3 10 8 6 7 2
6
5 6 4 1 2 3
10
4 2 9 5 10 1 7 3 6 8
6
3 6 4 2 1 5
7
5 1 4 3 2 7 6
4
3 4 2 1
3
2 1 3
7
5 6 7 3 4 2 1
11
8 2 4 9 11 5 1 6 10 3 7
10
2 5 10 9 8 6 4 7 3 1
5
5 4 1 2 3
2
2 1
7
7 2 3 5 6 4 1
3
1 3 2
5
2 4 1 3 5
7
7 4 3 6 5 1 2
11
10 11 7 3 6 8 4 2 1 5 9
11
9 11 10 8 5 2 6 3 7 4 1
3
3 2 1
11
10 11 5 7 9 6 8 3 1 4 2
11
8 5 3 1 11 10 9 6 2 7 4
8
5 3 7 6 2 1 8 4
1
1
9
4 6 7 5 2 3 9 1 8
4
4 3 1 2
6
6 1 2 5 3 4
10
7 8 1 3 10 5 4 6 2 9
4
2 1 3 4
4
2 1 4 3
4
1 3 2 4
6
6 1 3 4 2 5
11
4 2 10 3 6 9 8 11 1 7 5
3
3 1 2
9
6 5 9 7 2 3 1 8 4
4
2 3 4 1
6
4 5 1 6 2 3
8
1 7 6 5 4 3 8 2
8
1 3 5 7 2 4 6 8
4
1 2 3 4
6
5 3 4 1 2 6
9
1 2 5 6 4 3 7 9 8
8
4 2 6 7 1 8 5 3
8
3 6 8 7 5 4 2 1
3
1 3 2
1
1
8
4 1 8 7 3 5 6 2
11
3 2 6 10 4 9 1 7 5 8 11
11
3 10 4 9 8 5 7 6 2 11 1
7
4 3 6 7 2 5 1
8
3 8 4 2 1 6 5 7
5
5 2 3 4 1
3
1 3 2
5
3 2 4 5 1
2
2 1
7
3 2 1 4 6 5 7
8
1 3 6 2 4 5 8 7
6
2 1 3 6 5 4
5
3 1 4 2 5
10
5 8 3 2 4 7 10 9 6 1
4
2 1 4 3
1
1
7
7 5 1 6 3 4 2
9
2 3 5 4 6 7 9 1 8
11
6 10 5 2 7 3 4 8 11 9 1
11
11 4 1 10 7 9 8 5 6 3 2
3
3 1 2
5
5 2 1 4 3
6
6 1 5 3 2 4
5
1 3 2 4 5
4
2 1 4 3
10
10 1 7 5 4 3 2 8 9 6
7
1 4 3 2 7 5 6
6
1 4 2 6 5 3
1
1
10
5 6 2 7 3 4 9 1 8 10
4
1 4 2 3
8
7 8 5 6 1 4 3 2
10
9 8 7 3 4 10 2 5 6 1
6
6 4 1 5 3 2
11
6 7 5 2 1 4 10 9 8 3 11
7
7 3 1 2 4 5 6
8
7 3 1 2 4 5 6 8
11
5 7 8 6 4 2 9 3 1 10 11
10
7 10 4 3 2 5 6 1 8 9
8
7 6 8 2 1 3 4 5
9
8 9 7 3 6 4 5 2 1
1
1
4
3 2 4 1
8
1 6 3 4 8 2 5 7
7
7 5 6 4 2 1 3
7
3 5 1 4 7 6 2
11
4 9 5 3 7 11 2 6 10 8 1
11
1 8 5 9 10 11 6 4 2 3 7
7
2 5 1 3 6 7 4
7
2 3 4 7 5 6 1
1
1
5
5 2 1 3 4
6
5 2 4 1 6 3
4
2 4 1 3
11
11 8 9 6 5 4 7 1 3 2 10
7
6 3 2 4 7 1 5
10
9 5 7 6 2 4 8 1 10 3
2
1 2
4
3 2 1 4
9
1 2 3 9 8 5 6 4 7
4
2 1 4 3
1
1
1
1
10
6 1 3 4 2 5 9 10 7 8
5
3 1 4 2 5
9
3 6 4 9 7 8 2 5 1
3
1 3 2
1
1
6
6 2 1 3 4 5
8
2 3 1 7 5 6 4 8
10
6 9 10 5 7 4 1 3 8 2
3
1 3 2
6
1 6 3 5 2 4
1
1
8
1 4 2 3 7 6 8 5
9
6 4 1 3 5 2 8 9 7
6
5 6 4 2 3 1
8
2 7 1 3 4 5 6 8
8
1 3 2 4 5 6 8 7
5
2 3 5 4 1
5
1 2 3 4 5
9
4 8 1 3 7 6 5 2 9
7
1 2 4 5 6 7 3
4
1 4 2 3
3
3 2 1
7
2 3 4 1 5 6 7
8
2 3 1 4 5 7 6 8
`

func isSorted(a []int) bool {
	for i := 0; i+1 < len(a); i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}

// solve mirrors 1558F.go (odd-even sort step count).
func solve(tc testCase) string {
	arr := append([]int(nil), tc.arr...)
	ans := 0
	for !isSorted(arr) {
		if ans%2 == 0 {
			for i := 0; i+1 < tc.n; i += 2 {
				if arr[i] > arr[i+1] {
					arr[i], arr[i+1] = arr[i+1], arr[i]
				}
			}
		} else {
			for i := 1; i+1 < tc.n; i += 2 {
				if arr[i] > arr[i+1] {
					arr[i], arr[i+1] = arr[i+1], arr[i]
				}
			}
		}
		ans++
	}
	return strconv.Itoa(ans)
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	cases := make([]testCase, 0, t)
	idx := 1
	for caseNum := 0; caseNum < t; caseNum++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("case %d missing n", caseNum+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", caseNum+1, err)
		}
		idx++
		if idx >= len(lines) {
			return nil, fmt.Errorf("case %d missing array", caseNum+1)
		}
		vals := strings.Fields(lines[idx])
		if len(vals) != n {
			return nil, fmt.Errorf("case %d expected %d numbers got %d", caseNum+1, n, len(vals))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(vals[i])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value: %v", caseNum+1, err)
			}
			arr[i] = v
		}
		idx++
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
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
