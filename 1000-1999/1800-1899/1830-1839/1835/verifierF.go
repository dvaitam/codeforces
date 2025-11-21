package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	from int
	to   int
}

type bitset []uint64

func newBitset(size int) bitset {
	return make([]uint64, (size+63)>>6)
}

func (b bitset) set(pos int) {
	b[pos>>6] |= 1 << (pos & 63)
}

func (b bitset) has(pos int) bool {
	return (b[pos>>6]>>(pos&63))&1 == 1
}

func (b bitset) or(other bitset) {
	for i := range b {
		b[i] |= other[i]
	}
}

func (b bitset) copy() bitset {
	res := make(bitset, len(b))
	copy(res, b)
	return res
}

func (b bitset) equal(other bitset) bool {
	if len(b) != len(other) {
		return false
	}
	for i := range b {
		if b[i] != other[i] {
			return false
		}
	}
	return true
}

type hopcroftKarp struct {
	nLeft, nRight int
	adj           [][]int
	dist          []int
	pairU         []int
	pairV         []int
}

func newHopcroftKarp(n, m int) *hopcroftKarp {
	h := &hopcroftKarp{
		nLeft:  n,
		nRight: m,
		adj:    make([][]int, n),
		dist:   make([]int, n),
		pairU:  make([]int, n),
		pairV:  make([]int, m),
	}
	for i := 0; i < n; i++ {
		h.pairU[i] = -1
	}
	for i := 0; i < m; i++ {
		h.pairV[i] = -1
	}
	return h
}

func (h *hopcroftKarp) addEdge(u, v int) {
	h.adj[u] = append(h.adj[u], v)
}

func (h *hopcroftKarp) bfs() bool {
	queue := make([]int, 0)
	for u := 0; u < h.nLeft; u++ {
		if h.pairU[u] == -1 {
			h.dist[u] = 0
			queue = append(queue, u)
		} else {
			h.dist[u] = -1
		}
	}
	found := false
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range h.adj[u] {
			matchedU := h.pairV[v]
			if matchedU == -1 {
				found = true
				continue
			}
			if h.dist[matchedU] == -1 {
				h.dist[matchedU] = h.dist[u] + 1
				queue = append(queue, matchedU)
			}
		}
	}
	return found
}

func (h *hopcroftKarp) dfs(u int) bool {
	for _, v := range h.adj[u] {
		matchedU := h.pairV[v]
		if matchedU == -1 || (h.dist[matchedU] == h.dist[u]+1 && h.dfs(matchedU)) {
			h.pairU[u] = v
			h.pairV[v] = u
			return true
		}
	}
	h.dist[u] = -1
	return false
}

func (h *hopcroftKarp) maxMatching() int {
	matching := 0
	for h.bfs() {
		for u := 0; u < h.nLeft; u++ {
			if h.pairU[u] == -1 && h.dfs(u) {
				matching++
			}
		}
	}
	return matching
}

func sccTarjan(n int, adj [][]int) (compID []int, comps [][]int) {
	compID = make([]int, n)
	for i := range compID {
		compID[i] = -1
	}
	low := make([]int, n)
	dfn := make([]int, n)
	timeStamp := 0
	stack := make([]int, 0)
	inStack := make([]bool, n)

	var dfs func(int)
	dfs = func(v int) {
		timeStamp++
		dfn[v] = timeStamp
		low[v] = timeStamp
		stack = append(stack, v)
		inStack[v] = true
		for _, to := range adj[v] {
			if dfn[to] == 0 {
				dfs(to)
				if low[to] < low[v] {
					low[v] = low[to]
				}
			} else if inStack[to] && dfn[to] < low[v] {
				low[v] = dfn[to]
			}
		}
		if low[v] == dfn[v] {
			var comp []int
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[w] = false
				compID[w] = len(comps)
				comp = append(comp, w)
				if w == v {
					break
				}
			}
			comps = append(comps, comp)
		}
	}

	for i := 0; i < n; i++ {
		if dfn[i] == 0 {
			dfs(i)
		}
	}
	return
}

func topoSort(adj [][]int) []int {
	n := len(adj)
	indeg := make([]int, n)
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			indeg[v]++
		}
	}
	queue := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	order := make([]int, 0, n)
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		order = append(order, u)
		for _, v := range adj[u] {
			indeg[v]--
			if indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return order
}

func buildImplicationGraph(n int, edges []edge, matchL, matchR []int) [][]int {
	g := make([][]int, n)
	seen := make(map[[2]int]struct{})
	for _, e := range edges {
		if matchL[e.from] == e.to {
			continue
		}
		to := matchR[e.to]
		key := [2]int{e.from, to}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		g[e.from] = append(g[e.from], to)
	}
	return g
}

