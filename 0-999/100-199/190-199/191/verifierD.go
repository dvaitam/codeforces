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

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	for u := 1; u <= n; u++ {
		deg[u] = len(adj[u])
	}
	type frame struct{ u, parent, next int }
	state := make([]int8, n+1)
	pos := make([]int, n+1)
	for i := range pos {
		pos[i] = -1
	}
	stack := []int{}
	cycleMark := make([]bool, n+1)
	cycles := 0
	dfsStack := []frame{{1, -1, 0}}
	state[1] = 1
	pos[1] = 0
	stack = append(stack, 1)
	for len(dfsStack) > 0 {
		f := &dfsStack[len(dfsStack)-1]
		u := f.u
		if f.next < len(adj[u]) {
			v := adj[u][f.next]
			f.next++
			if v == f.parent {
				continue
			}
			if state[v] == 0 {
				state[v] = 1
				pos[v] = len(stack)
				stack = append(stack, v)
				dfsStack = append(dfsStack, frame{v, u, 0})
			} else if state[v] == 1 {
				cycles++
				for i := pos[v]; i < len(stack); i++ {
					cycleMark[stack[i]] = true
				}
			}
		} else {
			state[u] = 2
			dfsStack = dfsStack[:len(dfsStack)-1]
			stack = stack[:len(stack)-1]
			pos[u] = -1
		}
	}
	odd := 0
	for u := 1; u <= n; u++ {
		d := deg[u]
		if cycleMark[u] {
			d -= 2
		}
		if d%2 == 1 {
			odd++
		}
	}
	minLines := cycles + odd/2
	maxLines := len(edges)
	return fmt.Sprintf("%d %d", minLines, maxLines)
}

func generateTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func addCycle(rng *rand.Rand, edges [][2]int, n int) [][2]int {
	u := rng.Intn(n) + 1
	v := rng.Intn(n) + 1
	for v == u {
		v = rng.Intn(n) + 1
	}
	edges = append(edges, [2]int{u, v})
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	edges := generateTree(rng, n)
	if rng.Intn(2) == 0 && n >= 3 {
		edges = addCycle(rng, edges, n)
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&in, "%d %d\n", e[0], e[1])
	}
	exp := solve(n, edges)
	return in.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
