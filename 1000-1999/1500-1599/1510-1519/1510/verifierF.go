package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const tolerance = 1e-6

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := generateTests()
	for idx, tc := range tests {
		refVal := solve(tc.input)

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVal, err := parseFloat(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if !almostEqual(refVal, gotVal, tolerance) {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %.12f got %.12f\ninput:\n%sparticipant output:\n%s\n",
				idx+1, refVal, gotVal, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

// ---------- Embedded CF-accepted solver for 1510F ----------

type sPoint struct {
	x, y float64
}

type sRay struct {
	angle float64
	orig  sPoint
	dir   sPoint
}

type sSector struct {
	F1       sPoint
	F2       sPoint
	L_hidden float64
}

func solve(input string) float64 {
	reader := bufio.NewReader(strings.NewReader(input))

	var n int
	var L float64
	fmt.Fscan(reader, &n, &L)

	V := make([]sPoint, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &V[i].x, &V[i].y)
	}

	pref := make([]float64, n)
	cur := 0.0
	for i := 0; i < n; i++ {
		pref[i] = cur
		p1 := V[i]
		p2 := V[(i+1)%n]
		cur += math.Hypot(p1.x-p2.x, p1.y-p2.y)
	}
	P_poly := cur

	ccwDist := func(i, j int) float64 {
		d := pref[j] - pref[i]
		if d < 0 {
			d += P_poly
		}
		return d
	}

	rays := make([]sRay, 0, 2*n)
	for i := 0; i < n; i++ {
		p1 := V[i]
		p2 := V[(i+1)%n]
		E := sPoint{p2.x - p1.x, p2.y - p1.y}
		LE := math.Hypot(E.x, E.y)
		U := sPoint{E.x / LE, E.y / LE}

		angFwd := math.Atan2(U.y, U.x)
		if angFwd < 0 {
			angFwd += 2 * math.Pi
		}
		rays = append(rays, sRay{angle: angFwd, orig: p2, dir: U})

		Ubwd := sPoint{-U.x, -U.y}
		angBwd := math.Atan2(Ubwd.y, Ubwd.x)
		if angBwd < 0 {
			angBwd += 2 * math.Pi
		}
		rays = append(rays, sRay{angle: angBwd, orig: p1, dir: Ubwd})
	}

	sort.Slice(rays, func(i, j int) bool {
		return rays[i].angle < rays[j].angle
	})

	numRays := len(rays)
	sectors := make([]sSector, numRays)
	for k := 0; k < numRays; k++ {
		ang1 := rays[k].angle
		ang2 := rays[(k+1)%numRays].angle
		if k == numRays-1 {
			ang2 += 2 * math.Pi
		}
		mid := (ang1 + ang2) / 2.0
		ux := math.Cos(mid)
		uy := math.Sin(mid)

		maxVal := -1e18
		minVal := 1e18
		idxLeft := 0
		idxRight := 0
		for i, p := range V {
			val := -uy*p.x + ux*p.y
			if val > maxVal {
				maxVal = val
				idxLeft = i
			}
			if val < minVal {
				minVal = val
				idxRight = i
			}
		}
		sectors[k] = sSector{
			F1:       V[idxLeft],
			F2:       V[idxRight],
			L_hidden: ccwDist(idxLeft, idxRight),
		}
	}

	Pk := make([]sPoint, numRays)
	for k := 0; k < numRays; k++ {
		sec := sectors[k]
		F1 := sec.F1
		F2 := sec.F2
		D := L - sec.L_hidden

		Vorig := rays[k].orig
		U := rays[k].dir

		A := sPoint{Vorig.x - F1.x, Vorig.y - F1.y}
		B := sPoint{Vorig.x - F2.x, Vorig.y - F2.y}

		P := 2.0 * (U.x*(A.x-B.x) + U.y*(A.y-B.y))
		Q := (A.x*A.x + A.y*A.y) - (B.x*B.x + B.y*B.y) - D*D

		a_quad := P*P - 4.0*D*D
		b_quad := 2.0*P*Q - 8.0*D*D*(U.x*B.x+U.y*B.y)
		c_quad := Q*Q - 4.0*D*D*(B.x*B.x+B.y*B.y)

		disc := b_quad*b_quad - 4.0*a_quad*c_quad
		if disc < 0 {
			disc = 0
		}
		r := (-b_quad - math.Sqrt(disc)) / (2.0 * a_quad)

		Pk[k] = sPoint{Vorig.x + r*U.x, Vorig.y + r*U.y}
	}

	totalArea := 0.0
	for k := 0; k < numRays; k++ {
		sec := sectors[k]
		F1 := sec.F1
		F2 := sec.F2
		D := L - sec.L_hidden

		a_el := D / 2.0
		dxF := F1.x - F2.x
		dyF := F1.y - F2.y
		distF := math.Hypot(dxF, dyF)
		c_el := distF / 2.0
		b_el := math.Sqrt(math.Max(0, a_el*a_el-c_el*c_el))

		C := sPoint{(F1.x + F2.x) / 2.0, (F1.y + F2.y) / 2.0}

		var Uel, Vel sPoint
		if distF < 1e-9 {
			Uel = sPoint{1.0, 0.0}
		} else {
			Uel = sPoint{dxF / distF, dyF / distF}
		}
		Vel = sPoint{-Uel.y, Uel.x}

		pStart := Pk[k]
		pEnd := Pk[(k+1)%numRays]

		dx0 := pStart.x - C.x
		dy0 := pStart.y - C.y
		cos0 := (dx0*Uel.x + dy0*Uel.y) / a_el
		sin0 := (dx0*Vel.x + dy0*Vel.y) / b_el
		t0 := math.Atan2(sin0, cos0)

		dx1 := pEnd.x - C.x
		dy1 := pEnd.y - C.y
		cos1 := (dx1*Uel.x + dy1*Uel.y) / a_el
		sin1 := (dx1*Vel.x + dy1*Vel.y) / b_el
		t1 := math.Atan2(sin1, cos1)

		diff := t1 - t0
		for diff < 0 {
			diff += 2.0 * math.Pi
		}
		for diff >= 2.0*math.Pi {
			diff -= 2.0 * math.Pi
		}

		term1 := a_el * b_el * diff
		term2 := a_el * (math.Cos(t1) - math.Cos(t0)) * (C.x*Uel.y - C.y*Uel.x)
		term3 := b_el * (math.Sin(t1) - math.Sin(t0)) * (C.x*Vel.y - C.y*Vel.x)

		totalArea += 0.5 * (term1 + term2 + term3)
	}

	return totalArea
}

// ---------- Utility functions ----------

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseFloat(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single number, got %d tokens", len(fields))
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", fields[0], err)
	}
	return val, nil
}

func almostEqual(a, b, tol float64) bool {
	diff := math.Abs(a - b)
	if diff <= tol {
		return true
	}
	div := math.Max(1.0, math.Abs(a))
	return diff <= tol*div
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomPolygonTests(rng, 80)...)
	tests = append(tests, stressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	polygons := []struct {
		points [][2]float64
		l      float64
	}{
		{
			points: [][2]float64{{0, 0}, {1, 0}, {0, 1}},
			l:      4,
		},
		{
			points: [][2]float64{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			l:      5,
		},
		{
			points: [][2]float64{{0, 0}, {2, -1}, {3, 0}, {4, 3}, {-1, 4}},
			l:      17,
		},
	}
	var tests []testCase
	for _, poly := range polygons {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %.6f\n", len(poly.points), poly.l))
		for _, pt := range poly.points {
			sb.WriteString(fmt.Sprintf("%.6f %.6f\n", pt[0], pt[1]))
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}

func randomPolygonTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for i := 0; i < batches; i++ {
		n := rng.Intn(50) + 3
		radius := rng.Float64()*90 + 10
		points := generateConvexPolygon(rng, n, radius)
		perimeter := polygonPerimeter(points)
		l := perimeter + rng.Float64()*500 + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %.10f\n", len(points), l))
		for _, pt := range points {
			sb.WriteString(fmt.Sprintf("%.10f %.10f\n", pt[0], pt[1]))
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}

func stressTests(rng *rand.Rand) []testCase {
	var tests []testCase
	for i := 0; i < 5; i++ {
		n := 10000
		points := generateConvexPolygon(rng, n, 100000)
		perimeter := polygonPerimeter(points)
		l := perimeter + rng.Float64()*8e5 + 1e-3
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %.10f\n", len(points), l))
		for _, pt := range points {
			sb.WriteString(fmt.Sprintf("%.10f %.10f\n", pt[0], pt[1]))
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}

func generateConvexPolygon(rng *rand.Rand, n int, maxRadius float64) [][2]float64 {
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	points := make([][2]float64, n)
	for i := 0; i < n; i++ {
		r := rng.Float64()*maxRadius + 1
		points[i] = [2]float64{r * math.Cos(angles[i]), r * math.Sin(angles[i])}
	}
	return convexHull(points)
}

func polygonPerimeter(points [][2]float64) float64 {
	per := 0.0
	for i := 0; i < len(points); i++ {
		j := (i + 1) % len(points)
		per += math.Hypot(points[i][0]-points[j][0], points[i][1]-points[j][1])
	}
	return per
}

func convexHull(points [][2]float64) [][2]float64 {
	sort.Slice(points, func(i, j int) bool {
		if points[i][0] == points[j][0] {
			return points[i][1] < points[j][1]
		}
		return points[i][0] < points[j][0]
	})
	nn := len(points)
	if nn <= 1 {
		return points
	}
	half := func(points [][2]float64) [][2]float64 {
		h := make([][2]float64, 0, nn)
		for _, p := range points {
			for len(h) >= 2 && crossP(h[len(h)-2], h[len(h)-1], p) <= 0 {
				h = h[:len(h)-1]
			}
			h = append(h, p)
		}
		return h
	}
	lower := half(points)
	upper := half(reversePoints(points))
	lower = lower[:len(lower)-1]
	upper = upper[:len(upper)-1]
	return append(lower, upper...)
}

func reversePoints(points [][2]float64) [][2]float64 {
	res := make([][2]float64, len(points))
	copy(res, points)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func crossP(a, b, c [2]float64) float64 {
	return (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])
}
