package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Line struct {
	points []int
}

type testCase struct {
	in  string
	out string
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func matMul(a, b [][]float64) [][]float64 {
	n := len(a)
	c := make([][]float64, n)
	for i := 0; i < n; i++ {
		c[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c
}

func matVecMul(a [][]float64, v []float64) []float64 {
	n := len(a)
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := 0.0
		row := a[i]
		for j := 0; j < n; j++ {
			sum += row[j] * v[j]
		}
		res[i] = sum
	}
	return res
}

func expected(n int, x, y []int, queries [][2]int) string {
	type key struct{ dx, dy, c int }
	lineMap := make(map[key]map[int]struct{})
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := x[j] - x[i]
			dy := y[j] - y[i]
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			c := dy*x[i] - dx*y[i]
			k := key{dx, dy, c}
			if lineMap[k] == nil {
				lineMap[k] = make(map[int]struct{})
			}
			lineMap[k][i] = struct{}{}
			lineMap[k][j] = struct{}{}
		}
	}
	lines := make([]Line, 0, len(lineMap))
	for _, m := range lineMap {
		pts := make([]int, 0, len(m))
		for p := range m {
			pts = append(pts, p)
		}
		if len(pts) >= 2 {
			lines = append(lines, Line{pts})
		}
	}
	pointLines := make([][]int, n)
	for idx, ln := range lines {
		for _, p := range ln.points {
			pointLines[p] = append(pointLines[p], idx)
		}
	}
	P := make([][]float64, n)
	for i := 0; i < n; i++ {
		P[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		deg := len(pointLines[i])
		if deg == 0 {
			continue
		}
		for _, li := range pointLines[i] {
			pts := lines[li].points
			prob := 1.0 / float64(deg) / float64(len(pts))
			for _, j := range pts {
				P[i][j] += prob
			}
		}
	}
	maxM := 10000
	mats := make([][][]float64, 0)
	mats = append(mats, P)
	for (1 << uint(len(mats))) <= maxM {
		next := matMul(mats[len(mats)-1], mats[len(mats)-1])
		mats = append(mats, next)
	}
	var sb strings.Builder
	for _, q := range queries {
		t := q[0] - 1
		m := q[1]
		step := m - 1
		w := make([]float64, n)
		w[t] = 1
		bit := 0
		for step > 0 {
			if step&1 == 1 {
				w = matVecMul(mats[bit], w)
			}
			step >>= 1
			bit++
		}
		ans := 0.0
		for _, ln := range lines {
			sum := 0.0
			for _, p := range ln.points {
				sum += w[p]
			}
			prob := sum / float64(len(ln.points))
			if prob > ans {
				ans = prob
			}
		}
		sb.WriteString(fmt.Sprintf("%.9f\n", ans))
	}
	return strings.TrimSpace(sb.String())
}

func generate() []testCase {
	const T = 100
	rng := rand.New(rand.NewSource(5))
	cases := make([]testCase, T)
	for i := 0; i < T; i++ {
		n := rng.Intn(4) + 2
		pts := make(map[[2]int]struct{})
		x := make([]int, n)
		y := make([]int, n)
		for j := 0; j < n; j++ {
			for {
				xx := rng.Intn(11) - 5
				yy := rng.Intn(11) - 5
				key := [2]int{xx, yy}
				if _, ok := pts[key]; !ok {
					pts[key] = struct{}{}
					x[j] = xx
					y[j] = yy
					break
				}
			}
		}
		qn := rng.Intn(3) + 1
		queries := make([][2]int, qn)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", x[j], y[j])
		}
		fmt.Fprintf(&sb, "%d\n", qn)
		for j := 0; j < qn; j++ {
			t := rng.Intn(n) + 1
			m := rng.Intn(10) + 1
			queries[j] = [2]int{t, m}
			fmt.Fprintf(&sb, "%d %d\n", t, m)
		}
		cases[i] = testCase{
			in:  sb.String(),
			out: expected(n, x, y, queries),
		}
	}
	return cases
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
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.out {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
