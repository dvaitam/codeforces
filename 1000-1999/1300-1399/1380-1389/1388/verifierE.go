package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type line struct{ k, m float64 }

type cht struct {
	lines []line
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
	// Binary search for best line (no monotone pointer needed)
	lo, hi := 0, len(h.lines)-1
	for lo < hi {
		mid := (lo + hi) / 2
		if h.lines[mid].k*x+h.lines[mid].m <= h.lines[mid+1].k*x+h.lines[mid+1].m {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return h.lines[lo].k*x + h.lines[lo].m
}

func solve(n int, xl, xr, y []int) float64 {
	segs := make([][2]float64, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if y[i] == y[j] {
				continue
			}
			a := float64(xr[j]-xl[i]) / float64(y[i]-y[j])
			b := float64(xl[j]-xr[i]) / float64(y[i]-y[j])
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
	// Build convex hull trick for min and max
	// lmin: lines (-y[i], -xl[i]) and (-y[i], -xr[i]), query gives max => negate for min
	// lmax: lines (y[i], xl[i]) and (y[i], xr[i]), query gives max
	// width(c) = lmax.query(c) + lmin.query(c)
	//   because lmin.query(c) = max(-y*c - val) = -(min(y*c + val))
	//   and width = max(y*c + xr) - min(y*c + xl) = lmax.query(c) - (-lmin.query(c))
	//   wait: lmin stores lines for the MIN side. Let's be precise:
	//   lmin.query(c) = max over all i of (-y[i]*c - xl[i])
	//   -lmin.query(c) = min over all i of (y[i]*c + xl[i]) = leftmost projection
	//   lmax.query(c) = max over all i of (y[i]*c + xr[i]) = rightmost projection
	//   width = lmax.query(c) - (-lmin.query(c)) = lmax.query(c) + lmin.query(c)

	type slopeLine struct {
		k float64
		m float64
	}
	minMap := make(map[float64]float64)
	maxMap := make(map[float64]float64)
	for i := 0; i < n; i++ {
		k1 := float64(-y[i])
		// For same slope, keep the line with highest intercept (dominates for max query)
		m1 := float64(-xl[i])
		if v, ok := minMap[k1]; !ok || m1 > v {
			minMap[k1] = m1
		}
		m2 := float64(-xr[i])
		if v, ok := minMap[k1]; !ok || m2 > v {
			minMap[k1] = m2
		}

		k2 := float64(y[i])
		m3 := float64(xl[i])
		if v, ok := maxMap[k2]; !ok || m3 > v {
			maxMap[k2] = m3
		}
		m4 := float64(xr[i])
		if v, ok := maxMap[k2]; !ok || m4 > v {
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
			v := hmin.query(x) + hmax.query(x)
			if v < ans {
				ans = v
			}
		}
	}
	return ans
}

type testCase struct {
	n  int
	xl []int
	xr []int
	y  []int
}

func generateTestCase(rng *rand.Rand, maxCoord, maxY, maxN int) testCase {
	n := rng.Intn(maxN) + 1
	xl := make([]int, n)
	xr := make([]int, n)
	y := make([]int, n)
	// Generate non-overlapping segments at distinct heights for simplicity
	for i := 0; i < n; i++ {
		for {
			a := rng.Intn(2*maxCoord+1) - maxCoord
			b := a + rng.Intn(10) + 1
			if b > maxCoord {
				b = maxCoord
			}
			if a >= b {
				continue
			}
			yi := rng.Intn(maxY) + 1
			// Check no overlap with existing segments at same y
			ok := true
			for j := 0; j < i; j++ {
				if y[j] == yi && !(xr[j] <= a || b <= xl[j]) {
					ok = false
					break
				}
			}
			if ok {
				xl[i] = a
				xr[i] = b
				y[i] = yi
				break
			}
		}
	}
	return testCase{n: n, xl: xl, xr: xr, y: y}
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

	rng := rand.New(rand.NewSource(42))

	// Generate test cases with integer coordinates
	tests := make([]testCase, 0, 100)
	// Add some small deterministic cases
	tests = append(tests, testCase{n: 1, xl: []int{1}, xr: []int{5}, y: []int{3}})
	tests = append(tests, testCase{n: 2, xl: []int{1, 3}, xr: []int{2, 6}, y: []int{1, 2}})
	tests = append(tests, testCase{n: 2, xl: []int{0, 0}, xr: []int{2, 2}, y: []int{1, 2}})
	// Random small cases
	for i := 0; i < 97; i++ {
		tests = append(tests, generateTestCase(rng, 20, 10, 3))
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.xl[i], tc.xr[i], tc.y[i]))
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		res, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output %q\n", idx+1, got)
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
