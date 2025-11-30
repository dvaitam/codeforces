package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// solve mirrors 1486B.go.
func solve(points [][2]int64) int64 {
	n := len(points)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i, p := range points {
		xs[i] = p[0]
		ys[i] = p[1]
	}
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	if n%2 == 1 {
		return 1
	}
	xRange := xs[n/2] - xs[n/2-1] + 1
	yRange := ys[n/2] - ys[n/2-1] + 1
	return xRange * yRange
}

type testCase struct {
	points [][2]int64
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
3 9 1 4 1 7 7
8 10 6 3 1 7 0 6 6 9 0 7 4 3 9 1 5
1 0 0
9 0 6 10 3 6 0 8 3 7 7 8 3 5 3 10 3 7 4
1 6 8
2 2 10 4 1
6 8 6 8 10 3 4 4 9 7 8 6 9
1 7 3
7 6 10 2 5 8 10 5 1 7 10 8 1 2 8
7 5 7 0 7 0 4 9 9 9 6 10 2 2 8
4 0 3 8 8 3 6 8 5
10 5 7 4 10 8 9 0 6 8 2 8 8 3 6 0 7 5 9 8 3
9 6 7 5 6 5 0 8 8 9 9 5 7 9 0 3 10 2 8
10 2 1 8 4 0 10 1 1 0 7 0 4 3 4 1 9 2 5 4 1
3 2 4 8 2 10 4
5 7 5 7 7 1 0 4 6 5 6
4 4 1 4 8 3 9 6 0
4 0 6 2 0 2 7 8 10
7 8 3 10 8 7 3 8 10 0 6 10 9 5 10
7 0 4 2 3 0 4 1 1 4 4 2 6 9 4
3 0 8 0 9 3 9
8 2 9 8 0 6 3 5 1 3 9 10 6 9 3 7 1
7 4 8 7 0 5 9 6 4 0 2 3 5 9 2
6 6 3 4 10 1 6 8 5 10 8 7 8
4 1 0 1 2 2 2 8 3
5 5 9 8 4 5 5 5 1 4 3
10 7 2 9 8 1 5 0 6 1 6 2 2 5 1 9 9 6 1 9 8
4 9 1 4 5 4 9 8 1
8 4 1 0 4 0 9 10 0 1 6 1 0 3 3 9 6
3 1 7 2 10 3 2
2 6 6 8 4
9 4 7 5 1 3 10 5 0 0 0 4 9 5 7 6 5 6 1
2 5 9 7 1
5 3 9 8 7 10 5 4 2 8 3
5 3 3 5 1 4 1 7 1 10 9
6 3 6 4 0 5 2 5 9 4 3 5 1
9 9 9 9 1 3 3 0 3 6 1 4 8 1 1 0 10 0 4
6 7 7 2 1 8 5 1 8 10 2 2 2
3 5 4 1 8 9 4
3 3 2 8 0 5 9
9 3 2 4 6 8 2 0 10 3 4 1 10 7 6 8 4 8 7
9 7 0 6 5 2 4 7 0 10 6 9 0 0 5 9 2 9 2
3 4 4 6 9 6 2
10 1 3 7 0 2 8 5 8 10 7 10 10 3 3 5 7 10 7 3 6
6 8 9 10 4 10 3 0 1 8 10 5 2
9 3 4 4 4 8 5 2 7 9 1 1 9 8 9 6 2 2 4
7 3 9 0 7 10 6 10 5 6 8 2 8 0 8
2 4 10 1 4
2 2 9 10 10
2 7 3 6 6
7 2 5 7 2 9 7 3 1 6 9 8 6 1 10
5 4 3 6 8 0 3 8 7 9 0
1 10 9
4 4 3 2 4 2 8 3 4
5 9 4 10 7 2 8 5 7 6 1
4 9 6 3 4 1 0 1 9
1 8 4
3 1 8 5 9 4 6
9 10 5 8 5 0 1 7 7 5 4 8 6 5 10 9 7 1 10
7 6 3 8 0 4 10 9 8 3 7 9 8 6 4
3 7 9 10 8 3 5
9 0 10 6 9 6 6 5 9 9 1 7 3 10 10 4 10 0 6
3 10 6 4 2 1 9
1 5 4
7 10 8 4 2 7 4 7 2 7 8 0 4 8 1
10 6 1 5 1 10 7 0 2 8 2 1 6 10 4 9 4 3 8 3 3
6 4 1 1 8 10 5 7 8 8 0 2 4
9 4 5 9 3 6 8 6 2 7 4 9 5 3 4 9 3 10 0
10 6 5 6 3 4 3 1 10 2 9 7 9 2 9 4 7 8 2 2 2
8 5 4 6 3 1 3 10 4 1 1 3 6 5 7 1 2
1 0 9
1 3 10
1 7 8
10 7 5 10 4 1 9 2 1 3 6 3 7 7 6 2 3 3 4 7 8
10 6 3 7 4 5 7 9 1 3 1 0 0 0 7 5 6 9 4 3 6
3 10 2 0 0 6 2
9 0 9 6 4 2 1 7 10 4 0 0 8 0 8 2 0 4 1
7 1 3 0 7 10 2 4 10 3 10 7 6 5 10
5 4 10 10 3 3 0 9 9 2 5
7 9 8 10 8 0 5 8 6 8 3 8 6 10 1
5 9 1 4 2 1 2 0 3 6 0
1 10 1
9 7 8 5 1 5 0 2 8 0 7 10 2 6 7 0 8 4 1
5 5 1 4 0 6 0 4 5 2 4
7 1 10 4 1 6 3 8 8 3 5 5 8 6 9
8 1 2 10 7 8 8 9 8 8 0 4 2 3 5 6 8
6 1 6 5 2 9 1 0 4 10 8 5 6
5 5 5 4 5 8 8 0 8 1 2
6 5 5 9 1 7 4 7 7 5 6 1 9
1 2 0
9 7 9 4 3 9 5 5 10 5 6 4 7 9 5 8 8 2 0
3 4 10 3 9 2 1
3 6 9 0 1 8 10
5 1 3 4 1 10 9 8 10 1 1
4 10 2 8 6 0 9 5 7
5 3 3 9 7 3 6 7 10 5 8
4 7 1 4 6 3 0 8 6
9 7 1 6 9 8 9 9 6 0 5 7 0 3 4 10 0 8 1
5 8 5 8 10 9 8 4 8 6 8
9 6 9 10 9 4 7 4 2 8 7 9 2 8 2 4 10 0 6
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	pos := 0
	res := make([]testCase, 0, 100)
	for pos < len(fields) {
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("bad n at token %d: %v", pos+1, err)
		}
		pos++
		if pos+2*n > len(fields) {
			return nil, fmt.Errorf("not enough coordinates for n=%d", n)
		}
		pts := make([][2]int64, n)
		for i := 0; i < n; i++ {
			x, err := strconv.ParseInt(fields[pos+2*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad x at case starting %d point %d: %v", pos, i+1, err)
			}
			y, err := strconv.ParseInt(fields[pos+2*i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad y at case starting %d point %d: %v", pos, i+1, err)
			}
			pts[i] = [2]int64{x, y}
		}
		pos += 2 * n
		res = append(res, testCase{points: pts})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.points)))
		for _, p := range tc.points {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		expected := strconv.FormatInt(solve(tc.points), 10)
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
