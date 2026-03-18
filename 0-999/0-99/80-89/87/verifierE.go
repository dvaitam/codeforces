package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type point struct{ x, y int }

const embeddedRefSource = `package main

import (
	"bufio"
	"io"
	"os"
)

type Point struct {
	x, y int64
}

type FastScanner struct {
	data []byte
	idx  int
}

func (fs *FastScanner) nextInt() int64 {
	n := len(fs.data)
	for fs.idx < n && fs.data[fs.idx] <= ' ' {
		fs.idx++
	}
	sign := int64(1)
	if fs.data[fs.idx] == '-' {
		sign = -1
		fs.idx++
	}
	var val int64
	for fs.idx < n {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int64(c-'0')
		fs.idx++
	}
	return sign * val
}

func add(a, b Point) Point {
	return Point{a.x + b.x, a.y + b.y}
}

func sub(a, b Point) Point {
	return Point{a.x - b.x, a.y - b.y}
}

func crossSignXY(ax, ay, bx, by int64) int {
	p1 := ax * by
	p2 := ay * bx
	if p1 >= 0 && p2 < 0 {
		return 1
	}
	if p1 < 0 && p2 >= 0 {
		return -1
	}
	if p1 > p2 {
		return 1
	}
	if p1 < p2 {
		return -1
	}
	return 0
}

func crossSign(a, b Point) int {
	return crossSignXY(a.x, a.y, b.x, b.y)
}

func orient(a, b, c Point) int {
	return crossSignXY(b.x-a.x, b.y-a.y, c.x-a.x, c.y-a.y)
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func onSegment(a, b, p Point) bool {
	return orient(a, b, p) == 0 &&
		min64(a.x, b.x) <= p.x && p.x <= max64(a.x, b.x) &&
		min64(a.y, b.y) <= p.y && p.y <= max64(a.y, b.y)
}

func normalize(poly []Point) []Point {
	n := len(poly)
	idx := 0
	for i := 1; i < n; i++ {
		if poly[i].y < poly[idx].y || (poly[i].y == poly[idx].y && poly[i].x < poly[idx].x) {
			idx = i
		}
	}
	if idx == 0 {
		return poly
	}
	res := make([]Point, n)
	copy(res, poly[idx:])
	copy(res[n-idx:], poly[:idx])
	return res
}

func minkowski(a, b []Point) []Point {
	a = normalize(a)
	b = normalize(b)

	na, nb := len(a), len(b)
	ea := make([]Point, na)
	eb := make([]Point, nb)

	for i := 0; i < na; i++ {
		ea[i] = sub(a[(i+1)%na], a[i])
	}
	for i := 0; i < nb; i++ {
		eb[i] = sub(b[(i+1)%nb], b[i])
	}

	res := make([]Point, 0, na+nb)
	cur := add(a[0], b[0])
	res = append(res, cur)

	i, j := 0, 0
	for i < na || j < nb {
		if i == na {
			cur = add(cur, eb[j])
			j++
		} else if j == nb {
			cur = add(cur, ea[i])
			i++
		} else {
			s := crossSign(ea[i], eb[j])
			if s > 0 {
				cur = add(cur, ea[i])
				i++
			} else if s < 0 {
				cur = add(cur, eb[j])
				j++
			} else {
				cur = add(cur, add(ea[i], eb[j]))
				i++
				j++
			}
		}
		res = append(res, cur)
	}

	return normalize(res[:len(res)-1])
}

func containsConvex(poly []Point, q Point) bool {
	n := len(poly)
	if n == 1 {
		return poly[0] == q
	}
	if n == 2 {
		return onSegment(poly[0], poly[1], q)
	}

	s1 := orient(poly[0], poly[1], q)
	if s1 < 0 {
		return false
	}
	s2 := orient(poly[0], poly[n-1], q)
	if s2 > 0 {
		return false
	}
	if s1 == 0 {
		return onSegment(poly[0], poly[1], q)
	}
	if s2 == 0 {
		return onSegment(poly[0], poly[n-1], q)
	}

	l, r := 1, n-1
	for r-l > 1 {
		m := (l + r) >> 1
		if orient(poly[0], poly[m], q) >= 0 {
			l = m
		} else {
			r = m
		}
	}

	return orient(poly[l], poly[l+1], q) >= 0
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	fs := FastScanner{data: data}

	polys := make([][]Point, 3)
	for k := 0; k < 3; k++ {
		n := int(fs.nextInt())
		poly := make([]Point, n)
		for i := 0; i < n; i++ {
			poly[i] = Point{fs.nextInt(), fs.nextInt()}
		}
		polys[k] = poly
	}

	sum := minkowski(minkowski(polys[0], polys[1]), polys[2])

	m := int(fs.nextInt())
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	for i := 0; i < m; i++ {
		q := Point{fs.nextInt() * 3, fs.nextInt() * 3}
		if containsConvex(sum, q) {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
	out.Flush()
}
`

func buildRef() (string, error) {
	ref := "./refE.bin"
	tmpGo := "refE_src.go"
	if err := os.WriteFile(tmpGo, []byte(embeddedRefSource), 0644); err != nil {
		return "", fmt.Errorf("write embedded source: %v", err)
	}
	defer os.Remove(tmpGo)
	cmd := exec.Command("go", "build", "-o", ref, tmpGo)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTriangle(rng *rand.Rand) []point {
	x0 := rng.Intn(10)
	y0 := rng.Intn(10)
	dx1 := rng.Intn(5) + 1
	dx2 := rng.Intn(dx1)
	dy := rng.Intn(5) + 1
	p1 := point{x0, y0}
	p2 := point{x0 + dx1, y0}
	p3 := point{x0 + dx2, y0 + dy}
	return []point{p1, p2, p3}
}

func generateCase(rng *rand.Rand) string {
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		pts := genTriangle(rng)
		fmt.Fprintf(&sb, "%d\n", len(pts))
		for _, p := range pts {
			fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
		}
	}
	m := rng.Intn(5) + 1
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		x := rng.Intn(15)
		y := rng.Intn(15)
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
