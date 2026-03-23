package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ---------- embedded solver (correct Edmonds/Chu-Liu with edge recovery) ----------

func solveEmbedded(input string) string {
	sc := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(sc[idx])
		idx++
		return v
	}

	n := nextInt()
	m := nextInt()

	type edge struct {
		u, v, w int
	}
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		w := nextInt()
		edges[i] = edge{u, v, w}
	}

	// Check reachability from 1 in original directed graph
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
	}
	vis := make([]bool, n+1)
	vis[1] = true
	q := []int{1}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				q = append(q, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return "-1"
		}
	}

	// Edmonds' algorithm with edge recovery
	// We work with 0-indexed nodes internally
	root := 0
	curEdges := make([]wEdge, m)
	for i, e := range edges {
		curEdges[i] = wEdge{e.u - 1, e.v - 1, e.w, i}
	}

	// iterative Edmonds with edge id tracking
	// chosen[v] = original edge index chosen as incoming edge for v
	chosen := edmonds(n, root, curEdges)
	if chosen == nil {
		return "-1"
	}

	var result []int
	for _, eidx := range chosen {
		if edges[eidx].w == 1 {
			result = append(result, eidx+1)
		}
	}

	var out bytes.Buffer
	fmt.Fprintln(&out, len(result))
	if len(result) > 0 {
		for i, id := range result {
			if i > 0 {
				fmt.Fprint(&out, " ")
			}
			fmt.Fprint(&out, id)
		}
		fmt.Fprintln(&out)
	}
	return strings.TrimSpace(out.String())
}

