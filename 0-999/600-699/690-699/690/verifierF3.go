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
	"time"
)

type testCase struct {
	id    string
	input string
}

type instance struct {
	n        int
	drawings [][][2]int // 1-indexed for output
}

type caseInfo struct {
	n        int
	drawings [][][2]int // zero-indexed edges
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF3.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF3.go /path/to/candidate")
		os.Exit(1)
	}

	baseDir := currentDir()
	refBin, err := buildReference(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	executed := 0
	for _, tc := range tests {
		cases, err := parseInputCases(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated input for %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		refVerdicts, err := parseReferenceOutputs(refOut, cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on %s: %v\noutput:\n%s", tc.id, err, refOut)
			os.Exit(1)
		}
		if err := verifyCandidateOutput(refOut, cases, refVerdicts); err != nil {
			fmt.Fprintf(os.Stderr, "skipping test %s due to reference self-check failure: %v\n", tc.id, err)
			continue
		}
		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if err := verifyCandidateOutput(candOut, cases, refVerdicts); err != nil {
			fmt.Fprintf(os.Stderr, "test %s failed: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		executed++
		if executed%5 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d tests...\n", executed)
		}
	}
	fmt.Printf("All %d tests passed.\n", executed)
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine current file path")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref690F3.bin")
	cmd := exec.Command("go", "build", "-o", out, "690F3.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseInputCases(input string) ([]caseInfo, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read number of test cases: %v", err)
	}
	cases := make([]caseInfo, t)
	for i := 0; i < t; i++ {
		var n, k int
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return nil, fmt.Errorf("failed to read n,k for case %d: %v", i+1, err)
		}
		drawings := make([][][2]int, k)
		for j := 0; j < k; j++ {
			var m int
			if _, err := fmt.Fscan(reader, &m); err != nil {
				return nil, fmt.Errorf("failed to read m for case %d drawing %d: %v", i+1, j+1, err)
			}
			drawings[j] = make([][2]int, m)
			for e := 0; e < m; e++ {
				var u, v int
				if _, err := fmt.Fscan(reader, &u, &v); err != nil {
					return nil, fmt.Errorf("failed to read edge %d for case %d drawing %d: %v", e+1, i+1, j+1, err)
				}
				drawings[j][e] = [2]int{u - 1, v - 1}
			}
		}
		cases[i] = caseInfo{n: n, drawings: drawings}
	}
	return cases, nil
}

func parseReferenceOutputs(output string, cases []caseInfo) ([]bool, error) {
	ts := newTokenStream(output)
	verdicts := make([]bool, len(cases))
	for i, cs := range cases {
		token, err := ts.nextToken()
		if err != nil {
			return nil, fmt.Errorf("missing verdict for case %d: %v", i+1, err)
		}
		switch strings.ToLower(token) {
		case "no":
			verdicts[i] = false
		case "yes":
			verdicts[i] = true
			required := cs.n - 1
			for e := 0; e < required; e++ {
				if _, err := ts.nextInt(); err != nil {
					return nil, fmt.Errorf("reference case %d missing edge endpoint: %v", i+1, err)
				}
				if _, err := ts.nextInt(); err != nil {
					return nil, fmt.Errorf("reference case %d missing edge endpoint: %v", i+1, err)
				}
			}
		default:
			return nil, fmt.Errorf("reference case %d invalid verdict token %q", i+1, token)
		}
	}
	if token, err := ts.nextToken(); err == nil {
		return nil, fmt.Errorf("reference output has extra token %q", token)
	}
	return verdicts, nil
}

