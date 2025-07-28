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

type Point struct {
	x, y float64
}

type Segment struct {
	a, b Point
}

func sub(a, b Point) Point     { return Point{a.x - b.x, a.y - b.y} }
func dot(a, b Point) float64   { return a.x*b.x + a.y*b.y }
func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func norm(a Point) float64     { return math.Hypot(a.x, a.y) }
func normalize(a Point) Point {
	r := norm(a)
	if r == 0 {
		return Point{0, 0}
	}
	return Point{a.x / r, a.y / r}
}
func angleDiff(a, b Point) float64 {
	a = normalize(a)
	b = normalize(b)
	c := math.Abs(dot(a, b))
	if c > 1 {
		c = 1
	}
	return math.Acos(c)
}
func isOnRay(p, d, t Point) (bool, float64) {
	v := sub(t, p)
	if math.Abs(cross(v, d)) > 1e-9 {
		return false, 0
	}
	k := dot(v, d) / dot(d, d)
	if k >= -1e-9 {
		return true, k
	}
	return false, 0
}
func intersectRaySegment(p, d Point, s Segment) (bool, float64) {
	a := s.a
	b := s.b
	v := sub(b, a)
	w := sub(a, p)
	denom := cross(d, v)
	if math.Abs(denom) < 1e-9 {
		if math.Abs(cross(w, d)) < 1e-9 {
			t1 := dot(sub(a, p), d) / dot(d, d)
			t2 := dot(sub(b, p), d) / dot(d, d)
			best := math.Inf(1)
			if t1 >= -1e-9 {
				best = math.Min(best, t1)
			}
			if t2 >= -1e-9 {
				best = math.Min(best, t2)
			}
			if best < math.Inf(1) {
				return true, best
			}
		}
		return false, 0
	}
	t := cross(w, v) / denom
	u := cross(w, d) / denom
	if t >= -1e-9 && u >= -1e-9 && u <= 1+1e-9 {
		return true, t
	}
	return false, 0
}
func attempt(start, d Point, segs []Segment) bool {
	p := start
	visited := make(map[[2]int]bool)
	for iter := 0; iter <= len(segs); iter++ {
		hasT, distT := isOnRay(p, d, Point{0, 0})
		bestDist := math.Inf(1)
		bestIdx := -1
		for i, s := range segs {
			ok, dist := intersectRaySegment(p, d, s)
			if ok && dist > 1e-9 && dist < bestDist {
				bestDist = dist
				bestIdx = i
			}
		}
		if hasT && (bestIdx == -1 || distT <= bestDist+1e-9) {
			return true
		}
		if bestIdx == -1 {
			return false
		}
		seg := segs[bestIdx]
		v := sub(seg.b, seg.a)
		endpoint := seg.b
		slideVec := v
		if dot(d, v) < 0 {
			endpoint = seg.a
			slideVec = Point{-v.x, -v.y}
		}
		if angleDiff(d, slideVec) >= math.Pi/4-1e-9 {
			return false
		}
		key := [2]int{int(endpoint.x*1000 + 0.5), int(endpoint.y*1000 + 0.5)}
		if visited[key] {
			return false
		}
		visited[key] = true
		p = endpoint
	}
	return false
}

func solve(segs []Segment, queries []Point) []string {
	var res []string
	for _, start := range queries {
		dirs := []Point{{-start.x, -start.y}}
		for _, s := range segs {
			dirs = append(dirs, Point{-s.a.x, -s.a.y}, Point{-s.b.x, -s.b.y})
		}
		found := false
		for _, d := range dirs {
			if d.x == 0 && d.y == 0 {
				continue
			}
			if attempt(start, d, segs) {
				found = true
				break
			}
		}
		if found {
			res = append(res, "YES")
		} else {
			res = append(res, "NO")
		}
	}
	return res
}

func buildInput(segs []Segment, queries []Point) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(segs)))
	for _, s := range segs {
		sb.WriteString(fmt.Sprintf("%f %f %f %f\n", s.a.x, s.a.y, s.b.x, s.b.y))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%f %f\n", q.x, q.y))
	}
	return sb.String()
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		segs := make([]Segment, n)
		for j := 0; j < n; j++ {
			ax := rng.Float64()*6 - 3
			ay := rng.Float64()*6 - 3
			bx := ax + rng.Float64()*2 - 1
			by := ay + rng.Float64()*2 - 1
			segs[j] = Segment{Point{ax, ay}, Point{bx, by}}
		}
		q := rng.Intn(3) + 1
		queries := make([]Point, q)
		for j := range queries {
			queries[j] = Point{rng.Float64()*6 - 3, rng.Float64()*6 - 3}
		}
		input := buildInput(segs, queries)
		res := solve(segs, queries)
		exp := strings.Join(res, "\n")
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("case %d wrong answer\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
