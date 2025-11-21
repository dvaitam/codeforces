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
	"strings"
)

const (
	minShapeSize = 15
	separation   = 10
)

type box struct {
	minX, maxX int
	minY, maxY int
}

func (b box) intersects(o box) bool {
	return b.minX <= o.maxX && o.minX <= b.maxX && b.minY <= o.maxY && o.minY <= b.maxY
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "178E1.go")
	tmp, err := os.CreateTemp("", "oracle178E1")
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

func run(bin, input string) (string, error) {
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

func randRange(r *rand.Rand, lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + r.Intn(hi-lo+1)
}

func fits(candidate box, boxes []box, n int) bool {
	if candidate.minX < 0 || candidate.minY < 0 || candidate.maxX >= n || candidate.maxY >= n {
		return false
	}
	for _, b := range boxes {
		if candidate.intersects(b) {
			return false
		}
	}
	return true
}

func drawCircle(grid [][]byte, cx, cy, radius int) {
	n := len(grid)
	r2 := radius * radius
	for x := cx - radius; x <= cx+radius; x++ {
		if x < 0 || x >= n {
			continue
		}
		dx := x - cx
		maxDy := int(math.Sqrt(float64(r2 - dx*dx)))
		row := grid[x]
		for y := cy - maxDy; y <= cy+maxDy; y++ {
			if y >= 0 && y < n {
				row[y] = 1
			}
		}
	}
}

func drawAxisSquare(grid [][]byte, x0, y0, side int) {
	for x := x0; x < x0+side; x++ {
		row := grid[x]
		for y := y0; y < y0+side; y++ {
			row[y] = 1
		}
	}
}

func drawRotatedSquare(grid [][]byte, cx, cy, side int, angle float64) {
	n := len(grid)
	half := float64(side) / 2
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	radius := int(math.Ceil(half*(math.Abs(cosA)+math.Abs(sinA)))) + 1
	for x := cx - radius; x <= cx+radius; x++ {
		if x < 0 || x >= n {
			continue
		}
		dx := float64(x - cx)
		row := grid[x]
		for y := cy - radius; y <= cy+radius; y++ {
			if y < 0 || y >= n {
				continue
			}
			dy := float64(y - cy)
			xr := dx*cosA + dy*sinA
			yr := -dx*sinA + dy*cosA
			if math.Abs(xr) <= half && math.Abs(yr) <= half {
				row[y] = 1
			}
		}
	}
}

func tryPlaceCircle(r *rand.Rand, grid [][]byte, boxes *[]box) bool {
	n := len(grid)
	minR := minShapeSize
	maxBound := n/2 - separation - 1
	if maxBound < minR {
		return false
	}
	maxR := n / 6
	if maxR < minR {
		maxR = minR
	}
	if maxR > maxBound {
		maxR = maxBound
	}
	for attempt := 0; attempt < 300; attempt++ {
		radius := randRange(r, minR, maxR)
		minC := radius + separation
		maxC := n - radius - separation - 1
		if minC > maxC {
			continue
		}
		cx := randRange(r, minC, maxC)
		cy := randRange(r, minC, maxC)
		candidate := box{
			minX: cx - radius - separation,
			maxX: cx + radius + separation,
			minY: cy - radius - separation,
			maxY: cy + radius + separation,
		}
		if !fits(candidate, *boxes, n) {
			continue
		}
		drawCircle(grid, cx, cy, radius)
		*boxes = append(*boxes, candidate)
		return true
	}
	return false
}

func tryPlaceSquare(r *rand.Rand, grid [][]byte, boxes *[]box) bool {
	n := len(grid)
	maxSide := n / 4
	if maxSide < minShapeSize {
		maxSide = minShapeSize
	}
	for attempt := 0; attempt < 400; attempt++ {
		side := randRange(r, minShapeSize, maxSide)
		if r.Intn(2) == 0 {
			minCoord := separation
			maxCoord := n - side - separation
			if minCoord > maxCoord {
				continue
			}
			x0 := randRange(r, minCoord, maxCoord)
			y0 := randRange(r, minCoord, maxCoord)
			candidate := box{
				minX: x0 - separation,
				maxX: x0 + side - 1 + separation,
				minY: y0 - separation,
				maxY: y0 + side - 1 + separation,
			}
			if !fits(candidate, *boxes, n) {
				continue
			}
			drawAxisSquare(grid, x0, y0, side)
			*boxes = append(*boxes, candidate)
			return true
		}
		angle := r.Float64() * (math.Pi / 2)
		half := float64(side) / 2
		span := half * (math.Abs(math.Cos(angle)) + math.Abs(math.Sin(angle)))
		dx := int(math.Ceil(span))
		dy := dx
		cxMin := dx + separation
		cxMax := n - dx - separation - 1
		cyMin := dy + separation
		cyMax := n - dy - separation - 1
		if cxMin > cxMax || cyMin > cyMax {
			continue
		}
		cx := randRange(r, cxMin, cxMax)
		cy := randRange(r, cyMin, cyMax)
		candidate := box{
			minX: cx - dx - separation,
			maxX: cx + dx + separation,
			minY: cy - dy - separation,
			maxY: cy + dy + separation,
		}
		if !fits(candidate, *boxes, n) {
			continue
		}
		drawRotatedSquare(grid, cx, cy, side, angle)
		*boxes = append(*boxes, candidate)
		return true
	}
	return false
}

func applyNoise(r *rand.Rand, grid [][]byte) {
	noise := r.Float64() * 0.2
	if noise == 0 {
		return
	}
	for i := range grid {
		row := grid[i]
		for j := range row {
			if r.Float64() < noise {
				row[j] ^= 1
			}
		}
	}
}

func formatGrid(grid [][]byte) string {
	n := len(grid)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		row := grid[i]
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			if row[j] == 1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func tryGenCase(r *rand.Rand) (string, bool) {
	n := r.Intn(151) + 200
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, n)
	}
	boxes := make([]box, 0, 10)
	target := r.Intn(4) + 3
	circleCount, squareCount := 0, 0
	attempts := 0
	for len(boxes) < target && attempts < target*400 {
		attempts++
		if r.Intn(2) == 0 {
			if tryPlaceCircle(r, grid, &boxes) {
				circleCount++
			}
		} else {
			if tryPlaceSquare(r, grid, &boxes) {
				squareCount++
			}
		}
	}
	if circleCount == 0 && !tryPlaceCircle(r, grid, &boxes) {
		return "", false
	}
	if squareCount == 0 && !tryPlaceSquare(r, grid, &boxes) {
		return "", false
	}
	applyNoise(r, grid)
	return formatGrid(grid), true
}

func genCase(r *rand.Rand) string {
	for attempt := 0; attempt < 20; attempt++ {
		if input, ok := tryGenCase(r); ok {
			return input
		}
	}
	panic("failed to generate test case")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
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
	const tests = 20
	for i := 0; i < tests; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
