package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	q int
	c [][]int
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
2 2 2 4 4 3
3 1 5 4 1 1 2 1 5 5 3
3 1 4 2 3 2 2 4 2 5 5
4 1 1 1 4 5 1 2 5 1 4 3 3 5 2 1 3 2
3 2 4 4 5 3 3 1 4 2 1
3 2 5 4 4 2 5 4 2 1 5
1 1 4
3 1 4 3 3 1 3 5 5 1 1
1 1 2
1 2 4
4 1 3 2 1 4 5 2 1 2 4 1 1 5 4 2 5 3
4 2 5 1 4 5 2 5 2 2 5 2 1 5 4 3 3 3
2 2 1 1 5 1
1 2 2
2 1 1 3 1 3
2 2 5 1 1 3
2 2 4 3 4 4
1 1 3
4 2 4 5 4 5 4 5 1 2 2 4 5 3 2 1 2 2
1 2 5
1 1 5
3 2 3 4 4 5 3 5 4 5 3
2 1 4 3 3 5
1 1 5
1 2 1
3 2 5 2 5 1 5 5 1 1 5
3 1 3 2 3 3 4 4 4 1 2
4 2 2 1 4 5 2 4 3 5 4 2 1 4 5 5 2 5
4 2 5 3 4 5 2 5 4 4 2 3 4 3 2 2 3 2
3 2 2 5 3 5 2 3 2 2 3
2 1 4 4 2 2
3 1 1 5 4 2 3 4 3 5 2
3 2 3 1 1 2 2 5 2 5 4
1 2 5
3 2 4 3 4 1 1 1 3 2 3
4 2 5 4 4 3 2 1 1 1 5 1 5 5 3 5 1 5
2 1 1 1 3 1
4 1 5 5 5 2 3 3 2 2 5 2 2 2 2 3 4 3
2 2 2 1 2 4
2 1 5 1 3 5
4 1 3 1 5 3 5 2 1 2 2 5 1 4 2 2 4 5
1 2 2
2 2 3 4 2 3
2 1 3 4 3 3
3 1 3 5 4 1 2 4 1 5 1
2 2 2 4 3 5
2 1 1 3 4 3
4 1 3 4 4 5 2 5 5 2 1 3 5 4 5 5 4 2
4 1 4 1 1 5 4 1 5 4 2 1 2 2 2 1 5 4
4 1 5 2 5 2 2 4 1 2 1 5 2 5 2 5 5 5
4 2 3 4 5 1 3 3 2 1 2 4 2 5 2 1 2 2
4 2 1 5 3 3 4 3 3 1 5 1 4 3 5 5 1 2
2 2 2 1 4 3
2 2 2 1 5 3
3 1 3 2 2 3 5 4 4 4 1
2 2 1 3 2 5
2 2 3 2 5 2
4 2 2 1 2 2 4 5 4 2 1 3 4 5 1 3 1 1
4 1 5 1 1 5 1 4 3 5 3 5 4 2 2 1 5 5
2 2 2 1 2 2
2 1 5 5 2 1
1 2 1
3 2 5 3 5 4 1 1 2 4 1
3 2 3 4 4 1 3 2 2 1 5
2 1 5 2 4 3
1 2 4
2 2 5 5 4 2
3 1 4 3 2 3 3 1 3 2 4
1 2 4
1 2 5
4 2 5 5 5 2 2 1 3 1 4 3 2 4 4 3 5 4
2 2 3 2 1 2
3 2 1 1 4 3 3 5 2 3 4
2 1 2 3 4 4
2 2 4 4 4 3
1 1 5
4 1 1 1 5 4 4 2 4 1 4 1 3 5 2 3 3 5
1 2 1
2 1 1 3 3 4
2 1 3 5 4 3
2 2 4 1 1 1
1 2 2
3 1 5 3 1 1 5 5 3 1 2
1 1 4
2 1 4 3 2 4
1 2 1
1 1 3
1 2 2
3 2 2 4 5 5 5 3 2 4 1
4 1 3 3 5 3 4 4 3 3 3 2 2 2 4 3 1 3
4 2 2 1 1 4 5 2 2 4 5 3 5 2 1 4 4 4
4 2 1 2 5 2 1 4 2 1 5 5 4 5 3 2 1 4
1 2 3
1 1 5
3 2 4 5 3 1 1 1 3 4 4
4 1 3 4 4 5 1 4 3 4 2 4 4 4 5 5 3 4
3 2 4 5 2 3 1 4 4 4 3
1 1 3
2 2 1 5 3 2
2 2 1 1 2 5
`

func parseTestcases() ([]testCase, error) {
	data := strings.TrimSpace(testcaseData)
	if data == "" {
		return nil, fmt.Errorf("no test data")
	}
	lines := strings.Split(data, "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d missing n/q", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		q, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad q: %v", i+1, err)
		}
		expect := 2 + n*n
		if len(fields) != expect {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, expect, len(fields))
		}
		grid := make([][]int, n)
		pos := 2
		for r := 0; r < n; r++ {
			grid[r] = make([]int, n)
			for c := 0; c < n; c++ {
				v, err := strconv.Atoi(fields[pos])
				if err != nil {
					return nil, fmt.Errorf("case %d bad cell %d: %v", i+1, pos-1, err)
				}
				grid[r][c] = v
				pos++
			}
		}
		res = append(res, testCase{n: n, q: q, c: grid})
	}
	return res, nil
}

// solve mirrors 1500D.go and returns the expected output string.
func solve(tc testCase) string {
	n, q := tc.n, tc.q
	res := make([]int, n)
	for k := 1; k <= n; k++ {
		count := 0
		for i := 0; i <= n-k; i++ {
			for j := 0; j <= n-k; j++ {
				colors := make(map[int]struct{})
				ok := true
				for x := 0; x < k && ok; x++ {
					for y := 0; y < k && ok; y++ {
						colors[tc.c[i+x][j+y]] = struct{}{}
						if len(colors) > q {
							ok = false
						}
					}
				}
				if ok {
					count++
				}
			}
		}
		res[k-1] = count
	}
	parts := make([]string, n)
	for i, v := range res {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, " ")
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", tc.c[i][j])
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	output, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("runtime error: %v\n%s", err, string(ee.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func validateOutput(tc testCase, out string, expect string) bool {
	fields := strings.Fields(out)
	if len(fields) != tc.n {
		return false
	}
	for _, f := range fields {
		if _, err := strconv.Atoi(f); err != nil {
			return false
		}
	}
	return out == expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		if !validateOutput(tc, got, expect) {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
