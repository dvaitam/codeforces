package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type point struct {
	x, y float64
}

type circle struct {
	x, y, r float64
}

func runProgram(bin, input string) (string, error) {
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
	return out.String(), nil
}

func parseInput(data []byte) (int, int, []point, error) {
	reader := bytes.NewReader(data)
	var k, n int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to read k: %v", err)
	}
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to read n: %v", err)
	}
	total := k * n
	pts := make([]point, total)
	for i := 0; i < total; i++ {
		var xi, yi int
		if _, err := fmt.Fscan(reader, &xi, &yi); err != nil {
			return 0, 0, nil, fmt.Errorf("failed to read sample %d: %v", i+1, err)
		}
		pts[i] = point{float64(xi), float64(yi)}
	}
	return k, n, pts, nil
}

func parseCandidate(output string, k int) ([]circle, error) {
	fields := strings.Fields(output)
	if len(fields) < 3*k {
		return nil, fmt.Errorf("expected at least %d numbers, got %d", 3*k, len(fields))
	}
	res := make([]circle, 0, k)
	idx := 0
	for i := 0; i < k; i++ {
		if idx+3 > len(fields) {
			return nil, fmt.Errorf("insufficient tokens for circle %d", i+1)
		}
		x, err := strconv.ParseFloat(fields[idx], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid x for circle %d: %v", i+1, err)
		}
		y, err := strconv.ParseFloat(fields[idx+1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid y for circle %d: %v", i+1, err)
		}
		r, err := strconv.ParseFloat(fields[idx+2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid r for circle %d: %v", i+1, err)
		}
		if !(r > 0 && !math.IsNaN(r) && !math.IsInf(r, 0)) {
			return nil, fmt.Errorf("radius for circle %d must be positive finite", i+1)
		}
		if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
			return nil, fmt.Errorf("center for circle %d must be finite", i+1)
		}
		res = append(res, circle{x, y, r})
		idx += 3
	}
	return res, nil
}

func recoverRings(k int, pts []point) []circle {
	remain := make([]point, len(pts))
	copy(remain, pts)
	res := make([]circle, 0, k)
	for ring := 0; ring < k && len(remain) >= 3; ring++ {
		bestCnt := 0
		var best circle
		var bestInliers []int
		for idx := range remain {
			i1, i2 := nearestTwo(remain, idx)
			if i1 == -1 || i2 == -1 {
				continue
			}
			cx, cy, cr, ok := fitCircle(remain[idx], remain[i1], remain[i2])
			if !ok || cr <= 0 {
				continue
			}
			cnt, inliers := countInliers(remain, cx, cy, cr)
			if cnt > bestCnt {
				bestCnt = cnt
				best = circle{cx, cy, cr}
				bestInliers = inliers
			}
		}
		if len(bestInliers) > 3 {
			sub := make([]point, len(bestInliers))
			for i, id := range bestInliers {
				sub[i] = remain[id]
			}
			cx, cy, cr := fitCircleKasa(sub)
			best = circle{cx, cy, cr}
		}
		res = append(res, best)
		mark := make([]bool, len(remain))
		for _, id := range bestInliers {
			if id >= 0 && id < len(mark) {
				mark[id] = true
			}
		}
		tmp := remain[:0]
		for i, p := range remain {
			if !mark[i] {
				tmp = append(tmp, p)
			}
		}
		remain = tmp
	}
	return res
}

func nearestTwo(pts []point, idx int) (int, int) {
	min1, min2 := math.MaxFloat64, math.MaxFloat64
	id1, id2 := -1, -1
	pi := pts[idx]
	for j, pj := range pts {
		if j == idx {
			continue
		}
		d := math.Hypot(pi.x-pj.x, pi.y-pj.y)
		if d < min1 {
			min2, id2 = min1, id1
			min1, id1 = d, j
		} else if d < min2 {
			min2 = d
			id2 = j
		}
	}
	return id1, id2
}

func countInliers(pts []point, cx, cy, cr float64) (int, []int) {
	tol := 0.12*cr + 5000
	if tol < 60000 {
		tol = 60000
	}
	var ids []int
	for i, p := range pts {
		d := math.Hypot(p.x-cx, p.y-cy)
		if math.Abs(d-cr) <= tol {
			ids = append(ids, i)
		}
	}
	return len(ids), ids
}

func fitCircle(p1, p2, p3 point) (cx, cy, r float64, ok bool) {
	x1, y1 := p1.x, p1.y
	x2, y2 := p2.x, p2.y
	x3, y3 := p3.x, p3.y
	d := 2 * (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2))
	if math.Abs(d) < 1e-8 {
		return 0, 0, 0, false
	}
	sq1 := x1*x1 + y1*y1
	sq2 := x2*x2 + y2*y2
	sq3 := x3*x3 + y3*y3
	ux := (sq1*(y2-y3) + sq2*(y3-y1) + sq3*(y1-y2)) / d
	uy := (sq1*(x3-x2) + sq2*(x1-x3) + sq3*(x2-x1)) / d
	r = math.Hypot(x1-ux, y1-uy)
	return ux, uy, r, true
}

