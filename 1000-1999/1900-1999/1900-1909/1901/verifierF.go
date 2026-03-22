package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ── embedded reference solver ──────────────────────────────────────────────

type Point struct {
	x int64
	y int64
}

func cross(a, b, c Point) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func buildHull(l, r int, h []int64) []Point {
	hull := make([]Point, 0, r-l+1)
	for i := l; i <= r; i++ {
		p := Point{int64(i), h[i]}
		for len(hull) >= 2 && cross(hull[len(hull)-2], hull[len(hull)-1], p) >= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	return hull
}

func pairValue(lx, ly, rx, ry, n int64) float64 {
	u := n - 1 - 2*lx
	v := 2*rx - (n - 1)
	num := 2.0 * (float64(v)*float64(ly) + float64(u)*float64(ry))
	den := float64(u + v)
	return num / den
}

func queryToRight(hull []Point, x, y, n int64) float64 {
	if len(hull) == 1 {
		return pairValue(x, y, hull[0].x, hull[0].y, n)
	}
	q := Point{x, y}
	lo, hi := 0, len(hull)-2
	ans := -1
	for lo <= hi {
		mid := (lo + hi) >> 1
		if cross(q, hull[mid], hull[mid+1]) >= 0 {
			ans = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	idx := 0
	if ans != -1 {
		idx = ans + 1
	}
	p := hull[idx]
	return pairValue(x, y, p.x, p.y, n)
}

func queryToLeft(hull []Point, x, y, n int64) float64 {
	if len(hull) == 1 {
		return pairValue(hull[0].x, hull[0].y, x, y, n)
	}
	q := Point{x, y}
	lo, hi := 0, len(hull)-2
	idx := len(hull) - 1
	for lo <= hi {
		mid := (lo + hi) >> 1
		if cross(q, hull[mid], hull[mid+1]) >= 0 {
			idx = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	p := hull[idx]
	return pairValue(p.x, p.y, x, y, n)
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func referenceSolve(n int, a, b []int64) []float64 {
	N := int64(n)
	negInf := math.Inf(-1)
	ans := make([]float64, n)

	if n%2 == 0 {
		leftCnt := n / 2
		leftEnd := leftCnt - 1
		rightStart := leftCnt

		rightHullA := buildHull(rightStart, n-1, a)
		leftValA := make([]float64, leftCnt)
		leftValB := make([]float64, leftCnt)
		for i := 0; i < leftCnt; i++ {
			x := int64(i)
			leftValA[i] = queryToRight(rightHullA, x, a[i], N)
			leftValB[i] = queryToRight(rightHullA, x, b[i], N)
		}
		prefB := make([]float64, leftCnt)
		prefB[0] = leftValB[0]
		for i := 1; i < leftCnt; i++ {
			prefB[i] = maxf(prefB[i-1], leftValB[i])
		}
		sufA := make([]float64, leftCnt+1)
		sufA[leftCnt] = negInf
		for i := leftCnt - 1; i >= 0; i-- {
			sufA[i] = maxf(leftValA[i], sufA[i+1])
		}
		for i := 0; i <= leftEnd; i++ {
			ans[i] = maxf(prefB[i], sufA[i+1])
		}

		leftHullB := buildHull(0, leftEnd, b)
		rightCnt := n - rightStart
		rightValA := make([]float64, rightCnt)
		rightValB := make([]float64, rightCnt)
		for idx, k := 0, rightStart; k < n; k, idx = k+1, idx+1 {
			x := int64(k)
			rightValA[idx] = queryToLeft(leftHullB, x, a[k], N)
			rightValB[idx] = queryToLeft(leftHullB, x, b[k], N)
		}
		prefB2 := make([]float64, rightCnt)
		prefB2[0] = rightValB[0]
		for i := 1; i < rightCnt; i++ {
			prefB2[i] = maxf(prefB2[i-1], rightValB[i])
		}
		sufA2 := make([]float64, rightCnt+1)
		sufA2[rightCnt] = negInf
		for i := rightCnt - 1; i >= 0; i-- {
			sufA2[i] = maxf(rightValA[i], sufA2[i+1])
		}
		for i := rightStart; i < n; i++ {
			idx := i - rightStart
			ans[i] = maxf(prefB2[idx], sufA2[idx+1])
		}
	} else {
		c := n / 2
		rightStart := c + 1

		rightHullA := buildHull(rightStart, n-1, a)
		leftValA := make([]float64, c)
		leftValB := make([]float64, c)
		for i := 0; i < c; i++ {
			x := int64(i)
			leftValA[i] = queryToRight(rightHullA, x, a[i], N)
			leftValB[i] = queryToRight(rightHullA, x, b[i], N)
		}
		prefB := make([]float64, c)
		prefB[0] = leftValB[0]
		for i := 1; i < c; i++ {
			prefB[i] = maxf(prefB[i-1], leftValB[i])
		}
		sufA := make([]float64, c+1)
		sufA[c] = negInf
		for i := c - 1; i >= 0; i-- {
			sufA[i] = maxf(leftValA[i], sufA[i+1])
		}
		centerA := 2.0 * float64(a[c])
		for i := 0; i < c; i++ {
			v := centerA
			v = maxf(v, prefB[i])
			v = maxf(v, sufA[i+1])
			ans[i] = v
		}

		centerB := 2.0 * float64(b[c])
		ans[c] = maxf(centerB, prefB[c-1])

		leftHullB := buildHull(0, c-1, b)
		rightCnt := n - rightStart
		rightValA := make([]float64, rightCnt)
		rightValB := make([]float64, rightCnt)
		for idx, k := 0, rightStart; k < n; k, idx = k+1, idx+1 {
			x := int64(k)
			rightValA[idx] = queryToLeft(leftHullB, x, a[k], N)
			rightValB[idx] = queryToLeft(leftHullB, x, b[k], N)
		}
		prefB2 := make([]float64, rightCnt)
		prefB2[0] = rightValB[0]
		for i := 1; i < rightCnt; i++ {
			prefB2[i] = maxf(prefB2[i-1], rightValB[i])
		}
		sufA2 := make([]float64, rightCnt+1)
		sufA2[rightCnt] = negInf
		for i := rightCnt - 1; i >= 0; i-- {
			sufA2[i] = maxf(rightValA[i], sufA2[i+1])
		}
		for i := rightStart; i < n; i++ {
			idx := i - rightStart
			v := centerB
			v = maxf(v, prefB2[idx])
			v = maxf(v, sufA2[idx+1])
			ans[i] = v
		}
	}
	return ans
}

// ── verifier harness ───────────────────────────────────────────────────────

type runResult struct {
	out string
	err error
}

func runBinary(path, input string) runResult {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, stderr.String())
	}
	return runResult{strings.TrimSpace(out.String()), err}
}

func genTest() (string, int, []int64, []int64) {
	n := rand.Intn(8) + 3
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 1; i < n-1; i++ {
		a[i] = rand.Int63n(1000)
		b[i] = rand.Int63n(1000)
	}
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), n, a, b
}

func formatRef(vals []float64) string {
	var out bytes.Buffer
	for i, v := range vals {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.FormatFloat(v, 'f', 15, 64))
	}
	return out.String()
}

func parseFloats(s string) ([]float64, error) {
	fields := strings.Fields(s)
	out := make([]float64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func main() {
	_ = io.Discard // keep io imported for potential use
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(6)

	for i := 0; i < 100; i++ {
		tc, n, a, b := genTest()
		expVals := referenceSolve(n, a, b)
		expStr := formatRef(expVals)

		got := runBinary(binary, tc)
		if got.err != nil {
			fmt.Fprintf(os.Stderr, "binary failed on test %d: %v\n", i+1, got.err)
			os.Exit(1)
		}

		// Compare with tolerance
		gotVals, err := parseFloats(got.out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error on test %d: %v\noutput: %s\n", i+1, err, got.out)
			os.Exit(1)
		}
		if len(gotVals) != len(expVals) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d values, got %d\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n",
				i+1, len(expVals), len(gotVals), tc, expStr, got.out)
			os.Exit(1)
		}
		for j := range expVals {
			if math.Abs(expVals[j]-gotVals[j]) > 1e-6 {
				fmt.Fprintf(os.Stderr, "mismatch on test %d position %d: expected %.15f got %.15f\ninput:\n%s\n",
					i+1, j, expVals[j], gotVals[j], tc)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
