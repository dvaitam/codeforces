package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceF1  = "0-999/300-399/310-319/316/316F1.go"
	randomTrials = 60
	minRadius    = 20
	maxRadius    = 80
	minRayLen    = 10
	maxRayLen    = 30
)

var rayDirections = [][2]int{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

type point struct {
	x int
	y int
}

type raySpec struct {
	dir    [2]int
	length int
}

type img struct {
	h, w int
	cell [][]byte
}

func newImg(h, w int) *img {
	grid := make([][]byte, h)
	for i := range grid {
		grid[i] = make([]byte, w)
	}
	return &img{h: h, w: w, cell: grid}
}

func (g *img) toInput() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", g.h, g.w)
	for i := 0; i < g.h; i++ {
		for j := 0; j < g.w; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(int(g.cell[i][j])))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceF1)
	if err != nil {
		fatal("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	tests = append(tests, buildScene(rand.New(rand.NewSource(11)), caseConfig{hMin: 1200, hMax: 1600, wMin: 1200, wMax: 1600, minSuns: 2, maxSuns: 5}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for extra := 0; extra < randomTrials; extra++ {
		var cfg caseConfig
		if rng.Intn(6) == 0 {
			cfg = caseConfig{hMin: 800, hMax: 1600, wMin: 800, wMax: 1600, minSuns: 1, maxSuns: 5}
		} else {
			cfg = caseConfig{hMin: 120, hMax: 500, wMin: 120, wMax: 500, minSuns: 0, maxSuns: 4}
		}
		tests = append(tests, buildScene(rng, cfg))
	}

	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fatal("reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fatal("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		if normalize(expect) != normalize(got) {
			fatal("case %d mismatch\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "316F1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func normalize(out string) string {
	return strings.Join(strings.Fields(out), " ")
}

type caseConfig struct {
	hMin, hMax int
	wMin, wMax int
	minSuns    int
	maxSuns    int
}

func deterministicCases() []string {
	cases := []string{blankCase(5, 5)}
	seeds := []int64{3, 7, 13}
	configs := []caseConfig{
		{hMin: 200, hMax: 300, wMin: 220, wMax: 320, minSuns: 1, maxSuns: 2},
		{hMin: 400, hMax: 500, wMin: 400, wMax: 500, minSuns: 2, maxSuns: 3},
		{hMin: 600, hMax: 800, wMin: 600, wMax: 800, minSuns: 3, maxSuns: 4},
	}
	for i, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		cases = append(cases, buildScene(rng, configs[i]))
	}
	return cases
}

func blankCase(h, w int) string {
	return newImg(h, w).toInput()
}

func buildScene(rng *rand.Rand, cfg caseConfig) string {
	h := randRange(rng, cfg.hMin, cfg.hMax)
	w := randRange(rng, cfg.wMin, cfg.wMax)
	grid := newImg(h, w)
	minSuns := max(0, cfg.minSuns)
	maxSuns := max(minSuns, cfg.maxSuns)
	target := minSuns
	if maxSuns > minSuns {
		target += rng.Intn(maxSuns - minSuns + 1)
	}
	placed := 0
	attempts := 0
	limit := max(200, target*200)
	for attempts < limit && placed < target {
		attempts++
		radius := randRange(rng, minRadius, maxRadius)
		margin := radius + maxRayLen + 5
		if margin >= grid.h || margin >= grid.w || grid.h <= 2*margin || grid.w <= 2*margin {
			continue
		}
		cx := margin + rng.Intn(h-2*margin)
		cy := margin + rng.Intn(w-2*margin)
		if grid.placeSun(cx, cy, radius, randomRays(rng)) {
			placed++
		}
	}
	if placed < target {
		placed += placeFallback(grid, target-placed)
	}
	return grid.toInput()
}

func placeFallback(g *img, need int) int {
	placed := 0
	radius := (minRadius + maxRadius) / 2
	margin := radius + maxRayLen + 5
	for cx := margin; cx < g.h-margin && placed < need; cx += radius + maxRayLen + 5 {
		for cy := margin; cy < g.w-margin && placed < need; cy += radius + maxRayLen + 5 {
			if g.placeSun(cx, cy, radius, nil) {
				placed++
			}
		}
	}
	return placed
}

func randRange(rng *rand.Rand, lo, hi int) int {
	if hi < lo {
		lo, hi = hi, lo
	}
	if lo == hi {
		return lo
	}
	return lo + rng.Intn(hi-lo+1)
}

func randomRays(rng *rand.Rand) []raySpec {
	var count int
	switch rng.Intn(5) {
	case 0:
		count = 0
	case 1:
		count = 1
	default:
		count = 2 + rng.Intn(5)
	}
	if count > len(rayDirections) {
		count = len(rayDirections)
	}
	dirs := make([][2]int, len(rayDirections))
	copy(dirs, rayDirections)
	rng.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})
	rays := make([]raySpec, count)
	for i := 0; i < count; i++ {
		rays[i] = raySpec{
			dir:    dirs[i],
			length: randRange(rng, minRayLen, maxRayLen),
		}
	}
	return rays
}

func (g *img) placeSun(cx, cy, radius int, rays []raySpec) bool {
	circle := collectCircleCells(cx, cy, radius)
	pending := make([]point, 0, len(circle)+len(rays)*maxRayLen*3)
	seen := make(map[point]struct{}, len(circle))

	addCells := func(cells []point) bool {
		for _, c := range cells {
			if c.x < 0 || c.x >= g.h || c.y < 0 || c.y >= g.w {
				return false
			}
			if _, ok := seen[c]; ok {
				continue
			}
			if g.cell[c.x][c.y] == 1 {
				return false
			}
			seen[c] = struct{}{}
			pending = append(pending, c)
		}
		return true
	}

	if !addCells(circle) {
		return false
	}
	for _, ray := range rays {
		cells, ok := collectRayCells(cx, cy, radius, ray)
		if !ok {
			return false
		}
		if !addCells(cells) {
			return false
		}
	}
	for _, c := range pending {
		g.cell[c.x][c.y] = 1
	}
	return true
}

func collectCircleCells(cx, cy, radius int) []point {
	var cells []point
	r2 := radius * radius
	for x := cx - radius - 1; x <= cx+radius+1; x++ {
		for y := cy - radius - 1; y <= cy+radius+1; y++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= r2 {
				cells = append(cells, point{x: x, y: y})
			}
		}
	}
	return cells
}

func collectRayCells(cx, cy, radius int, ray raySpec) ([]point, bool) {
	dx, dy := ray.dir[0], ray.dir[1]
	if dx == 0 && dy == 0 {
		return nil, false
	}
	startStep := radius + 1
	sx := cx + dx*startStep
	sy := cy + dy*startStep
	perpX := sign(-dy)
	perpY := sign(dx)
	offsets := []point{{0, 0}}
	if perpX != 0 || perpY != 0 {
		offsets = append(offsets, point{x: perpX, y: perpY})
		offsets = append(offsets, point{x: -perpX, y: -perpY})
	}
	var cells []point
	for step := 0; step < ray.length; step++ {
		px := sx + dx*step
		py := sy + dy*step
		for _, off := range offsets {
			cells = append(cells, point{x: px + off.x, y: py + off.y})
		}
	}
	return cells, true
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
