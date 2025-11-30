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
	m int
	a [][]uint64
}

// Embedded logic from 1704D.go.
func solve(tc testCase) (int, uint64) {
	sums := make([]uint64, tc.n)
	for i := 0; i < tc.n; i++ {
		var s uint64
		for j := 0; j < tc.m; j++ {
			s += tc.a[i][j] * uint64(j+1)
		}
		sums[i] = s
	}
	base := majorityValue(sums)
	var idx int
	var diff uint64
	for i, v := range sums {
		if v != base {
			idx = i + 1
			if v > base {
				diff = v - base
			} else {
				diff = base - v
			}
			break
		}
	}
	return idx, diff
}

func majorityValue(arr []uint64) uint64 {
	var cand uint64
	cnt := 0
	for _, v := range arr {
		if cnt == 0 {
			cand = v
			cnt = 1
		} else if v == cand {
			cnt++
		} else {
			cnt--
		}
	}
	return cand
}

// Embedded copy of testcasesD.txt (each line: n m then n*m values).
const testcaseData = `
3 3 0 3 0 0 5 5 4 5 4
3 4 4 2 2 0 5 3 1 5 4 3 4 4
4 4 4 3 4 5 5 4 3 4 3 1 2 5 4 2 4 3
4 4 2 2 2 0 4 0 3 3 2 1 4 3 1 5 3 2
4 4 1 3 3 0 5 0 2 0 2 4 5 0 2 3 0 2
3 3 4 0 1 5 0 2 0 5 5
2 2 5 2 3 4
3 2 0 5 5 5 1 4
4 4 5 2 2 2 3 3 4 3 0 1 3 1 4 0 4 1
4 2 1 5 3 3 1 4 1 5
3 4 5 2 0 5 5 5 2 3 3 1 5 1
2 3 3 4 2 4 3 2
4 2 4 5 2 4 5 2 3 0
3 2 5 3 0 4 1 4
4 4 5 2 3 2 1 0 0 4 4 4 4 1 1 2 0 0
2 2 5 0 3 5
4 2 0 0 0 0 4 2 2 0
4 2 4 1 3 1 2 2 4 4
4 3 1 1 1 3 0 1 4 5 3 0 2 2
3 2 0 4 1 4 5 0
2 2 1 1 2 0
2 3 5 1 0 3 1 1
2 4 2 2 0 4 0 3 1 1
4 2 0 0 1 0 5 5 0 0
4 2 2 5 2 4 3 1 5 0
4 3 1 2 2 2 3 5 4 3 5 3 0 3
2 3 2 1 4 5 0 1
2 3 2 4 2 2 2 3
3 3 2 2 3 3 4 0 2 1 2
2 3 4 1 4 5 5 1
2 3 5 5 1 1 1 0
4 3 1 2 5 2 1 2 3 2 0 3 1 4
3 3 0 1 5 2 0 5 1 3 4
2 2 5 1 5 2
4 3 4 2 4 5 5 2 2 5 0 1 3 0
3 4 5 1 0 1 2 2 4 3 1 2 0 1
3 4 4 0 1 0 1 1 2 4 2 1 1 1
2 3 4 4 4 4 4 4
3 3 4 4 3 1 4 2 1 3 0
3 2 1 1 1 1 0 1
3 3 1 0 2 0 4 1 5 5 1
2 3 5 5 1 4 2 1
4 4 5 5 0 1 2 3 4 0 0 2 3 4 2 1 3 0
4 3 5 5 4 5 2 4 2 4 0 3 2 3
2 3 0 5 5 2 0 1
3 2 2 2 2 2 3 3
2 3 1 5 2 0 0 3
3 4 1 2 2 5 5 4 3 4 2 4 2 5
2 3 1 2 0 3 2 4
4 4 5 1 3 3 1 3 5 1 0 5 5 1 0 5 0 0
4 4 3 4 3 5 2 4 1 4 5 3 3 0 3 4 5 4
4 3 1 4 4 2 0 5 2 2 3 1 5 5
4 4 2 0 3 2 1 1 1 3 0 2 4 2 1 4 3 3
2 4 1 4 0 3 5 1 4 1
3 3 0 2 2 1 1 2 1 1 5
4 4 2 1 4 5 3 3 3 4 4 2 1 4 4 4 2 4
2 3 5 1 5 0 2 0
3 3 5 5 4 5 1 4 3 3 4
3 3 1 0 0 5 0 0 4 3 1
3 3 1 3 3 4 4 0 4 1 4
3 3 3 2 2 1 4 2 1 0 4
2 4 0 4 1 3 2 3 2 5
2 3 0 5 3 2 3 3
3 4 0 1 3 4 3 4 5 3 4 1 2 1
3 2 4 1 1 1 3 5
2 2 5 0 3 0
3 2 2 4 2 2 3 0
3 3 5 3 2 5 5 2 5 4 3
3 3 1 0 3 1 3 1 3 0 4
2 4 2 2 3 5 4 5 2 5
4 3 4 4 3 2 3 5 3 4 1 1 1 4
2 4 1 0 2 4 0 2 3 0
3 3 2 4 4 5 3 1 2 1 2
3 4 3 5 1 0 3 5 2 4 4 1 1 0
4 3 3 1 5 2 0 3 0 2 0 3 1 0
4 4 1 2 1 3 3 2 4 5 5 4 2 1 1 1 1 4
2 2 0 3 4 4
2 3 1 2 4 5 5 5
3 4 1 0 2 1 1 1 5 3 4 3 0 5
2 2 4 3 4 1
2 3 5 1 2 5 2 0
3 3 0 4 3 1 5 1 1 5 0
4 4 2 0 2 4 1 0 0 0 3 2 0 3 5 4 4 5
4 3 5 2 1 3 2 5 2 3 3 3 2 4
2 2 0 1 1 1
2 2 5 3 3 4
3 4 0 0 1 4 0 0 3 2 3 1 1 5
4 3 3 2 4 5 0 4 3 1 5 4 2 4
2 3 0 1 5 5 0 3
2 2 2 1 1 2
2 3 5 0 3 0 4 3
2 3 4 4 4 0 1 4
4 4 3 5 4 4 3 1 4 3 3 5 2 3 0 0 1 4
4 4 2 0 0 2 0 1 0 5 5 4 0 0 4 3 1 0
3 3 4 0 4 3 4 2 3 5 0
4 4 0 2 4 3 3 1 2 2 3 3 4 0 2 4 1 5
4 3 3 2 4 4 1 2 4 5 3 5 5 3
4 4 4 3 3 5 1 2 0 0 1 3 0 2 2 2 4 2
2 2 4 1 3 0
3 3 5 4 2 4 2 1 5 0 4
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: not enough data", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		if len(fields) != 2+n*m {
			return nil, fmt.Errorf("line %d: expected %d matrix values got %d", i+1, n*m, len(fields)-2)
		}
		a := make([][]uint64, n)
		idx := 2
		for r := 0; r < n; r++ {
			a[r] = make([]uint64, m)
			for c := 0; c < m; c++ {
				val, err := strconv.ParseUint(fields[idx], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: bad value at %d: %v", i+1, idx-1, err)
				}
				a[r][c] = val
				idx++
			}
		}
		tests = append(tests, testCase{n: n, m: m, a: a})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expectedIdx int, expectedDiff uint64) error {
	var input strings.Builder
	input.WriteString("1\n")
	fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
	for r := 0; r < tc.n; r++ {
		for c := 0; c < tc.m; c++ {
			if c > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatUint(tc.a[r][c], 10))
		}
		input.WriteByte('\n')
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
	got := strings.Fields(strings.TrimSpace(out.String()))
	if len(got) != 2 {
		return fmt.Errorf("expected two outputs got %d (%v)", len(got), got)
	}
	idxVal, err := strconv.Atoi(got[0])
	if err != nil {
		return fmt.Errorf("bad idx output %q", got[0])
	}
	diffVal, err := strconv.ParseUint(got[1], 10, 64)
	if err != nil {
		return fmt.Errorf("bad diff output %q", got[1])
	}
	if idxVal != expectedIdx || diffVal != expectedDiff {
		return fmt.Errorf("expected %d %d got %d %d", expectedIdx, expectedDiff, idxVal, diffVal)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		idx, diff := solve(tc)
		if err := runCase(bin, tc, idx, diff); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