func verifyCandidateOutput(output string, cases []caseInfo, refVerdicts []bool) error {
	ts := newTokenStream(output)
	for i, cs := range cases {
		token, err := ts.nextToken()
		if err != nil {
			return fmt.Errorf("missing verdict for case %d: %v", i+1, err)
		}
		switch strings.ToLower(token) {
		case "no":
			if refVerdicts[i] {
				return fmt.Errorf("case %d: expected YES but got NO", i+1)
			}
		case "yes":
			if !refVerdicts[i] {
				return fmt.Errorf("case %d: expected NO but got YES", i+1)
			}
			edges := make([][2]int, cs.n-1)
			for e := 0; e < cs.n-1; e++ {
				u, err := ts.nextInt()
				if err != nil {
					return fmt.Errorf("case %d: missing edge endpoint: %v", i+1, err)
				}
				v, err := ts.nextInt()
				if err != nil {
					return fmt.Errorf("case %d: missing edge endpoint: %v", i+1, err)
				}
				edges[e] = [2]int{u, v}
			}
			if err := validateTree(cs.n, edges); err != nil {
				return fmt.Errorf("case %d: invalid tree: %v", i+1, err)
			}
			if err := checkDrawings(cs, edges); err != nil {
				return fmt.Errorf("case %d: tree does not match drawings: %v", i+1, err)
			}
		default:
			return fmt.Errorf("case %d: unexpected verdict token %q", i+1, token)
		}
	}
	if token, err := ts.nextToken(); err == nil {
		return fmt.Errorf("candidate produced extra token %q", token)
	}
	return nil
}

func validateTree(n int, edges [][2]int) error {
	if len(edges) != n-1 {
		return fmt.Errorf("expected %d edges, got %d", n-1, len(edges))
	}
	adj := make([][]int, n)
	seen := make(map[[2]int]struct{})
	for _, e := range edges {
		u, v := e[0], e[1]
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge (%d,%d) uses vertex outside 1..%d", u, v, n)
		}
		if u == v {
			return fmt.Errorf("self-loop at vertex %d", u)
		}
		a, b := u-1, v-1
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate edge between %d and %d", u, v)
		}
		seen[key] = struct{}{}
		adj[u-1] = append(adj[u-1], v-1)
		adj[v-1] = append(adj[v-1], u-1)
	}
	visited := make([]bool, n)
	queue := []int{0}
	visited[0] = true
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, y := range adj[x] {
			if !visited[y] {
				visited[y] = true
				queue = append(queue, y)
			}
		}
	}
	for i, v := range visited {
		if !v {
			return fmt.Errorf("vertex %d is disconnected", i+1)
		}
	}
	return nil
}

func checkDrawings(cs caseInfo, edges [][2]int) error {
	treeZero := make([][2]int, len(edges))
	for i, e := range edges {
		treeZero[i] = [2]int{e[0] - 1, e[1] - 1}
	}
	can := newCanonizer()
	available := make(map[int]struct{})
	for v := 0; v < cs.n; v++ {
		forest := make([][2]int, 0, len(treeZero))
		for _, e := range treeZero {
			if e[0] != v && e[1] != v {
				forest = append(forest, e)
			}
		}
		h := can.canonizeForest(forest)
		available[h] = struct{}{}
	}

	for idx, drawing := range cs.drawings {
		h := can.canonizeForest(cloneEdges(drawing))
		if _, ok := available[h]; !ok {
			return fmt.Errorf("drawing %d cannot be reproduced by removing any vertex", idx+1)
		}
	}
	return nil
}

type canonizer struct {
	mp map[string]int
}

func newCanonizer() *canonizer {
	return &canonizer{mp: make(map[string]int)}
}

func (c *canonizer) mapT(ch []int) int {
	key := keyOfSlice(ch)
	if v, ok := c.mp[key]; ok {
		return v
	}
	id := len(c.mp)
	c.mp[key] = id
	return id
}

func keyOfSlice(ch []int) string {
	if len(ch) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(ch[0]))
	for i := 1; i < len(ch); i++ {
		sb.WriteByte(',')
		sb.WriteString(strconv.Itoa(ch[i]))
	}
	return sb.String()
}

func canonizeIndices(e [][2]int) [][2]int {
	if len(e) == 0 {
		return e
	}
	mapping := make(map[int]int, len(e)*2)
	next := 0
	for i := range e {
		u := e[i][0]
		v := e[i][1]
		if _, ok := mapping[u]; !ok {
			mapping[u] = next
			next++
		}
		if _, ok := mapping[v]; !ok {
			mapping[v] = next
			next++
		}
		e[i][0] = mapping[u]
		e[i][1] = mapping[v]
	}
	return e
}

