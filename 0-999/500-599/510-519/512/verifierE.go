package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type triData struct {
	n    int
	orig [][2]int
	goal [][2]int
}

type testCase struct {
	input string
	data  triData
}

type diagStore struct {
	list [][2]int
	pos  map[[2]int]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if err := verifyOutput(tc.data, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%s\nReference output:\n%s\n", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if err := verifyOutput(tc.data, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%s\nCandidate output:\n%s\n", idx+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"512E.go",
		filepath.Join("0-999", "500-599", "510-519", "512", "512E.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 512E.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref512E_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func verifyOutput(data triData, out string) error {
	k, flips, err := parseOutput(out)
	if err != nil {
		return err
	}
	if len(flips) != k {
		return fmt.Errorf("declared %d steps but parsed %d", k, len(flips))
	}
	return simulate(data, flips)
}

func parseOutput(out string) (int, [][2]int, error) {
	reader := strings.NewReader(out)
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return 0, nil, fmt.Errorf("failed to read number of steps: %w", err)
	}
	if k < 0 {
		return 0, nil, fmt.Errorf("number of steps is negative")
	}
	if k > 20000 {
		return 0, nil, fmt.Errorf("number of steps exceeds 20000")
	}
	flips := make([][2]int, 0, k)
	for i := 0; i < k; i++ {
		var a, b int
		if _, err := fmt.Fscan(reader, &a, &b); err != nil {
			return 0, nil, fmt.Errorf("failed to read move %d: %w", i+1, err)
		}
		flips = append(flips, [2]int{a, b})
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return 0, nil, fmt.Errorf("extraneous output detected: %s", extra)
	}
	return k, flips, nil
}

func simulate(data triData, flips [][2]int) error {
	n := data.n
	adj := buildAdjacency(n, data.orig)
	for step, mv := range flips {
		a, b := mv[0], mv[1]
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("move %d: vertex out of range (%d,%d)", step+1, a, b)
		}
		if a == b {
			return fmt.Errorf("move %d: diagonal endpoints coincide", step+1)
		}
		if !adj[a][b] {
			return fmt.Errorf("move %d: diagonal (%d,%d) does not exist", step+1, a, b)
		}
		if isEdge(a, b, n) {
			return fmt.Errorf("move %d: cannot flip polygon edge (%d,%d)", step+1, a, b)
		}
		c, d, err := findOpposite(adj, n, a, b)
		if err != nil {
			return fmt.Errorf("move %d: %v", step+1, err)
		}
		c, d = normalizePair(c, d)
		if adj[c][d] {
			return fmt.Errorf("move %d: diagonal (%d,%d) already exists", step+1, c, d)
		}
		adj[a][b], adj[b][a] = false, false
		adj[c][d], adj[d][c] = true, true
	}
	goalAdj := buildAdjacency(n, data.goal)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if adj[i][j] != goalAdj[i][j] {
				if isEdge(i, j, n) {
					return fmt.Errorf("edge (%d,%d) mismatch", i, j)
				}
				return fmt.Errorf("final triangulation differs at diagonal (%d,%d)", i, j)
			}
		}
	}
	return nil
}

func buildAdjacency(n int, diagonals [][2]int) [][]bool {
	adj := make([][]bool, n+1)
	for i := range adj {
		adj[i] = make([]bool, n+1)
	}
	for i := 1; i <= n; i++ {
		j := i%n + 1
		adj[i][j], adj[j][i] = true, true
	}
	for _, d := range diagonals {
		a, b := d[0], d[1]
		adj[a][b], adj[b][a] = true, true
	}
	return adj
}

func isEdge(a, b, n int) bool {
	if a > b {
		a, b = b, a
	}
	if b-a == 1 {
		return true
	}
	return a == 1 && b == n
}

func findOpposite(adj [][]bool, n, a, b int) (int, int, error) {
	first, second := 0, 0
	for v := 1; v <= n; v++ {
		if v == a || v == b {
			continue
		}
		if adj[a][v] && adj[b][v] {
			if first == 0 {
				first = v
			} else if second == 0 {
				second = v
			} else {
				return 0, 0, fmt.Errorf("diagonal (%d,%d) belongs to more than two triangles", a, b)
			}
		}
	}
	if second == 0 {
		return 0, 0, fmt.Errorf("diagonal (%d,%d) is not shared by two triangles", a, b)
	}
	return first, second, nil
}

func normalizePair(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTest(4, [][2]int{{1, 3}}, [][2]int{{1, 3}}))
	tests = append(tests, manualTest(4, [][2]int{{1, 3}}, [][2]int{{2, 4}}))
	tests = append(tests, manualTest(5, [][2]int{{1, 3}, {1, 4}}, [][2]int{{1, 3}, {3, 5}}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		var n int
		switch {
		case len(tests) < 15:
			n = 4 + rng.Intn(5) // 4..8
		case len(tests) < 35:
			n = 9 + rng.Intn(40) // 9..48
		case len(tests) < 50:
			n = 50 + rng.Intn(150) // 50..199
		default:
			n = 200 + rng.Intn(801) // 200..1000
		}
		tests = append(tests, randomTestCase(n, rng))
	}
	return tests
}