func computeReachability(n int, edges []edge, matchL, matchR []int) ([]bitset, error) {
	if len(matchL) != n || len(matchR) != n {
		return nil, fmt.Errorf("invalid matching sizes")
	}
	imp := buildImplicationGraph(n, edges, matchL, matchR)
	compID, comps := sccTarjan(n, imp)
	compCnt := len(comps)
	compAdj := make([][]int, compCnt)
	seen := make(map[[2]int]struct{})
	for u := 0; u < n; u++ {
		for _, v := range imp[u] {
			cu, cv := compID[u], compID[v]
			if cu == cv {
				continue
			}
			key := [2]int{cu, cv}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			compAdj[cu] = append(compAdj[cu], cv)
		}
	}
	topo := topoSort(compAdj)
	compReach := make([]bitset, compCnt)
	for i := 0; i < compCnt; i++ {
		compReach[i] = newBitset(compCnt)
	}
	for i := len(topo) - 1; i >= 0; i-- {
		u := topo[i]
		for _, v := range compAdj[u] {
			compReach[u].or(compReach[v])
		}
		for _, v := range compAdj[u] {
			compReach[u].set(v)
		}
		compReach[u].set(u)
	}
	compBits := make([]bitset, compCnt)
	for idx, comp := range comps {
		bs := newBitset(n)
		for _, v := range comp {
			bs.set(v)
		}
		compBits[idx] = bs
	}
	compReachVertices := make([]bitset, compCnt)
	for i := 0; i < compCnt; i++ {
		bs := newBitset(n)
		for c := 0; c < compCnt; c++ {
			if compReach[i].has(c) {
				bs.or(compBits[c])
			}
		}
		compReachVertices[i] = bs
	}
	vertexReach := make([]bitset, n)
	for v := 0; v < n; v++ {
		vertexReach[v] = compReachVertices[compID[v]]
	}
	return vertexReach, nil
}

func minimizeCondensationEdges(compAdj [][]int) [][2]int {
	c := len(compAdj)
	topo := topoSort(compAdj)
	reach := make([]bitset, c)
	for i := 0; i < c; i++ {
		reach[i] = newBitset(c)
	}
	var edges [][2]int
	for i := len(topo) - 1; i >= 0; i-- {
		u := topo[i]
		for _, v := range compAdj[u] {
			reach[u].or(reach[v])
		}
		for _, v := range compAdj[u] {
			if !reach[u].has(v) {
				edges = append(edges, [2]int{u, v})
			}
		}
		for _, v := range compAdj[u] {
			reach[u].set(v)
		}
	}
	return edges
}

func minimalEdgeCount(n int, edges []edge, matchL, matchR []int) int {
	imp := buildImplicationGraph(n, edges, matchL, matchR)
	compID, comps := sccTarjan(n, imp)
	compCnt := len(comps)
	compAdj := make([][]int, compCnt)
	seen := make(map[[2]int]struct{})
	for u := 0; u < n; u++ {
		for _, v := range imp[u] {
			cu, cv := compID[u], compID[v]
			if cu == cv {
				continue
			}
			key := [2]int{cu, cv}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			compAdj[cu] = append(compAdj[cu], cv)
		}
	}

	edgeSet := make(map[[2]int]struct{})
	for u := 0; u < n; u++ {
		edgeSet[[2]int{u, matchL[u]}] = struct{}{}
	}
	for _, comp := range comps {
		if len(comp) <= 1 {
			continue
		}
		for i := 0; i < len(comp); i++ {
			u := comp[i]
			v := comp[(i+1)%len(comp)]
			edgeSet[[2]int{u, matchL[v]}] = struct{}{}
		}
	}
	minEdges := minimizeCondensationEdges(compAdj)
	rep := make([]int, compCnt)
	for idx, comp := range comps {
		rep[idx] = comp[0]
	}
	for _, e := range minEdges {
		fromComp, toComp := e[0], e[1]
		u := rep[fromComp]
		v := rep[toComp]
		edgeSet[[2]int{u, matchL[v]}] = struct{}{}
	}
	return len(edgeSet)
}

type testCase struct {
	n            int
	edges        []edge
	input        string
	good         bool
	matchL       []int
	matchR       []int
	reach        []bitset
	minEdgeCount int
	neighborBits []bitset
}

func buildNeighbors(n int, edges []edge) []bitset {
	neighbors := make([]bitset, n)
	for i := 0; i < n; i++ {
		neighbors[i] = newBitset(n)
	}
	for _, e := range edges {
		neighbors[e.from].set(e.to)
	}
	return neighbors
}