func (c *canonizer) encodeSubtree(x, parent int, adj [][]int) int {
	children := make([]int, 0, len(adj[x]))
	for _, y := range adj[x] {
		if y != parent {
			children = append(children, c.encodeSubtree(y, x, adj))
		}
	}
	sort.Ints(children)
	return c.mapT(children)
}

func (c *canonizer) dfsSize(x, parent int, adj [][]int, sz, bal []int) {
	sz[x] = 1
	bal[x] = 0
	for _, y := range adj[x] {
		if y == parent {
			continue
		}
		c.dfsSize(y, x, adj, sz, bal)
		if sz[y] > bal[x] {
			bal[x] = sz[y]
		}
		sz[x] += sz[y]
	}
}

func (c *canonizer) canonizeTree(e [][2]int) int {
	if len(e) == 0 {
		return c.mapT(nil)
	}
	e = canonizeIndices(e)
	n := 0
	for _, p := range e {
		if p[0]+1 > n {
			n = p[0] + 1
		}
		if p[1]+1 > n {
			n = p[1] + 1
		}
	}
	if n == 0 {
		return c.mapT(nil)
	}
	adj := make([][]int, n)
	for _, p := range e {
		adj[p[0]] = append(adj[p[0]], p[1])
		adj[p[1]] = append(adj[p[1]], p[0])
	}
	sz := make([]int, n)
	bal := make([]int, n)
	c.dfsSize(0, -1, adj, sz, bal)
	best := -1
	for v := 0; v < n; v++ {
		rem := n - sz[v]
		cur := bal[v]
		if rem > cur {
			cur = rem
		}
		if 2*cur <= n {
			val := c.encodeSubtree(v, -1, adj)
			if best == -1 || val < best {
				best = val
			}
		}
	}
	if best == -1 {
		return c.encodeSubtree(0, -1, adj)
	}
	return best
}

func (c *canonizer) canonizeForest(e [][2]int) int {
	if len(e) == 0 {
		return c.mapT(nil)
	}
	e = canonizeIndices(e)
	n := 0
	for _, p := range e {
		if p[0]+1 > n {
			n = p[0] + 1
		}
		if p[1]+1 > n {
			n = p[1] + 1
		}
	}
	if n == 0 {
		return c.mapT(nil)
	}
	adj := make([][]int, n)
	for _, p := range e {
		adj[p[0]] = append(adj[p[0]], p[1])
		adj[p[1]] = append(adj[p[1]], p[0])
	}
	comp := make([]int, n)
	var hashes []int
	curID := 0
	for i := 0; i < n; i++ {
		if comp[i] != 0 {
			continue
		}
		curID++
		stack := []int{i}
		comp[i] = curID
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, w := range adj[v] {
				if comp[w] == 0 {
					comp[w] = curID
					stack = append(stack, w)
				}
			}
		}
		var subEdges [][2]int
		for _, p := range e {
			if comp[p[0]] == curID && comp[p[1]] == curID {
				subEdges = append(subEdges, p)
			}
		}
		hashes = append(hashes, c.canonizeTree(subEdges))
	}
	sort.Ints(hashes)
	return c.mapT(hashes)
}

func cloneEdges(e [][2]int) [][2]int {
	if len(e) == 0 {
		return nil
	}
	res := make([][2]int, len(e))
	copy(res, e)
	return res
}

type tokenStream struct {
	r *bufio.Reader
}

func newTokenStream(s string) *tokenStream {
	return &tokenStream{r: bufio.NewReader(strings.NewReader(s))}
}

func (ts *tokenStream) nextToken() (string, error) {
	var sb strings.Builder
	for {
		b, err := ts.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				if sb.Len() == 0 {
					return "", io.EOF
				}
				return sb.String(), nil
			}
			return "", err
		}
		if isSpace(b) {
			if sb.Len() == 0 {
				continue
			}
			return sb.String(), nil
		}
		sb.WriteByte(b)
		for {
			b, err = ts.r.ReadByte()
			if err != nil {
				if err == io.EOF {
					return sb.String(), nil
				}
				return "", err
			}
			if isSpace(b) {
				return sb.String(), nil
			}
			sb.WriteByte(b)
		}
	}
}

