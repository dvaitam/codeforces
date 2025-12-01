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

const refSource = "./324D1.go"

type segment struct {
	vertical     bool
	fixed        int64
	low, high    int64
	headCoord    int64
	headDir      byte
	fromX, fromY int64
	toX, toY     int64
}

type query struct {
	x, y int64
	dir  byte
	t    int64
}

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := compareOutputs(tc.input, refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-324D1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref324D1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func compareOutputs(input, refOut, candOut string) error {
	refLines := splitLines(refOut)
	candLines := splitLines(candOut)
	if len(refLines) != len(candLines) {
		return fmt.Errorf("output line count mismatch: expected %d got %d", len(refLines), len(candLines))
	}
	reader := strings.NewReader(input)
	var n int
	var b int64
	if _, err := fmt.Fscan(reader, &n, &b); err != nil {
		return fmt.Errorf("failed to parse input header: %v", err)
	}
	var segments []segment
	for i := 0; i < n; i++ {
		var x0, y0, x1, y1 int64
		fmt.Fscan(reader, &x0, &y0, &x1, &y1)
		segments = append(segments, buildSegment(x0, y0, x1, y1))
	}
	var q int
	fmt.Fscan(reader, &q)
	if len(refLines) != q {
		return fmt.Errorf("output lines %d do not match q=%d", len(refLines), q)
	}
	for i := 0; i < q; i++ {
		var x, y, t int64
		var dirStr string
		fmt.Fscan(reader, &x, &y, &dirStr, &t)
		expX, expY, err := parsePoint(refLines[i])
		if err != nil {
			return fmt.Errorf("reference line %d invalid: %v", i+1, err)
		}
		gotX, gotY, err := parsePoint(candLines[i])
		if err != nil {
			return fmt.Errorf("candidate line %d invalid: %v", i+1, err)
		}
		if err := validateTrajectory(b, segments, x, y, dirStr[0], t, gotX, gotY); err != nil {
			return fmt.Errorf("line %d invalid: %v", i+1, err)
		}
		if gotX != expX || gotY != expY {
			return fmt.Errorf("line %d mismatch: expected %d %d got %d %d", i+1, expX, expY, gotX, gotY)
		}
	}
	return nil
}

func buildSegment(x0, y0, x1, y1 int64) segment {
	if x0 == x1 {
		var low, high, head int64
		dir := byte('U')
		if y0 < y1 {
			low, high = y0, y1
			head = y1
			dir = 'U'
		} else {
			low, high = y1, y0
			head = y1
			dir = 'D'
		}
		return segment{
			vertical: true, fixed: x0,
			low: low, high: high,
			headCoord: head, headDir: dir,
			fromX: x0, fromY: y0, toX: x1, toY: y1,
		}
	}

	// horizontal
	var low, high, head int64
	dir := byte('R')
	if x0 < x1 {
		low, high = x0, x1
		head = x1
		dir = 'R'
	} else {
		low, high = x1, x0
		head = x1
		dir = 'L'
	}
	return segment{
		vertical: false, fixed: y0,
		low: low, high: high,
		headCoord: head, headDir: dir,
		fromX: x0, fromY: y0, toX: x1, toY: y1,
	}
}

func parsePoint(line string) (int64, int64, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("expected two integers, got %q", line)
	}
	x, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid x %q", fields[0])
	}
	y, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid y %q", fields[1])
	}
	return x, y, nil
}

func validateTrajectory(b int64, segs []segment, sx, sy int64, dir byte, t, gx, gy int64) error {
	x, y := sx, sy
	timeLeft := t
	for timeLeft > 0 {
		dx, dy := directionVector(dir)
		if dx == 0 && dy == 0 {
			return fmt.Errorf("unknown direction %c", dir)
		}
		distBound := boundaryDistance(b, x, y, dx, dy)
		bestDist := distBound + 1
		var hit *segment
		for i := range segs {
			s := &segs[i]
			d := travelToSegment(x, y, dx, dy, s)
			if d > 0 && d < bestDist {
				bestDist = d
				hit = s
			}
		}
		if hit == nil || bestDist > distBound {
			move := timeLeft
			if move > distBound {
				move = distBound
			}
			x += dx * move
			y += dy * move
			break
		}
		if timeLeft < bestDist {
			x += dx * timeLeft
			y += dy * timeLeft
			timeLeft = 0
			break
		}
		x += dx * bestDist
		y += dy * bestDist
		timeLeft -= bestDist
		var travel int64
		if hit.vertical {
			travel = abs64(hit.headCoord - y)
			if timeLeft < travel {
				y += sign64(hit.headCoord-y) * timeLeft
				timeLeft = 0
				break
			}
			y = hit.headCoord
		} else {
			travel = abs64(hit.headCoord - x)
			if timeLeft < travel {
				x += sign64(hit.headCoord-x) * timeLeft
				timeLeft = 0
				break
			}
			x = hit.headCoord
		}
		timeLeft -= travel
		dir = hit.headDir
	}
	if x != gx || y != gy {
		return fmt.Errorf("reached %d %d, expected %d %d via trajectory", x, y, gx, gy)
	}
	return nil
}