func manualTest(n int, orig, goal [][2]int) testCase {
	return newTestCase(n, orig, goal)
}

func randomTestCase(n int, rng *rand.Rand) testCase {
	adjOrig, storeOrig := buildFan(n)
	randomFlips(adjOrig, storeOrig, n, rng.Intn(n*3+1), rng)
	origDiags := diagonalsFromAdj(adjOrig, n)

	adjGoal := cloneAdj(adjOrig)
	storeGoal := newDiagStore(diagonalsFromAdj(adjGoal, n))
	flipCount := rng.Intn(n*3 + 1)
	if flipCount == 0 && n > 4 {
		flipCount = 1
	}
	randomFlips(adjGoal, storeGoal, n, flipCount, rng)
	goalDiags := diagonalsFromAdj(adjGoal, n)
	if sameDiagonals(origDiags, goalDiags) && n > 4 {
		randomFlips(adjGoal, storeGoal, n, 1, rng)
		goalDiags = diagonalsFromAdj(adjGoal, n)
	}

	return newTestCase(n, origDiags, goalDiags)
}

func buildFan(n int) ([][]bool, *diagStore) {
	adj := buildAdjacency(n, nil)
	var diags [][2]int
	for v := 3; v <= n-1; v++ {
		adj[1][v], adj[v][1] = true, true
		diags = append(diags, [2]int{1, v})
	}
	return adj, newDiagStore(diags)
}

func randomFlips(adj [][]bool, store *diagStore, n, cnt int, rng *rand.Rand) {
	if cnt <= 0 {
		return
	}
	for t := 0; t < cnt; t++ {
		if len(store.list) == 0 {
			return
		}
		idx := rng.Intn(len(store.list))
		a := store.list[idx][0]
		b := store.list[idx][1]
		c, d, err := findOpposite(adj, n, a, b)
		if err != nil {
			continue
		}
		c, d = normalizePair(c, d)
		if store.contains(c, d) {
			continue
		}
		adj[a][b], adj[b][a] = false, false
		adj[c][d], adj[d][c] = true, true
		store.remove(a, b)
		store.add(c, d)
	}
}

func diagonalsFromAdj(adj [][]bool, n int) [][2]int {
	var diags [][2]int
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if isEdge(i, j, n) {
				continue
			}
			if adj[i][j] {
				diags = append(diags, [2]int{i, j})
			}
		}
	}
	return diags
}

func sameDiagonals(a, b [][2]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i][0] != b[i][0] || a[i][1] != b[i][1] {
			return false
		}
	}
	return true
}

func newTestCase(n int, orig, goal [][2]int) testCase {
	if len(orig) != n-3 || len(goal) != n-3 {
		panic("invalid triangulation size")
	}
	origCopy := append([][2]int(nil), orig...)
	goalCopy := append([][2]int(nil), goal...)
	sort.Slice(origCopy, func(i, j int) bool {
		if origCopy[i][0] == origCopy[j][0] {
			return origCopy[i][1] < origCopy[j][1]
		}
		return origCopy[i][0] < origCopy[j][0]
	})
	sort.Slice(goalCopy, func(i, j int) bool {
		if goalCopy[i][0] == goalCopy[j][0] {
			return goalCopy[i][1] < goalCopy[j][1]
		}
		return goalCopy[i][0] < goalCopy[j][0]
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, d := range origCopy {
		sb.WriteString(fmt.Sprintf("%d %d\n", d[0], d[1]))
	}
	for _, d := range goalCopy {
		sb.WriteString(fmt.Sprintf("%d %d\n", d[0], d[1]))
	}
	return testCase{
		input: sb.String(),
		data: triData{
			n:    n,
			orig: origCopy,
			goal: goalCopy,
		},
	}
}

func cloneAdj(adj [][]bool) [][]bool {
	n := len(adj) - 1
	cp := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		cp[i] = make([]bool, n+1)
		copy(cp[i], adj[i])
	}
	return cp
}

func newDiagStore(diags [][2]int) *diagStore {
	ds := &diagStore{
		list: make([][2]int, len(diags)),
		pos:  make(map[[2]int]int, len(diags)),
	}
	copy(ds.list, diags)
	for i, d := range ds.list {
		key := [2]int{d[0], d[1]}
		ds.pos[key] = i
	}
	return ds
}

func (ds *diagStore) remove(a, b int) {
	key := [2]int{a, b}
	idx, ok := ds.pos[key]
	if !ok {
		return
	}
	last := len(ds.list) - 1
	if idx != last {
		mv := ds.list[last]
		ds.list[idx] = mv
		ds.pos[[2]int{mv[0], mv[1]}] = idx
	}
	ds.list = ds.list[:last]
	delete(ds.pos, key)
}

func (ds *diagStore) add(a, b int) {
	key := [2]int{a, b}
	ds.pos[key] = len(ds.list)
	ds.list = append(ds.list, [2]int{a, b})
}

func (ds *diagStore) contains(a, b int) bool {
	_, ok := ds.pos[[2]int{a, b}]
	return ok
}
