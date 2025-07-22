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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

const eps = 1e-9

func cut(x, y []float64, aa, bb, cc float64) ([]float64, []float64) {
	n := len(x)
	if n == 0 {
		return x, y
	}
	var nx, ny []float64
	for i := 0; i < n; i++ {
		xi, yi := x[i], y[i]
		zi := aa*xi + bb*yi + cc
		if zi < eps {
			nx = append(nx, xi)
			ny = append(ny, yi)
		}
		j := i + 1
		if j == n {
			j = 0
		}
		xj, yj := x[j], y[j]
		zj := aa*xj + bb*yj + cc
		if (zi < -eps && zj > eps) || (zi > eps && zj < -eps) {
			a := yj - yi
			b := xi - xj
			c := -a*xi - b*yi
			d := a*bb - b*aa
			ix := (b*cc - c*bb) / d
			iy := (c*aa - a*cc) / d
			nx = append(nx, ix)
			ny = append(ny, iy)
		}
	}
	return nx, ny
}

func solveE(r *bufio.Reader) string {
	var w, h, n int
	fmt.Fscan(r, &w, &h, &n)
	pts := make([]struct{ x, y int }, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &pts[i].x, &pts[i].y)
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	xs := make([]struct{ x, y int }, 0, n)
	ws := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if i > 0 && pts[i] == pts[i-1] {
			ws[len(ws)-1]++
		} else {
			xs = append(xs, pts[i])
			ws = append(ws, 1)
		}
	}
	pts = xs
	n = len(pts)
	ansSq := 0.0
	for st := 0; st < n; st++ {
		x := []float64{0, float64(w), float64(w), 0}
		y := []float64{0, 0, float64(h), float64(h)}
		for i := 0; i < n; i++ {
			dx := float64(pts[i].x - pts[st].x)
			dy := float64(pts[i].y - pts[st].y)
			midx := float64(pts[i].x+pts[st].x) * 0.5
			midy := float64(pts[i].y+pts[st].y) * 0.5
			cc := -dx*midx - dy*midy
			x, y = cut(x, y, dx, dy, cc)
			if len(x) == 0 {
				break
			}
		}
		if ws[st] >= 2 {
			for i := 0; i < len(x); i++ {
				dx := x[i] - float64(pts[st].x)
				dy := y[i] - float64(pts[st].y)
				d := dx*dx + dy*dy
				if d > ansSq {
					ansSq = d
				}
			}
			continue
		}
		m := len(x)
		ptsIdx := make([]int, m)
		for i := 0; i < m; i++ {
			j := (i + 1) % m
			mx := (x[i] + x[j]) * 0.5
			my := (y[i] + y[j]) * 0.5
			best := -1
			bestD := 1e100
			for k := 0; k < n; k++ {
				if k == st {
					continue
				}
				dx := mx - float64(pts[k].x)
				dy := my - float64(pts[k].y)
				d := dx*dx + dy*dy
				if d < bestD {
					bestD = d
					best = k
				}
			}
			ptsIdx[i] = best
		}
		ox := make([]float64, len(x))
		oy := make([]float64, len(y))
		copy(ox, x)
		copy(oy, y)
		sz := len(ptsIdx)
		for jj := 0; jj < sz; jj++ {
			pt := ptsIdx[jj]
			x = make([]float64, len(ox))
			y = make([]float64, len(oy))
			copy(x, ox)
			copy(y, oy)
			for u := 0; u < sz; u++ {
				if u == jj {
					continue
				}
				i := ptsIdx[u]
				dx := float64(pts[i].x - pts[pt].x)
				dy := float64(pts[i].y - pts[pt].y)
				midx := float64(pts[i].x+pts[pt].x) * 0.5
				midy := float64(pts[i].y+pts[pt].y) * 0.5
				cc := -dx*midx - dy*midy
				x, y = cut(x, y, dx, dy, cc)
				if len(x) == 0 {
					break
				}
			}
			for i := 0; i < len(x); i++ {
				dx := x[i] - float64(pts[pt].x)
				dy := y[i] - float64(pts[pt].y)
				d := dx*dx + dy*dy
				if d > ansSq {
					ansSq = d
				}
			}
		}
	}
	return fmt.Sprintf("%.17f\n", math.Sqrt(ansSq))
}

func generateCase(rng *rand.Rand) string {
	w := rng.Intn(10) + 1
	h := rng.Intn(10) + 1
	n := rng.Intn(4) + 2
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", w, h, n)
	for i := 0; i < n; i++ {
		x := rng.Intn(w + 1)
		y := rng.Intn(h + 1)
		fmt.Fprintf(&b, "%d %d\n", x, y)
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solveE(bufio.NewReader(strings.NewReader(tc)))
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output\n", i+1)
			os.Exit(1)
		}
		exp, _ := strconv.ParseFloat(strings.TrimSpace(expect), 64)
		if math.Abs(got-exp) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", i+1, exp, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
