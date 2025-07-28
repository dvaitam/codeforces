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

func solveCase(n, m int, v, tarr []int64, edges [][2]int) string {
	diff := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		diff[i] = tarr[i] - v[i]
		sum += diff[i]
	}
	adj := make([][]int, n)
	for _, e := range edges {
		u, w := e[0], e[1]
		adj[u] = append(adj[u], w)
		adj[w] = append(adj[w], u)
	}
	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	bip := true
	queue := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if color[i] != -1 {
			continue
		}
		color[i] = 0
		queue = append(queue, i)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, w := range adj[u] {
				if color[w] == -1 {
					color[w] = color[u] ^ 1
					queue = append(queue, w)
				} else if color[w] == color[u] {
					bip = false
				}
			}
		}
	}
	if !bip {
		if sum%2 == 0 {
			return "YES\n"
		}
		return "NO\n"
	}
	var diffSum int64
	for i := 0; i < n; i++ {
		if color[i] == 0 {
			diffSum += diff[i]
		} else {
			diffSum -= diff[i]
		}
	}
	if diffSum == 0 {
		return "YES\n"
	}
	return "NO\n"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateConnectedGraph(rng *rand.Rand, n int) [][2]int {
	// start with tree
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		u := rng.Intn(i)
		edges = append(edges, [2]int{u, i})
	}
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	exist := make(map[[2]int]bool)
	for _, e := range edges {
		a, b := e[0], e[1]
		if a > b {
			a, b = b, a
		}
		exist[[2]int{a, b}] = true
	}
	for len(edges) < m {
		a := rng.Intn(n)
		b := rng.Intn(n)
		if a == b {
			continue
		}
		x, y := a, b
		if x > y {
			x, y = y, x
		}
		if exist[[2]int{x, y}] {
			continue
		}
		exist[[2]int{x, y}] = true
		edges = append(edges, [2]int{a, b})
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	edges := generateConnectedGraph(rng, n)
	m := len(edges)
	v := make([]int64, n)
	tarr := make([]int64, n)
	for i := 0; i < n; i++ {
		v[i] = int64(rng.Intn(21) - 10)
		tarr[i] = int64(rng.Intn(21) - 10)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", tarr[i])
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	expect := solveCase(n, m, v, tarr, edges)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
