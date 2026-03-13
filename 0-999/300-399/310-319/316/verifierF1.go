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
	randomTrials = 20
	minRadius    = 20
	maxRadius    = 50
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
		fmt.Fprintln(os.Stderr, "usage: verifierF1 /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	// Build reference binary from REFERENCE_SOURCE_PATH
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		fatal("REFERENCE_SOURCE_PATH not set")
	}

	refBin := filepath.Join(os.TempDir(), "ref_316F1")
	cmd := exec.Command("go", "build", "-o", refBin, refSrc)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()

	for idx, input := range tests {
		expected, err := runBinary(refBin, input)
		if err != nil {
			fatal("reference failed on case %d: %v", idx+1, err)
		}
		got, err := runBinary(candidate, input)
		if err != nil {
			fatal("candidate failed on case %d: %v", idx+1, err)
		}
		if normalizeOutput(got) != normalizeOutput(expected) {
			fatal("case %d mismatch\nexpected: %s\n     got: %s", idx+1, expected, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func normalizeOutput(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func buildTests() []string {
	var tests []string

	// Blank case: no suns
	tests = append(tests, newImg(5, 5).toInput())

	// Build known scenes
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, buildKnownScene(rng))
	}
	return tests
}

func buildKnownScene(rng *rand.Rand) string {
	h := randRange(rng, 200, 400)
	w := randRange(rng, 200, 400)
	grid := newImg(h, w)

	target := randRange(rng, 1, 3)
	placed := 0
	radius := randRange(rng, minRadius, maxRadius)
	margin := radius + maxRayLen + 5

	for attempt := 0; attempt < 200 && placed < target; attempt++ {
		if margin >= h/2 || margin >= w/2 {
			break
		}
		cx := margin + rng.Intn(max(1, h-2*margin))
		cy := margin + rng.Intn(max(1, w-2*margin))
		numRays := rng.Intn(7)
		rays := makeRays(rng, numRays)
		if grid.placeSun(cx, cy, radius, rays) {
			placed++
		}
	}

	return grid.toInput()
}

func makeRays(rng *rand.Rand, count int) []raySpec {
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

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func runBinary(target, input string) (string, error) {
	cmd := exec.Command(target)
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

func randRange(rng *rand.Rand, lo, hi int) int {
	if hi < lo {
		lo, hi = hi, lo
	}
	if lo == hi {
		return lo
	}
	return lo + rng.Intn(hi-lo+1)
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

// suppress unused import
var _ = time.Now
