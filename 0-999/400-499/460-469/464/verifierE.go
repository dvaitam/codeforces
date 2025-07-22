package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD = 1000000007

type Edge struct {
	to, x int
}

func dijkstra(n int, edges [][]Edge, s, t int) (*big.Int, []int) {
	dist := make([]*big.Int, n+1)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = nil
		parent[i] = -1
	}
	dist[s] = big.NewInt(0)
	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, &Item{node: s, dist: big.NewInt(0)})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(*Item)
		u := it.node
		d := it.dist
		if dist[u] == nil || d.Cmp(dist[u]) != 0 {
			continue
		}
		if u == t {
			break
		}
		for _, e := range edges[u] {
			w := new(big.Int).Lsh(big.NewInt(1), uint(e.x))
			nd := new(big.Int).Add(d, w)
			if dist[e.to] == nil || nd.Cmp(dist[e.to]) < 0 {
				dist[e.to] = nd
				parent[e.to] = u
				heap.Push(pq, &Item{node: e.to, dist: nd})
			}
		}
	}
	if dist[t] == nil {
		return nil, nil
	}
	path := []int{}
	for v := t; v != -1; v = parent[v] {
		path = append(path, v)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return dist[t], path
}

func modInt(b *big.Int) int64 {
	return new(big.Int).Mod(b, big.NewInt(MOD)).Int64()
}

// priority queue
type Item struct {
	node int
	dist *big.Int
}

type PQ []*Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist.Cmp(pq[j].dist) < 0 }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
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
	return strings.TrimSpace(out.String()), nil
}

type caseE struct {
	n, m  int
	edges [][3]int
	s, t  int
}

func parseCases(path string) ([]caseE, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	cases := []caseE{}
	for {
		var n, m int
		if !scan.Scan() {
			break
		}
		n, _ = strconv.Atoi(scan.Text())
		if !scan.Scan() {
			return nil, fmt.Errorf("bad file")
		}
		m, _ = strconv.Atoi(scan.Text())
		edges := make([][3]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			edges[i] = [3]int{u, v, x}
		}
		scan.Scan()
		s, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		t, _ := strconv.Atoi(scan.Text())
		cases = append(cases, caseE{n: n, m: m, edges: edges, s: s, t: t})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesE.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		adj := make([][]Edge, tc.n+1)
		for _, e := range tc.edges {
			u, v, x := e[0], e[1], e[2]
			adj[u] = append(adj[u], Edge{v, x})
			adj[v] = append(adj[v], Edge{u, x})
		}
		dist, _ := dijkstra(tc.n, adj, tc.s, tc.t)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
		}
		fmt.Fprintf(&sb, "%d %d\n", tc.s, tc.t)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if dist == nil {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: no output\n", idx+1)
				os.Exit(1)
			}
			if scan.Text() != "-1" {
				fmt.Fprintf(os.Stderr, "case %d: expected -1\n", idx+1)
				os.Exit(1)
			}
			if scan.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: missing distance\n", idx+1)
			os.Exit(1)
		}
		modValExp := modInt(dist)
		modValGot, err := strconv.Atoi(scan.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad distance\n", idx+1)
			os.Exit(1)
		}
		if int64(modValGot) != modValExp {
			fmt.Fprintf(os.Stderr, "case %d: expected dist %d got %d\n", idx+1, modValExp, modValGot)
			os.Exit(1)
		}
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: missing path len\n", idx+1)
			os.Exit(1)
		}
		L, _ := strconv.Atoi(scan.Text())
		path := make([]int, L)
		for i := 0; i < L; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: incomplete path\n", idx+1)
				os.Exit(1)
			}
			path[i], _ = strconv.Atoi(scan.Text())
		}
		if scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
		if path[0] != tc.s || path[len(path)-1] != tc.t {
			fmt.Fprintf(os.Stderr, "case %d: wrong endpoints\n", idx+1)
			os.Exit(1)
		}
		// verify path cost
		cost := big.NewInt(0)
		for i := 0; i < len(path)-1; i++ {
			u := path[i]
			v := path[i+1]
			found := false
			for _, e := range adj[u] {
				if e.to == v {
					cost.Add(cost, new(big.Int).Lsh(big.NewInt(1), uint(e.x)))
					found = true
					break
				}
			}
			if !found {
				fmt.Fprintf(os.Stderr, "case %d: invalid path\n", idx+1)
				os.Exit(1)
			}
		}
		if cost.Cmp(dist) != 0 {
			fmt.Fprintf(os.Stderr, "case %d: path is not shortest\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
