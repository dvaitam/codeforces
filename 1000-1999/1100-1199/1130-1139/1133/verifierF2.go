package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DSU struct {
	p []int
}

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) {
	fa := d.Find(a)
	fb := d.Find(b)
	if fa != fb {
		d.p[fa] = fb
	}
}

type edge struct{ u, v int }

func solveF2(n, D int, edges []edge) (bool, [][2]int) {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	if len(adj[1]) < D {
		return false, nil
	}
	comp := make([]int, n+1)
	compID := 0
	for i := 2; i <= n; i++ {
		if comp[i] == 0 {
			compID++
			stack := []int{i}
			comp[i] = compID
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, v := range adj[u] {
					if v == 1 || comp[v] != 0 {
						continue
					}
					comp[v] = compID
					stack = append(stack, v)
				}
			}
		}
	}
	if compID > D {
		return false, nil
	}
	compEdges := make(map[int][]int)
	for _, y := range adj[1] {
		compEdges[comp[y]] = append(compEdges[comp[y]], y)
	}
	chosen := make(map[int]bool)
	initial := make([]int, 0, D)
	for cid := 1; cid <= compID; cid++ {
		list := compEdges[cid]
		if len(list) == 0 {
			return false, nil
		}
		y := list[0]
		chosen[y] = true
		initial = append(initial, y)
	}
	left := D - compID
	for _, y := range adj[1] {
		if left == 0 {
			break
		}
		if chosen[y] {
			continue
		}
		chosen[y] = true
		initial = append(initial, y)
		left--
	}
	visited := make([]bool, n+1)
	visited[1] = true
	edgesOut := make([][2]int, 0, n-1)
	queue := make([]int, len(initial))
	for i, v := range initial {
		edgesOut = append(edgesOut, [2]int{1, v})
		visited[v] = true
		queue[i] = v
	}
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, v := range adj[u] {
			if v == 1 || visited[v] {
				continue
			}
			visited[v] = true
			edgesOut = append(edgesOut, [2]int{u, v})
			queue = append(queue, v)
		}
	}
	if len(edgesOut) != n-1 {
		return false, nil
	}
	return true, edgesOut
}

func generateCase(rng *rand.Rand) (string, int, []edge, int) {
	n := rng.Intn(6) + 2
	edges := make([]edge, 0, n*(n-1)/2)
	used := make(map[[2]int]bool)
	for i := 2; i <= n; i++ {
		edges = append(edges, edge{i - 1, i})
		used[[2]int{i - 1, i}] = true
	}
	total := n * (n - 1) / 2
	extras := rng.Intn(total - (n - 1) + 1)
	for len(edges) < n-1+extras {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if used[[2]int{u, v}] {
			continue
		}
		used[[2]int{u, v}] = true
		edges = append(edges, edge{u, v})
	}
	D := rng.Intn(n-1) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), D))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String(), D, edges, n
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseOutput(out string, n, D int, edgesIn map[[2]int]bool) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(scanner.Text())
	if first == "NO" {
		if scanner.Scan() {
			return fmt.Errorf("extra data")
		}
		return fmt.Errorf("NO")
	}
	if first != "YES" {
		return fmt.Errorf("missing YES/NO")
	}
	edgesOut := make([][2]int, 0, n-1)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var u, v int
		if _, err := fmt.Sscan(line, &u, &v); err != nil {
			return fmt.Errorf("bad edge")
		}
		edgesOut = append(edgesOut, [2]int{u, v})
	}
	if len(edgesOut) != n-1 {
		return fmt.Errorf("wrong edge count")
	}
	dsu := NewDSU(n)
	deg1 := 0
	for _, e := range edgesOut {
		a, b := e[0], e[1]
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("vertex range")
		}
		x, y := a, b
		if x > y {
			x, y = y, x
		}
		if !edgesIn[[2]int{x, y}] {
			return fmt.Errorf("edge not in input")
		}
		if dsu.Find(a) == dsu.Find(b) {
			return fmt.Errorf("not tree")
		}
		dsu.Union(a, b)
		if a == 1 || b == 1 {
			deg1++
		}
	}
	root := dsu.Find(1)
	for i := 2; i <= n; i++ {
		if dsu.Find(i) != root {
			return fmt.Errorf("not tree")
		}
	}
	if deg1 != D {
		return fmt.Errorf("degree not %d", D)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, D, edges, n := generateCase(rng)
		expectedOk, _ := solveF2(n, D, edges)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		edgeMap := make(map[[2]int]bool)
		for _, e := range edges {
			a, b := e.u, e.v
			if a > b {
				a, b = b, a
			}
			edgeMap[[2]int{a, b}] = true
		}
		err = parseOutput(out, n, D, edgeMap)
		if expectedOk {
			if err != nil && err.Error() == "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected YES but got NO\ninput:\n%s", i+1, in)
				os.Exit(1)
			}
			if err != nil && err.Error() != "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
				os.Exit(1)
			}
			if err == nil {
				continue
			}
		} else {
			if err == nil {
				fmt.Fprintf(os.Stderr, "case %d failed: expected NO but output tree\ninput:\n%s", i+1, in)
				os.Exit(1)
			}
			if err.Error() != "NO" {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
