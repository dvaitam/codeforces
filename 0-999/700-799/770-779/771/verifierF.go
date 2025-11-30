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

const testcases = `
100
3 2 4
3 -4
4 3
3 -1 -1
3 1
1 4
5 1 -5
-5 -4
1 0
4 0
-4 5
3 -5 3
-3 -5
1 -3
3 4 -4
3 0
-3 5
4 1 -3
3 0
5 5
-1 2
4 -4 -1
0 5
3 2
-5 3
3 -3 5
-1 -2
1 -4
3 -5 -5
0 -1
2 -4
4 4 3
4 -4
3 1
-1 5
4 5 -1
4 -1
0 4
2 1
3 0 -5
4 -4
3 1
3 -3 5
-5 4
3 0
4 4 0
-2 -1
5 -5
-1 2
4 -2 3
0 5
-1 -5
1 3
4 -1 5
-2 1
4 3
2 5
4 1 3
1 4
1 -1
5 -5
5 -4 -2
-4 3
-5 -5
-1 5
-3 -3
5 4 -3
5 -3
1 -1
-1 -2
3 -4
5 -1 -1
-3 3
-1 4
4 -5
-5 5
3 -3 2
-4 -3
-5 -2
4 3 -3
5 3
5 4
1 1
5 -3 1
-3 4
-2 2
-5 5
4 -3
4 2 -2
-1 4
-4 -2
-4 -4
4 1 2
-3 2
-4 1
2 3
4 -1 3
-4 -3
4 -4
-2 5
3 -4 -2
3 1
-1 -2
3 5 2
-2 0
-1 -4
5 0 -5
-1 3
-2 -4
-1 5
-2 5
4 1 2
-3 5
-3 5
-2 5
3 5 5
5 1
-2 0
3 5 5
2 3
2 3
5 3 2
1 -3
2 -1
-4 2
2 -3
3 -3 -1
-2 3
5 0
3 3 -5
0 3
-1 0
5 4 2
5 -2
0 -5
-2 4
0 0
4 3 3
4 -3
-3 -3
-4 4
3 -4 2
5 2
3 4
4 2 -4
2 -1
-5 0
-1 4
5 -3 -3
0 -1
-3 -1
-1 -3
-2 -1
5 -5 3
-3 2
4 -3
-2 -2
0 -3
4 0 1
-1 5
-3 -2
3 4
4 -3 -5
-1 1
1 0
5 4
5 2 4
2 -3
1 3
2 -2
3 0
3 -1 2
3 3
-5 -1
4 4 5
1 1
4 -1
-4 3
5 2 -5
-5 2
2 1
0 -1
2 4
5 4 5
-1 -4
-1 -5
-3 -3
-2 4
4 -2 0
-5 4
4 -5
3 5
4 0 4
2 -3
-4 5
3 4
3 -2 2
5 1
0 -2
5 -1 -5
-4 3
-5 1
0 5
-2 -3
5 -4 -3
-1 1
3 -5
1 4
-3 5
4 0 -1
-3 1
3 1
-3 4
3 3 -2
1 1
4 -5
3 -5 -1
5 -1
2 -3
4 4 4
1 -3
-3 1
0 4
4 3 3
4 5
0 -3
-1 5
5 2 4
-5 -3
0 -5
4 -2
-2 1
3 3 -3
4 5
-4 1
4 -2 -1
3 3
-5 5
-2 -1
3 4 1
5 -3
3 2
5 -1 4
-2 0
-1 -2
4 -1
-1 2
3 5 2
0 -3
4 1
4 -4 3
3 -1
-3 4
-2 1
4 1 2
-2 -4
-4 3
-2 0
3 -1 5
3 -1
3 -1
5 -3 2
-2 0
1 2
3 -2
-4 -1
5 -2 -2
-1 -2
3 -2
0 -2
5 3
4 -4 4
-5 2
2 4
-5 -3
5 -5 2
0 2
-4 3
2 -3
2 -5
4 -4 -1
1 0
-3 -1
4 5
4 -4 -3
-4 3
-3 5
-3 4
3 -2 2
-5 5
-1 -3
5 5 3
3 -4
-2 4
-1 -1
0 0
5 -5 -3
2 -2
-1 -5
-5 -2
-4 5
3 0 -1
-1 2
5 -4
3 -3 -4
5 1
-5 3
3 5 3
4 2
-5 -3
4 -5 1
-2 1
-5 2
4 -2
3 5 4
-5 -3
-5 -3
3 -2 -4
1 -1
3 -4
4 -4 -2
-3 -1
2 1
3 4
4 1 1
3 3
-3 -3
-4 3
5 2 4
-2 1
-2 3
-1 5
3 -1
4 5 -2
2 0
4 -4
3 -3
4 -3 -5
-3 4
4 3
-3 5
3 -5 0
3 3
-4 2
5 4 2
-5 -5
-3 -4
2 5
5 -5
3 -3 -3
-4 5
2 5
5 4 -1
-4 0
1 -3
4 0
-2 5
3 1 -1
-4 -4
-4 0
3 4 0
4 3
-3 2
3 0 -2
-5 -4
-1 1
5 1 1
0 2
-3 -1
5 -3
-4 -3
4 -2 -5
-5 -4
2 5
1 5
3 5 -5
-2 1
5 -4
4 3 -4
1 3
-3 3
-2 0
3 -1 2
0 -1
-3 -1
5 -3 0
0 1
-5 -3
-3 -5
-1 -5

`
const eps = 1e-9

