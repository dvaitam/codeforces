package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// gcd mirrors 1516C.go.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type testCase struct {
	n   int
	arr []int
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
2 3 3
4 6 10 9 20
3 20 2 19
3 14 13 17
4 18 15 17 9
2 1 12
5 11 13 14 17 6
3 8 8 1
3 11 6 5
4 17 18 6 15
5 17 12 19 12 12
5 6 13 15 17 8
5 9 16 17 17 12
5 15 12 19 18 15
5 8 11 6 20 9
5 10 10 17 18 17
5 10 7 16 17 12
2 11 1
3 4 2 19
2 9 19
3 4 17 5
4 8 7 2 14
2 2 12
4 6 8 1 3
2 3 1
2 1 12
4 5 6 6 17
2 13 19
2 8 5
2 1 12
2 10 11
5 1 10 15 18 20
2 9 13
3 16 8 3
4 4 1 15 5
5 16 17 11 5 11
4 9 20 14 1
3 2 9 2
3 6 6 4
5 8 17 2 8 8
5 3 9 3 19 8
4 9 14 9 17
2 5 2
5 14 6 4 17 3
3 4 4 1
3 8 4 7
2 17 15
5 10 18 13 7 7
5 14 17 1 19 19
2 14 17
3 4 16 12
2 17 4
4 10 12 10 1
5 4 4 10 7 1
5 2 14 16 15 7
2 1 10
2 12 10
2 8 16
3 4 19 12
5 15 5 12 13 4
4 4 4 3 20
4 13 7 4 1
5 2 16 10 12 15
3 12 9 16
5 14 16 10 13 8
3 16 20 9
5 3 19 19 4 3
4 6 18 5 14
2 3 2
3 10 13 8
4 15 6 17 10
2 5 18
5 4 11 17 8 17
4 6 6 15 8
5 12 19 5 15 15
2 20 13
3 13 17 2
5 9 13 9 14 16
4 18 11 3 8
3 13 13 1
4 15 17 15 6
2 1 13
3 19 20 13
3 4 13 18
3 9 19 19
3 16 20 5
2 20 14
5 9 17 19 6 15
3 3 12 1
5 18 3 19 16 11
5 9 17 15 1 3
4 6 13 9 5
2 6 16
5 15 10 5 1 10
5 1 12 2 18 13
5 7 10 16 5 16
4 3 9 11 10
4 10 13 17 3
3 13 20 17
3 17 3 10
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
		res = append(res, testCase{n: n, arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1516C.go for one test case.
func solve(tc testCase) string {
	sum := 0
	for _, v := range tc.arr {
		sum += v
	}
	half := sum / 2
	dp := make([]bool, half+1)
	dp[0] = true
	for _, x := range tc.arr {
		for s := half; s >= x; s-- {
			if dp[s-x] {
				dp[s] = true
			}
		}
	}
	if sum%2 == 1 || !dp[half] {
		return "0"
	}
	g := tc.arr[0]
	for i := 1; i < tc.n; i++ {
		g = gcd(g, tc.arr[i])
	}
	best := 31
	idx := -1
	for i, v := range tc.arr {
		x := v / g
		cnt := 0
		for x%2 == 0 {
			x /= 2
			cnt++
		}
		if cnt < best {
			best = cnt
			idx = i
		}
	}
	return fmt.Sprintf("1\n%d", idx+1)
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
