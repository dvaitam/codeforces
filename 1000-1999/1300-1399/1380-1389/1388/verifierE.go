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

func runCase(bin string, n int, xl, xr, y []float64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f\n", xl[i], xr[i], y[i]))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got, err := strconv.ParseFloat(strings.TrimSpace(out.String()), 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := solve(n, xl, xr, y)
	diff := math.Abs(got-exp) / math.Max(1, math.Abs(exp))
	if diff > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("bad test line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+3*n {
			fmt.Printf("bad test line %d field count\n", idx+1)
			os.Exit(1)
		}
		xl := make([]float64, n)
		xr := make([]float64, n)
		y := make([]float64, n)
		for i := 0; i < n; i++ {
			xl[i], _ = strconv.ParseFloat(fields[1+3*i], 64)
			xr[i], _ = strconv.ParseFloat(fields[1+3*i+1], 64)
			y[i], _ = strconv.ParseFloat(fields[1+3*i+2], 64)
		}
		idx++
		if err := runCase(bin, n, xl, xr, y); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