type Pt struct{ x, y float64 }

func (a Pt) Add(b Pt) Pt        { return Pt{a.x + b.x, a.y + b.y} }
func (a Pt) Sub(b Pt) Pt        { return Pt{a.x - b.x, a.y - b.y} }
func (a Pt) Mul(f float64) Pt   { return Pt{a.x * f, a.y * f} }
func (a Pt) Dot(b Pt) float64   { return a.x*b.x + a.y*b.y }
func (a Pt) Cross(b Pt) float64 { return a.x*b.y - a.y*b.x }
func (a Pt) Abs() float64       { return math.Hypot(a.x, a.y) }

func (a Pt) Up() bool {
	if math.Abs(a.y) < eps {
		return a.x > 0
	}
	return a.y > 0
}

type line struct {
	v Pt
	c float64
}

func newLinePoints(p1, p2 Pt) line {
	dir := p2.Sub(p1)
	n := Pt{-dir.y, dir.x}
	d := n.Abs()
	if d > 0 {
		n = n.Mul(1 / d)
	}
	return line{v: n, c: n.Dot(p1)}
}

func (l line) signedDist(p Pt) float64 { return l.v.Dot(p) - l.c }

func cmpAngle(a, b Pt) bool {
	au, bu := a.Up(), b.Up()
	if au != bu {
		return au
	}
	return a.Cross(b) > eps
}

func eqLine(a, b line) bool {
	if a.v.Up() != b.v.Up() {
		return false
	}
	return math.Abs(a.v.Cross(b.v)) < eps
}

func cmpLine(a, b line) bool {
	au, bu := a.v.Up(), b.v.Up()
	if au != bu {
		return au
	}
	cr := a.v.Cross(b.v)
	if math.Abs(cr) > eps {
		return cr > 0
	}
	return a.c > b.c
}

func det3x3(a, b, c line) float64 {
	return a.c*(b.v.Cross(c.v)) + b.c*(c.v.Cross(a.v)) + c.c*(a.v.Cross(b.v))
}

func intersect(l1, l2 line) Pt {
	d := l1.v.x*l2.v.y - l1.v.y*l2.v.x
	if math.Abs(d) < eps {
		return Pt{1e18, 1e18}
	}
	dx := l1.c*l2.v.y - l1.v.y*l2.c
	dy := l1.v.x*l2.c - l1.c*l2.v.x
	return Pt{dx / d, dy / d}
}

