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

func computeColors(n, k int, adj [][]int, leafColor []int) []int {
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, u := range adj[v] {
			if u != parent[v] {
				parent[u] = v
				stack = append(stack, u)
			}
		}
	}
	color := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if v != 1 && len(adj[v]) == 1 {
			color[v] = leafColor[v]
		} else {
			blue, red := 0, 0
			for _, u := range adj[v] {
				if u == parent[v] {
					continue
				}
				if color[u] == 1 {
					blue++
				} else {
					red++
				}
			}
			if blue-red >= k {
				color[v] = 1
			} else {
				color[v] = 0
			}
		}
	}
	return color
}

func solveCase(n, k int, edges [][2]int, s []int, queries [][]int) []int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	leafColor := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if s[i] == -1 {
			leafColor[i] = 0
		} else {
			leafColor[i] = s[i]
		}
	}
	curK := k
	color := computeColors(n, curK, adj, leafColor)
	var res []int
	for _, q := range queries {
		switch q[0] {
		case 1:
			v := q[1]
			color = computeColors(n, curK, adj, leafColor)
			res = append(res, color[v])
		case 2:
			v := q[1]
			c := q[2]
			leafColor[v] = c
			color = computeColors(n, curK, adj, leafColor)
		case 3:
			curK = q[1]
			color = computeColors(n, curK, adj, leafColor)
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	k := rng.Intn(2*n+1) - n
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		u := i + 2
		v := rng.Intn(i+1) + 1
		edges[i] = [2]int{u, v}
	}
	s := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if rng.Intn(2) == 0 {
			s[i] = -1
		} else {
			s[i] = rng.Intn(2)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", s[i]))
	}
	sb.WriteByte('\n')
	q := rng.Intn(5) + 1
	queries := make([][]int, q)
	sb.WriteString(fmt.Sprintf("%d\n", q))
	leaves := []int{}
	for i := 2; i <= n; i++ {
		if s[i] != -1 {
			leaves = append(leaves, i)
		}
	}
	if len(leaves) == 0 {
		leaves = append(leaves, 2)
	}
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		if t == 1 {
			v := rng.Intn(n) + 1
			queries[i] = []int{1, v}
			sb.WriteString(fmt.Sprintf("1 %d\n", v))
		} else if t == 2 {
			v := leaves[rng.Intn(len(leaves))]
			c := rng.Intn(2)
			queries[i] = []int{2, v, c}
			sb.WriteString(fmt.Sprintf("2 %d %d\n", v, c))
		} else {
			h := rng.Intn(2*n+1) - n
			queries[i] = []int{3, h}
			sb.WriteString(fmt.Sprintf("3 %d\n", h))
		}
	}
	answers := solveCase(n, k, edges, s, queries)
	var exp strings.Builder
	for i, v := range answers {
		if i > 0 {
			exp.WriteByte('\n')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
