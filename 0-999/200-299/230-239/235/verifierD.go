package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveCase(n int, edges [][2]int) float64 {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n)
	pd := make([]bool, n)
	dArr := make([]int, n)
	uArr := make([]int, n)
	var cycle []int
	var findCycle func(u, parent int, path *[]int) []int
	findCycle = func(u, parent int, path *[]int) []int {
		visited[u] = true
		*path = append(*path, u)
		for _, v := range adj[u] {
			if v == parent {
				continue
			}
			if visited[v] {
				idx := 0
				for i, x := range *path {
					if x == v {
						idx = i
						break
					}
				}
				cyc := make([]int, len(*path)-idx)
				copy(cyc, (*path)[idx:])
				return cyc
			}
			if res := findCycle(v, u, path); len(res) > 0 {
				return res
			}
		}
		*path = (*path)[:len(*path)-1]
		return nil
	}
	path := make([]int, 0, n)
	cycle = findCycle(0, -1, &path)
	lenC := len(cycle)
	for i, node := range cycle {
		pd[node] = true
		dArr[node] = i
		uArr[node] = 0
	}
	maxA := 2*n + 3
	a := make([]float64, maxA)
	for i := 1; i < maxA; i++ {
		a[i] = 1.0 / float64(i)
	}
	var dfs2 func(u, parent, compIdx, depth int)
	dfs2 = func(u, parent, compIdx, depth int) {
		dArr[u] = compIdx
		uArr[u] = depth
		for _, v := range adj[u] {
			if v == parent || pd[v] {
				continue
			}
			dfs2(v, u, compIdx, depth+1)
		}
	}
	for i, node := range cycle {
		dfs2(node, -1, i, 0)
	}
	var ans float64
	var dfs3 func(u, parent, depth int)
	dfs3 = func(u, parent, depth int) {
		ans += a[depth]
		for _, v := range adj[u] {
			if v != parent && !pd[v] {
				dfs3(v, u, depth+1)
			}
		}
	}
	for i := 0; i < n; i++ {
		root := cycle[dArr[i]]
		pd[root] = false
		dfs3(i, -1, 1)
		pd[root] = true
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if dArr[i] != dArr[j] {
				k1 := uArr[i] + uArr[j]
				diff := abs(dArr[i]-dArr[j]) + 1
				k2 := diff
				k3 := lenC - diff + 2
				ans += a[k1+k2] + a[k1+k3] - a[k1+lenC]
			}
		}
	}
	return ans
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 3 // 3..6
	edges := make([][2]int, n)
	perm := rng.Perm(n)
	for i := 0; i < n; i++ {
		u := perm[i] + 1
		v := perm[(i+1)%n] + 1
		edges[i] = [2]int{u, v}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	input := sb.String()
	expect := fmt.Sprintf("%.11f", solveCase(n, edges))
	return input, expect
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
