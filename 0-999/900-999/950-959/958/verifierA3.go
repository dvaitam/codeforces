package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	testCount     = 60
	requiredRatio = 0.9
	distTolerance = 0.35
)

type point struct {
	x float64
	y float64
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "958A3.go")
	tmp, err := os.CreateTemp("", "oracle958A3")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func randomPoint(r *rand.Rand) point {
	return point{
		x: r.Float64()*20000 - 10000,
		y: r.Float64()*20000 - 10000,
	}
}

func clonePoints(src []point) []point {
	dst := make([]point, len(src))
	copy(dst, src)
	return dst
}

func randomPoints(count int, r *rand.Rand) []point {
	res := make([]point, count)
	for i := 0; i < count; i++ {
		res[i] = randomPoint(r)
	}
	return res
}

func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

func applyTransform(pts []point, angle, dx, dy float64) []point {
	out := make([]point, len(pts))
	c := math.Cos(angle)
	s := math.Sin(angle)
	for i, p := range pts {
		x := c*p.x - s*p.y + dx
		y := s*p.x + c*p.y + dy
		out[i] = point{x: round2(x), y: round2(y)}
	}
	return out
}

func shufflePoints(pts []point, r *rand.Rand) {
	r.Shuffle(len(pts), func(i, j int) {
		pts[i], pts[j] = pts[j], pts[i]
	})
}

func genCase(r *rand.Rand) (string, []point, []point, int, int, int) {
	N := 1000 + r.Intn(1001) // [1000,2000]
	n1 := N + r.Intn(N/2+1)
	n2 := N + r.Intn(N/2+1)
	base := randomPoints(N, r)

	arr1Base := append(clonePoints(base), randomPoints(n1-N, r)...)
	arr2Base := append(clonePoints(base), randomPoints(n2-N, r)...)

	angle1 := r.Float64() * 2 * math.Pi
	angle2 := r.Float64() * 2 * math.Pi
	dx1 := r.Float64()*20000 - 10000
	dy1 := r.Float64()*20000 - 10000
	dx2 := r.Float64()*20000 - 10000
	dy2 := r.Float64()*20000 - 10000

	arr1 := applyTransform(arr1Base, angle1, dx1, dy1)
	arr2 := applyTransform(arr2Base, angle2, dx2, dy2)

	shufflePoints(arr1, r)
	shufflePoints(arr2, r)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", N)
	fmt.Fprintf(&sb, "%d\n", n1)
	for _, p := range arr1 {
		fmt.Fprintf(&sb, "%.2f %.2f\n", p.x, p.y)
	}
	fmt.Fprintf(&sb, "%d\n", n2)
	for _, p := range arr2 {
		fmt.Fprintf(&sb, "%.2f %.2f\n", p.x, p.y)
	}
	return sb.String(), arr1, arr2, N, n1, n2
}

func parsePairs(out string, N, n1, n2 int) ([][2]int, error) {
	fields := strings.Fields(out)
	if len(fields) != 2*N {
		return nil, fmt.Errorf("expected %d pairs, got %d tokens", N, len(fields)/2)
	}
	pairs := make([][2]int, N)
	used1 := make([]bool, n1)
	used2 := make([]bool, n2)
	for i := 0; i < N; i++ {
		a, err := strconv.Atoi(fields[2*i])
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i*2+1, err)
		}
		b, err := strconv.Atoi(fields[2*i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i*2+2, err)
		}
		if a < 1 || a > n1 || b < 1 || b > n2 {
			return nil, fmt.Errorf("pair %d indices out of range", i+1)
		}
		if used1[a-1] {
			return nil, fmt.Errorf("first map index %d used multiple times", a)
		}
		if used2[b-1] {
			return nil, fmt.Errorf("second map index %d used multiple times", b)
		}
		used1[a-1] = true
		used2[b-1] = true
		pairs[i] = [2]int{a - 1, b - 1}
	}
	return pairs, nil
}

func estimateTransform(pPts, qPts []point) (cosv, sinv, tx, ty float64) {
	n := len(pPts)
	var sumPx, sumPy, sumQx, sumQy float64
	for i := 0; i < n; i++ {
		sumPx += pPts[i].x
		sumPy += pPts[i].y
		sumQx += qPts[i].x
		sumQy += qPts[i].y
	}
	cpx := sumPx / float64(n)
	cpy := sumPy / float64(n)
	cqx := sumQx / float64(n)
	cqy := sumQy / float64(n)

	var sxx, sxy, syx, syy float64
	for i := 0; i < n; i++ {
		px := pPts[i].x - cpx
		py := pPts[i].y - cpy
		qx := qPts[i].x - cqx
		qy := qPts[i].y - cqy
		sxx += px * qx
		sxy += px * qy
		syx += py * qx
		syy += py * qy
	}
	c := sxx + syy
	s := syx - sxy
	norm := math.Hypot(c, s)
	if norm < 1e-9 {
		cosv = 1
		sinv = 0
	} else {
		cosv = c / norm
		sinv = s / norm
	}
	tx = cqx - (cosv*cpx - sinv*cpy)
	ty = cqy - (sinv*cpx + cosv*cpy)
	return
}

func countGoodPairs(pPts, qPts []point) int {
	cosv, sinv, tx, ty := estimateTransform(pPts, qPts)
	good := 0
	for i := 0; i < len(pPts); i++ {
		x := cosv*pPts[i].x - sinv*pPts[i].y + tx
		y := sinv*pPts[i].x + cosv*pPts[i].y + ty
		if math.Hypot(x-qPts[i].x, y-qPts[i].y) <= distTolerance {
			good++
		}
	}
	return good
}

func validateOutput(out string, N, n1, n2 int, pts1, pts2 []point) error {
	pairs, err := parsePairs(out, N, n1, n2)
	if err != nil {
		return err
	}
	pPts := make([]point, N)
	qPts := make([]point, N)
	for i := 0; i < N; i++ {
		pPts[i] = pts1[pairs[i][0]]
		qPts[i] = pts2[pairs[i][1]]
	}
	good := countGoodPairs(pPts, qPts)
	ratio := float64(good) / float64(N)
	if ratio+1e-9 < requiredRatio {
		return fmt.Errorf("only %.2f%% pairs consistent with a rigid transform", ratio*100)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA3.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input, pts1, pts2, N, n1, n2 := genCase(r)
		expectStr, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if err := validateOutput(expectStr, N, n1, n2, pts1, pts2); err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := runProgram(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if err := validateOutput(gotStr, N, n1, n2, pts1, pts2); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
