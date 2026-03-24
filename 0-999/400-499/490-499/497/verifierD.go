package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type IPoint struct{ x, y int64 }

func dist2(a, b IPoint) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return dx*dx + dy*dy
}

func cross3(a, b, c IPoint) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func between(a, b, x int64) bool {
	if a > b {
		a, b = b, a
	}
	return a <= x && x <= b
}

func onSegment(a, b, p IPoint) bool {
	return cross3(a, b, p) == 0 && between(a.x, b.x, p.x) && between(a.y, b.y, p.y)
}

func segmentsIntersect(a, b, c, d IPoint) bool {
	o1 := cross3(a, b, c)
	o2 := cross3(a, b, d)
	o3 := cross3(c, d, a)
	o4 := cross3(c, d, b)

	if o1 == 0 && onSegment(a, b, c) {
		return true
	}
	if o2 == 0 && onSegment(a, b, d) {
		return true
	}
	if o3 == 0 && onSegment(c, d, a) {
		return true
	}
	if o4 == 0 && onSegment(c, d, b) {
		return true
	}

	return (o1 > 0) != (o2 > 0) && (o3 > 0) != (o4 > 0)
}

func pointSegWithin(p, a, b IPoint, r2 int64) bool {
	dx := b.x - a.x
	dy := b.y - a.y
	len2 := dx*dx + dy*dy
	if len2 == 0 {
		return dist2(p, a) <= r2
	}

	apx := p.x - a.x
	apy := p.y - a.y
	dot := apx*dx + apy*dy

	if dot <= 0 {
		return apx*apx+apy*apy <= r2
	}
	if dot >= len2 {
		bpx := p.x - b.x
		bpy := p.y - b.y
		return bpx*bpx+bpy*bpy <= r2
	}

	cr := dx*apy - dy*apx
	if cr < 0 {
		cr = -cr
	}
	return cr*cr <= r2*len2
}

type fastScanner struct {
	data []byte
	idx  int
}

func (fs *fastScanner) nextInt() int64 {
	n := len(fs.data)
	for fs.idx < n && fs.data[fs.idx] <= ' ' {
		fs.idx++
	}
	if fs.idx >= n {
		return 0
	}
	sign := int64(1)
	if fs.data[fs.idx] == '-' {
		sign = -1
		fs.idx++
	}
	var v int64
	for fs.idx < n {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int64(c-'0')
		fs.idx++
	}
	return sign * v
}

func solveD(input string) string {
	data := []byte(input)
	fs := &fastScanner{data: data}

	px := fs.nextInt()
	py := fs.nextInt()

	n := int(fs.nextInt())
	A := make([]IPoint, n)
	for i := 0; i < n; i++ {
		x := fs.nextInt()
		y := fs.nextInt()
		A[i] = IPoint{x - px, y - py}
	}

	qx := fs.nextInt()
	qy := fs.nextInt()

	m := int(fs.nextInt())
	B := make([]IPoint, m)
	for i := 0; i < m; i++ {
		x := fs.nextInt()
		y := fs.nextInt()
		B[i] = IPoint{x - qx, y - qy}
	}

	dx := qx - px
	dy := qy - py
	r2 := dx*dx + dy*dy

	for i := 0; i < n; i++ {
		a0 := A[i]
		a1 := A[(i+1)%n]
		for j := 0; j < m; j++ {
			b0 := B[j]
			b1 := B[(j+1)%m]

			maxd := dist2(a0, b0)
			if d := dist2(a0, b1); d > maxd {
				maxd = d
			}
			if d := dist2(a1, b0); d > maxd {
				maxd = d
			}
			if d := dist2(a1, b1); d > maxd {
				maxd = d
			}
			if maxd < r2 {
				continue
			}

			if segmentsIntersect(a0, a1, b0, b1) ||
				pointSegWithin(a0, b0, b1, r2) ||
				pointSegWithin(a1, b0, b1, r2) ||
				pointSegWithin(b0, a0, a1, r2) ||
				pointSegWithin(b1, a0, a1, r2) {
				return "YES"
			}
		}
	}

	return "NO"
}

func randPoly(rng *rand.Rand, n int) []IPoint {
	pts := make([]IPoint, n)
	for i := 0; i < n; i++ {
		pts[i] = IPoint{int64(rng.Intn(21) - 10), int64(rng.Intn(21) - 10)}
	}
	cx, cy := 0.0, 0.0
	for _, p := range pts {
		cx += float64(p.x)
		cy += float64(p.y)
	}
	cx /= float64(n)
	cy /= float64(n)
	sort.Slice(pts, func(i, j int) bool {
		ai := math.Atan2(float64(pts[i].y)-cy, float64(pts[i].x)-cx)
		aj := math.Atan2(float64(pts[j].y)-cy, float64(pts[j].x)-cx)
		return ai < aj
	})
	return pts
}

func genCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3
	m := rng.Intn(3) + 3
	P := IPoint{int64(rng.Intn(11) - 5), int64(rng.Intn(11) - 5)}
	Q := IPoint{int64(rng.Intn(11) - 5), int64(rng.Intn(11) - 5)}
	if P.x == Q.x && P.y == Q.y {
		Q.x++
	}
	A := randPoly(rng, n)
	B := randPoly(rng, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", P.x, P.y)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", A[i].x, A[i].y)
	}
	fmt.Fprintf(&sb, "%d %d\n", Q.x, Q.y)
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", B[i].x, B[i].y)
	}
	inp := sb.String()
	exp := solveD(inp)
	return inp, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