// edmonds finds minimum cost arborescence rooted at root.
// Returns slice where result[v] = original edge index for v's incoming edge (v != root).
// Returns nil if no arborescence exists.
func edmonds(n, root int, edges []wEdge) []int {
	const INF = int(1e18)

	// We'll track which original edge is chosen for each vertex through contractions.
	// id[v] = supernode id for vertex v
	// For each contraction, we record what happened so we can expand later.

	type contraction struct {
		cycleNodes []int          // nodes in the cycle (in supernode space before contraction)
		cycleEdge  map[int]int    // cycleEdge[v] = edge index in 'edges' that was the min incoming for v (in the cycle)
		newNode    int            // the supernode id
	}

	curN := n
	// id[v] maps original vertex v to current supernode
	id := make([]int, n)
	for i := range id {
		id[i] = i
	}

	var contractions []contraction
	curEdges := make([]wEdge, len(edges))
	copy(curEdges, edges)

	for {
		// Find minimum incoming edge for each non-root node
		minIn := make([]int, curN)    // min incoming cost
		minEdge := make([]int, curN)  // index in curEdges of min incoming edge
		pre := make([]int, curN)      // predecessor node
		for i := range minIn {
			minIn[i] = INF
			minEdge[i] = -1
			pre[i] = -1
		}

		for i, e := range curEdges {
			if e.from == e.to {
				continue
			}
			if e.cost < minIn[e.to] {
				minIn[e.to] = e.cost
				minEdge[e.to] = i
				pre[e.to] = e.from
			}
		}

		// Check if all non-root nodes are reachable
		for v := 0; v < curN; v++ {
			if v == root && v < curN {
				continue
			}
			if v == root {
				continue
			}
			if minIn[v] == INF {
				return nil
			}
		}

		// Find cycles
		visited := make([]int, curN) // 0=unvisited, -1=done, >0=visit id
		cycleID := make([]int, curN) // which cycle a node belongs to (-1 = none)
		for i := range cycleID {
			cycleID[i] = -1
		}

		numCycles := 0
		for v := 0; v < curN; v++ {
			if v == root || visited[v] != 0 {
				continue
			}
			// Walk back through predecessors
			path := []int{}
			u := v
			for u != root && visited[u] == 0 {
				visited[u] = v + 1 // mark with unique id
				path = append(path, u)
				u = pre[u]
			}
			if u != root && visited[u] == v+1 {
				// Found a cycle, mark all nodes from u back to u
				cid := numCycles
				numCycles++
				w := u
				for {
					cycleID[w] = cid
					w = pre[w]
					if w == u {
						break
					}
				}
			}
			// Mark all visited nodes as done
			for _, w := range path {
				visited[w] = -1
			}
		}

		if numCycles == 0 {
			// No cycles, we have our arborescence
			// Recover original edge indices
			// minEdge[v] points to the edge in curEdges chosen for each v
			result := make([]int, n) // result[original_v] = original edge index
			for i := range result {
				result[i] = -1
			}

			// First, record choices at current level
			chosenAtLevel := make([]int, curN) // chosenAtLevel[supernode] = curEdges index
			for v := 0; v < curN; v++ {
				if v == root {
					continue
				}
				chosenAtLevel[v] = minEdge[v]
			}

			// Now expand contractions in reverse
			// We need to map supernodes back to original vertices
			// and figure out which original edges are chosen

			// Build final chosen map: supernode -> original edge index
			superChosen := make(map[int]int)
			for v := 0; v < curN; v++ {
				if v == root {
					continue
				}
				superChosen[v] = curEdges[chosenAtLevel[v]].origIdx
			}

			// Expand contractions in reverse order
			for i := len(contractions) - 1; i >= 0; i-- {
				c := contractions[i]
				// The supernode c.newNode was chosen to have incoming edge superChosen[c.newNode]
				// This edge's target (in original space) determines which cycle node receives external edge
				incomingOrigIdx := superChosen[c.newNode]

				// Add all cycle edges to superChosen
				for _, cNode := range c.cycleNodes {
					superChosen[cNode] = c.cycleEdge[cNode]
				}
				// Override the one that receives the external edge
				// We need to find which cycle node contains incomingTarget
				// But after multiple contractions, incomingTarget might be inside a nested supernode
				// We need to find which cycle node of THIS contraction contains incomingTarget
				// Actually, incomingTarget is an original vertex (0-indexed). We need to figure out
				// which cycleNode it mapped to at the time of this contraction.

				// The cycle nodes are supernode IDs at the level before this contraction.
				// We need to know which cycle node contained the original vertex incomingTarget.
				// We can determine this by looking at the edge: curEdges[chosenAtLevel].to was the
				// supernode that incomingTarget mapped to.

				// Actually, a cleaner approach: the edge in edges[incomingOrigIdx] targets original vertex
				// incomingTarget+1. At the time of contraction i, this vertex was part of one of the
				// cycle nodes. The cycle edge for that node gets replaced by the external edge.

				// But we don't easily know the mapping at contraction time. Let me use a different approach.

				// Actually let's track this differently.
				delete(superChosen, c.newNode)

				// The external incoming edge replaces the cycle edge for one cycle node.
				// That cycle node is the one whose cycle edge target is "replaced" by the external edge.
				// In Edmonds', the external edge enters the cycle at the vertex it points to.
				// At contraction level i, the edge pointed to c.newNode (the supernode).
				// Before contraction, it pointed to one of the cycle nodes.

				// We need to figure out which cycle node the external edge originally entered.
				// The origIdx tells us the original edge, and edges[origIdx].to-1 = original target vertex.
				// We need to trace this original vertex through contractions 0..i-1 to find its supernode
				// at level i.

				targetSuper := edges[incomingOrigIdx].to
				for j := 0; j < i; j++ {
					cj := contractions[j]
					for _, cn := range cj.cycleNodes {
						if cn == targetSuper {
							targetSuper = cj.newNode
							break
						}
					}
				}

				// targetSuper is now the cycle node that receives the external edge
				superChosen[targetSuper] = incomingOrigIdx
			}

			// Now superChosen maps original vertices (0-indexed) to their incoming original edge index
			for v, eidx := range superChosen {
				if v < n {
					result[v] = eidx
				}
			}

			return result
		}

		// Contract cycles
		// Build new node mapping
		newID := make([]int, curN)
		nextNode := 0
		// First assign IDs to cycle supernodes
		cycleNewNode := make([]int, numCycles)
		for i := 0; i < numCycles; i++ {
			cycleNewNode[i] = nextNode
			nextNode++
		}
		// Then assign IDs to non-cycle nodes
		for v := 0; v < curN; v++ {
			if cycleID[v] >= 0 {
				newID[v] = cycleNewNode[cycleID[v]]
			} else {
				newID[v] = nextNode
				nextNode++
			}
		}
		newRoot := newID[root]

		// Record contraction info
		for cid := 0; cid < numCycles; cid++ {
			var cycleNodes []int
			cycleEdge := make(map[int]int)
			for v := 0; v < curN; v++ {
				if cycleID[v] == cid {
					cycleNodes = append(cycleNodes, v)
					cycleEdge[v] = curEdges[minEdge[v]].origIdx
				}
			}
			contractions = append(contractions, contraction{
				cycleNodes: cycleNodes,
				cycleEdge:  cycleEdge,
				newNode:    cycleNewNode[cid],
			})
		}

		// Build new edge set
		var newEdges []wEdge
		for _, e := range curEdges {
			nu := newID[e.from]
			nv := newID[e.to]
			if nu == nv {
				continue
			}
			newCost := e.cost - minIn[e.to]
			newEdges = append(newEdges, wEdge{nu, nv, newCost, e.origIdx})
		}

		curEdges = newEdges
		curN = nextNode
		root = newRoot
	}
}