func (ts *tokenStream) nextInt() (int, error) {
	tok, err := ts.nextToken()
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer %q: %v", tok, err)
	}
	return val, nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}

func generateTests() []testCase {
	var tests []testCase
	detTree := [][2]int{{1, 2}, {2, 3}, {3, 4}}
	inst := instance{
		n: 4,
		drawings: [][][2]int{
			buildDrawing(detTree, 4, 1, nil),
			buildDrawing(detTree, 4, 3, nil),
		},
	}
	tests = append(tests, testCase{
		id:    "det-small-yes",
		input: formatInstances([]instance{inst}),
	})
	rngDet := rand.New(rand.NewSource(7))
	tests = append(tests, testCase{
		id:    "det-small-no",
		input: formatInstances([]instance{makeInvalidInstance(rngDet, 4)}),
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 25; i++ {
		caseCount := rng.Intn(3) + 1
		insts := make([]instance, caseCount)
		for j := 0; j < caseCount; j++ {
			n := rng.Intn(30) + 2
			if rng.Float64() < 0.7 {
				insts[j] = makeValidInstance(rng, n)
			} else {
				insts[j] = makeInvalidInstance(rng, n)
			}
		}
		tests = append(tests, testCase{
			id:    fmt.Sprintf("rand-%02d", i+1),
			input: formatInstances(insts),
		})
	}

	largeSizes := []int{200, 500, 1000}
	for i, sz := range largeSizes {
		tests = append(tests, testCase{
			id:    fmt.Sprintf("large-yes-%d", i+1),
			input: formatInstances([]instance{makeValidInstance(rng, sz)}),
		})
	}
	tests = append(tests, testCase{
		id:    "large-no-1",
		input: formatInstances([]instance{makeInvalidInstance(rng, 200)}),
	})
	return tests
}

func formatInstances(instances []instance) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(instances))
	for _, inst := range instances {
		fmt.Fprintf(&sb, "%d %d\n", inst.n, 2)
		for _, drawing := range inst.drawings {
			fmt.Fprintf(&sb, "%d\n", len(drawing))
			for _, e := range drawing {
				fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
			}
		}
	}
	return sb.String()
}

func makeValidInstance(rng *rand.Rand, n int) instance {
	tree := randomTree(n, rng)
	r1 := rng.Intn(n) + 1
	r2 := rng.Intn(n-1) + 1
	if r2 >= r1 {
		r2++
	}
	return instance{
		n: n,
		drawings: [][][2]int{
			buildDrawing(tree, n, r1, rng),
			buildDrawing(tree, n, r2, rng),
		},
	}
}

func makeInvalidInstance(rng *rand.Rand, n int) instance {
	inst := makeValidInstance(rng, n)
	idx := rng.Intn(len(inst.drawings))
	inst.drawings[idx] = corruptDrawing(inst.drawings[idx], n, rng)
	return inst
}

func randomTree(n int, rng *rand.Rand) [][2]int {
	if n == 1 {
		return nil
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func buildDrawing(tree [][2]int, n, removed int, rng *rand.Rand) [][2]int {
	var filtered [][2]int
	for _, e := range tree {
		if e[0] == removed || e[1] == removed {
			continue
		}
		filtered = append(filtered, e)
	}
	nodes := make([]int, 0, n-1)
	for v := 1; v <= n; v++ {
		if v != removed {
			nodes = append(nodes, v)
		}
	}
	if rng != nil {
		rng.Shuffle(len(nodes), func(i, j int) {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		})
	}
	mapping := make(map[int]int, len(nodes))
	for idx, val := range nodes {
		mapping[val] = idx + 1
	}
	res := make([][2]int, len(filtered))
	for i, e := range filtered {
		res[i] = [2]int{mapping[e[0]], mapping[e[1]]}
	}
	return res
}

func corruptDrawing(edges [][2]int, n int, rng *rand.Rand) [][2]int {
	if len(edges) == 0 {
		node := 1
		if n > 1 {
			node = rng.Intn(n-1) + 1
		}
		return append(edges, [2]int{node, node})
	}
	res := make([][2]int, len(edges)+1)
	copy(res, edges)
	res[len(edges)] = edges[rng.Intn(len(edges))]
	return res
}
