package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `
1 8.80 13.77 4.12
2 5.12 9.94 9.96 8.13 11.55 2.39
1 5.83 9.01 9.83
3 1.07 4.96 7.28 5.96 8.97 5.62 2.67 6.33 2.52
1 7.93 7.99 9.03
3 4.81 5.35 5.07 5.84 7.11 5.38 7.76 12.37 6.05
1 7.81 10.35 8.97
1 8.27 12.52 8.91
3 4.09 7.39 4.94 6.84 8.16 6.29 9.20 10.05 4.97
1 9.16 13.81 1.78
3 4.05 5.22 8.86 0.39 3.49 3.42 5.08 7.83 4.55
2 1.81 4.52 2.44 8.52 12.68 2.29
1 9.62 12.01 6.10
2 5.56 6.89 3.07 1.11 1.81 8.31
1 6.58 6.87 6.99
3 5.59 5.62 8.76 5.58 9.36 5.41 6.90 11.56 6.04
2 8.73 12.15 3.18 4.98 9.90 3.47
1 4.96 9.69 5.58
2 0.92 2.44 4.80 7.39 12.21 8.75
1 2.32 7.00 7.68
2 6.60 10.51 3.55 6.12 10.28 2.08
2 5.00 7.36 6.91 3.74 8.32 4.89
2 9.73 13.56 7.43 2.83 3.36 9.13
1 4.46 6.18 3.19
1 9.30 10.99 3.06
3 3.84 5.78 5.62 9.80 14.68 6.09 6.18 9.56 5.52
2 2.33 4.92 1.01 4.72 6.69 9.55
2 9.84 12.72 1.36 0.93 1.94 3.94
1 1.81 5.30 1.25
1 0.08 2.75 2.17
3 6.06 10.62 1.73 4.66 8.90 5.68 6.21 8.29 5.83
1 2.27 4.35 5.27
1 4.36 6.26 7.89
3 9.55 9.69 4.12 3.67 5.36 5.09 1.34 3.93 3.27
1 7.15 8.51 7.23
3 7.72 8.40 2.86 9.76 11.36 4.79 5.00 7.80 7.24
2 1.25 5.71 5.23 4.55 6.25 4.75
2 0.79 3.31 8.38 8.63 10.90 4.16
2 0.24 3.86 1.05 4.05 7.87 5.01
2 2.41 4.80 2.27 4.41 6.21 5.38
1 9.69 12.30 1.79
2 9.33 12.48 2.10 0.59 2.30 4.01
3 8.05 11.88 5.22 6.78 8.84 2.73 3.91 7.84 8.22
2 8.88 12.29 5.69 7.24 8.16 9.31
3 8.51 11.43 2.37 6.89 9.35 6.97 9.29 12.28 1.30
2 9.05 13.27 6.00 1.96 2.18 2.21
2 8.74 11.79 6.95 1.09 5.92 5.64
3 6.34 7.77 7.15 4.93 5.90 4.35 8.43 12.02 4.52
3 3.18 8.13 1.61 2.20 2.76 2.40 9.36 14.12 7.37
1 8.92 12.44 1.19
1 0.88 2.94 6.13
2 5.41 8.58 9.65 6.69 9.34 9.16
2 3.80 6.49 1.70 9.70 12.16 1.14
2 6.17 7.97 1.75 0.23 4.27 6.52
1 5.88 10.67 4.31
2 2.58 5.55 3.51 8.38 9.48 4.46
3 9.27 10.42 6.31 9.36 9.58 3.84 5.84 8.27 9.60
3 4.79 6.39 7.56 0.24 2.41 6.98 9.62 13.43 8.97
1 7.05 8.29 6.66
2 3.84 5.56 4.36 8.03 8.98 8.42
3 0.60 3.17 1.13 9.08 13.18 6.73 2.58 3.27 9.96
1 6.13 8.22 8.12
1 2.52 3.04 1.68
3 1.28 1.46 8.57 2.23 2.72 8.53 2.17 3.06 6.97
1 3.99 7.03 7.71
2 1.58 5.18 9.05 9.74 12.47 2.81
2 8.28 13.20 3.49 6.64 10.49 1.75
2 3.08 6.62 9.55 0.35 3.41 3.63
1 1.66 3.04 2.81
3 3.46 5.71 4.73 5.32 7.36 1.72 9.79 14.78 2.57
1 6.28 9.65 3.58
1 8.35 11.55 3.42
3 3.13 8.02 6.88 2.87 6.79 3.24 8.88 11.47 3.34
1 7.86 9.95 6.24
2 9.51 14.06 5.80 8.89 12.44 4.81
1 5.21 6.36 8.27
2 0.31 0.71 1.89 1.83 2.47 4.94
2 5.25 5.88 3.22 2.82 4.84 5.24
1 1.64 6.08 5.72
1 2.49 3.60 3.71
1 1.86 6.78 5.90
1 2.33 6.43 4.76
3 9.44 10.65 6.04 8.81 11.72 2.51 2.48 7.42 3.69
3 7.42 11.03 8.11 8.47 8.79 2.51 5.06 6.12 5.80
2 1.03 5.31 5.80 8.99 13.46 1.09
3 9.84 13.57 5.05 2.76 4.82 4.11 3.96 7.59 9.03
1 5.86 9.56 8.72
1 5.11 5.45 5.02
2 7.10 11.59 6.95 6.00 10.40 3.60
3 1.48 4.81 4.30 9.64 12.15 7.19 1.34 3.73 7.61
1 6.09 7.78 3.71
1 4.75 6.23 8.28
2 3.49 6.68 4.43 5.79 9.27 5.51
3 8.15 10.97 3.12 4.31 5.64 2.51 1.40 4.39 7.81
3 7.78 10.31 8.44 4.76 6.47 4.90 4.56 7.81 1.47
3 1.94 3.23 7.22 0.18 3.17 8.70 9.79 14.51 8.43
1 9.40 10.32 2.56
3 3.60 7.66 1.08 9.91 9.99 6.47 9.28 13.44 3.79
1 3.93 6.43 4.27
2 0.88 3.04 2.04 0.25 1.11 6.34
`

