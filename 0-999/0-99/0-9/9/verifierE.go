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

type pair struct{ u, v int }

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n, m int, u, v []int) string {
	dsuP := make([]int, n)
	for i := range dsuP {
		dsuP[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if dsuP[x] != x {
			dsuP[x] = find(dsuP[x])
		}
		return dsuP[x]
	}
	unite := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra != rb {
			dsuP[rb] = ra
		}
	}
	for i := 0; i < m; i++ {
		unite(u[i], v[i])
	}
	compVerts := make([][]int, n)
	compEdges := make([]int, n)
	for i := 0; i < n; i++ {
		compVerts[find(i)] = append(compVerts[find(i)], i)
	}
	for i := 0; i < m; i++ {
		compEdges[find(u[i])]++
	}
	for r, verts := range compVerts {
		if len(verts) == 0 {
			continue
		}
		if compEdges[r]-len(verts)+1 > 1 {
			return "NO"
		}
	}
	deg := make([]int, n)
	inCycle := make([]bool, n)
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		a, b := u[i], v[i]
		if a == b {
			inCycle[a] = true
		} else {
			deg[a]++
			deg[b]++
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
	}
	queue := []int{}
	for i := 0; i < n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, y := range adj[x] {
			if deg[y] > 0 {
				deg[y]--
				if deg[y] == 1 {
					queue = append(queue, y)
				}
			}
		}
		deg[x] = 0
	}
	for i := 0; i < n; i++ {
		if deg[i] > 0 {
			inCycle[i] = true
		}
	}
	var add []pair
	for i := 0; i < n; i++ {
		if !inCycle[i] {
			add = append(add, pair{i, i})
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(add)))
	for _, p := range add {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.u+1, p.v+1))
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	maxEdges := n + rng.Intn(4)
	m := rng.Intn(maxEdges + 1)
	u := make([]int, m)
	v := make([]int, m)
	for i := 0; i < m; i++ {
		u[i] = rng.Intn(n)
		v[i] = rng.Intn(n)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", u[i]+1, v[i]+1))
	}
	expect := solveCase(n, m, u, v)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