func fitCircleKasa(pts []point) (cx, cy, r float64) {
	m := float64(len(pts))
	if m == 0 {
		return 0, 0, 0
	}
	var mx, my float64
	for _, p := range pts {
		mx += p.x
		my += p.y
	}
	mx /= m
	my /= m
	var Suu, Suv, Svv, Suz, Svz float64
	for _, p := range pts {
		X := p.x - mx
		Y := p.y - my
		Z := X*X + Y*Y
		Suu += X * X
		Suv += X * Y
		Svv += Y * Y
		Suz += X * Z
		Svz += Y * Z
	}
	det := Suu*Svv - Suv*Suv
	var a, b float64
	if math.Abs(det) > 1e-12 {
		a = (Svv*Suz - Suv*Svz) / det
		b = (Suu*Svz - Suv*Suz) / det
	}
	cx = mx + a/2
	cy = my + b/2
	var sum float64
	for _, p := range pts {
		sum += math.Hypot(p.x-cx, p.y-cy)
	}
	r = sum / m
	return
}

func hausdorff(a, b circle) float64 {
	d := math.Hypot(a.x-b.x, a.y-b.y)
	r1, r2 := a.r, b.r
	dpp := math.Abs(d + r1 + r2)
	dpm := math.Abs(d + r1 - r2)
	dmp := math.Abs(d - r1 + r2)
	dmm := math.Abs(d - r1 - r2)
	return math.Max(
		math.Max(math.Min(dmm, dmp), math.Min(dpm, dpp)),
		math.Max(math.Min(dmm, dpm), math.Min(dmp, dpp)),
	)
}

func checkMatching(candidate, reference []circle, tol float64) error {
	if len(candidate) != len(reference) {
		return fmt.Errorf("mismatched circle counts")
	}
	k := len(candidate)
	used := make([]bool, k)
	var dfs func(int) bool
	dfs = func(idx int) bool {
		if idx == k {
			return true
		}
		for j := 0; j < k; j++ {
			if used[j] {
				continue
			}
			if hausdorff(candidate[idx], reference[j]) > tol {
				continue
			}
			used[j] = true
			if dfs(idx + 1) {
				return true
			}
			used[j] = false
		}
		return false
	}
	if dfs(0) {
		return nil
	}
	return fmt.Errorf("no permutation matches recovered rings within tolerance %.0f", tol)
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC3.go /path/to/binary")
		os.Exit(1)
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	k, _, pts, err := parseInput(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	output, err := runProgram(args[0], string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	candidate, err := parseCandidate(output, k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "output parse error: %v\n", err)
		os.Exit(1)
	}
	reference := recoverRings(k, pts)
	const tolerance = 150000.0
	if err := checkMatching(candidate, reference, tolerance); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