type wEdge struct {
	from, to, cost, origIdx int
}

// ---------- verifier infrastructure ----------

type edgeData struct {
	from, to   int
	needRepair bool
}

type verTestCase struct {
	input  string
	expect int
	n      int
	edges  []edgeData
}

type refEdge struct {
	from, to int
	cost     int
}

const inf = int(1e9)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkAnswer(tc, out); err != nil {
			// Also run embedded solver for comparison
			refOut := solveEmbedded(tc.input)
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\nexpected minimal repairs: %d\nreference output:\n%s\n", i+1, err, tc.input, out, tc.expect, refOut)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkAnswer(tc verTestCase, output string) error {
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		return fmt.Errorf("empty output")
	}
	tokens := strings.Fields(output)
	if tokens[0] == "-1" {
		if tc.expect != -1 {
			return fmt.Errorf("expected %d repairs but got -1", tc.expect)
		}
		if len(tokens) != 1 {
			return fmt.Errorf("unexpected extra tokens after -1")
		}
		return nil
	}
	if tc.expect == -1 {
		return fmt.Errorf("expected -1 but got %s", tokens[0])
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil || k < 0 {
		return fmt.Errorf("invalid number of repairs %q", tokens[0])
	}
	if k == 0 {
		if tc.expect != 0 {
			return fmt.Errorf("expected %d repairs but reported 0", tc.expect)
		}
		if len(tokens) != 1 {
			return fmt.Errorf("extra tokens after 0")
		}
		return nil
	}
	if len(tokens) < k+1 {
		return fmt.Errorf("reported %d roads but provided %d identifiers", k, len(tokens)-1)
	}
	if k != tc.expect {
		return fmt.Errorf("expected %d repairs but reported %d", tc.expect, k)
	}
	selected := make(map[int]struct{}, k)
	for i := 0; i < k; i++ {
		id, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return fmt.Errorf("invalid road index %q", tokens[i+1])
		}
		if id < 1 || id > len(tc.edges) {
			return fmt.Errorf("road index %d out of range", id)
		}
		if _, ok := selected[id]; ok {
			return fmt.Errorf("road index %d listed multiple times", id)
		}
		if !tc.edges[id-1].needRepair {
			return fmt.Errorf("road %d is already good and cannot be repaired", id)
		}
		selected[id] = struct{}{}
	}
	if !allReachable(tc.n, tc.edges, selected) {
		return fmt.Errorf("not all cities reachable from the capital with reported repairs")
	}
	return nil
}

