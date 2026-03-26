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

// ── embedded solver (CF-accepted 598F) ──────────────────────────────

type Point struct{ x, y float64 }

func solveCase(input string) []float64 {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	readFloat := func() float64 {
		scanner.Scan()
		f, _ := strconv.ParseFloat(scanner.Text(), 64)
		return f
	}
	readInt := func() int {
		scanner.Scan()
		i, _ := strconv.Atoi(scanner.Text())
		return i
	}

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	m := readInt()

	origV := make([]Point, n)
	Vx := make([]int64, n)
	Vy := make([]int64, n)

	for i := 0; i < n; i++ {
		origV[i].x = readFloat()
		origV[i].y = readFloat()
		Vx[i] = int64(math.Round(origV[i].x * 100))
		Vy[i] = int64(math.Round(origV[i].y * 100))
	}

	results := make([]float64, m)
	for q := 0; q < m; q++ {
		origA := Point{readFloat(), readFloat()}
		origB := Point{readFloat(), readFloat()}

		Ax := int64(math.Round(origA.x * 100))
		Ay := int64(math.Round(origA.y * 100))
		Bx := int64(math.Round(origB.x * 100))
		By := int64(math.Round(origB.y * 100))

		a := Ay - By
		b := Bx - Ax
		c := Ax*By - Bx*Ay

		E := make([]int64, n)
		for j := 0; j < n; j++ {
			E[j] = a*Vx[j] + b*Vy[j] + c
		}

		L := math.Hypot(origB.x-origA.x, origB.y-origA.y)
		ux := (origB.x - origA.x) / L
		uy := (origB.y - origA.y) / L

		var pts []float64

		for j := 0; j < n; j++ {
			if E[j] == 0 {
				p := (origV[j].x-origA.x)*ux + (origV[j].y-origA.y)*uy
				pts = append(pts, p)
			}
			j1 := (j + 1) % n
			if (E[j] > 0 && E[j1] < 0) || (E[j] < 0 && E[j1] > 0) {
				t := float64(E[j]) / float64(E[j]-E[j1])
				Qx := origV[j].x + t*(origV[j1].x-origV[j].x)
				Qy := origV[j].y + t*(origV[j1].y-origV[j].y)
				p := (Qx-origA.x)*ux + (Qy-origA.y)*uy
				pts = append(pts, p)
			}
		}

		sortFloat64s(pts)

		totalLength := 0.0

		for j := 0; j < len(pts)-1; j++ {
			p1 := pts[j]
			p2 := pts[j+1]
			if p2-p1 < 1e-11 {
				continue
			}
			mp := (p1 + p2) / 2.0

			onBoundary := false
			for k := 0; k < n; k++ {
				k1 := (k + 1) % n
				if E[k] == 0 && E[k1] == 0 {
					pk := (origV[k].x-origA.x)*ux + (origV[k].y-origA.y)*uy
					pk1 := (origV[k1].x-origA.x)*ux + (origV[k1].y-origA.y)*uy
					minp := math.Min(pk, pk1)
					maxp := math.Max(pk, pk1)
					if mp > minp+1e-7 && mp < maxp-1e-7 {
						onBoundary = true
						break
					}
				}
			}

			if onBoundary {
				totalLength += (p2 - p1)
			} else {
				Mx := origA.x + mp*ux
				My := origA.y + mp*uy

				inside := false
				for k := 0; k < n; k++ {
					k1 := (k + 1) % n
					xi, yi := origV[k].x, origV[k].y
					xj, yj := origV[k1].x, origV[k1].y

					if (yi > My) != (yj > My) {
						intersect := xi + (My-yi)*(xj-xi)/(yj-yi)
						if Mx < intersect {
							inside = !inside
						}
					}
				}

				if inside {
					totalLength += (p2 - p1)
				}
			}
		}

		results[q] = totalLength
	}
	return results
}

func sortFloat64s(a []float64) {
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

// ── test generation (convex polygons to guarantee simple polygon) ───

func generateConvexPoly(rng *rand.Rand, n int) []Point {
	// Generate random angles, sort, place on random radii
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	// Check uniqueness (regenerate if duplicates)
	poly := make([]Point, n)
	cx, cy := rng.Float64()*10-5, rng.Float64()*10-5
	for i := 0; i < n; i++ {
		r := 1 + rng.Float64()*9
		poly[i] = Point{cx + r*math.Cos(angles[i]), cy + r*math.Sin(angles[i])}
	}
	return poly
}

func generateCase(rng *rand.Rand) (string, []float64) {
	n := rng.Intn(8) + 3
	m := rng.Intn(5) + 1
	poly := generateConvexPoly(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%.2f %.2f\n", poly[i].x, poly[i].y))
	}
	for i := 0; i < m; i++ {
		x1 := rng.Float64()*40 - 20
		y1 := rng.Float64()*40 - 20
		x2 := rng.Float64()*40 - 20
		y2 := rng.Float64()*40 - 20
		for math.Hypot(x2-x1, y2-y1) < 1e-6 {
			x2 = rng.Float64()*40 - 20
			y2 = rng.Float64()*40 - 20
		}
		sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f\n", x1, y1, x2, y2))
	}
	input := sb.String()
	expected := solveCase(input)
	return input, expected
}

func runExe(exe, input string) (string, error) {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		in, expected := generateCase(rng)
		out, err := runExe(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		gotFields := strings.Fields(out)
		if len(gotFields) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d values, got %d\ninput:\n%s", i+1, len(expected), len(gotFields), in)
			os.Exit(1)
		}
		for j := range gotFields {
			g, err1 := strconv.ParseFloat(gotFields[j], 64)
			e := expected[j]
			if err1 != nil {
				fmt.Fprintf(os.Stderr, "case %d line %d: parse error: %v\n", i+1, j+1, err1)
				os.Exit(1)
			}
			rel := math.Max(1, math.Abs(e))
			if math.Abs(g-e) > 1e-4*rel {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected %.10f got %.10f\ninput:\n%s", i+1, j+1, e, g, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
