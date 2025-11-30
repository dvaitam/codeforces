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
	a [][]float64
}

// solve is the embedded logic from 16E.go.
func solve(tc testCase) []float64 {
	n := tc.n
	a := tc.a
	size := 1 << n
	d := make([]float64, size)
	ans := make([]float64, n)
	full := size - 1
	d[full] = 1.0
	s := make([]int, n)
	for mask := full; mask > 0; mask-- {
		probCur := d[mask]
		if probCur == 0 {
			continue
		}
		cnt := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				s[cnt] = i
				cnt++
			}
		}
		if cnt == 1 {
			ans[s[0]] = probCur
		} else {
			totalPairs := float64(cnt * (cnt - 1))
			fac := 2.0 / totalPairs
			for x := 0; x < cnt-1; x++ {
				i := s[x]
				for y := x + 1; y < cnt; y++ {
					j := s[y]
					d[mask^(1<<i)] += probCur * fac * a[j][i]
					d[mask^(1<<j)] += probCur * fac * a[i][j]
				}
			}
		}
	}
	return ans
}

// Embedded copy of testcasesE.txt (each line: n then n*n probabilities).
const testcaseData = `
2 0.00 0.34 0.66 0.00
1 0.00
1 0.00
4 0.00 0.31 0.62 0.92 0.69 0.00 0.33 0.78 0.38 0.67 0.00 0.13 0.08 0.22 0.87 0.00
1 0.00
2 0.00 0.36 0.64 0.00
4 0.00 0.44 0.69 0.35 0.56 0.00 0.85 0.41 0.31 0.15 0.00 0.58 0.65 0.59 0.42 0.00
1 0.00
4 0.00 0.13 0.02 0.80 0.87 0.00 0.86 0.39 0.98 0.14 0.00 0.48 0.20 0.61 0.52 0.00
4 0.00 0.64 0.10 0.06 0.36 0.00 0.84 0.60 0.90 0.16 0.00 0.30 0.94 0.40 0.70 0.00
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
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d: missing n", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(fields) != 1+n*n {
			return nil, fmt.Errorf("line %d: expected %d values, got %d", i+1, 1+n*n, len(fields))
		}
		a := make([][]float64, n)
		idx := 1
		for r := 0; r < n; r++ {
			a[r] = make([]float64, n)
			for c := 0; c < n; c++ {
				val, err := strconv.ParseFloat(fields[idx], 64)
				if err != nil {
					return nil, fmt.Errorf("line %d: bad float at %d: %v", i+1, idx, err)
				}
				a[r][c] = val
				idx++
			}
		}
		tests = append(tests, testCase{n: n, a: a})
	}
	return tests, nil
}

func formatAns(ans []float64) string {
	var b strings.Builder
	for i, v := range ans {
		fmt.Fprintf(&b, "%.12f", v)
		if i+1 < len(ans) {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%.2f", tc.a[i][j])
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := formatAns(solve(tc))
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
