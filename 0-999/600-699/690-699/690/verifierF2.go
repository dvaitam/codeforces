package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type testCaseInput struct {
	name  string
	input string
	info  caseInfo
}

type caseInfo struct {
	n        int
	drawings [][][2]int
}

type canonizer struct {
	memo map[string]int
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(690))
	tests := generateTests(rng)

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refYes, _, err := parseSingleCaseOutput(refOut, tc.info.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candYes, candEdges, err := parseSingleCaseOutput(candOut, tc.info.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse solution output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if !refYes {
			if candYes {
				fmt.Fprintf(os.Stderr, "test %d (%s) expected NO but solution printed YES\ninput:\n%soutput:\n%s", idx+1, tc.name, tc.input, candOut)
				os.Exit(1)
			}
			continue
		}

		if !candYes {
			fmt.Fprintf(os.Stderr, "test %d (%s) expected YES but solution printed NO\ninput:\n%s", idx+1, tc.name, tc.input)
			os.Exit(1)
		}

		canon := newCanonizer()
		expected, err := computeExpectedCounts(canon, tc.info)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on test %d (%s): %v", idx+1, tc.name, err)
			os.Exit(1)
		}
		if err := verifyCandidateTree(tc.info.n, candEdges, canon, expected); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-690F2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "690F2.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSingleCaseOutput(output string, n int) (bool, [][2]int, error) {
	reader := bufio.NewReader(strings.NewReader(output))
	var verdict string
	if _, err := fmt.Fscan(reader, &verdict); err != nil {
		return false, nil, fmt.Errorf("unable to read verdict: %v", err)
	}
	verdict = strings.ToUpper(verdict)
	if verdict == "NO" {
		if err := ensureEOF(reader); err != nil {
			return false, nil, err
		}
		return false, nil, nil
	}
	if verdict != "YES" {
		return false, nil, fmt.Errorf("expected YES or NO, got %q", verdict)
	}
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return false, nil, fmt.Errorf("failed to read edge %d: %v", i+1, err)
		}
		edges[i] = [2]int{u, v}
	}
	if err := ensureEOF(reader); err != nil {
		return false, nil, err
	}
	return true, edges, nil
}

func ensureEOF(r *bufio.Reader) error {
	var extra string
	if _, err := fmt.Fscan(r, &extra); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("error while checking extra output: %v", err)
	}
	return fmt.Errorf("unexpected extra token %q", extra)
}

func generateTests(rng *rand.Rand) []testCaseInput {
	var tests []testCaseInput

	tests = append(tests, makeTestFromTree("n2_line", 2, lineTree(2), rng))
	tests = append(tests, makeTestFromTree("n3_path", 3, lineTree(3), rng))
	tests = append(tests, makeTestFromTree("star5", 5, starTree(5), rng))
	tests = append(tests, randomCase("random7", 7, rng))
	tests = append(tests, randomCase("random10", 10, rng))
	tests = append(tests, randomCase("random15", 15, rng))
	tests = append(tests, randomCase("random25", 25, rng))
	tests = append(tests, randomCase("random40", 40, rng))
	tests = append(tests, randomCase("random60", 60, rng))
	tests = append(tests, randomCase("random80", 80, rng))
	tests = append(tests, randomCase("random100", 100, rng))
	tests = append(tests, makeAllEmptyCase(5))

	return tests
}

func randomCase(name string, n int, rng *rand.Rand) testCaseInput {
	return makeTestFromTree(name, n, randomTree(n, rng), rng)
}

func makeTestFromTree(name string, n int, edges [][2]int, rng *rand.Rand) testCaseInput {
	drawings := make([][][2]int, n)
	for removed := 0; removed < n; removed++ {
		drawings[removed] = makeDrawing(n, edges, removed, rng)
	}
	rng.Shuffle(len(drawings), func(i, j int) {
		drawings[i], drawings[j] = drawings[j], drawings[i]
	})
	info := caseInfo{n: n, drawings: drawings}
	return testCaseInput{
		name:  name,
		info:  info,
		input: caseInfoToInput(info),
	}
}

