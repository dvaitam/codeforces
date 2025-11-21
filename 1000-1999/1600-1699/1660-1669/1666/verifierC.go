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
	refSource      = "1000-1999/1600-1699/1660-1669/1666/1666C.go"
	maxSegments    = 100
	coordLimit     = 1_000_000_000
	requiredPoints = 3
)

type point struct {
	x, y int64
}

type segment struct {
	x1, y1, x2, y2 int64
}

type solution struct {
	segments []segment
}

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		pts, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		minLen := minimalLength(pts)

		refRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refRaw)
			os.Exit(1)
		}
		refSol, err := parseSolution(refRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refRaw)
			os.Exit(1)
		}
		if err := validateSolution(pts, refSol, minLen); err != nil {
			fmt.Fprintf(os.Stderr, "reference solution invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refRaw)
			os.Exit(1)
		}

		candRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candRaw)
			os.Exit(1)
		}
		candSol, err := parseSolution(candRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candRaw)
			os.Exit(1)
		}
		if err := validateSolution(pts, candSol, minLen); err != nil {
			fmt.Fprintf(os.Stderr, "candidate solution invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candRaw)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1666C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1666C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseInput(data string) ([]point, error) {
	reader := strings.NewReader(data)
	pts := make([]point, requiredPoints)
	for i := 0; i < requiredPoints; i++ {
		if _, err := fmt.Fscan(reader, &pts[i].x, &pts[i].y); err != nil {
			return nil, fmt.Errorf("failed to read point %d: %v", i+1, err)
		}
	}
	return pts, nil
}

func minimalLength(pts []point) int64 {
	xs := []int64{pts[0].x, pts[1].x, pts[2].x}
	ys := []int64{pts[0].y, pts[1].y, pts[2].y}
	sortInt64s(xs)
	sortInt64s(ys)
	minX, maxX := xs[0], xs[2]
	minY, maxY := ys[0], ys[2]
	medianX, medianY := xs[1], ys[1]

	var horiz int64 = maxX - minX
	for _, p := range pts {
		horiz += abs64(p.y - medianY)
	}

	var vert int64 = maxY - minY
	for _, p := range pts {
		vert += abs64(p.x - medianX)
	}

	if horiz < vert {
		return horiz
	}
	return vert
}

func sortInt64s(arr []int64) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func parseSolution(output string) (*solution, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	if fields[0] == "-1" {
		return nil, fmt.Errorf("problem always solvable; '-1' is invalid")
	}

	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid segment count %q", fields[0])
	}
	if n < 0 || n > maxSegments {
		return nil, fmt.Errorf("segment count %d outside [0,%d]", n, maxSegments)
	}
	if len(fields) != 1+n*4 {
		return nil, fmt.Errorf("expected %d integers, got %d", 1+n*4, len(fields))
	}
	segs := make([]segment, n)
	for i := 0; i < n; i++ {
		var coords [4]int64
		for j := 0; j < 4; j++ {
			val, err := strconv.ParseInt(fields[1+i*4+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[1+i*4+j])
			}
			if abs64(val) > coordLimit {
				return nil, fmt.Errorf("coordinate %d exceeds limit %d", val, coordLimit)
			}
			coords[j] = val
		}
		segs[i] = segment{coords[0], coords[1], coords[2], coords[3]}
		if !isAxisAligned(segs[i]) {
			return nil, fmt.Errorf("segment %d is not axis-aligned or has zero length", i+1)
		}
	}
	return &solution{segments: segs}, nil
}

func isAxisAligned(seg segment) bool {
	if seg.x1 == seg.x2 && seg.y1 != seg.y2 {
		return true
	}
	if seg.y1 == seg.y2 && seg.x1 != seg.x2 {
		return true
	}
	return false
}

func validateSolution(pts []point, sol *solution, minLen int64) error {
	if len(sol.segments) > maxSegments {
		return fmt.Errorf("uses %d segments, exceeds limit %d", len(sol.segments), maxSegments)
	}
	totalLen := int64(0)
	for _, seg := range sol.segments {
		totalLen += segmentLength(seg)
	}
	if totalLen != minLen {
		return fmt.Errorf("total length %d differs from optimal %d", totalLen, minLen)
	}
	pointCovered := make([]bool, len(pts))
	for pi, pt := range pts {
		for _, seg := range sol.segments {
			if pointOnSegment(pt, seg) {
				pointCovered[pi] = true
				break
			}
		}
		if !pointCovered[pi] {
			return fmt.Errorf("point %d (%d,%d) is not on any segment", pi+1, pt.x, pt.y)
		}
	}
	if !checkConnectivity(pts, sol.segments) {
		return fmt.Errorf("segments do not connect all points")
	}
	return nil
}

