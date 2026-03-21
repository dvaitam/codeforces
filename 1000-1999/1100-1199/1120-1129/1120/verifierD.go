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

const embeddedSolverD = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v, id int
	w        int64
}

var (
	in    []int
	out   []int
	adj   [][]int
	timer int
)

func dfs(u, p int) {
	isLeaf := true
	in[u] = 1e9
	out[u] = -1e9
	for _, v := range adj[u] {
		if v != p {
			isLeaf = false
			dfs(v, u)
			if in[v] < in[u] {
				in[u] = in[v]
			}
			if out[v] > out[u] {
				out[u] = out[v]
			}
		}
	}
	if isLeaf && u != 1 {
		timer++
		in[u] = timer
		out[u] = timer
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &c[i])
	}

	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	in = make([]int, n+1)
	out = make([]int, n+1)
	dfs(1, 0)

	edges := make([]Edge, 0, n)
	for i := 1; i <= n; i++ {
		if in[i] <= out[i] {
			edges = append(edges, Edge{
				u:  in[i],
				v:  out[i] + 1,
				w:  c[i],
				id: i,
			})
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})

	parent := make([]int, timer+2)
	for i := 1; i <= timer+1; i++ {
		parent[i] = i
	}

	var find func(int) int
	find = func(i int) int {
		if parent[i] == i {
			return i
		}
		parent[i] = find(parent[i])
		return parent[i]
	}

	union := func(i, j int) {
		rootI := find(i)
		rootJ := find(j)
		if rootI != rootJ {
			parent[rootI] = rootJ
		}
	}

	current_components := timer + 1
	var total_cost int64
	var valid_edges []int

	for i := 0; i < len(edges); {
		j := i
		for j < len(edges) && edges[j].w == edges[i].w {
			j++
		}

		comps_before := current_components

		for k := i; k < j; k++ {
			if find(edges[k].u) != find(edges[k].v) {
				valid_edges = append(valid_edges, edges[k].id)
			}
		}

		for k := i; k < j; k++ {
			if find(edges[k].u) != find(edges[k].v) {
				union(edges[k].u, edges[k].v)
				current_components--
			}
		}

		total_cost += int64(comps_before-current_components) * edges[i].w
		i = j
	}

	sort.Ints(valid_edges)

	fmt.Fprintln(writer, total_cost, len(valid_edges))
	for i, id := range valid_edges {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, id)
	}
	fmt.Fprintln(writer)
}
`

func buildEmbeddedOracle() (string, func(), error) {
	tmpSrc, err := os.CreateTemp("", "oracle1120D-*.go")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpSrc.WriteString(embeddedSolverD); err != nil {
		tmpSrc.Close()
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpSrc.Close()

	tmpBin, err := os.CreateTemp("", "oracle1120D-bin-*")
	if err != nil {
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpBin.Close()

	cmd := exec.Command("go", "build", "-o", tmpBin.Name(), tmpSrc.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpSrc.Name())
		os.Remove(tmpBin.Name())
		return "", nil, fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	os.Remove(tmpSrc.Name())
	return tmpBin.Name(), func() { os.Remove(tmpBin.Name()) }, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp, err := os.CreateTemp("", "1120D-cand-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), abs)
		cmd.Dir = filepath.Dir(abs)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProgram(path string, input []byte) (string, string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func parseOutput(out string, n int) (int64, []int, error) {
	tokens := strings.Fields(out)
	if len(tokens) < 2 {
		return 0, nil, fmt.Errorf("expected at least 2 tokens, got %d", len(tokens))
	}
	cost, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid minimum cost %q", tokens[0])
	}
	k, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid vertex count %q", tokens[1])
	}
	if k < 0 {
		return 0, nil, fmt.Errorf("negative vertex count %d", k)
	}
	if len(tokens) != 2+k {
		return 0, nil, fmt.Errorf("expected %d tokens, got %d", 2+k, len(tokens))
	}
	verts := make([]int, k)
	prev := 0
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(tokens[2+i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid vertex index %q", tokens[2+i])
		}
		if v < 1 || v > n {
			return 0, nil, fmt.Errorf("vertex %d out of range [1,%d]", v, n)
		}
		if i > 0 && v <= prev {
			return 0, nil, fmt.Errorf("vertices not in strictly increasing order")
		}
		prev = v
		verts[i] = v
	}
	return cost, verts, nil
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(100)+1))
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, i))
	}
	return sb.String(), n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	candBin, cleanupCand, err := prepareCandidate(candidatePath)
	if err != nil {
		fail("failed to prepare candidate: %v", err)
	}
	defer cleanupCand()

	refBin, cleanupRef, err := buildEmbeddedOracle()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer cleanupRef()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 200; i++ {
		input, n := generateCase(rng)
		inputBytes := []byte(input)

		refOut, refErr, err := runProgram(refBin, inputBytes)
		if err != nil {
			fail("reference runtime error on case %d: %v\n%s", i, err, refErr)
		}
		expCost, expVerts, err := parseOutput(refOut, n)
		if err != nil {
			fail("failed to parse reference output on case %d: %v\noutput:\n%s", i, err, refOut)
		}

		candOut, candErr, err := runProgram(candBin, inputBytes)
		if err != nil {
			fail("candidate runtime error on case %d: %v\n%s", i, err, candErr)
		}
		gotCost, gotVerts, err := parseOutput(candOut, n)
		if err != nil {
			fail("invalid candidate output on case %d: %v\noutput:\n%s", i, err, candOut)
		}

		if gotCost != expCost {
			fail("case %d: wrong minimum cost: expected %d got %d\ninput:\n%s", i, expCost, gotCost, input)
		}
		if len(gotVerts) != len(expVerts) {
			fail("case %d: wrong vertex count: expected %d got %d", i, len(expVerts), len(gotVerts))
		}
		for j, v := range expVerts {
			if gotVerts[j] != v {
				fail("case %d: vertices mismatch at position %d: expected %d got %d", i, j+1, v, gotVerts[j])
			}
		}
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