func makeAllEmptyCase(n int) testCaseInput {
	drawings := make([][][2]int, n)
	info := caseInfo{n: n, drawings: drawings}
	return testCaseInput{
		name:  "all_empty",
		info:  info,
		input: caseInfoToInput(info),
	}
}

func caseInfoToInput(info caseInfo) string {
	var b strings.Builder
	fmt.Fprintf(&b, "1\n%d %d\n", info.n, info.n)
	for _, drawing := range info.drawings {
		fmt.Fprintf(&b, "%d\n", len(drawing))
		for _, e := range drawing {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
	}
	return b.String()
}

func makeDrawing(n int, edges [][2]int, removed int, rng *rand.Rand) [][2]int {
	perm := rng.Perm(n)
	mapping := make([]int, n)
	idx := 0
	for v := 0; v < n; v++ {
		if v == removed {
			continue
		}
		mapping[v] = perm[idx] + 1
		idx++
	}
	var drawing [][2]int
	for _, e := range edges {
		u, v := e[0], e[1]
		if u == removed || v == removed {
			continue
		}
		drawing = append(drawing, [2]int{mapping[u], mapping[v]})
	}
	return drawing
}

func lineTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i - 1})
	}
	return edges
}

func starTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{0, i})
	}
	return edges
}

func randomTree(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 1; v < n; v++ {
		parent := rng.Intn(v)
		edges = append(edges, [2]int{v, parent})
	}
	return edges
}

func newCanonizer() *canonizer {
	return &canonizer{memo: make(map[string]int)}
}

func (c *canonizer) mapSlice(vals []int) int {
	var sb strings.Builder
	sb.Grow(len(vals)*4 + 2)
	sb.WriteByte('[')
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte(']')
	key := sb.String()
	if id, ok := c.memo[key]; ok {
		return id
	}
	id := len(c.memo)
	c.memo[key] = id
	return id
}

func (c *canonizer) encodeRooted(adj [][]int, node, parent int) int {
	var children []int
	for _, nb := range adj[node] {
		if nb == parent {
			continue
		}
		children = append(children, c.encodeRooted(adj, nb, node))
	}
	sort.Ints(children)
	return c.mapSlice(children)
}

func (c *canonizer) canonizeTree(adj [][]int) int {
	n := len(adj)
	if n == 0 {
		return c.mapSlice(nil)
	}
	if n == 1 {
		return c.mapSlice(nil)
	}
	centers := findTreeCenters(adj)
	best := -1
	for _, center := range centers {
		label := c.encodeRooted(adj, center, -1)
		if best == -1 || label < best {
			best = label
		}
	}
	return best
}

func findTreeCenters(adj [][]int) []int {
	n := len(adj)
	if n <= 2 {
		centers := make([]int, n)
		for i := 0; i < n; i++ {
			centers[i] = i
		}
		return centers
	}
	degree := make([]int, n)
	var leaves []int
	for i := 0; i < n; i++ {
		degree[i] = len(adj[i])
		if degree[i] == 1 {
			leaves = append(leaves, i)
		}
	}
	remaining := n
	for remaining > 2 {
		remaining -= len(leaves)
		var next []int
		for _, leaf := range leaves {
			for _, nb := range adj[leaf] {
				degree[nb]--
				if degree[nb] == 1 {
					next = append(next, nb)
				}
			}
		}
		leaves = next
	}
	return leaves
}

func (c *canonizer) canonizeForest(totalNodes int, edges [][2]int) int {
	adj := make([][]int, totalNodes)
	for _, e := range edges {
		u, v := e[0], e[1]
		if u < 0 || u >= totalNodes || v < 0 || v >= totalNodes {
			panic("edge endpoint out of range")
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, totalNodes)
	var comps []int
	for i := 0; i < totalNodes; i++ {
		if visited[i] {
			continue
		}
		var nodes []int
		stack := []int{i}
		visited[i] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nodes = append(nodes, v)
			for _, nb := range adj[v] {
				if !visited[nb] {
					visited[nb] = true
					stack = append(stack, nb)
				}
			}
		}
		compAdj := make([][]int, len(nodes))
		idx := make(map[int]int, len(nodes))
		for j, v := range nodes {
			idx[v] = j
		}
		for _, v := range nodes {
			for _, nb := range adj[v] {
				compAdj[idx[v]] = append(compAdj[idx[v]], idx[nb])
			}
		}
		comps = append(comps, c.canonizeTree(compAdj))
	}
	sort.Ints(comps)
	return c.mapSlice(comps)
}