func segmentLength(seg segment) int64 {
	if seg.x1 == seg.x2 {
		return abs64(seg.y1 - seg.y2)
	}
	return abs64(seg.x1 - seg.x2)
}

func pointOnSegment(pt point, seg segment) bool {
	if seg.x1 == seg.x2 {
		if pt.x != seg.x1 {
			return false
		}
		minY, maxY := ordered(seg.y1, seg.y2)
		return pt.y >= minY && pt.y <= maxY
	}
	if pt.y != seg.y1 {
		return false
	}
	minX, maxX := ordered(seg.x1, seg.x2)
	return pt.x >= minX && pt.x <= maxX
}

func ordered(a, b int64) (int64, int64) {
	if a <= b {
		return a, b
	}
	return b, a
}

func checkConnectivity(pts []point, segs []segment) bool {
	if len(segs) == 0 {
		return false
	}
	totalNodes := len(segs) + len(pts)
	ds := newDSU(totalNodes)

	for i := range segs {
		for j := i + 1; j < len(segs); j++ {
			if segmentsIntersect(segs[i], segs[j]) {
				ds.union(i, j)
			}
		}
	}

	for pi, pt := range pts {
		attached := false
		for si, seg := range segs {
			if pointOnSegment(pt, seg) {
				ds.union(len(segs)+pi, si)
				attached = true
			}
		}
		if !attached {
			return false
		}
	}

	root := ds.find(len(segs))
	for pi := range pts {
		if ds.find(len(segs)+pi) != root {
			return false
		}
	}
	return true
}

func segmentsIntersect(a, b segment) bool {
	if a.x1 == a.x2 && b.x1 == b.x2 {
		if a.x1 != b.x1 {
			return false
		}
		minA, maxA := ordered(a.y1, a.y2)
		minB, maxB := ordered(b.y1, b.y2)
		return max(minA, minB) <= min(maxA, maxB)
	}
	if a.y1 == a.y2 && b.y1 == b.y2 {
		if a.y1 != b.y1 {
			return false
		}
		minA, maxA := ordered(a.x1, a.x2)
		minB, maxB := ordered(b.x1, b.x2)
		return max(minA, minB) <= min(maxA, maxB)
	}
	if a.x1 == a.x2 {
		return crossIntersect(a, b)
	}
	return crossIntersect(b, a)
}

func crossIntersect(vertical, horizontal segment) bool {
	if vertical.x1 != vertical.x2 || horizontal.y1 != horizontal.y2 {
		return false
	}
	minY, maxY := ordered(vertical.y1, vertical.y2)
	minX, maxX := ordered(horizontal.x1, horizontal.x2)
	return horizontal.y1 >= minY && horizontal.y1 <= maxY &&
		vertical.x1 >= minX && vertical.x1 <= maxX
}

type dsu struct {
	parent []int
	rank   []int
}

func newDSU(n int) *dsu {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &dsu{parent: parent, rank: rank}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++
	}
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample", []point{{1, 1}, {3, 5}, {8, 6}}),
		newTestCase("colinear_horizontal", []point{{-5, 0}, {0, 0}, {10, 0}}),
		newTestCase("colinear_vertical", []point{{3, -7}, {3, 0}, {3, 20}}),
		newTestCase("mixed", []point{{0, 0}, {0, 10}, {5, 5}}),
		newTestCase("negative_coords", []point{{-10, -10}, {-5, 2}, {7, -3}}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTestCase(rng, i))
	}
	return tests
}

func newTestCase(name string, pts []point) testCase {
	var sb strings.Builder
	for i, pt := range pts {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d %d", pt.x, pt.y))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	pts := make([]point, 0, requiredPoints)
	for len(pts) < requiredPoints {
		x := rng.Int63n(2*coordLimit+1) - coordLimit
		y := rng.Int63n(2*coordLimit+1) - coordLimit
		p := point{x, y}
		unique := true
		for _, existing := range pts {
			if existing == p {
				unique = false
				break
			}
		}
		if unique {
			pts = append(pts, p)
		}
	}
	name := fmt.Sprintf("random_%d", idx+1)
	return newTestCase(name, pts)
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
