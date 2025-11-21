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

const refSource = "0-999/300-399/330-339/331/331D2.go"

type arrow struct {
	x0, y0 int
	x1, y1 int
}

type query struct {
	x, y int
	dir  byte
	t    int64
}

type testCase struct {
	name  string
	input string
	q     int
}

type point struct {
	x, y int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
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
		expectRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}
		expect, err := parseAnswers(expectRaw, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		got, err := parseAnswers(gotRaw, tc.q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, gotRaw)
			os.Exit(1)
		}
		if !equalPoints(expect, got) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: answers differ\ninput:\n%sreference:\n%s\ncandidate:\n%s\n",
				idx+1, tc.name, tc.input, formatPoints(expect), formatPoints(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-331D2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref331D2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
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

func parseAnswers(output string, q int) ([]point, error) {
	fields := strings.Fields(output)
	if len(fields) != q*2 {
		return nil, fmt.Errorf("expected %d integers, got %d", q*2, len(fields))
	}
	ans := make([]point, q)
	for i := 0; i < q; i++ {
		xVal, err := strconv.ParseInt(fields[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[2*i])
		}
		yVal, err := strconv.ParseInt(fields[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[2*i+1])
		}
		ans[i] = point{x: xVal, y: yVal}
	}
	return ans, nil
}

func equalPoints(a, b []point) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func formatPoints(arr []point) string {
	var sb strings.Builder
	for i, p := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("(%d,%d)", p.x, p.y))
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := []testCase{
		manualNoArrowsTest(),
		manualVerticalChainTest(),
		manualSeparatedDirectionsTest(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, i, 30, 25, 40))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestCase(rng, i+40, 1000, 200, 200))
	}
	for i := 0; i < 15; i++ {
		tests = append(tests, randomTestCase(rng, i+70, 100000, 400, 400))
	}
	tests = append(tests,
		buildStripTest("horizontal_strip_stress", true, 1500, 800),
		buildStripTest("vertical_strip_stress", false, 1500, 800),
		buildRandomLargeQTest(rng),
	)
	return tests
}

func manualNoArrowsTest() testCase {
	b := 5
	arrows := []arrow{}
	queries := []query{
		{0, 0, 'R', 3},
		{5, 5, 'L', 10},
		{2, 5, 'D', 2},
		{4, 1, 'U', 0},
	}
	return buildTestCase("no_arrows_basic", b, arrows, queries)
}

func manualVerticalChainTest() testCase {
	b := 30
	arrows := []arrow{
		{10, 2, 10, 8},
		{10, 11, 10, 7},
		{10, 15, 10, 22},
		{10, 25, 10, 18},
	}
	queries := []query{
		{0, 5, 'R', 40},
		{10, 3, 'U', 50},
		{10, 24, 'D', 7},
		{8, 26, 'R', 2},
	}
	return buildTestCase("vertical_chain_turns", b, arrows, queries)
}

func manualSeparatedDirectionsTest() testCase {
	b := 80
	arrows := []arrow{
		{5, 5, 60, 5},
		{60, 40, 60, 70},
		{5, 15, 60, 15},
		{60, 70, 60, 50},
	}
	queries := []query{
		{0, 5, 'R', 100},
		{60, 45, 'U', 10},
		{30, 0, 'U', 200},
		{60, 75, 'D', 100},
	}
	return buildTestCase("separated_regions", b, arrows, queries)
}

func randomTestCase(rng *rand.Rand, idx, bLimit, nLimit, qLimit int) testCase {
	b := rng.Intn(bLimit) + 1
	targetN := 0
	if nLimit > 0 {
		targetN = rng.Intn(nLimit + 1)
	}
	arrows := generateArrows(rng, b, targetN)
	q := rng.Intn(qLimit) + 1
	queries := randomQueries(rng, b, q)
	name := fmt.Sprintf("random_%d_b%d_n%d_q%d", idx+1, b, len(arrows), q)
	return buildTestCase(name, b, arrows, queries)
}

func buildRandomLargeQTest(rng *rand.Rand) testCase {
	b := 100000
	arrows := generateArrows(rng, b, 300)
	q := 1000
	queries := randomQueries(rng, b, q)
	return buildTestCase("large_coordinates_random", b, arrows, queries)
}

