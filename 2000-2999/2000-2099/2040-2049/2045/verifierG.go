package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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
	expect  []answer
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

	tests := buildTests()

	for idx, tc := range tests {
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

		if !equalAnswers(tc.expect, candAns) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed\ninput:\n%sexpected:\n%s\ncandidate:\n%s", idx+1, tc.name, tc.input, formatAnswers(tc.expect), candOut)
			return 1
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
	return 0
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

func formatAnswers(ans []answer) string {
	var sb strings.Builder
	for i, a := range ans {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if a.invalid {
			sb.WriteString("INVALID")
		} else {
			sb.WriteString(strconv.FormatInt(a.val, 10))
		}
	}
	if len(ans) > 0 {
		sb.WriteByte('\n')
	}
	return sb.String()
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
	gc, err := parseGridCase(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, queries: len(gc.queries), expect: solveCase(gc)}
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
	return testCase{name: name, input: input, queries: len(gc.queries), expect: solveCase(gc)}
}

func parseGridCase(input string) (gridCase, error) {
	reader := strings.NewReader(input)
	var gc gridCase
	if _, err := fmt.Fscan(reader, &gc.R, &gc.C, &gc.X); err != nil {
		return gridCase{}, fmt.Errorf("failed to read R C X: %v", err)
	}
	gc.cells = make([][]int, gc.R)
	for i := 0; i < gc.R; i++ {
		var row string
		if _, err := fmt.Fscan(reader, &row); err != nil {
			return gridCase{}, fmt.Errorf("failed to read row %d: %v", i+1, err)
		}
		if len(row) != gc.C {
			return gridCase{}, fmt.Errorf("row %d has length %d, expected %d", i+1, len(row), gc.C)
		}
		gc.cells[i] = make([]int, gc.C)
		for j := 0; j < gc.C; j++ {
			if row[j] < '0' || row[j] > '9' {
				return gridCase{}, fmt.Errorf("invalid digit in row %d", i+1)
			}
			gc.cells[i][j] = int(row[j] - '0')
		}
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return gridCase{}, fmt.Errorf("failed to read Q: %v", err)
	}
	gc.queries = make([]query, q)
	for i := 0; i < q; i++ {
		if _, err := fmt.Fscan(reader, &gc.queries[i].rs, &gc.queries[i].cs, &gc.queries[i].rf, &gc.queries[i].cf); err != nil {
			return gridCase{}, fmt.Errorf("failed to read query %d: %v", i+1, err)
		}
	}
	return gc, nil
}

func solveCase(gc gridCase) []answer {
	wDiff := precomputeDiff(gc.X)
	n := gc.R * gc.C
	type edge struct {
		u, v int
		w    int64
	}
	edges := make([]edge, 0, n*4)
	adj := make([][]int, n)
	for r := 0; r < gc.R; r++ {
		for c := 0; c < gc.C; c++ {
			u := r*gc.C + c
			hu := gc.cells[r][c]
			if r > 0 {
				v := (r-1)*gc.C + c
				w := wDiff[hu-gc.cells[r-1][c]+9]
				edges = append(edges, edge{u: u, v: v, w: w})
				adj[u] = append(adj[u], v)
			}
			if r+1 < gc.R {
				v := (r+1)*gc.C + c
				w := wDiff[hu-gc.cells[r+1][c]+9]
				edges = append(edges, edge{u: u, v: v, w: w})
				adj[u] = append(adj[u], v)
			}
			if c > 0 {
				v := r*gc.C + (c - 1)
				w := wDiff[hu-gc.cells[r][c-1]+9]
				edges = append(edges, edge{u: u, v: v, w: w})
				adj[u] = append(adj[u], v)
			}
			if c+1 < gc.C {
				v := r*gc.C + (c + 1)
				w := wDiff[hu-gc.cells[r][c+1]+9]
				edges = append(edges, edge{u: u, v: v, w: w})
				adj[u] = append(adj[u], v)
			}
		}
	}

	reachable := make([][]bool, n)
	for s := 0; s < n; s++ {
		reachable[s] = make([]bool, n)
		queue := []int{s}
		reachable[s][s] = true
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				if !reachable[s][v] {
					reachable[s][v] = true
					queue = append(queue, v)
				}
			}
		}
	}

	const inf int64 = 1 << 60
	res := make([]answer, len(gc.queries))
	for i, q := range gc.queries {
		s := (q.rs-1)*gc.C + (q.cs - 1)
		t := (q.rf-1)*gc.C + (q.cf - 1)

		dist := make([]int64, n)
		for j := range dist {
			dist[j] = inf
		}
		dist[s] = 0
		for it := 0; it < n-1; it++ {
			changed := false
			for _, e := range edges {
				if dist[e.u] == inf {
					continue
				}
				cand := dist[e.u] + e.w
				if cand < dist[e.v] {
					dist[e.v] = cand
					changed = true
				}
			}
			if !changed {
				break
			}
		}

		neg := make([]bool, n)
		queue := make([]int, 0, n)
		for _, e := range edges {
			if dist[e.u] == inf {
				continue
			}
			if dist[e.u]+e.w < dist[e.v] {
				if !neg[e.v] {
					neg[e.v] = true
					queue = append(queue, e.v)
				}
			}
		}
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				if !neg[v] {
					neg[v] = true
					queue = append(queue, v)
				}
			}
		}

		bad := false
		for v := 0; v < n; v++ {
			if neg[v] && reachable[v][t] {
				bad = true
				break
			}
		}
		if bad {
			res[i] = answer{invalid: true}
		} else {
			res[i] = answer{val: dist[t]}
		}
	}
	return res
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