func directionVector(dir byte) (int64, int64) {
	switch dir {
	case 'R':
		return 1, 0
	case 'L':
		return -1, 0
	case 'U':
		return 0, 1
	case 'D':
		return 0, -1
	default:
		return 0, 0
	}
}

func boundaryDistance(b, x, y, dx, dy int64) int64 {
	if dx > 0 {
		return b - x
	}
	if dx < 0 {
		return x
	}
	if dy > 0 {
		return b - y
	}
	return y
}

func travelToSegment(x, y, dx, dy int64, s *segment) int64 {
	if dx != 0 && s.vertical {
		if y < s.low || y > s.high {
			return -1
		}
		dist := (s.fixed - x) * dx
		if dist <= 0 {
			return -1
		}
		return dist
	} else if dy != 0 && !s.vertical {
		if x < s.low || x > s.high {
			return -1
		}
		dist := (s.fixed - y) * dy
		if dist <= 0 {
			return -1
		}
		return dist
	}
	return -1
}

func splitLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) == 1 && len(lines[0]) == 0 {
		return []string{}
	}
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{name: "no_arrows", input: formatCase(0, 5, nil, []query{
		{0, 0, 'R', 3},
		{5, 5, 'L', 2},
	})})
	tests = append(tests, testCase{name: "single_vertical", input: formatCase(1, 5, []segment{
		buildSegment(2, 0, 2, 5),
	}, []query{
		{0, 3, 'R', 10},
		{5, 4, 'L', 10},
	})})
	tests = append(tests, testCase{name: "single_horizontal", input: formatCase(1, 5, []segment{
		buildSegment(0, 2, 5, 2),
	}, []query{
		{1, 0, 'U', 5},
		{4, 5, 'D', 5},
	})})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tc := randomCase(rng, i)
		tests = append(tests, tc)
	}
	return tests
}

func formatCase(n int, b int64, segs []segment, qs []query) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, b)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d %d %d\n", segs[i].fromX, segs[i].fromY, segs[i].toX, segs[i].toY)
	}
	fmt.Fprintf(&sb, "%d\n", len(qs))
	for _, q := range qs {
		fmt.Fprintf(&sb, "%d %d %c %d\n", q.x, q.y, q.dir, q.t)
	}
	return sb.String()
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(3)
	b := int64(rng.Intn(10) + 5)
	var segs []segment
	used := make(map[string]bool)
	for len(segs) < n {
		x0 := int64(rng.Intn(int(b + 1)))
		y0 := int64(rng.Intn(int(b + 1)))
		if rng.Intn(2) == 0 {
			x1 := x0
			y1 := int64(rng.Intn(int(b + 1)))
			if y1 == y0 {
				continue
			}
			key := fmt.Sprintf("v_%d_%d_%d", x0, y0, y1)
			if used[key] {
				continue
			}
			used[key] = true
			segs = append(segs, buildSegment(x0, y0, x1, y1))
		} else {
			y1 := y0
			x1 := int64(rng.Intn(int(b + 1)))
			if x1 == x0 {
				continue
			}
			key := fmt.Sprintf("h_%d_%d_%d", y0, x0, x1)
			if used[key] {
				continue
			}
			used[key] = true
			segs = append(segs, buildSegment(x0, y0, x1, y1))
		}
	}
	qc := rng.Intn(4) + 1
	var qs []query
	dirs := []byte{'U', 'D', 'L', 'R'}
	for i := 0; i < qc; i++ {
		x := int64(rng.Intn(int(b + 1)))
		y := int64(rng.Intn(int(b + 1)))
		dir := dirs[rng.Intn(len(dirs))]
		t := int64(rng.Intn(20) + 1)
		qs = append(qs, query{x, y, dir, t})
	}
	name := fmt.Sprintf("random_%d", idx+1)
	return testCase{name: name, input: formatCase(n, b, segs, qs)}
}

func splitInput(input string) (b int64, segs []segment, qs []query, err error) {
	reader := strings.NewReader(input)
	var n int
	if _, err = fmt.Fscan(reader, &n, &b); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var x0, y0, x1, y1 int64
		if _, err = fmt.Fscan(reader, &x0, &y0, &x1, &y1); err != nil {
			return
		}
		segs = append(segs, buildSegment(x0, y0, x1, y1))
	}
	var q int
	if _, err = fmt.Fscan(reader, &q); err != nil {
		return
	}
	for i := 0; i < q; i++ {
		var x, y, t int64
		var dir string
		if _, err = fmt.Fscan(reader, &x, &y, &dir, &t); err != nil {
			return
		}
		qs = append(qs, query{x: x, y: y, dir: dir[0], t: t})
	}
	return
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func sign64(x int64) int64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
