package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"5",
	"1 1 2 4",
	"1 2 2 3",
	"2",
	"",
	"6",
	"1 2 1 1 5",
	"1 1 2 1 1",
	"3",
	"",
	"5",
	"1 2 2 3",
	"1 2 2 3",
	"1",
	"",
	"6",
	"1 1 3 4 1",
	"1 2 1 3 1",
	"2",
	"",
	"6",
	"1 1 1 4 1",
	"1 2 3 4 1",
	"3",
	"",
	"6",
	"1 1 3 3 5",
	"1 2 2 1 5",
	"4",
	"",
	"4",
	"1 2 1",
	"1 1 1",
	"3",
	"",
	"4",
	"1 1 1",
	"1 1 1",
	"1",
	"",
	"6",
	"1 2 3 2 2",
	"1 2 2 4 3",
	"1",
	"",
	"4",
	"1 2 3",
	"1 1 1",
	"1",
	"",
	"4",
	"1 1 2",
	"1 2 2",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 1 1 2 5",
	"1 2 1 3 1",
	"1",
	"",
	"6",
	"1 1 1 1 4",
	"1 1 3 1 5",
	"4",
	"",
	"6",
	"1 2 1 2 1",
	"1 2 2 2 1",
	"5",
	"",
	"5",
	"1 1 3 4",
	"1 2 2 4",
	"2",
	"",
	"3",
	"1 1",
	"1 2",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 2 3 3",
	"1 2 1 1",
	"4",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 2",
	"2",
	"",
	"6",
	"1 2 2 4 1",
	"1 1 3 3 2",
	"2",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"6",
	"1 1 1 4 1",
	"1 1 2 4 5",
	"5",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 1 2",
	"1 1 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 1 1 2 3",
	"1 1 1 4 4",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 1",
	"1 2 1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 3",
	"1 2 2",
	"3",
	"",
	"3",
	"1 2",
	"1 1",
	"1",
	"",
	"3",
	"1 2",
	"1 2",
	"1",
	"",
	"4",
	"1 1 2",
	"1 1 3",
	"2",
	"",
	"4",
	"1 1 2",
	"1 2 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 2",
	"1 1 1",
	"3",
	"",
	"5",
	"1 2 2 1",
	"1 1 3 4",
	"4",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 1",
	"1 1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"6",
	"1 2 2 2 4",
	"1 2 1 3 5",
	"2",
	"",
	"5",
	"1 2 1 1",
	"1 2 3 2",
	"1",
	"",
	"5",
	"1 2 1 1",
	"1 1 1 2",
	"4",
	"",
	"5",
	"1 1 1 4",
	"1 2 3 4",
	"4",
	"",
	"4",
	"1 2 3",
	"1 1 1",
	"2",
	"",
	"4",
	"1 1 3",
	"1 2 3",
	"1",
	"",
	"5",
	"1 2 1 1",
	"1 1 1 2",
	"1",
	"",
	"4",
	"1 2 2",
	"1 1 3",
	"2",
	"",
	"4",
	"1 1 3",
	"1 1 3",
	"2",
	"",
	"3",
	"1 2",
	"1 2",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 2 1 4 3",
	"1 2 2 4 1",
	"2",
	"",
	"3",
	"1 2",
	"1 1",
	"1",
	"",
	"5",
	"1 1 2 1",
	"1 2 2 1",
	"4",
	"",
	"5",
	"1 1 1 4",
	"1 1 1 2",
	"3",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 1 3 4",
	"1 1 3 3",
	"2",
	"",
	"5",
	"1 1 3 2",
	"1 2 2 3",
	"2",
	"",
	"5",
	"1 1 1 1",
	"1 2 1 2",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 2 2",
	"1 1 1",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 1",
	"1 2",
	"1",
	"",
	"6",
	"1 1 1 4 1",
	"1 2 3 3 2",
	"1",
	"",
	"3",
	"1 2",
	"1 2",
	"2",
	"",
	"6",
	"1 1 3 1 1",
	"1 1 2 2 2",
	"2",
	"",
	"4",
	"1 1 1",
	"1 2 2",
	"3",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"4",
	"1 1 3",
	"1 2 3",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"6",
	"1 2 2 1 4",
	"1 2 1 3 3",
	"4",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"4",
	"1 2 2",
	"1 2 1",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 1",
	"1 1",
	"2",
	"",
	"5",
	"1 2 1 4",
	"1 1 1 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 2 2 3",
	"1 1 1 1",
	"3",
	"",
	"4",
	"1 1 3",
	"1 1 1",
	"2",
	"",
	"4",
	"1 1 3",
	"1 1 2",
	"3",
	"",
	"6",
	"1 2 2 4 3",
	"1 1 1 1 4",
	"4",
	"",
	"6",
	"1 1 3 3 3",
	"1 2 1 4 4",
	"5",
	"",
	"4",
	"1 1 2",
	"1 1 3",
	"2",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"5",
	"1 1 3 1",
	"1 2 1 3",
	"1",
	"",
	"2",
	"1",
	"1",
	"1",
	"",
	"3",
	"1 2",
	"1 1",
	"2",
	"",
	"5",
	"1 2 1 3",
	"1 2 1 4",
	"2",
	"",
	"6",
	"1 1 1 2 4",
	"1 2 2 4 4",
	"2",
	"",
}

