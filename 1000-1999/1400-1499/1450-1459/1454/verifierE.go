package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveCaseE(n int, edges [][2]int) int64 {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	removed := make([]bool, n+1)
	queue := make([]int, 0)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	qi := 0
	for qi < len(queue) {
		u := queue[qi]
		qi++
		removed[u] = true
		for _, v := range adj[u] {
			if removed[v] {
				continue
			}
			deg[v]--
			if deg[v] == 1 {
				queue = append(queue, v)
			}
		}
	}
	iscycle := make([]bool, n+1)
	cycleRoots := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !removed[i] {
			iscycle[i] = true
			cycleRoots = append(cycleRoots, i)
		}
	}
	visited := make([]bool, n+1)
	var sumS2 int64
	for _, u := range cycleRoots {
		var cnt int64 = 1
		stack := make([]int, 0)
		for _, v := range adj[u] {
			if !iscycle[v] && !visited[v] {
				visited[v] = true
				stack = append(stack, v)
				cnt++
			}
		}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, w := range adj[v] {
				if !iscycle[w] && !visited[w] {
					visited[w] = true
					stack = append(stack, w)
					cnt++
				}
			}
		}
		sumS2 += cnt * cnt
	}
	nn := int64(n)
	t1 := nn * (nn - 1) / 2
	t2 := (nn*nn - sumS2) / 2
	return t1 + t2
}

func generateGraph(n int, r *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n)
	// build random tree first
	for i := 2; i <= n; i++ {
		v := r.Intn(i-1) + 1
		edges = append(edges, [2]int{i, v})
	}
	// add one extra edge to create single cycle
	for {
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		if a == b {
			continue
		}
		exists := false
		for _, e := range edges {
			if (e[0] == a && e[1] == b) || (e[0] == b && e[1] == a) {
				exists = true
				break
			}
		}
		if !exists {
			edges = append(edges, [2]int{a, b})
			break
		}
	}
	return edges
}

func generateTests() ([]int, [][][2]int, string) {
	const t = 100
	r := rand.New(rand.NewSource(5))
	ns := make([]int, t)
	graphs := make([][][2]int, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(8) + 3
		edges := generateGraph(n, r)
		ns[i] = n
		graphs[i] = edges
		fmt.Fprintf(&sb, "%d\n", n)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
	}
	return ns, graphs, sb.String()
}

func verify(ns []int, graphs [][][2]int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for idx, n := range ns {
		if !scanner.Scan() {
			return fmt.Errorf("case %d: missing output", idx+1)
		}
		var ans int64
		fmt.Sscan(scanner.Text(), &ans)
		expected := solveCaseE(n, graphs[idx])
		if ans != expected {
			return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, ans)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	ns, graphs, input := generateTests()
	out, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if err := verify(ns, graphs, out); err != nil {
		fmt.Fprintln(os.Stderr, "verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed for problem E")
}
