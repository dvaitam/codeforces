package main

import (
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

type Pnt struct{ x, y float64 }

func get(a, b Pnt, s, w float64) float64 {
	dir := Pnt{b.x - a.x, b.y - a.y}
	length := math.Hypot(dir.x, dir.y)
	best := 0.0
	twoPi := 2 * math.Pi
	for i := 0; i <= 100; i++ {
		t := float64(i) / 100.0
		x := a.x + dir.x*t
		y := a.y + dir.y*t
		r := math.Atan2(y, x) - s
		for r < 0 {
			r += twoPi
		}
		for r > twoPi {
			r -= twoPi
		}
		minr := math.Min(r, twoPi-r)
		val := length * t * w / minr
		if val > best {
			best = val
		}
	}
	return best
}

func solve(ax, ay, bx, by float64, tanks []struct{ x, y, s, w float64 }, k int) float64 {
	a := Pnt{ax, ay}
	b := Pnt{bx, by}
	v := make([]float64, 0, len(tanks)+1)
	v = append(v, 0.0)
	for _, t := range tanks {
		p := Pnt{t.x, t.y}
		ap := Pnt{a.x - p.x, a.y - p.y}
		bp := Pnt{b.x - p.x, b.y - p.y}
		v = append(v, get(ap, bp, t.s, t.w))
	}
	sort.Float64s(v)
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
	return v[k]
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		ax := float64(rng.Intn(11) - 5)
		ay := float64(rng.Intn(11) - 5)
		bx := float64(rng.Intn(11) - 5)
		by := float64(rng.Intn(11) - 5)
		if ax == bx && ay == by {
			bx++
		}
		n := rng.Intn(3) + 1
		tanks := make([]struct{ x, y, s, w float64 }, n)
		for j := 0; j < n; j++ {
			tanks[j].x = float64(rng.Intn(11) - 5)
			tanks[j].y = float64(rng.Intn(11) - 5)
			tanks[j].s = rng.Float64() * 2 * math.Pi
			tanks[j].w = rng.Float64()*2 + 0.1
		}
		k := rng.Intn(n + 1)
		expected := solve(ax, ay, bx, by, tanks, k)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%.0f %.0f %.0f %.0f\n", ax, ay, bx, by))
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, t := range tanks {
			sb.WriteString(fmt.Sprintf("%.0f %.0f %.5f %.5f\n", t.x, t.y, t.s, t.w))
		}
		sb.WriteString(fmt.Sprintf("%d\n", k))
		input := sb.String()
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(gotStr), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output\n", i+1)
			os.Exit(1)
		}
		if math.Abs(got-expected) > 1e-4*math.Max(1.0, math.Abs(expected)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
