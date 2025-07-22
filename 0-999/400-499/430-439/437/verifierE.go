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

const mod = 1000000007

type Point struct {
	x, y int64
}

func orient(a, b, c Point) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func onSegment(a, b, c Point) bool {
	return c.x >= minInt64(a.x, b.x) && c.x <= maxInt64(a.x, b.x) &&
		c.y >= minInt64(a.y, b.y) && c.y <= maxInt64(a.y, b.y)
}

func intersect(a, b, c, d Point) bool {
	o1 := orient(a, b, c)
	o2 := orient(a, b, d)
	o3 := orient(c, d, a)
	o4 := orient(c, d, b)
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
	return (o1 > 0 && o2 < 0 || o1 < 0 && o2 > 0) &&
		(o3 > 0 && o4 < 0 || o3 < 0 && o4 > 0)
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func pointInPoly(poly []Point, p struct{ x, y float64 }) bool {
	cnt := 0
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		ay := float64(a.y)
		by := float64(b.y)
		if ay > by {
			ay, by = by, ay
		}
		if p.y <= ay || p.y > by {
			continue
		}
		ax := float64(a.x)
		bx := float64(b.x)
		xint := ax + (bx-ax)*((p.y-ay)/(by-ay))
		if xint > p.x {
			cnt++
		}
	}
	return cnt%2 == 1
}

func runBinary(bin, input string) (string, error) {
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

func solveE(poly []Point) string {
	n := len(poly)
	var area2 int64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area2 += poly[i].x*poly[j].y - poly[i].y*poly[j].x
	}
	ccw := area2 > 0
	valid := make([][]bool, n)
	for i := range valid {
		valid[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		valid[i][(i+1)%n] = true
		valid[(i+1)%n][i] = true
	}
	for i := 0; i < n; i++ {
		for j := i + 2; j < n; j++ {
			if i == 0 && j == n-1 {
				valid[i][j] = true
				valid[j][i] = true
				continue
			}
			a := poly[i]
			b := poly[j]
			ok := true
			for k := 0; k < n; k++ {
				k2 := (k + 1) % n
				if k == i || k == j || k2 == i || k2 == j {
					continue
				}
				if intersect(a, b, poly[k], poly[k2]) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			mx := (float64(a.x) + float64(b.x)) * 0.5
			my := (float64(a.y) + float64(b.y)) * 0.5
			if !pointInPoly(poly, struct{ x, y float64 }{mx, my}) {
				continue
			}
			valid[i][j] = true
			valid[j][i] = true
		}
	}
	dp := make([][]int64, n)
	for i := range dp {
		dp[i] = make([]int64, n)
		dp[i][i] = 1
		if i+1 < n {
			dp[i][i+1] = 1
		}
	}
	for length := 2; length < n; length++ {
		for i := 0; i+length < n; i++ {
			j := i + length
			if !valid[i][j] {
				dp[i][j] = 0
				continue
			}
			var ways int64
			for k := i + 1; k < j; k++ {
				if !valid[i][k] || !valid[k][j] {
					continue
				}
				o := orient(poly[i], poly[k], poly[j])
				if (ccw && o <= 0) || (!ccw && o >= 0) {
					continue
				}
				ways = (ways + dp[i][k]*dp[k][j]) % mod
			}
			dp[i][j] = ways
		}
	}
	return fmt.Sprintf("%d", dp[0][n-1])
}

func genPolygon(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 3
	angles := make([]float64, n)
	for i := range angles {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	poly := make([]Point, n)
	for i, ang := range angles {
		r := rng.Float64()*50 + 1
		x := int64(math.Round(r * math.Cos(ang) * 100))
		y := int64(math.Round(r * math.Sin(ang) * 100))
		poly[i] = Point{x, y}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, p := range poly {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	input := sb.String()
	exp := solveE(poly)
	return input, exp
}

func manualCase() (string, string) {
	poly := []Point{{0, 0}, {100, 0}, {100, 100}, {0, 100}}
	var sb strings.Builder
	fmt.Fprintf(&sb, "4\n")
	for _, p := range poly {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	input := sb.String()
	exp := solveE(poly)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	in1, ex1 := manualCase()
	cases = append(cases, [2]string{in1, ex1})
	for i := 0; i < 100; i++ {
		in, exp := genPolygon(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for idx, tc := range cases {
		out, err := runBinary(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
