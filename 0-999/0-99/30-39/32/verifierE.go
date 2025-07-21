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

type Point struct{ x, y float64 }

func sub(a, b Point) Point         { return Point{a.x - b.x, a.y - b.y} }
func add(a, b Point) Point         { return Point{a.x + b.x, a.y + b.y} }
func mul(a Point, k float64) Point { return Point{a.x * k, a.y * k} }
func dot(a, b Point) float64       { return a.x*b.x + a.y*b.y }
func cross(a, b Point) float64     { return a.x*b.y - a.y*b.x }
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func orient(a, b, c Point) float64 {
	return cross(sub(b, a), sub(c, a))
}

func onSegment(a, b, c Point) bool {
	if abs(orient(a, b, c)) > 1e-9 {
		return false
	}
	minx, maxx := math.Min(a.x, b.x), math.Max(a.x, b.x)
	miny, maxy := math.Min(a.y, b.y), math.Max(a.y, b.y)
	return c.x >= minx-1e-9 && c.x <= maxx+1e-9 && c.y >= miny-1e-9 && c.y <= maxy+1e-9
}

func segIntersect(a, b, c, d Point) bool {
	o1 := orient(a, b, c)
	o2 := orient(a, b, d)
	o3 := orient(c, d, a)
	o4 := orient(c, d, b)
	if o1*o2 < 0 && o3*o4 < 0 {
		return true
	}
	if abs(o1) < 1e-9 && onSegment(a, b, c) {
		return true
	}
	if abs(o2) < 1e-9 && onSegment(a, b, d) {
		return true
	}
	if abs(o3) < 1e-9 && onSegment(c, d, a) {
		return true
	}
	if abs(o4) < 1e-9 && onSegment(c, d, b) {
		return true
	}
	return false
}

func lineIntersect(p, r, q, s Point) (Point, bool, float64) {
	rxs := cross(r, s)
	if abs(rxs) < 1e-9 {
		return Point{}, false, 0
	}
	qp := sub(q, p)
	t := cross(qp, s) / rxs
	u := cross(qp, r) / rxs
	ip := add(p, mul(r, t))
	return ip, true, u
}

func solveCase(V, P, W1, W2, M1, M2 Point) string {
	if !segIntersect(V, P, W1, W2) && !segIntersect(V, P, M1, M2) {
		return "YES"
	}
	oV := orient(M1, M2, V)
	oP := orient(M1, M2, P)
	if oV*oP <= 0 {
		return "NO"
	}
	ap := sub(P, M1)
	ab := sub(M2, M1)
	ab2 := dot(ab, ab)
	t := dot(ap, ab) / ab2
	proj := add(M1, mul(ab, t))
	Pp := Point{2*proj.x - P.x, 2*proj.y - P.y}
	rVec := sub(Pp, V)
	sVec := ab
	R, ok, u := lineIntersect(V, rVec, M1, sVec)
	if !ok || u < -1e-9 || u > 1+1e-9 {
		return "NO"
	}
	if segIntersect(V, R, W1, W2) || segIntersect(R, P, W1, W2) {
		return "NO"
	}
	return "YES"
}

func randPoint(rng *rand.Rand) Point {
	return Point{float64(rng.Intn(41) - 20), float64(rng.Intn(41) - 20)}
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		V := randPoint(rng)
		P := randPoint(rng)
		if V == P {
			continue
		}
		W1 := randPoint(rng)
		W2 := randPoint(rng)
		if W1 == W2 {
			continue
		}
		M1 := randPoint(rng)
		M2 := randPoint(rng)
		if M1 == M2 {
			continue
		}
		if segIntersect(W1, W2, M1, M2) {
			continue
		}
		if onSegment(W1, W2, V) || onSegment(W1, W2, P) || onSegment(M1, M2, V) || onSegment(M1, M2, P) {
			continue
		}
		input := fmt.Sprintf("%.0f %.0f\n%.0f %.0f\n%.0f %.0f %.0f %.0f\n%.0f %.0f %.0f %.0f\n",
			V.x, V.y, P.x, P.y, W1.x, W1.y, W2.x, W2.y, M1.x, M1.y, M2.x, M2.y)
		exp := solveCase(V, P, W1, W2, M1, M2)
		return input, exp
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
