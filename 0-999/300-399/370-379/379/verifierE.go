package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Line struct {
	a, b   int
	x0, x1 float64
}

func eval(a, b int, x float64) float64 { return float64(a)*x + float64(b) }
func P(a, b int, x0, x1 float64) float64 {
	return (eval(a, b, x0) + eval(a, b, x1)) * (x1 - x0) / 2.0
}
func lowerBound(lines []Line, a int) int {
	return sort.Search(len(lines), func(i int) bool { return lines[i].a >= a })
}
func removeAt(lines []Line, idx int) []Line { return append(lines[:idx], lines[idx+1:]...) }
func insertAt(lines []Line, idx int, ln Line) []Line {
	lines = append(lines, Line{})
	copy(lines[idx+1:], lines[idx:])
	lines[idx] = ln
	return lines
}
func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func solveCase(n, k int, vals [][]int) []float64 {
	S := make([][]Line, k)
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		y0 := vals[i][0]
		area := 0.0
		for j := 0; j < k; j++ {
			y1 := vals[i][j+1]
			a := y1 - y0
			b := y0
			segs := S[j]
			if len(segs) == 0 {
				segs = []Line{{a: a, b: b, x0: 0, x1: 1}}
				area += P(a, b, 0, 1)
				S[j] = segs
				y0 = y1
				continue
			}
			left, right := 1e18, -1e18
			idx := lowerBound(segs, a)
			for idx < len(segs) {
				l := segs[idx]
				if eval(a, b, l.x0) >= eval(l.a, l.b, l.x0) && eval(a, b, l.x1) >= eval(l.a, l.b, l.x1) {
					area += P(a, b, l.x0, l.x1) - P(l.a, l.b, l.x0, l.x1)
					left = minF(left, l.x0)
					right = maxF(right, l.x1)
					segs = removeAt(segs, idx)
					continue
				}
				if eval(a, b, l.x0) > eval(l.a, l.b, l.x0) {
					xsr := float64(b-l.b) / float64(l.a-a)
					area += P(a, b, l.x0, xsr) - P(l.a, l.b, l.x0, xsr)
					left = minF(left, l.x0)
					right = maxF(right, xsr)
					segs[idx] = Line{a: l.a, b: l.b, x0: xsr, x1: l.x1}
				}
				break
			}
			idx = lowerBound(segs, a)
			if idx > 0 {
				idx--
				for idx >= 0 {
					l := segs[idx]
					if eval(a, b, l.x0) >= eval(l.a, l.b, l.x0) && eval(a, b, l.x1) >= eval(l.a, l.b, l.x1) {
						area += P(a, b, l.x0, l.x1) - P(l.a, l.b, l.x0, l.x1)
						left = minF(left, l.x0)
						right = maxF(right, l.x1)
						segs = removeAt(segs, idx)
						idx--
						continue
					}
					if eval(a, b, l.x1) > eval(l.a, l.b, l.x1) {
						xsr := float64(b-l.b) / float64(l.a-a)
						area += P(a, b, xsr, l.x1) - P(l.a, l.b, xsr, l.x1)
						left = minF(left, xsr)
						right = maxF(right, l.x1)
						segs[idx] = Line{a: l.a, b: l.b, x0: l.x0, x1: xsr}
					}
					break
				}
			}
			if left <= right {
				insertPos := lowerBound(segs, a)
				segs = insertAt(segs, insertPos, Line{a: a, b: b, x0: left, x1: right})
			}
			S[j] = segs
			y0 = y1
		}
		res[i] = area
	}
	return res
}

func runCase(exe string, n, k int, vals [][]int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		for j := 0; j <= k; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", vals[i][j]))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	expected := solveCase(n, k, vals)
	for i := 0; i < n; i++ {
		var v float64
		if _, err := fmt.Sscan(fields[i], &v); err != nil {
			return fmt.Errorf("parse output: %v", err)
		}
		diff := v - expected[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > 1e-4 {
			return fmt.Errorf("value %d mismatch: expected %.5f got %.5f", i+1, expected[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		n, k int
		vals [][]int
	}{}
	cases = append(cases, struct {
		n, k int
		vals [][]int
	}{2, 2, [][]int{{1, 1, 1}, {2, 2, 2}}})
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		k := rng.Intn(3) + 1
		vals := make([][]int, n)
		for j := 0; j < n; j++ {
			vals[j] = make([]int, k+1)
			for t := 0; t <= k; t++ {
				vals[j][t] = rng.Intn(5) + 1
			}
		}
		cases = append(cases, struct {
			n, k int
			vals [][]int
		}{n, k, vals})
	}
	for idx, c := range cases {
		if err := runCase(exe, c.n, c.k, c.vals); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
