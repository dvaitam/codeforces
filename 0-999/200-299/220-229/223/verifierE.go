package main

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Correct solver for 223E (planar graph, count vertices inside/on a cycle).

type Edge223 struct {
	to  int
	id  int
	dir int
}

var area_b, term1_b, term2_b, xu_b, yv_b, xv_b, yu_b *big.Int

func initBig223() {
	area_b = new(big.Int)
	term1_b = new(big.Int)
	term2_b = new(big.Int)
	xu_b = new(big.Int)
	yv_b = new(big.Int)
	xv_b = new(big.Int)
	yu_b = new(big.Int)
}

func faceAreaSign223(edges []int, u_of_dir, v_of_dir []int, x, y []int64) int {
	area_b.SetInt64(0)
	for _, e := range edges {
		u := u_of_dir[e]
		v := v_of_dir[e]
		xu_b.SetInt64(x[u])
		yv_b.SetInt64(y[v])
		xv_b.SetInt64(x[v])
		yu_b.SetInt64(y[u])
		term1_b.Mul(xu_b, yv_b)
		term2_b.Mul(xv_b, yu_b)
		area_b.Add(area_b, term1_b.Sub(term1_b, term2_b))
	}
	return area_b.Sign()
}

func polyAreaSign223(nodes []int, x, y []int64) int {
	area_b.SetInt64(0)
	k := len(nodes)
	for i := 0; i < k; i++ {
		u := nodes[i]
		v := nodes[(i+1)%k]
		xu_b.SetInt64(x[u])
		yv_b.SetInt64(y[v])
		xv_b.SetInt64(x[v])
		yu_b.SetInt64(y[u])
		term1_b.Mul(xu_b, yv_b)
		term2_b.Mul(xv_b, yu_b)
		area_b.Add(area_b, term1_b.Sub(term1_b, term2_b))
	}
	return area_b.Sign()
}

func half223(dx, dy int64) int {
	if dy > 0 || (dy == 0 && dx > 0) {
		return 1
	}
	return 2
}

func solve223E(input string) string {
	data := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(data) && data[pos] <= ' ' {
			pos++
		}
		if pos >= len(data) {
			return 0
		}
		sign := 1
		if data[pos] == '-' {
			sign = -1
			pos++
		}
		res := 0
		for pos < len(data) && data[pos] > ' ' {
			res = res*10 + int(data[pos]-'0')
			pos++
		}
		return res * sign
	}

	n := nextInt()
	m := nextInt()

	adj := make([][]Edge223, n+1)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		adj[u] = append(adj[u], Edge223{v, i, 0})
		adj[v] = append(adj[v], Edge223{u, i, 1})
	}

	x := make([]int64, n+1)
	y := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		x[i] = int64(nextInt())
		y[i] = int64(nextInt())
	}

	for i := 1; i <= n; i++ {
		u := i
		sort.Slice(adj[u], func(j, k int) bool {
			e1 := adj[u][j]
			e2 := adj[u][k]
			dx1 := x[e1.to] - x[u]
			dy1 := y[e1.to] - y[u]
			dx2 := x[e2.to] - x[u]
			dy2 := y[e2.to] - y[u]
			h1 := half223(dx1, dy1)
			h2 := half223(dx2, dy2)
			if h1 != h2 {
				return h1 < h2
			}
			cross := dx1*dy2 - dx2*dy1
			return cross > 0
		})
	}

	pos_in_adj := make([]int, 2*m)
	u_of_dir_edge := make([]int, 2*m)
	v_of_dir_edge := make([]int, 2*m)
	for i := 1; i <= n; i++ {
		for k, e := range adj[i] {
			dir_edge := 2*e.id + e.dir
			pos_in_adj[dir_edge] = k
			u_of_dir_edge[dir_edge] = i
			v_of_dir_edge[dir_edge] = e.to
		}
	}

	initBig223()

	visited := make([]bool, 2*m)
	type Face struct {
		edges []int
	}
	faces := []Face{}

	for i := 0; i < 2*m; i++ {
		if !visited[i] {
			curr := i
			face_edges := []int{}
			for {
				visited[curr] = true
				face_edges = append(face_edges, curr)
				v := v_of_dir_edge[curr]
				rev := curr ^ 1
				p := pos_in_adj[rev]
				next_idx := (p + 1) % len(adj[v])
				next_e := adj[v][next_idx]
				curr = 2*next_e.id + next_e.dir
				if curr == i {
					break
				}
			}
			faces = append(faces, Face{face_edges})
		}
	}

	f_out := -1
	for f := 0; f < len(faces); f++ {
		if faceAreaSign223(faces[f].edges, u_of_dir_edge, v_of_dir_edge, x, y) < 0 {
			f_out = f
			break
		}
	}

	face_of_dir := make([]int, 2*m)
	for f, face := range faces {
		for _, e := range face.edges {
			face_of_dir[e] = f
		}
	}

	visited_face := make([]bool, len(faces))
	visited_face[f_out] = true
	queue := []int{f_out}
	parent_edge := make([]int, len(faces))
	for i := range parent_edge {
		parent_edge[i] = -1
	}

	post_order := make([]int, 0, len(faces))

	for len(queue) > 0 {
		curr_f := queue[0]
		queue = queue[1:]
		post_order = append(post_order, curr_f)
		for _, e := range faces[curr_f].edges {
			neighbor := face_of_dir[e^1]
			if !visited_face[neighbor] {
				visited_face[neighbor] = true
				parent_edge[neighbor] = e ^ 1
				queue = append(queue, neighbor)
			}
		}
	}

	g_prime := make([]int64, 2*m)
	for i := len(post_order) - 1; i >= 1; i-- {
		f := post_order[i]
		p_edge := parent_edge[f]
		sum := int64(0)
		for _, e := range faces[f].edges {
			if e != p_edge {
				sum += g_prime[e]
			}
		}
		W := int64(len(faces[f].edges) - 2)
		val := W - sum
		g_prime[p_edge] = val
		g_prime[p_edge^1] = -val
	}

	edge_id_map := make(map[uint64]int)
	for i := 0; i < 2*m; i++ {
		u := u_of_dir_edge[i]
		v := v_of_dir_edge[i]
		edge_id_map[uint64(u)<<32|uint64(v)] = i
	}

	q := nextInt()
	var outBuf bytes.Buffer
	for i := 0; i < q; i++ {
		k := nextInt()
		a := make([]int, k)
		for j := 0; j < k; j++ {
			a[j] = nextInt()
		}

		edges := make([]int, k)
		for j := 0; j < k; j++ {
			u := a[j]
			v := a[(j+1)%k]
			edges[j] = edge_id_map[uint64(u)<<32|uint64(v)]
		}

		sign := polyAreaSign223(a, x, y)

		S_given := int64(0)
		for j := 0; j < k; j++ {
			S_given += g_prime[edges[j]]
		}

		S_ccw := S_given
		if sign < 0 {
			S_ccw = -S_given
		}

		V_in := (int64(2) + int64(k) + S_ccw) / 2
		outBuf.WriteString(strconv.FormatInt(V_in, 10))
		if i < q-1 {
			outBuf.WriteByte(' ')
		}
	}

	return outBuf.String()
}

