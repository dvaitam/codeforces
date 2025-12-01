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

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2045G.go"

type answer struct {
	invalid bool
	val     int64
}

type query struct {
	rs, cs, rf, cf int
}

type gridCase struct {
	R, C, X int
	cells   [][]int
	queries []query
}

type testCase struct {
	name    string
	input   string
	queries int
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		return 1
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			return 1
		}
		refAns, err := parseOutput(refOut, tc.queries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			return 1
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			return 1
		}
		candAns, err := parseOutput(candOut, tc.queries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			return 1
		}

		if !equalAnswers(refAns, candAns) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\ninput:\n%sreference:\n%s\ncandidate:\n%s", idx+1, tc.name, tc.input, refOut, candOut)
			return 1
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
	return 0
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2045G-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2045G.bin")
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
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, expected int) ([]answer, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]answer, expected)
	for i, token := range fields {
		if token == "INVALID" {
			res[i] = answer{invalid: true}
			continue
		}
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = answer{val: val}
	}
	return res, nil
}

func equalAnswers(a, b []answer) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].invalid != b[i].invalid {
			return false
		}
		if !a[i].invalid && a[i].val != b[i].val {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample1", "3 4 1\n3359\n4294\n3681\n5\n1 1 3 4\n3 3 2 1\n2 2 1 4\n1 3 3 2\n1 1 1 1\n"),
		newTestCase("sample2", "2 4 5\n1908\n2023\n2\n1 1 2 4\n1 1 1 1\n"),
		newTestCase("sample3", "3 3 9\n135\n357\n579\n2\n3 3 1 1\n2 2 2 2\n"),
	}

	tests = append(tests, buildSingleCellCase(), buildLineCase(), buildForcedInvalidCase())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, buildModerateCase(rng))
	for i := 0; i < 20; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func newTestCase(name, input string) testCase {
	q, err := extractQueryCount(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, queries: q}
}

func extractQueryCount(input string) (int, error) {
	reader := strings.NewReader(input)
	var R, C, X int
	if _, err := fmt.Fscan(reader, &R, &C, &X); err != nil {
		return 0, fmt.Errorf("failed to read R C X: %v", err)
	}
	for i := 0; i < R; i++ {
		var row string
		if _, err := fmt.Fscan(reader, &row); err != nil {
			return 0, fmt.Errorf("failed to read row %d: %v", i+1, err)
		}
		if len(row) != C {
			return 0, fmt.Errorf("row %d has length %d, expected %d", i+1, len(row), C)
		}
	}
	var Q int
	if _, err := fmt.Fscan(reader, &Q); err != nil {
		return 0, fmt.Errorf("failed to read Q: %v", err)
	}
	if Q <= 0 {
		return 0, fmt.Errorf("non-positive Q: %d", Q)
	}
	return Q, nil
}

func buildSingleCellCase() testCase {
	gc := gridCase{
		R: 1, C: 1, X: 9,
		cells: [][]int{{7}},
		queries: []query{
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
		},
	}
	return structuredTest("single_cell", gc)
}

func buildLineCase() testCase {
	cells := [][]int{{1, 9, 0, 8, 4}}
	gc := gridCase{
		R: 1, C: len(cells[0]), X: 5,
		cells: cells,
		queries: []query{
			{1, 1, 1, 5},
			{1, 5, 1, 1},
			{1, 3, 1, 3},
			{1, 2, 1, 4},
		},
	}
	return structuredTest("single_row", gc)
}

func buildForcedInvalidCase() testCase {
	x := 3
	square, ok := findNegativeSquare(x)
	if !ok {
		panic("failed to find a negative cycle for X=3")
	}
	gc := gridCase{
		R: 2,
		C: 2,
		X: x,
		cells: [][]int{
			{square[0], square[1]},
			{square[2], square[3]},
		},
		queries: []query{
			{1, 1, 2, 2},
			{2, 1, 1, 2},
			{1, 2, 1, 1},
		},
	}
	return structuredTest("forced_invalid", gc)
}

func buildModerateCase(rng *rand.Rand) testCase {
	r, c := 8, 9
	grid := make([][]int, r)
	for i := 0; i < r; i++ {
		grid[i] = make([]int, c)
		for j := 0; j < c; j++ {
			grid[i][j] = rng.Intn(10)
		}
	}
	qCount := 25
	queries := make([]query, qCount)
	for i := 0; i < qCount; i++ {
		queries[i] = query{
			rs: rng.Intn(r) + 1,
			cs: rng.Intn(c) + 1,
			rf: rng.Intn(r) + 1,
			cf: rng.Intn(c) + 1,
		}
	}
	gc := gridCase{
		R:       r,
		C:       c,
		X:       7,
		cells:   grid,
		queries: queries,
	}
	return structuredTest("moderate_grid", gc)
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	dims := rng.Intn(10) + 2
	r := rng.Intn(dims) + 1
	c := dims - r + 1
	if c <= 0 {
		c = 1
	}
	odd := []int{1, 3, 5, 7, 9}
	x := odd[rng.Intn(len(odd))]

	grid := make([][]int, r)
	for i := 0; i < r; i++ {
		grid[i] = make([]int, c)
		for j := 0; j < c; j++ {
			grid[i][j] = rng.Intn(10)
		}
	}

	qCount := rng.Intn(60) + 1
	queries := make([]query, qCount)
	for i := 0; i < qCount; i++ {
		queries[i] = query{
			rs: rng.Intn(r) + 1,
			cs: rng.Intn(c) + 1,
			rf: rng.Intn(r) + 1,
			cf: rng.Intn(c) + 1,
		}
	}

	gc := gridCase{
		R:       r,
		C:       c,
		X:       x,
		cells:   grid,
		queries: queries,
	}
	name := fmt.Sprintf("random_%d", idx)
	return structuredTest(name, gc)
}

func structuredTest(name string, gc gridCase) testCase {
	input := buildInput(gc)
	return testCase{name: name, input: input, queries: len(gc.queries)}
}

func buildInput(gc gridCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", gc.R, gc.C, gc.X))
	for i := 0; i < gc.R; i++ {
		for j := 0; j < gc.C; j++ {
			sb.WriteByte(byte('0' + gc.cells[i][j]))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(gc.queries)))
	for _, q := range gc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", q.rs, q.cs, q.rf, q.cf))
	}
	return sb.String()
}

func findNegativeSquare(x int) ([4]int, bool) {
	wDiff := precomputeDiff(x)
	for tl := 0; tl <= 9; tl++ {
		for tr := 0; tr <= 9; tr++ {
			for bl := 0; bl <= 9; bl++ {
				for br := 0; br <= 9; br++ {
					val := wDiff[tl-bl+9] + wDiff[bl-br+9] + wDiff[br-tr+9] + wDiff[tr-tl+9]
					if val < 0 {
						return [4]int{tl, tr, bl, br}, true
					}
				}
			}
		}
	}
	return [4]int{}, false
}

func precomputeDiff(x int) []int64 {
	wDiff := make([]int64, 19)
	for d := -9; d <= 9; d++ {
		val := int64(1)
		for i := 0; i < x; i++ {
			val *= int64(abs(d))
		}
		if d < 0 {
			val = -val
		}
		wDiff[d+9] = val
	}
	return wDiff
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