type line struct{ k, m float64 }

type cht struct {
	lines []line
	ptr   int
}

func isBad(a, b, c line) bool {
	return (b.m-a.m)*(b.k-c.k) >= (c.m-b.m)*(a.k-b.k)
}

func (h *cht) add(k, m float64) {
	ln := line{k: k, m: m}
	l := h.lines
	for len(l) >= 2 && isBad(l[len(l)-2], l[len(l)-1], ln) {
		l = l[:len(l)-1]
	}
	h.lines = append(l, ln)
}

func (h *cht) query(x float64) float64 {
	for h.ptr+1 < len(h.lines) && h.lines[h.ptr].k*x+h.lines[h.ptr].m <= h.lines[h.ptr+1].k*x+h.lines[h.ptr+1].m {
		h.ptr++
	}
	return h.lines[h.ptr].k*x + h.lines[h.ptr].m
}

func solve(n int, xl, xr, y []float64) float64 {
	segs := make([][2]float64, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if y[i] == y[j] {
				continue
			}
			a := (xr[j] - xl[i]) / (y[i] - y[j])
			b := (xl[j] - xr[i]) / (y[i] - y[j])
			if a > b {
				a, b = b, a
			}
			segs = append(segs, [2]float64{a, b})
		}
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i][0] < segs[j][0] })
	mr := make([][2]float64, 0)
	for _, s := range segs {
		if len(mr) == 0 || mr[len(mr)-1][1] <= s[0] {
			mr = append(mr, s)
		} else if mr[len(mr)-1][1] < s[1] {
			mr[len(mr)-1][1] = s[1]
		}
	}
	if len(mr) == 0 {
		mr = append(mr, [2]float64{0, 0})
	}
	minMap := make(map[float64]float64)
	maxMap := make(map[float64]float64)
	for i := 0; i < n; i++ {
		k1 := -y[i]
		m1 := -xl[i]
		if v, ok := minMap[k1]; !ok || v < m1 {
			minMap[k1] = m1
		}
		m2 := -xr[i]
		if v, ok := minMap[k1]; !ok || v < m2 {
			minMap[k1] = m2
		}
		k2 := y[i]
		m3 := xl[i]
		if v, ok := maxMap[k2]; !ok || v < m3 {
			maxMap[k2] = m3
		}
		m4 := xr[i]
		if v, ok := maxMap[k2]; !ok || v < m4 {
			maxMap[k2] = m4
		}
	}
	minLines := make([]line, 0, len(minMap))
	for k, m := range minMap {
		minLines = append(minLines, line{k: k, m: m})
	}
	maxLines := make([]line, 0, len(maxMap))
	for k, m := range maxMap {
		maxLines = append(maxLines, line{k: k, m: m})
	}
	sort.Slice(minLines, func(i, j int) bool { return minLines[i].k < minLines[j].k })
	sort.Slice(maxLines, func(i, j int) bool { return maxLines[i].k < maxLines[j].k })
	hmin := &cht{lines: make([]line, 0, len(minLines))}
	for _, ln := range minLines {
		hmin.add(ln.k, ln.m)
	}
	hmax := &cht{lines: make([]line, 0, len(maxLines))}
	for _, ln := range maxLines {
		hmax.add(ln.k, ln.m)
	}
	ans := math.Inf(1)
	for _, seg := range mr {
		for _, x := range []float64{seg[0], seg[1]} {
			v := -hmin.query(x) + hmax.query(x)
			if v < ans {
				ans = v
			}
		}
	}
	return ans
}

type testCase struct {
	n  int
	xl []float64
	xr []float64
	y  []float64
}

func parseTestcases(raw string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	tests := make([]testCase, 0)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		fields := strings.Fields(sc.Text())
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n on line %d", lineNo)
		}
		expected := 1 + 3*n
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d values got %d", lineNo, expected, len(fields))
		}
		xl := make([]float64, n)
		xr := make([]float64, n)
		y := make([]float64, n)
		idx := 1
		for i := 0; i < n; i++ {
			var err1, err2, err3 error
			xl[i], err1 = strconv.ParseFloat(fields[idx], 64)
			xr[i], err2 = strconv.ParseFloat(fields[idx+1], 64)
			y[i], err3 = strconv.ParseFloat(fields[idx+2], 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("invalid value on line %d", lineNo)
			}
			idx += 3
		}
		tests = append(tests, testCase{n: n, xl: xl, xr: xr, y: y})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests parsed")
	}
	return tests, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
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

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			sb.WriteString(fmt.Sprintf("%.15g %.15g %.15g\n", tc.xl[i], tc.xr[i], tc.y[i]))
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		res, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx+1)
			os.Exit(1)
		}
		expect := solve(tc.n, tc.xl, tc.xr, tc.y)
		diff := math.Abs(res-expect) / math.Max(1, math.Abs(expect))
		if diff > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %.9f got %.9f\n", idx+1, expect, res)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