// Generate a small planar graph for testing.
// We create a triangulated grid which is guaranteed planar and 2-connected.
func genPlanarGraph(rng *rand.Rand) string {
	// Grid size: rows x cols of vertices
	rows := rng.Intn(3) + 2 // 2-4
	cols := rng.Intn(3) + 2 // 2-4
	n := rows * cols

	// Generate coordinates
	xs := make([]int64, n+1)
	ys := make([]int64, n+1)
	idx := func(r, c int) int { return r*cols + c + 1 }
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			v := idx(r, c)
			xs[v] = int64(c*10 + rng.Intn(3))
			ys[v] = int64(r*10 + rng.Intn(3))
		}
	}

	type edgePair struct{ u, v int }
	edgeSet := map[edgePair]bool{}
	addEdge := func(u, v int) {
		if u > v {
			u, v = v, u
		}
		edgeSet[edgePair{u, v}] = true
	}

	// Grid edges + diagonals for triangulation
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if c+1 < cols {
				addEdge(idx(r, c), idx(r, c+1))
			}
			if r+1 < rows {
				addEdge(idx(r, c), idx(r+1, c))
			}
			if r+1 < rows && c+1 < cols {
				addEdge(idx(r, c), idx(r+1, c+1))
			}
		}
	}

	edges := make([]edgePair, 0, len(edgeSet))
	for e := range edgeSet {
		edges = append(edges, e)
	}
	m := len(edges)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}

	// Generate queries: use small face cycles from the triangulation
	// Pick triangles from the grid
	type triangle struct{ a, b, c int }
	var triangles []triangle
	for r := 0; r < rows-1; r++ {
		for c := 0; c < cols-1; c++ {
			// Lower triangle: (r,c), (r+1,c), (r+1,c+1)
			triangles = append(triangles, triangle{idx(r, c), idx(r+1, c), idx(r+1, c+1)})
			// Upper triangle: (r,c), (r,c+1), (r+1,c+1)
			triangles = append(triangles, triangle{idx(r, c), idx(r, c+1), idx(r+1, c+1)})
		}
	}

	q := rng.Intn(3) + 1
	if q > len(triangles) {
		q = len(triangles)
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	perm := rng.Perm(len(triangles))
	for i := 0; i < q; i++ {
		t := triangles[perm[i]]
		sb.WriteString(fmt.Sprintf("3 %d %d %d\n", t.a, t.b, t.c))
	}

	return sb.String()
}

func runBin(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s\nstdout: %s", err, stderr.String(), stdout.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = io.Discard // suppress unused import
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genPlanarGraph(rng)
		exp := solve223E(input)
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expFields := strings.Fields(exp)
		gotFields := strings.Fields(got)
		if len(expFields) != len(gotFields) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s\nexpected: %s\ngot: %s\n",
				i+1, len(expFields), len(gotFields), input, exp, got)
			os.Exit(1)
		}
		for j := range expFields {
			if expFields[j] != gotFields[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expFields[j], gotFields[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
