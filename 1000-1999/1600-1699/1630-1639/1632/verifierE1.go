package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct {
	delta int
	depth int
	dist  int
}

func solveE1(n int, edges [][2]int) []int {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = -1
		}
		queue := []int{i}
		dist[i][i] = 0
		for head := 0; head < len(queue); head++ {
			cur := queue[head]
			for _, nb := range adj[cur] {
				if dist[i][nb] == -1 {
					dist[i][nb] = dist[i][cur] + 1
					queue = append(queue, nb)
				}
			}
		}
	}
	depth := dist[0]
	e1 := 0
	for _, d := range depth {
		if d > e1 {
			e1 = d
		}
	}
	ans := make([]int, n+1)
	for x := 1; x <= n; x++ {
		ans[x] = e1
	}
	arr := make([]pair, n)
	prefix := make([]int, n+1)
	suffix := make([]int, n+1)
	for v := 1; v < n; v++ {
		for u := 0; u < n; u++ {
			arr[u] = pair{
				delta: depth[u] - dist[v][u],
				depth: depth[u],
				dist:  dist[v][u],
			}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].delta < arr[j].delta })
		prefix[0] = 0
		for i := 0; i < n; i++ {
			if arr[i].depth > prefix[i] {
				prefix[i+1] = arr[i].depth
			} else {
				prefix[i+1] = prefix[i]
			}
		}
		suffix[n] = 0
		for i := n - 1; i >= 0; i-- {
			if arr[i].dist > suffix[i+1] {
				suffix[i] = arr[i].dist
			} else {
				suffix[i] = suffix[i+1]
			}
		}
		for x := 1; x <= n; x++ {
			idx := sort.Search(n, func(i int) bool { return arr[i].delta > x })
			val := prefix[idx]
			tmp := x + suffix[idx]
			if tmp > val {
				val = tmp
			}
			if val < ans[x] {
				ans[x] = val
			}
		}
	}
	return ans[1:]
}

func randTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges = append(edges, [2]int{i, p})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(8) + 2
	edges := randTree(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	ans := solveE1(n, edges)
	return sb.String(), ans
}

func runCase(bin, input string, exp []int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(parts))
	}
	for i, p := range parts {
		var v int
		if _, err := fmt.Sscan(p, &v); err != nil {
			return fmt.Errorf("bad int at pos %d: %v", i+1, err)
		}
		if v != exp[i] {
			return fmt.Errorf("pos %d expected %d got %d", i+1, exp[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
