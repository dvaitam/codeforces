package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// Embedded solver for 1610I - Mashtali and Hagh Trees (game theory on tree).
func solveI(input string) string {
	reader := strings.NewReader(input)
	readInt := func() int {
		var x int
		fmt.Fscan(reader, &x)
		return x
	}

	n := readInt()
	if n == 0 {
		return ""
	}

	head := make([]int, n+1)
	next := make([]int, 2*n)
	to := make([]int, 2*n)
	edgeCount := 1

	addEdge := func(u, v int) {
		to[edgeCount] = v
		next[edgeCount] = head[u]
		head[u] = edgeCount
		edgeCount++
	}

	for i := 0; i < n-1; i++ {
		u := readInt()
		v := readInt()
		addEdge(u, v)
		addEdge(v, u)
	}

	parent := make([]int, n+1)
	order := make([]int, 0, n)
	order = append(order, 1)

	visited := make([]bool, n+1)
	visited[1] = true

	headQueue := 0
	for headQueue < len(order) {
		u := order[headQueue]
		headQueue++
		for e := head[u]; e != 0; e = next[e] {
			v := to[e]
			if !visited[v] {
				visited[v] = true
				parent[v] = u
				order = append(order, v)
			}
		}
	}

	g := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		u := order[i]
		for e := head[u]; e != 0; e = next[e] {
			v := to[e]
			if v != parent[u] {
				g[u] ^= (g[v] + 1)
			}
		}
	}

	inH := make([]bool, n+1)
	inH[1] = true

	ans := make([]byte, n)
	X := g[1]

	for k := 1; k <= n; k++ {
		curr := k
		for !inH[curr] {
			inH[curr] = true
			X ^= g[curr] ^ (g[curr] + 1) ^ 1
			curr = parent[curr]
		}
		if X != 0 {
			ans[k-1] = '1'
		} else {
			ans[k-1] = '2'
		}
	}

	return string(ans)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// Generate a random tree on n vertices.
func genTree(rng *rand.Rand, n int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, i))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(10) + 1
		input := genTree(rng, n)
		exp := solveI(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", t, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
