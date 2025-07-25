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

const eps = 1e-9

// point structure
type Pt struct{ x, y float64 }

func (a Pt) Add(b Pt) Pt        { return Pt{a.x + b.x, a.y + b.y} }
func (a Pt) Sub(b Pt) Pt        { return Pt{a.x - b.x, a.y - b.y} }
func (a Pt) Mul(f float64) Pt   { return Pt{a.x * f, a.y * f} }
func (a Pt) Dot(b Pt) float64   { return a.x*b.x + a.y*b.y }
func (a Pt) Cross(b Pt) float64 { return a.x*b.y - a.y*b.x }
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
	d := math.Hypot(n.x, n.y)
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

func solveOne(n int, pivot Pt, pts []Pt) float64 {
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

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	var inputs strings.Builder
	inputs.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]float64, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		px, _ := strconv.ParseFloat(scan.Text(), 64)
		scan.Scan()
		py, _ := strconv.ParseFloat(scan.Text(), 64)
		pts := make([]Pt, n-1)
		for i := 1; i < n; i++ {
			scan.Scan()
			x, _ := strconv.ParseFloat(scan.Text(), 64)
			scan.Scan()
			y, _ := strconv.ParseFloat(scan.Text(), 64)
			pts[i-1] = Pt{x, y}
		}
		expected[caseIdx] = solveOne(n, Pt{px, py}, pts)
		inputs.WriteString(fmt.Sprintf("%d %.0f %.0f\n", n, px, py))
		for _, p := range pts {
			inputs.WriteString(fmt.Sprintf("%.0f %.0f\n", p.x, p.y))
		}
	}
	out, err := runCandidate(os.Args[1], []byte(inputs.String()))
	if err != nil {
		fmt.Println("candidate failed:", err)
		os.Exit(1)
	}
	outLines := strings.Fields(out)
	if len(outLines) != t {
		fmt.Printf("expected %d outputs got %d\n", t, len(outLines))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		val, err := strconv.ParseFloat(outLines[i], 64)
		if err != nil || math.Abs(val-expected[i]) > 1e-6 {
			fmt.Printf("case %d failed: expected %.6f got %s\n", i+1, expected[i], outLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
