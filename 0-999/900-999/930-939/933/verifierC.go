package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Circle struct{ x, y, r float64 }

type Point struct{ x, y float64 }

const eps = 1e-9

func intersections(a, b Circle) []Point {
	dx := b.x - a.x
	dy := b.y - a.y
	d := math.Hypot(dx, dy)
	if d > a.r+b.r+eps || d < math.Abs(a.r-b.r)-eps || d == 0 {
		return nil
	}
	alpha := (a.r*a.r - b.r*b.r + d*d) / (2 * d)
	h2 := a.r*a.r - alpha*alpha
	if h2 < eps {
		h2 = 0
	}
	xm := a.x + alpha*dx/d
	ym := a.y + alpha*dy/d
	if h2 == 0 {
		return []Point{{xm, ym}}
	}
	h := math.Sqrt(h2)
	rx := -dy * h / d
	ry := dx * h / d
	return []Point{{xm + rx, ym + ry}, {xm - rx, ym - ry}}
}

func addPoint(pts []Point, p Point) ([]Point, int) {
	for i, q := range pts {
		if (p.x-q.x)*(p.x-q.x)+(p.y-q.y)*(p.y-q.y) < eps*eps {
			return pts, i
		}
	}
	return append(pts, p), len(pts)
}

func find(par []int, x int) int {
	if par[x] != x {
		par[x] = find(par, par[x])
	}
	return par[x]
}

func union(par []int, x, y int) {
	rx := find(par, x)
	ry := find(par, y)
	if rx != ry {
		par[ry] = rx
	}
}

func expectedC(n int, circles []Circle) int {
	points := make([]Point, 0)
	sets := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		sets[i] = make(map[int]struct{})
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for _, p := range intersections(circles[i], circles[j]) {
				var idx int
				points, idx = addPoint(points, p)
				sets[i][idx] = struct{}{}
				sets[j][idx] = struct{}{}
			}
		}
	}
	V := len(points)
	E := 0
	loops := 0
	for i := 0; i < n; i++ {
		if len(sets[i]) == 0 {
			loops++
			E++
		} else {
			E += len(sets[i])
		}
	}
	V += loops
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for idx := range sets[i] {
				if _, ok := sets[j][idx]; ok {
					union(parent, i, j)
					break
				}
			}
		}
	}
	comp := make(map[int]struct{})
	for i := 0; i < n; i++ {
		comp[find(parent, i)] = struct{}{}
	}
	C := len(comp)
	return E - V + C + 1
}

func generateCaseC(rng *rand.Rand) (int, []Circle) {
	n := rng.Intn(3) + 1
	circles := make([]Circle, n)
	used := make(map[[3]int]struct{})
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(21) - 10
			y := rng.Intn(21) - 10
			r := rng.Intn(10) + 1
			key := [3]int{x, y, r}
			if _, ok := used[key]; !ok {
				used[key] = struct{}{}
				circles[i] = Circle{float64(x), float64(y), float64(r)}
				break
			}
		}
	}
	return n, circles
}

func runCaseC(bin string, n int, circles []Circle) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, c := range circles {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", int(c.x), int(c.y), int(c.r)))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expectedC(n, circles)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, circles := generateCaseC(rng)
		if err := runCaseC(bin, n, circles); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