func normalizeDrawing(n int, edges [][2]int) ([][2]int, error) {
	total := n - 1
	labelMap := make(map[int]int)
	next := 0
	converted := make([][2]int, len(edges))
	for i, e := range edges {
		u, v := e[0], e[1]
		if u < 1 || u > n || v < 1 || v > n {
			return nil, fmt.Errorf("edge endpoints out of range")
		}
		if u == v {
			return nil, fmt.Errorf("edge with identical endpoints")
		}
		mu, ok := labelMap[u]
		if !ok {
			if next >= total {
				return nil, fmt.Errorf("drawing references more than %d vertices", total)
			}
			mu = next
			labelMap[u] = mu
			next++
		}
		mv, ok := labelMap[v]
		if !ok {
			if next >= total {
				return nil, fmt.Errorf("drawing references more than %d vertices", total)
			}
			mv = next
			labelMap[v] = mv
			next++
		}
		converted[i] = [2]int{mu, mv}
	}
	return converted, nil
}

func computeExpectedCounts(canon *canonizer, info caseInfo) (map[int]int, error) {
	expected := make(map[int]int)
	for _, drawing := range info.drawings {
		edges, err := normalizeDrawing(info.n, drawing)
		if err != nil {
			return nil, err
		}
		code := canon.canonizeForest(info.n-1, edges)
		expected[code]++
	}
	return expected, nil
}

func verifyCandidateTree(n int, edges [][2]int, canon *canonizer, expected map[int]int) error {
	if len(edges) != n-1 {
		return fmt.Errorf("expected %d edges, got %d", n-1, len(edges))
	}
	adj := make([][]int, n)
	seen := make(map[[2]int]struct{})
	for idx, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		if u < 0 || u >= n || v < 0 || v >= n {
			return fmt.Errorf("edge %d endpoints out of range", idx+1)
		}
		if u == v {
			return fmt.Errorf("edge %d is a self-loop", idx+1)
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate edge between %d and %d", a+1, b+1)
		}
		seen[key] = struct{}{}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	if len(seen) != n-1 {
		return fmt.Errorf("edges do not form a tree")
	}
	if !isConnected(adj) {
		return fmt.Errorf("graph is disconnected")
	}

	remaining := copyMap(expected)
	for rem := 0; rem < n; rem++ {
		fEdges := buildEdgesWithout(adj, rem)
		code := canon.canonizeForest(n-1, fEdges)
		if remaining[code] == 0 {
			return fmt.Errorf("forest after removing vertex %d does not match any drawing", rem+1)
		}
		remaining[code]--
	}
	for code, cnt := range remaining {
		if cnt != 0 {
			return fmt.Errorf("missing %d drawings of type %d", cnt, code)
		}
	}
	return nil
}

func buildEdgesWithout(adj [][]int, removed int) [][2]int {
	n := len(adj)
	mapping := make([]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		if i == removed {
			mapping[i] = -1
		} else {
			mapping[i] = idx
			idx++
		}
	}
	var edges [][2]int
	for u := 0; u < n; u++ {
		if u == removed {
			continue
		}
		for _, v := range adj[u] {
			if v == removed || u > v {
				continue
			}
			edges = append(edges, [2]int{mapping[u], mapping[v]})
		}
	}
	return edges
}

func isConnected(adj [][]int) bool {
	n := len(adj)
	visited := make([]bool, n)
	stack := []int{0}
	visited[0] = true
	count := 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		for _, nb := range adj[v] {
			if !visited[nb] {
				visited[nb] = true
				stack = append(stack, nb)
			}
		}
	}
	return count == n
}

func copyMap(in map[int]int) map[int]int {
	out := make(map[int]int, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