func makeTestCase(rng *rand.Rand, forceGood bool) testCase {
	for {
		n := rng.Intn(7) + 2
		maxM := n * n
		m := rng.Intn(maxM + 1)
		seen := make(map[[2]int]struct{})
		var edges []edge
		if forceGood {
			for i := 0; i < n; i++ {
				key := [2]int{i, i}
				seen[key] = struct{}{}
				edges = append(edges, edge{i, i})
			}
		}
		for len(edges) < m {
			u := rng.Intn(n)
			v := rng.Intn(n)
			key := [2]int{u, v}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			edges = append(edges, edge{u, v})
		}
		hk := newHopcroftKarp(n, n)
		for _, e := range edges {
			hk.addEdge(e.from, e.to)
		}
		size := hk.maxMatching()
		good := size == n
		if forceGood && !good {
			continue
		}
		if !forceGood && good {
			// Keep some bad cases; if too well connected regenerate.
			if rng.Float64() < 0.5 {
				continue
			}
		}
		var matchL, matchR []int
		var reach []bitset
		minEdges := 0
		if good {
			matchL = append([]int{}, hk.pairU...)
			matchR = append([]int{}, hk.pairV...)
			var err error
			reach, err = computeReachability(n, edges, matchL, matchR)
			if err != nil {
				continue
			}
			minEdges = minimalEdgeCount(n, edges, matchL, matchR)
		}
		neighborBits := buildNeighbors(n, edges)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e.from+1, e.to+n+1)
		}
		return testCase{
			n:            n,
			edges:        edges,
			input:        sb.String(),
			good:         good,
			matchL:       matchL,
			matchR:       matchR,
			reach:        reach,
			minEdgeCount: minEdges,
			neighborBits: neighborBits,
		}
	}
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseInt(token string) (int, error) {
	return strconv.Atoi(token)
}

func formatEdge(e edge, n int) string {
	return fmt.Sprintf("%d %d", e.from+1, e.to+n+1)
}

func verifyNO(tc testCase, tokens []string) error {
	if len(tokens) < 2 {
		return fmt.Errorf("expected k after NO")
	}
	k, err := parseInt(tokens[1])
	if err != nil || k < 0 {
		return fmt.Errorf("invalid k after NO")
	}
	if len(tokens) < 2+k {
		return fmt.Errorf("not enough vertices for NO answer")
	}
	subset := make(map[int]struct{})
	for i := 0; i < k; i++ {
		v, err := parseInt(tokens[2+i])
		if err != nil || v < 1 || v > tc.n {
			return fmt.Errorf("invalid vertex %q in NO answer", tokens[2+i])
		}
		subset[v-1] = struct{}{}
	}
	if len(subset) == 0 {
		return fmt.Errorf("subset must be non-empty to violate Hall")
	}
	neighbor := newBitset(tc.n)
	for v := range subset {
		neighbor.or(tc.neighborBits[v])
	}
	neighborCnt := 0
	for _, w := range neighbor {
		neighborCnt += bits.OnesCount64(w)
	}
	if len(subset) <= neighborCnt {
		return fmt.Errorf("provided subset is not a counterexample (|S|=%d, |N(S)|=%d)", len(subset), neighborCnt)
	}
	return nil
}

func verifyYES(tc testCase, tokens []string) error {
	if len(tokens) < 2 {
		return fmt.Errorf("expected m' after YES")
	}
	m, err := parseInt(tokens[1])
	if err != nil || m < 0 {
		return fmt.Errorf("invalid m' value")
	}
	if len(tokens) < 2+2*m {
		return fmt.Errorf("not enough edge tokens for YES answer")
	}
	edgeSet := make(map[[2]int]struct{})
	var edges []edge
	for i := 0; i < m; i++ {
		lTok := tokens[2+2*i]
		rTok := tokens[3+2*i]
		l, err1 := parseInt(lTok)
		r, err2 := parseInt(rTok)
		if err1 != nil || err2 != nil || l < 1 || l > tc.n || r < tc.n+1 || r > 2*tc.n {
			return fmt.Errorf("invalid edge (%s,%s)", lTok, rTok)
		}
		e := edge{from: l - 1, to: r - tc.n - 1}
		key := [2]int{e.from, e.to}
		if _, ok := edgeSet[key]; ok {
			return fmt.Errorf("duplicate edge %s", formatEdge(e, tc.n))
		}
		edgeSet[key] = struct{}{}
		edges = append(edges, e)
	}
	if len(edgeSet) != m {
		return fmt.Errorf("declared m' does not match unique edges")
	}
	hk := newHopcroftKarp(tc.n, tc.n)
	for _, e := range edges {
		hk.addEdge(e.from, e.to)
	}
	if hk.maxMatching() != tc.n {
		return fmt.Errorf("output graph is not good (no perfect matching)")
	}
	matchL := hk.pairU
	matchR := hk.pairV
	reach, err := computeReachability(tc.n, edges, matchL, matchR)
	if err != nil {
		return fmt.Errorf("failed to compute reachability: %v", err)
	}
	for v := 0; v < tc.n; v++ {
		if !reach[v].equal(tc.reach[v]) {
			return fmt.Errorf("tight sets differ (reachability mismatch at vertex %d)", v+1)
		}
	}
	minEdges := tc.minEdgeCount
	if len(edges) != minEdges {
		return fmt.Errorf("edge count not minimal: got %d, expected %d", len(edges), minEdges)
	}
	return nil
}

func verifyCase(tc testCase, output string) error {
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	first := strings.ToUpper(tokens[0])
	if tc.good {
		if first != "YES" {
			return fmt.Errorf("expected YES for good graph")
		}
		return verifyYES(tc, tokens)
	}
	if first != "NO" {
		return fmt.Errorf("expected NO for bad graph")
	}
	return verifyNO(tc, tokens)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/1835F_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := makeTestCase(rng, i%2 == 0)
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := verifyCase(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