func allReachable(n int, edges []edgeData, repaired map[int]struct{}) bool {
	adj := make([][]int, n+1)
	for idx, e := range edges {
		if !e.needRepair {
			adj[e.from] = append(adj[e.from], e.to)
			continue
		}
		if _, ok := repaired[idx+1]; ok {
			adj[e.from] = append(adj[e.from], e.to)
		}
	}
	vis := make([]bool, n+1)
	queue := []int{1}
	vis[1] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				queue = append(queue, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return false
		}
	}
	return true
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), err
}

func genTests() []verTestCase {
	rand.Seed(42)
	var tests []verTestCase
	tests = append(tests, newTestCase(1, nil))
	tests = append(tests, newTestCase(2, []edgeData{
		{from: 1, to: 2, needRepair: false},
	}))
	tests = append(tests, newTestCase(2, nil))
	tests = append(tests, newTestCase(3, []edgeData{
		{from: 1, to: 2, needRepair: true},
		{from: 2, to: 3, needRepair: false},
	}))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTestCase())
	}
	return tests
}

func randomTestCase() verTestCase {
	n := rand.Intn(10) + 1
	maxPossible := n * (n - 1)
	limit := n*3 + rand.Intn(n+1)
	if limit > maxPossible {
		limit = maxPossible
	}
	m := 0
	if limit > 0 {
		m = rand.Intn(limit + 1)
	}
	edges := make([]edgeData, 0, m)
	used := make(map[int]struct{})
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		key := (u-1)*n + (v - 1)
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		edges = append(edges, edgeData{
			from:       u,
			to:         v,
			needRepair: rand.Intn(2) == 1,
		})
	}
	return newTestCase(n, edges)
}

func newTestCase(n int, edges []edgeData) verTestCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		c := 0
		if e.needRepair {
			c = 1
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.from, e.to, c))
	}
	edgesCopy := make([]edgeData, len(edges))
	copy(edgesCopy, edges)
	expect := calcMinRepairs(n, edgesCopy)
	return verTestCase{
		input:  sb.String(),
		expect: expect,
		n:      n,
		edges:  edgesCopy,
	}
}

func calcMinRepairs(n int, edges []edgeData) int {
	if n == 0 {
		return 0
	}
	refEdges := make([]refEdge, 0, len(edges))
	for _, e := range edges {
		c := 0
		if e.needRepair {
			c = 1
		}
		refEdges = append(refEdges, refEdge{
			from: e.from - 1,
			to:   e.to - 1,
			cost: c,
		})
	}
	cost, ok := directedMST(n, 0, refEdges)
	if !ok {
		return -1
	}
	return cost
}

func directedMST(n, root int, edges []refEdge) (int, bool) {
	total := 0
	for {
		in := make([]int, n)
		pre := make([]int, n)
		id := make([]int, n)
		vis := make([]int, n)
		for i := range in {
			in[i] = inf
			id[i] = -1
			vis[i] = -1
		}
		for _, e := range edges {
			if e.from == e.to {
				continue
			}
			if e.cost < in[e.to] {
				in[e.to] = e.cost
				pre[e.to] = e.from
			}
		}
		in[root] = 0
		for v := 0; v < n; v++ {
			if v == root {
				continue
			}
			if in[v] == inf {
				return 0, false
			}
		}
		cnt := 0
		for v := 0; v < n; v++ {
			total += in[v]
			u := v
			for vis[u] != v && id[u] == -1 && u != root {
				vis[u] = v
				u = pre[u]
			}
			if u != root && id[u] == -1 {
				for x := pre[u]; x != u; x = pre[x] {
					id[x] = cnt
				}
				id[u] = cnt
				cnt++
			}
		}
		if cnt == 0 {
			return total, true
		}
		for v := 0; v < n; v++ {
			if id[v] == -1 {
				id[v] = cnt
				cnt++
			}
		}
		newEdges := make([]refEdge, 0, len(edges))
		for _, e := range edges {
			u := id[e.from]
			v := id[e.to]
			if u != v {
				newEdges = append(newEdges, refEdge{
					from: u,
					to:   v,
					cost: e.cost - in[e.to],
				})
			}
		}
		root = id[root]
		n = cnt
		edges = newEdges
	}
}