func halfplanesIntersection(ls []line) []Pt {
	sort.Slice(ls, func(i, j int) bool { return cmpLine(ls[i], ls[j]) })
	uniq := make([]line, 0, len(ls))
	for i, l := range ls {
		if i == 0 || !eqLine(ls[i-1], l) {
			uniq = append(uniq, l)
		}
	}
	ls = uniq
	n := len(ls)
	st := make([]int, 0, n*2)
	for iter := 0; iter < 2; iter++ {
		for i := 0; i < n; i++ {
			for len(st) > 1 {
				j := st[len(st)-1]
				k := st[len(st)-2]
				if ls[k].v.Cross(ls[i].v) <= eps || det3x3(ls[k], ls[j], ls[i]) <= eps {
					break
				}
				st = st[:len(st)-1]
			}
			st = append(st, i)
		}
	}
	pos := make([]int, n)
	for i := range pos {
		pos[i] = -1
	}
	ok := false
	var seq []int
	for i, id := range st {
		if pos[id] != -1 {
			seq = st[pos[id]:i]
			ok = true
			break
		}
		pos[id] = i
	}
	if !ok {
		return nil
	}
	k := len(seq)
	res := make([]Pt, k)
	M := Pt{0, 0}
	for i := 0; i < k; i++ {
		l1 := ls[seq[i]]
		l2 := ls[seq[(i+1)%k]]
		p := intersect(l1, l2)
		res[i] = p
		M = M.Add(p)
	}
	M = M.Mul(1.0 / float64(k))
	for _, id := range seq {
		if ls[id].signedDist(M) < -eps {
			return nil
		}
	}
	return res
}

func solveArea(n int, pivot Pt, pts []Pt) float64 {
	ev := make([]Pt, 0, (n-1)*2)
	for _, p := range pts {
		v := p.Sub(pivot)
		ev = append(ev, v)
		ev = append(ev, Pt{-v.x, -v.y})
	}
	sort.Slice(ev, func(i, j int) bool { return cmpAngle(ev[i], ev[j]) })
	const B = 1e6
	hp := make([]line, 0, len(ev)+4)
	negPivot := pivot.Mul(-1)
	box := []Pt{negPivot.Add(Pt{B, B}), negPivot.Add(Pt{-B, B}), negPivot.Add(Pt{-B, -B}), negPivot.Add(Pt{B, -B})}
	for i := 0; i < 4; i++ {
		hp = append(hp, newLinePoints(box[i], box[(i+1)%4]))
	}
	m := len(ev)
	for i := 0; i < m; i++ {
		A := ev[i]
		Bv := ev[(i+1)%m]
		if math.Abs(A.Cross(Bv)) < eps {
			return 0
		}
		l := newLinePoints(A, Bv)
		if l.signedDist(Pt{0, 0}) < 0 {
			l = newLinePoints(Bv, A)
		}
		hp = append(hp, l)
	}
	poly := halfplanesIntersection(hp)
	if poly == nil {
		return 0
	}
	area := 0.0
	k := len(poly)
	for i := 0; i < k; i++ {
		area += poly[i].Cross(poly[(i+1)%k])
	}
	return math.Abs(area) / 2.0
}

type testCase struct {
	n   int
	piv Pt
	pts []Pt
}

func parseCases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scan.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("read t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", i+1, err)
		}
		px, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read px: %w", i+1, err)
		}
		py, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read py: %w", i+1, err)
		}
		pts := make([]Pt, n-1)
		for j := 0; j < n-1; j++ {
			x, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: read x%d: %w", i+1, j, err)
			}
			y, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: read y%d: %w", i+1, j, err)
			}
			pts[j] = Pt{float64(x), float64(y)}
		}
		cases = append(cases, testCase{n: n, piv: Pt{float64(px), float64(py)}, pts: pts})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "1\n")
		fmt.Fprintf(&input, "%d %.0f %.0f\n", tc.n, tc.piv.x, tc.piv.y)
		for _, p := range tc.pts {
			fmt.Fprintf(&input, "%.0f %.0f\n", p.x, p.y)
		}

		expected := solveArea(tc.n, tc.piv, tc.pts)
		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Fprintf(os.Stderr, "case %d: empty output\n", idx+1)
			os.Exit(1)
		}
		val, err := strconv.ParseFloat(fields[0], 64)
		if err != nil || math.Abs(val-expected) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %.6f\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
