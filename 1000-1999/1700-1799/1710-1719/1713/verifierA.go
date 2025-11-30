package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type testCase struct {
	n      int
	points []point
}

// solve contains the logic from 1713A.go.
func solve(tc testCase) string {
	maxRight, maxLeft := 0, 0
	maxUp, maxDown := 0, 0
	for _, p := range tc.points {
		if p.x > 0 && p.x > maxRight {
			maxRight = p.x
		} else if p.x < 0 && -p.x > maxLeft {
			maxLeft = -p.x
		}
		if p.y > 0 && p.y > maxUp {
			maxUp = p.y
		} else if p.y < 0 && -p.y > maxDown {
			maxDown = -p.y
		}
	}
	ans := 2 * (maxRight + maxLeft + maxUp + maxDown)
	return strconv.Itoa(ans)
}

// Embedded copy of testcasesA.txt.
const testcaseData = `
4 0 1 3 0 -1 0 0 0
5 0 -2 0 -1 -4 0 0 -1 0 3
5 0 -1 -4 0 0 5 3 0 1 0
3 0 -2 0 2 3 0
3 3 0 0 -4 0 1
1 0 0
2 0 -4 4 0
2 -3 0 0 2
1 0 0
5 0 2 3 0 -4 0 0 3 4 0
5 0 2 1 0 -2 0 -2 0 0 -5
5 0 -1 -4 0 0 -3 0 -5 0 3
4 0 3 -2 0 0 5 0 1
5 2 0 0 0 4 0 4 0 0 -2
2 -1 0 -2 0
3 0 0 -5 0 -3 0
2 4 0 0 3
5 0 -5 -2 0 0 4 -4 0 -4 0
1 0 -2
2 0 2 -5 0
1 0 4
1 0 -4
2 -1 0 -3 0
1 0 -5
5 1 0 0 0 0 2 0 4 5 0
2 0 -5 0 -3
2 -1 0 2 0
2 5 0 4 0
5 0 5 5 0 3 0 0 2 0 0
1 0 -3
2 0 2 -1 0
3 0 5 0 -3 0 1
4 0 -4 -2 0 0 -3 5 0
4 5 0 0 1 4 0 5 0
1 -4 0
3 0 2 0 2 0 4
1 0 2
3 2 0 1 0 5 0
1 0 -3
1 5 0
4 -2 0 -5 0 0 5 0 -4
2 5 0 -1 0
3 0 -3 1 0 0 -5
3 0 -4 0 -3 0 5
3 -3 0 -5 0 -2 0
3 0 0 0 -5 0 4
4 0 2 0 0 0 -3 4 0
3 -3 0 0 0 0 -4
3 0 -5 -3 0 4 0
3 3 0 -4 0 -2 0
1 3 0
1 0 0
3 -4 0 0 2 0 0
1 2 0
4 0 0 0 -3 0 5 0 5
1 -4 0
2 0 -5 -4 0
4 0 -1 2 0 0 5 -4 0
3 -1 0 0 -3 0 0
1 -5 0
5 5 0 2 0 -2 0 0 -2 0 -4
2 0 0 0 -3
1 0 2
2 0 5 0 3
4 0 0 0 2 0 -2 0 -2
1 0 0
3 -3 0 0 4 0 1
5 2 0 -4 0 0 -5 -3 0 -5 0
4 -3 0 0 5 0 0 0 3
5 -4 0 0 3 0 -4 0 -2 4 0
4 0 2 0 1 0 -2 0 -5
1 0 -1
5 0 0 -1 0 0 -1 0 1 0 -5
2 0 -3 0 0
1 0 2
4 4 0 0 5 0 0 4 0
4 -5 0 -3 0 4 0 0 5
3 3 0 0 -2 2 0
1 0 4
4 0 0 0 5 0 -1 0 -3
4 0 5 0 -4 0 -2 1 0
4 -2 0 5 0 -5 0 0 -5
4 -5 0 0 5 0 -2 5 0
5 0 1 0 -4 4 0 -4 0 0 1
1 -3 0
1 0 1
4 0 0 0 -4 -4 0 -5 0
3 -5 0 0 -2 0 -4
5 0 -2 5 0 0 -4 0 -1 -5 0
5 -3 0 -4 0 -1 0 -2 0 2 0
3 3 0 0 -3 0 -4
1 0 -1
2 -3 0 0 -2
3 0 -2 0 -1 5 0
1 0 4
1 -4 0
2 3 0 -3 0
5 5 0 -2 0 0 0 0 3 1 0
1 0 0
4 -1 0 0 -4 0 -4 3 0
`

var expectedOutputs = []string{
	"10",
	"18",
	"26",
	"14",
	"16",
	"0",
	"16",
	"10",
	"0",
	"22",
	"20",
	"24",
	"14",
	"12",
	"4",
	"10",
	"14",
	"26",
	"4",
	"14",
	"8",
	"8",
	"6",
	"10",
	"18",
	"10",
	"6",
	"10",
	"20",
	"6",
	"6",
	"16",
	"22",
	"12",
	"8",
	"8",
	"4",
	"10",
	"6",
	"10",
	"28",
	"12",
	"12",
	"18",
	"10",
	"18",
	"18",
	"14",
	"24",
	"14",
	"6",
	"0",
	"12",
	"4",
	"16",
	"8",
	"18",
	"24",
	"8",
	"10",
	"22",
	"6",
	"4",
	"10",
	"8",
	"0",
	"14",
	"24",
	"16",
	"30",
	"14",
	"2",
	"14",
	"6",
	"4",
	"18",
	"28",
	"10",
	"8",
	"16",
	"20",
	"30",
	"34",
	"26",
	"6",
	"2",
	"18",
	"18",
	"28",
	"12",
	"14",
	"2",
	"10",
	"14",
	"8",
	"8",
	"12",
	"20",
	"0",
	"16",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d: empty", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("line %d: expected %d coordinate values, got %d", i+1, 2*n, len(fields)-1)
		}
		points := make([]point, n)
		for j := 0; j < n; j++ {
			x, err := strconv.Atoi(fields[1+2*j])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad x%d: %v", i+1, j, err)
			}
			y, err := strconv.Atoi(fields[2+2*j])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad y%d: %v", i+1, j, err)
			}
			points[j] = point{x: x, y: y}
		}
		tests = append(tests, testCase{n: n, points: points})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for _, p := range tc.points {
		input.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := expectedOutputs[i]
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