type testcase struct {
	n     int
	a     []int
	b     []int
	idx   int
	input string
}

func parseInts(line string) []int {
	if strings.TrimSpace(line) == "" {
		return nil
	}
	fields := strings.Fields(line)
	res := make([]int, len(fields))
	for i, f := range fields {
		val, _ := strconv.Atoi(f)
		res[i] = val
	}
	return res
}

func parseCases() []testcase {
	var cases []testcase
	for pos := 0; pos < len(rawTestcases); {
		for pos < len(rawTestcases) && strings.TrimSpace(rawTestcases[pos]) == "" {
			pos++
		}
		if pos >= len(rawTestcases) {
			break
		}
		n, _ := strconv.Atoi(strings.TrimSpace(rawTestcases[pos]))
		pos++
		if pos+2 >= len(rawTestcases) {
			break
		}
		lineA := rawTestcases[pos]
		lineB := rawTestcases[pos+1]
		lineIdx := rawTestcases[pos+2]
		pos += 3

		aParents := parseInts(lineA)
		bParents := parseInts(lineB)
		idxVal, _ := strconv.Atoi(strings.TrimSpace(lineIdx))

		a := make([]int, n+1)
		b := make([]int, n+1)
		for i := 0; i < len(aParents) && i+2 <= n; i++ {
			a[i+2] = aParents[i]
		}
		for i := 0; i < len(bParents) && i+2 <= n; i++ {
			b[i+2] = bParents[i]
		}

		var sb strings.Builder
		fmt.Fprintln(&sb, n)
		fmt.Fprintln(&sb, strings.TrimSpace(lineA))
		fmt.Fprintln(&sb, strings.TrimSpace(lineB))
		fmt.Fprintln(&sb, strings.TrimSpace(lineIdx))

		cases = append(cases, testcase{
			n:     n,
			a:     a,
			b:     b,
			idx:   idxVal,
			input: sb.String(),
		})
	}
	return cases
}

// --- Correct solver from accepted solution, adapted as function ---

type Edge struct {
	to, id int
}

type Tree struct {
	n          int
	adj        [][]Edge
	parent     []int
	depth      []int
	heavy      []int
	head       []int
	pos        []int
	edgeNode   []int
	currentPos int
	seg        [][]int
}

func newTree(n int) *Tree {
	return &Tree{
		n:        n,
		adj:      make([][]Edge, n+1),
		parent:   make([]int, n+1),
		depth:    make([]int, n+1),
		heavy:    make([]int, n+1),
		head:     make([]int, n+1),
		pos:      make([]int, n+1),
		edgeNode: make([]int, n),
		seg:      make([][]int, 4*n+1),
	}
}

func (t *Tree) dfs1(u, p, d int) int {
	t.parent[u] = p
	t.depth[u] = d
	size := 1
	maxSub := 0
	for _, edge := range t.adj[u] {
		v := edge.to
		if v != p {
			t.edgeNode[edge.id] = v
			sub := t.dfs1(v, u, d+1)
			size += sub
			if sub > maxSub {
				maxSub = sub
				t.heavy[u] = v
			}
		}
	}
	return size
}

func (t *Tree) dfs2(u, h int) {
	t.head[u] = h
	t.currentPos++
	t.pos[u] = t.currentPos
	if t.heavy[u] != 0 {
		t.dfs2(t.heavy[u], h)
	}
	for _, edge := range t.adj[u] {
		v := edge.to
		if v != t.parent[u] && v != t.heavy[u] {
			t.dfs2(v, v)
		}
	}
}

