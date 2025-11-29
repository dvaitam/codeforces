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

const testcasesRaw = `100
5 3 80730 2 1 3 7 2 6 2 5 30 1 6 2 2 2 2 1 2 2 2 3 1 033 1 2 5 3 38323 4 4 1 5 4 2 3 3 154 3 9 1 9 3 5 2 3 54 2 5 1 7 2 2 1 3 0 1 6 1 8 1 5 7 5 0702931 7 3 7 7 3 8 3 8 3 7 2 5 54 1 6 1 3 2 8 2 2 2 4 2 3 54 1 1 1 4 2 2 1 1 9 1 0 4 5 5476 2 0 1 7 3 3 2 9 2 6 2 2 65 1 0 2 4 3 4 928 2 7 3 5 2 4 2 7 7 2 1682752 1 7 3 8 6 1 590458 6 4 8 3 45290784 6 4 8 4 6 5 5 3 65275 3 8 2 8 2 3 6 4 416299 5 6 3 9 5 4 6 0 4 2 9792 2 2 1 7 4 2 0215 2 7 2 8 1 4 7 1 6 1 9 1 3 1 0 6 4 461850 5 9 3 1 3 8 5 5 5 3 26698 3 7 2 2 5 6 8 2 29150615 3 5 7 6 7 2 7476629 5 4 6 4 8 3 60547427 1 1 8 3 5 0 3 4 078 3 4 1 7 1 3 2 4 3 3 479 2 8 3 9 3 1 1 2 4 1 8 1 9 5 5 07559 2 0 1 4 5 7 1 8 2 0 7 4 9976763 3 7 7 1 3 0 6 6 5 4 42278 4 5 5 2 4 9 5 0 2 2 41 1 0 2 6 2 4 70 1 1 1 9 1 2 2 7 3 2 926 1 1 3 3 1 5 1 1 1 0 1 6 1 0 1 1 6 5 775502 6 4 2 9 6 8 3 8 5 9 4 3 1834 3 8 2 3 3 7 7 2 2051901 7 1 5 9 8 2 67751808 7 0 3 6 4 1 1732 4 5 4 3 5956 4 2 3 4 4 5 1 5 9 1 2 1 1 1 0 1 2 1 3 1 4 7 1 0 1 7 1 9 1 5 8 4 13221567 3 8 5 1 5 2 5 3 1 4 0 1 5 1 0 1 7 1 2 2 3 47 1 2 1 3 1 9 4 1 9654 2 7 6 3 172323 1 9 4 4 1 7 1 4 6 1 0 1 8 1 1 1 3 8 1 70307600 5 6 8 3 00328638 4 1 6 1 2 8 3 5 139 1 7 3 5 2 6 3 6 3 2 1 3 7 1 8 1 2 1 7 2 3 65 1 2 1 7 1 1 7 5 3407303 7 4 3 1 1 6 7 5 4 8 1 3 2 1 7 1 0 1 6 5 1 68251 3 2 3 3 754 1 1 3 8 1 0 2 1 25 2 4 1 2 8 1 2 1 3 6 3 029921 3 8 5 2 6 2 8 4 61557379 7 2 8 0 3 8 8 7 1 1 2 1 0 4 1 0937 3 6 8 2 34348853 6 5 4 0 8 4 52237601 3 6 2 4 4 5 8 7 2 2 60 1 8 1 6 1 3 0 1 8 1 8 1 8 7 4 9808265 5 0 6 8 2 4 3 1 1 1 5 1 4 1 3 3 1 0 1 4 1 9 4 4 8113 3 8 4 6 1 6 4 2 4 5 5651 3 3 1 3 3 4 3 9 1 2 2 4 23 1 9 2 6 1 6 2 8 4 3 6419 3 2 4 0 2 1 5 5 89926 2 5 5 1 3 1 1 7 2 8 2 3 37 2 3 2 1 1 3 3 1 152 2 9 2 5 75 1 8 1 5 1 1 1 3 2 5 3 4 136 3 7 2 0 1 9 2 2 1 4 6 1 1 1 4 1 1 1 5 4 3 6247 1 9 1 4 4 5 2 5 60 1 6 1 4 1 2 2 1 2 6 3 5 631 2 5 1 1 1 2 2 0 1 4 8 4 81033707 2 4 8 8 1 0 3 6 7 5 9903066 7 6 4 0 7 2 3 2 3 3 7 1 5366342 4 5 5 1 12711 4 0 7 2 3075972 4 3 1 5 1 1 2 1 1 5 2 55324 5 8 1 4 2 1 65 2 5 8 5 97392408 1 6 3 3 8 0 1 0 4 6 4 5 9340 2 0 4 5 1 8 3 7 4 1 2 3 98 2 9 2 6 2 2 2 4 03 2 1 2 8 2 6 1 7 6 1 709939 4 0 1 1 1 1 2`

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
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("empty test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count")
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid n for case %d", i+1)
		}
		xl := make([]float64, n)
		xr := make([]float64, n)
		y := make([]float64, n)
		for j := 0; j < n; j++ {
			for k := 0; k < 3; k++ {
				if !sc.Scan() {
					return nil, fmt.Errorf("missing value for case %d line %d", i+1, j+1)
				}
				val, err := strconv.ParseFloat(sc.Text(), 64)
				if err != nil {
					return nil, fmt.Errorf("invalid value for case %d line %d", i+1, j+1)
				}
				switch k {
				case 0:
					xl[j] = val
				case 1:
					xr[j] = val
				case 2:
					y[j] = val
				}
			}
		}
		tests = append(tests, testCase{n: n, xl: xl, xr: xr, y: y})
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
			sb.WriteString(fmt.Sprintf("%.0f %.0f %.0f\n", tc.xl[i], tc.xr[i], tc.y[i]))
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
