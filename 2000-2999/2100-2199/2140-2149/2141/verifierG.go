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
	refSource     = "./2141G.go"
	maxTotalN     = 30000
	randomTests   = 120
	maxNPerRandom = 200
	coordLimit    = 200000
)

type point struct {
	x int64
	y int64
}

type testCase struct {
	pts []point
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference runtime error: %v\n%s", err, refOut)
	}
	expect, err := parseOutput(refOut, len(tests))
	if err != nil {
		fail("reference output invalid: %v\n%s", err, refOut)
	}

	candCmd := commandFor(candidate)
	candOut, err := runProgram(candCmd, input)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candOut)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fail("candidate output invalid: %v\n%s", err, candOut)
	}

	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch on test %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2141G-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2141G.bin")
	cmd := exec.Command("go", "build", "-o", binPath, filepath.Clean(refSource))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(raw string, t int) ([]int64, error) {
	fields := strings.Fields(raw)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", len(tc.pts))
		for _, p := range tc.pts {
			fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if totalN+len(tc.pts) > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += len(tc.pts)
	}

	// Minimal cases.
	add(testCase{pts: []point{{0, 0}}})
	add(testCase{pts: []point{{0, 0}, {1, 0}}})
	add(testCase{pts: []point{{0, 0}, {0, 1}}})

	// Small grids and lines.
	add(makeLine(0, 0, 5, true))  // horizontal length 5
	add(makeLine(0, 0, 6, false)) // vertical length 6
	add(makeGrid(0, 0, 2, 2))
	add(makeGrid(3, 3, 3, 2))
	add(makeGrid(-4, -4, 3, 3))

	// Sparse unique scattered points.
	add(testCase{pts: []point{{-1000, -1000}, {1000, 1000}, {0, 0}, {5, -3}, {-2, 7}}})

	// Consecutive x block with multiple y stacks to stress rectangle counting.
	add(makeBlock(0, 0, 5, 5))
	add(makeBlock(10, -3, 4, 6))

	for len(tests) < randomTests && totalN < maxTotalN {
		maxN := maxNPerRandom
		if rem := maxTotalN - totalN; rem < maxN {
			maxN = rem
		}
		if maxN < 1 {
			break
		}
		n := rng.Intn(maxN) + 1
		tc := randomCase(rng, n)
		add(tc)
	}

	return tests
}

func makeLine(x, y int64, length int, horizontal bool) testCase {
	pts := make([]point, length)
	for i := 0; i < length; i++ {
		if horizontal {
			pts[i] = point{x + int64(i), y}
		} else {
			pts[i] = point{x, y + int64(i)}
		}
	}
	return testCase{pts: pts}
}

func makeGrid(x0, y0 int64, w, h int) testCase {
	pts := make([]point, 0, w*h)
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < h; dy++ {
			pts = append(pts, point{x0 + int64(dx), y0 + int64(dy)})
		}
	}
	return testCase{pts: pts}
}

func makeBlock(x0, y0 int64, w, h int) testCase {
	// Consecutive x columns, each with consecutive y rows, but not full grid.
	pts := make([]point, 0, w*h)
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < h; dy++ {
			if (dx+dy)%2 == 0 {
				pts = append(pts, point{x0 + int64(dx), y0 + int64(dy)})
			}
		}
	}
	return testCase{pts: pts}
}

func randomCase(rng *rand.Rand, n int) testCase {
	pts := make([]point, 0, n)
	seen := make(map[int64]map[int64]struct{})
	for len(pts) < n {
		x := rng.Int63n(2*coordLimit+1) - coordLimit
		y := rng.Int63n(2*coordLimit+1) - coordLimit
		if row, ok := seen[x]; ok {
			if _, dup := row[y]; dup {
				continue
			}
			row[y] = struct{}{}
		} else {
			seen[x] = map[int64]struct{}{y: {}}
		}
		pts = append(pts, point{x: x, y: y})
	}
	return testCase{pts: pts}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