func buildStripTest(name string, horizontal bool, segments, q int) testCase {
	if segments <= 0 {
		return buildTestCase(name, 5, nil, randomQueries(rand.New(rand.NewSource(1)), 5, q))
	}
	b := segments*3 + 10
	if b > 100000 {
		b = 100000
	}
	arrows := make([]arrow, 0, segments)
	if horizontal {
		for i := 0; i < segments; i++ {
			y := i*3 + 1
			if y > b {
				break
			}
			x0, x1 := 1, b-1
			if i%2 == 0 {
				arrows = append(arrows, arrow{x0, y, x1, y})
			} else {
				arrows = append(arrows, arrow{x1, y, x0, y})
			}
		}
	} else {
		for i := 0; i < segments; i++ {
			x := i*3 + 1
			if x > b {
				break
			}
			y0, y1 := 1, b-1
			if i%2 == 0 {
				arrows = append(arrows, arrow{x, y0, x, y1})
			} else {
				arrows = append(arrows, arrow{x, y1, x, y0})
			}
		}
	}
	if len(arrows) == 0 {
		arrows = append(arrows, arrow{0, 0, b, 0})
	}
	queries := make([]query, q)
	for i := 0; i < q; i++ {
		x := (i*7 + 3) % (b + 1)
		y := (i*13 + 5) % (b + 1)
		dir := []byte{'L', 'R', 'U', 'D'}[i%4]
		t := int64(50000 + i*10)
		if i%5 == 0 {
			t += 1_000_000_000_000 // ensure large times are covered
		}
		queries[i] = query{x: x, y: y, dir: dir, t: t}
	}
	return buildTestCase(name, b, arrows, queries)
}

func buildTestCase(name string, b int, arrows []arrow, queries []query) testCase {
	var sb strings.Builder
	sb.Grow(len(arrows)*32 + len(queries)*32 + 64)
	sb.WriteString(fmt.Sprintf("%d %d\n", len(arrows), b))
	for _, ar := range arrows {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", ar.x0, ar.y0, ar.x1, ar.y1))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d %c %d\n", qu.x, qu.y, qu.dir, qu.t))
	}
	return testCase{name: name, input: sb.String(), q: len(queries)}
}

func randomQueries(rng *rand.Rand, b int, q int) []query {
	if q <= 0 {
		return nil
	}
	dirs := []byte{'L', 'R', 'U', 'D'}
	res := make([]query, q)
	for i := 0; i < q; i++ {
		x := 0
		y := 0
		if b > 0 {
			x = rng.Intn(b + 1)
			y = rng.Intn(b + 1)
		}
		dir := dirs[rng.Intn(len(dirs))]
		t := rng.Int63n(1_000_000)
		if rng.Intn(7) == 0 {
			t = rng.Int63n(1_000_000_000_000_001)
		}
		res[i] = query{x: x, y: y, dir: dir, t: t}
	}
	return res
}

func generateArrows(rng *rand.Rand, b int, target int) []arrow {
	if target <= 0 || b == 0 {
		return nil
	}
	arrows := make([]arrow, 0, target)
	maxAttempts := target*500 + 1000
	for len(arrows) < target && maxAttempts > 0 {
		maxAttempts--
		var ar arrow
		if rng.Intn(2) == 0 {
			y := rng.Intn(b + 1)
			x0 := rng.Intn(b + 1)
			x1 := rng.Intn(b + 1)
			if x0 == x1 {
				continue
			}
			ar = arrow{x0, y, x1, y}
		} else {
			x := rng.Intn(b + 1)
			y0 := rng.Intn(b + 1)
			y1 := rng.Intn(b + 1)
			if y0 == y1 {
				continue
			}
			ar = arrow{x, y0, x, y1}
		}
		if hasConflict(ar, arrows) {
			continue
		}
		arrows = append(arrows, ar)
	}
	return arrows
}

func hasConflict(candidate arrow, existing []arrow) bool {
	for _, ar := range existing {
		if segmentsConflict(candidate, ar) {
			return true
		}
	}
	return false
}

func segmentsConflict(a, b arrow) bool {
	aHoriz := a.y0 == a.y1
	bHoriz := b.y0 == b.y1
	if aHoriz && bHoriz {
		if a.y0 != b.y0 {
			return false
		}
		al, ar := ordered(a.x0, a.x1)
		bl, br := ordered(b.x0, b.x1)
		return max(al, bl) <= min(ar, br)
	}
	if !aHoriz && !bHoriz {
		if a.x0 != b.x0 {
			return false
		}
		al, ar := ordered(a.y0, a.y1)
		bl, br := ordered(b.y0, b.y1)
		return max(al, bl) <= min(ar, br)
	}
	var h, v arrow
	if aHoriz {
		h, v = a, b
	} else {
		h, v = b, a
	}
	hx0, hx1 := ordered(h.x0, h.x1)
	vy0, vy1 := ordered(v.y0, v.y1)
	return v.x0 >= hx0 && v.x0 <= hx1 && h.y0 >= vy0 && h.y0 <= vy1
}

func ordered(a, b int) (int, int) {
	if a <= b {
		return a, b
	}
	return b, a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
