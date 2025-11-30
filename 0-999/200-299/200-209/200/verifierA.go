package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pair struct{ x, y int }

type testCase struct {
	n, m, k int
	seats   []pair
}

// testcasesA contains the contents of testcasesA.txt inline to avoid external dependencies.
const testcasesA = `20
4 4 1
3 4
4 3 4
3 3
2 3
2 2
2 1
5 3 5
5 1
3 1
1 3
3 2
5 1
3 4 3
3 2
3 4
2 3
1 5 1
1 4
1 5 4
1 2
1 1
1 5
1 2
2 5 4
1 1
2 5
2 1
2 5
3 1 3
2 1
3 1
2 1
5 4 3
5 2
3 2
2 2
1 5 3
1 1
1 2
1 1
1 5 4
1 5
1 2
1 5
1 4
4 3 1
3 3
1 4 3
1 2
1 3
1 2
3 2 3
2 1
1 1
3 1
1 5 5
1 1
1 2
1 4
1 3
1 1
5 1 2
2 1
4 1
1 1 1
1 1
1 2 1
1 2
4 2 1
4 1
5 1 4
2 1
3 1
5 1
2 1
`

func loadTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesA))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("no test count found")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing n for test %d", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("bad n for test %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing m for test %d", i+1)
		}
		m, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("bad m for test %d: %w", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing k for test %d", i+1)
		}
		k, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("bad k for test %d: %w", i+1, err)
		}
		seats := make([]pair, k)
		for j := 0; j < k; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("missing x for test %d seat %d", i+1, j+1)
			}
			x, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("bad x for test %d seat %d: %w", i+1, j+1, err)
			}
			if !scan.Scan() {
				return nil, fmt.Errorf("missing y for test %d seat %d", i+1, j+1)
			}
			y, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("bad y for test %d seat %d: %w", i+1, j+1, err)
			}
			seats[j] = pair{x: x, y: y}
		}
		tests = append(tests, testCase{n: n, m: m, k: k, seats: seats})
	}
	return tests, nil
}

// assignSeats is the exact logic from 200A.go, embedded here so the verifier has no external dependencies.
func assignSeats(n, m int, seats []pair) []pair {
	occ := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		occ[i] = make([]bool, m+1)
	}
	res := make([]pair, len(seats))
	for idx, p := range seats {
		x1, y1 := p.x, p.y
		if !occ[x1][y1] {
			occ[x1][y1] = true
			res[idx] = pair{x1, y1}
			continue
		}
		found := false
		var rx, ry int
		for d := 1; !found; d++ {
			for dx := 0; dx <= d && !found; dx++ {
				dy := d - dx
				rows := []int{}
				if dx == 0 {
					if x1 >= 1 && x1 <= n {
						rows = append(rows, x1)
					}
				} else {
					if x1-dx >= 1 {
						rows = append(rows, x1-dx)
					}
					if x1+dx <= n {
						rows = append(rows, x1+dx)
					}
				}
				for _, x2 := range rows {
					cols := []int{}
					if dy == 0 {
						if y1 >= 1 && y1 <= m {
							cols = append(cols, y1)
						}
					} else {
						if y1-dy >= 1 {
							cols = append(cols, y1-dy)
						}
						if y1+dy <= m {
							cols = append(cols, y1+dy)
						}
					}
					for _, y2 := range cols {
						if !occ[x2][y2] {
							found = true
							rx, ry = x2, y2
							break
						}
					}
					if found {
						break
					}
				}
			}
		}
		occ[rx][ry] = true
		res[idx] = pair{rx, ry}
	}
	return res
}

func formatInput(tc testCase) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d %d\n", tc.n, tc.m, tc.k)
	for _, s := range tc.seats {
		fmt.Fprintf(&buf, "%d %d\n", s.x, s.y)
	}
	return buf.String()
}

func parseOutput(out []byte, expected int) ([]pair, error) {
	scan := bufio.NewScanner(bytes.NewReader(out))
	scan.Split(bufio.ScanWords)
	res := make([]pair, 0, expected)
	for i := 0; i < expected; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing output for line %d", i+1)
		}
		x, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid x on line %d", i+1)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("missing output for line %d", i+1)
		}
		y, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid y on line %d", i+1)
		}
		res = append(res, pair{x: x, y: y})
	}
	if scan.Scan() {
		return nil, fmt.Errorf("extra output detected")
	}
	return res, nil
}

func verifyCase(bin string, idx int, tc testCase) error {
	input := formatInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d execution failed: %v\n%s", idx+1, err, out)
	}
	got, err := parseOutput(out, len(tc.seats))
	if err != nil {
		return fmt.Errorf("case %d: %v", idx+1, err)
	}
	expected := assignSeats(tc.n, tc.m, tc.seats)
	for i := range expected {
		if expected[i] != got[i] {
			return fmt.Errorf("case %d line %d mismatch: expected %d %d got %d %d", idx+1, i+1, expected[i].x, expected[i].y, got[i].x, got[i].y)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load embedded tests: %v\n", err)
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range tests {
		if err := verifyCase(bin, i, tc); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