func (t *Tree) addSeg(node, L, R, ql, qr, id int) {
	if ql <= L && R <= qr {
		t.seg[node] = append(t.seg[node], id)
		return
	}
	mid := (L + R) / 2
	if ql <= mid {
		t.addSeg(node*2, L, mid, ql, qr, id)
	}
	if qr > mid {
		t.addSeg(node*2+1, mid+1, R, ql, qr, id)
	}
}

func (t *Tree) addPath(u, v, id int) {
	for t.head[u] != t.head[v] {
		if t.depth[t.head[u]] < t.depth[t.head[v]] {
			u, v = v, u
		}
		t.addSeg(1, 1, t.n, t.pos[t.head[u]], t.pos[u], id)
		u = t.parent[t.head[u]]
	}
	if t.depth[u] > t.depth[v] {
		u, v = v, u
	}
	if t.pos[u]+1 <= t.pos[v] {
		t.addSeg(1, 1, t.n, t.pos[u]+1, t.pos[v], id)
	}
}

func (t *Tree) queryAndClear(node, L, R, p int, result *[]int, isDeleted []bool) {
	if len(t.seg[node]) > 0 {
		for _, id := range t.seg[node] {
			if !isDeleted[id] {
				isDeleted[id] = true
				*result = append(*result, id)
			}
		}
		t.seg[node] = nil
	}
	if L == R {
		return
	}
	mid := (L + R) / 2
	if p <= mid {
		t.queryAndClear(node*2, L, mid, p, result, isDeleted)
	} else {
		t.queryAndClear(node*2+1, mid+1, R, p, result, isDeleted)
	}
}

func solve(n int, a, b []int, idx int) string {
	var out bytes.Buffer

	blueTree := newTree(n)
	for i := 1; i < n; i++ {
		u := i + 1
		v := a[i]
		blueTree.adj[u] = append(blueTree.adj[u], Edge{v, i})
		blueTree.adj[v] = append(blueTree.adj[v], Edge{u, i})
	}

	redTree := newTree(n)
	for i := 1; i < n; i++ {
		u := i + 1
		v := b[i]
		redTree.adj[u] = append(redTree.adj[u], Edge{v, i})
		redTree.adj[v] = append(redTree.adj[v], Edge{u, i})
	}

	blueTree.dfs1(1, 0, 0)
	blueTree.dfs2(1, 1)

	redTree.dfs1(1, 0, 0)
	redTree.dfs2(1, 1)

	for i := 1; i < n; i++ {
		u := i + 1
		v := b[i]
		blueTree.addPath(u, v, i)
	}

	for i := 1; i < n; i++ {
		u := i + 1
		v := a[i]
		redTree.addPath(u, v, i)
	}

	blueDeleted := make([]bool, n)
	redDeleted := make([]bool, n)

	blueStage := []int{idx}
	blueDeleted[idx] = true
	var redStage []int

	printSlice := func(s []int) {
		for i, v := range s {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(&out, v)
		}
		out.WriteByte('\n')
	}

	for len(blueStage) > 0 || len(redStage) > 0 {
		if len(blueStage) > 0 {
			out.WriteString("Blue\n")
			sort.Ints(blueStage)
			printSlice(blueStage)

			redStage = nil
			for _, id := range blueStage {
				p := blueTree.pos[blueTree.edgeNode[id]]
				blueTree.queryAndClear(1, 1, blueTree.n, p, &redStage, redDeleted)
			}
			blueStage = nil
		} else if len(redStage) > 0 {
			out.WriteString("Red\n")
			sort.Ints(redStage)
			printSlice(redStage)

			blueStage = nil
			for _, id := range redStage {
				p := redTree.pos[redTree.edgeNode[id]]
				redTree.queryAndClear(1, 1, redTree.n, p, &blueStage, blueDeleted)
			}
			redStage = nil
		}
	}

	return strings.TrimSpace(out.String())
}

// --- end correct solver ---

func runSolution(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for i, tc := range cases {
		// Convert a/b arrays to the format solve() expects (1-indexed parent arrays)
		// solve() expects a[i] = parent of node i+1 for i=1..n-1
		// tc.a[j] for j=2..n is the parent of node j in blue tree
		// We need to pass arrays where a[i] = tc.a[i+1] for i=1..n-1
		aArr := make([]int, tc.n)
		bArr := make([]int, tc.n)
		for j := 1; j < tc.n; j++ {
			aArr[j] = tc.a[j+1]
			bArr[j] = tc.b[j+1]
		}
		expected := solve(tc.n, aArr, bArr, tc.idx)
		got, err := runSolution(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d mismatch:\nexpected:\n%s\n\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
