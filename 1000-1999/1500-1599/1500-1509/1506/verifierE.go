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
	q []int
}

// Embedded testcases from testcasesE.txt.
const testcaseData = `
8 2 3 7 8 8 8 8 8
3 3 3 3
18 18 18 18 18 18 18 18 18 18 18 18 18 18 18 18 18 18 18
10 8 8 9 9 9 10 10 10 10 10
6 5 6 6 6 6 6
1 1
10 1 6 6 6 9 9 9 9 9 10
6 6 6 6 6 6 6
2 1 2
9 7 7 7 7 7 7 7 8 9
7 5 5 5 5 5 7 7
14 5 8 13 13 13 13 13 14 14 14 14 14 14 14
4 4 4 4 4
15 8 8 10 11 11 11 11 14 14 14 14 14 14 15 15
20 13 13 14 14 14 14 14 14 14 14 17 17 18 18 20 20 20 20 20 20
19 13 13 13 16 17 17 17 17 17 17 17 17 17 19 19 19 19 19 19
11 8 8 11 11 11 11 11 11 11 11 11
15 5 5 13 13 13 13 13 13 13 14 15 15 15 15 15
9 3 3 7 8 8 8 9 9 9
11 3 5 9 9 11 11 11 11 11 11 11
14 1 9 12 12 13 14 14 14 14 14 14 14 14 14
11 3 4 10 10 10 11 11 11 11 11 11
9 9 9 9 9 9 9 9 9 9
13 8 8 10 10 10 11 11 11 12 12 13 13 13
7 1 7 7 7 7 7 7
17 6 6 7 7 14 14 14 14 14 14 17 17 17 17 17 17 17
10 1 3 7 8 10 10 10 10 10 10
8 3 3 6 6 6 6 7 8
17 3 6 7 11 11 11 11 15 15 15 15 15 15 15 16 17 17
6 1 4 5 5 5 6
9 7 7 7 7 8 8 8 8 9
16 8 13 13 13 13 14 14 14 14 14 14 14 15 15 15 16
18 13 14 14 14 14 14 14 14 14 14 14 14 14 16 17 18 18 18
10 3 9 9 9 9 9 9 10 10 10
11 2 4 7 9 9 9 9 9 10 10 11
12 5 6 10 10 10 10 10 12 12 12 12 12
20 1 3 3 11 15 19 19 19 20 20 20 20 20 20 20 20 20 20 20 20
15 2 11 14 14 14 15 15 15 15 15 15 15 15 15 15
15 15 15 15 15 15 15 15 15 15 15 15 15 15 15 15
3 3 3 3
2 2 2
13 4 4 7 7 7 8 8 8 9 12 12 12 13
2 1 2
3 3 3 3
2 1 2
20 18 18 18 18 18 18 18 18 18 18 18 18 18 20 20 20 20 20 20 20
13 8 8 8 8 12 12 12 12 12 12 12 13 13
14 7 9 9 9 10 12 12 12 12 13 14 14 14 14
18 13 13 18 18 18 18 18 18 18 18 18 18 18 18 18 18 18 18
20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20
20 7 7 10 15 15 15 15 15 15 20 20 20 20 20 20 20 20 20 20 20
10 6 7 7 7 7 10 10 10 10 10
2 2 2
16 3 11 11 11 13 13 13 13 13 16 16 16 16 16 16 16
18 2 12 12 16 16 16 18 18 18 18 18 18 18 18 18 18 18 18
9 8 8 8 9 9 9 9 9 9
20 10 17 17 17 17 17 17 18 18 18 18 18 20 20 20 20 20 20 20 20
8 4 7 7 8 8 8 8 8
6 4 5 5 5 5 6
9 2 4 9 9 9 9 9 9 9
7 4 5 5 6 6 7 7
16 1 2 11 15 15 15 15 15 15 15 15 15 15 16 16 16
2 2 2
6 5 5 5 6 6 6
12 8 11 11 11 11 11 11 11 11 12 12 12
6 6 6 6 6 6 6
7 5 5 5 5 7 7 7
13 11 11 11 11 11 11 12 12 12 12 12 13 13
14 11 11 12 12 13 14 14 14 14 14 14 14 14 14
20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20
18 3 8 8 8 8 8 13 16 16 17 17 17 17 17 17 18 18 18
16 12 12 12 12 13 13 13 13 13 13 13 13 15 15 15 16
3 2 2 3
20 19 19 19 19 19 19 19 19 19 19 19 19 19 19 20 20 20 20 20 20
17 8 16 16 16 17 17 17 17 17 17 17 17 17 17 17 17 17
3 2 3 3
10 2 6 6 6 6 8 8 8 10 10
7 2 4 4 4 7 7 7
3 3 3 3
5 5 5 5 5 5
9 8 8 9 9 9 9 9 9 9
6 6 6 6 6 6 6
5 1 2 4 4 5
8 7 7 7 7 7 7 7 8
9 1 4 5 9 9 9 9 9 9
7 1 6 6 6 7 7 7
6 6 6 6 6 6 6
6 2 6 6 6 6 6
2 2 2
18 2 9 14 15 15 15 15 15 15 16 16 16 16 16 18 18 18 18
2 1 2
1 1
13 2 10 10 10 10 12 13 13 13 13 13 13 13
1 1
12 10 12 12 12 12 12 12 12 12 12 12 12
3 1 3 3
17 16 16 16 16 16 16 17 17 17 17 17 17 17 17 17 17 17
13 8 8 8 8 10 10 11 13 13 13 13 13 13
14 10 10 10 12 12 13 14 14 14 14 14 14 14 14
14 12 12 12 12 12 12 13 13 14 14 14 14 14 14
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, q: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1506E.go for one test case.
func solve(tc testCase) (string, string) {
	n := tc.n
	q := tc.q
	minAns := make([]int, n)
	maxAns := make([]int, n)
	usedMin := make([]bool, n+1)
	usedMax := make([]bool, n+1)

	nextMin := 1
	for i := 0; i < n; i++ {
		if i == 0 || q[i] != q[i-1] {
			minAns[i] = q[i]
			usedMin[q[i]] = true
		} else {
			for nextMin <= n && usedMin[nextMin] {
				nextMin++
			}
			minAns[i] = nextMin
			usedMin[nextMin] = true
		}
	}

	cur := 0
	stack := []int{}
	for i := 0; i < n; i++ {
		if i == 0 || q[i] != q[i-1] {
			maxAns[i] = q[i]
			for x := cur + 1; x < q[i]; x++ {
				if !usedMax[x] {
					stack = append(stack, x)
					usedMax[x] = true
				}
			}
			usedMax[q[i]] = true
			cur = q[i]
		} else {
			idx := len(stack) - 1
			maxAns[i] = stack[idx]
			stack = stack[:idx]
		}
	}

	var sb1 strings.Builder
	var sb2 strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb1.WriteByte(' ')
			sb2.WriteByte(' ')
		}
		sb1.WriteString(strconv.Itoa(minAns[i]))
		sb2.WriteString(strconv.Itoa(maxAns[i]))
	}
	return sb1.String(), sb2.String()
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.q {
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect1, expect2 := solve(tc)
		expect := expect1 + "\n" + expect2
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed:\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
