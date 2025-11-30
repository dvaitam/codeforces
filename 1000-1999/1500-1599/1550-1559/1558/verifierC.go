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

// Embedded testcases from testcasesC.txt.
const testcaseData = `
100
3
2 3 1
7
6 1 4 5 3 2 7
5
1 2 5 4 3
9
8 1 2 6 7 4 9 3 5
3
3 1 2
1
1
3
2 3 1
5
1 2 3 5 4
5
1 4 2 5 3
5
1 4 2 3 5
3
2 1 3
5
1 5 2 4 3
9
7 2 3 9 1 8 4 6 5
3
2 1 3
5
4 2 3 5 1
5
5 3 2 1 4
7
7 4 3 2 5 6 1
3
3 1 2
1
1
7
4 6 2 1 3 5 7
1
1
7
5 3 4 6 7 1 2
9
8 5 2 3 9 6 1 7 4
9
2 3 5 1 7 4 6 8 9
9
8 5 4 9 2 6 3 7 1
5
5 1 3 2 4
7
7 3 2 1 4 6 5
7
6 1 3 4 5 2 7
7
6 4 1 7 2 3 5
9
4 1 2 6 3 8 7 5 9
3
1 3 2
7
5 3 6 1 2 4 7
7
7 4 6 3 1 5 2
1
1
3
2 1 3
5
2 1 3 4 5
5
4 1 2 3 5
7
4 6 3 2 1 5 7
7
7 4 5 1 6 2 3
5
2 4 1 5 3
5
3 4 1 2 5
7
3 4 7 5 6 2 1
7
5 3 7 1 6 2 4
5
2 1 4 3 5
7
6 4 2 1 3 7 5
5
5 1 3 2 4
9
7 2 9 5 4 3 1 6 8
9
2 8 3 5 1 4 9 6 7
9
2 7 9 8 3 6 1 5 4
7
6 7 2 4 5 1 3
9
9 2 3 5 8 1 6 7 4
7
2 4 5 7 6 3 1
1
1
7
1 5 2 6 3 4 7
7
7 5 4 3 2 6 1
7
4 3 6 1 7 5 2
5
1 3 4 2 5
3
3 2 1
5
4 2 5 1 3
9
7 6 2 9 8 1 5 3 4
5
1 2 3 4 5
9
2 3 1 4 5 6 9 7 8
9
4 2 3 6 7 8 9 5 1
5
5 2 1 4 3
5
2 1 5 4 3
9
6 4 8 7 3 9 2 5 1
9
5 7 9 4 8 6 1 3 2
1
1
5
3 2 1 4 5
7
5 4 2 3 1 6 7
1
1
9
5 2 3 4 7 6 1 9 8
1
1
7
1 3 4 5 6 7 2
7
2 7 4 5 1 6 3
3
1 2 3
7
2 1 3 6 5 4 7
1
1
1
1
3
1 2 3
5
1 4 5 3 2
5
3 4 1 5 2
7
6 5 2 3 4 7 1
9
8 2 1 3 4 5 6 7 9
3
2 3 1
5
3 5 4 1 2
1
1
7
6 5 7 3 4 1 2
5
5 4 3 1 2
5
4 3 2 1 5
5
1 4 3 2 5
9
8 6 5 2 7 4 9 1 3
5
3 1 2 4 5
7
6 5 3 4 2 1 7
7
4 2 6 1 7 3 5
7
1 4 3 6 2 7 5
5
1 3 4 2 5
1
1
7
6 1 2 3 4 5 7
9
4 5 7 9 1 6 3 2 8
7
1 6 3 4 5 7 2
9
4 2 1 6 7 5 3 8 9
7
4 2 3 1 5 7 6
5
5 1 2 3 4
9
1 3 4 2 6 9 7 5 8
3
3 2 1
5
3 2 5 1 4
3
3 1 2
7
3 1 4 5 6 7 2
5
1 5 3 4 2
5
3 2 1 4 5
1
1
7
5 4 2 3 1 7 6
3
3 2 1
7
2 1 4 3 5 6 7
3
1 3 2
3
2 3 1
1
1
9
4 2 5 3 7 1 6 8 9
5
1 5 3 4 2
7
1 2 3 4 5 6 7
5
4 3 5 1 2
5
1 2 5 4 3
5
1 3 2 4 5
3
1 2 3
5
1 3 4 5 2
7
6 4 2 1 7 3 5
3
2 3 1
3
1 2 3
9
2 4 8 1 6 9 7 3 5
1
1
1
1
5
3 4 1 5 2
3
1 2 3
1
1
7
5 2 1 4 3 6 7
1
1
5
2 4 1 3 5
1
1
5
2 3 4 5 1
5
3 4 2 5 1
3
3 1 2
9
8 3 5 7 6 1 9 2 4
3
1 2 3
1
1
5
1 3 4 2 5
1
1
3
1 3 2
3
1 3 2
3
2 1 3
7
5 1 2 3 4 6 7
3
1 2 3
7
5 1 3 4 2 6 7
3
2 1 3
7
1 2 3 4 5 6 7
5
3 1 4 2 5
1
1
5
1 4 3 2 5
7
5 1 3 4 2 6 7
3
2 1 3
1
1
5
1 4 3 2 5
5
1 3 2 4 5
5
1 4 3 5 2
5
1 4 3 2 5
1
1
5
1 3 2 5 4
5
1 3 2 4 5
5
1 3 2 4 5
5
1 2 3 4 5
5
1 2 4 5 3
5
1 3 4 2 5
5
1 3 2 4 5
5
1 3 2 4 5
5
1 2 4 5 3
`

var (
	p   []int
	vec []int
	curN int
)

func getpos(x int) int {
	for i := 1; i <= curN; i++ {
		if p[i] == x {
			return i
		}
	}
	panic("position not found")
}

func applyOp(x int) {
	vec = append(vec, x)
	for i, j := 1, x; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
}

// solveCase mirrors 1558C.go.
func solveCase(arr []int) string {
	curN = len(arr)
	p = make([]int, curN+1)
	ok := true
	for i := 1; i <= curN; i++ {
		p[i] = arr[i-1]
		if (p[i]-i)%2 != 0 {
			ok = false
		}
	}
	if !ok {
		return "-1"
	}
	vec = make([]int, 0, curN*5)
	for k := curN / 2; k >= 1; k-- {
		pos := getpos(2*k + 1)
		applyOp(pos)
		pos = getpos(2 * k)
		applyOp(pos - 1)
		pos = getpos(2 * k)
		applyOp(pos + 1)
		applyOp(3)
		applyOp(2*k + 1)
	}

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(vec)))
	if len(vec) > 0 {
		sb.WriteByte('\n')
		for i, v := range vec {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	}
	return sb.String()
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
		numsStr := strings.Fields(lines[idx])
		if len(numsStr) != n {
			return nil, fmt.Errorf("case %d expected %d numbers got %d", caseNum+1, n, len(numsStr))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(numsStr[i])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value: %v", caseNum+1, err)
			}
			arr[i] = val
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

func tokensEqual(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
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
		expect := solveCase(tc.arr)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !tokensEqual(expect, got) {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
