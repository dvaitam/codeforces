package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct{ u, v int }

type caseE struct {
	n     int
	edges []edge
	root  int
	dist  [][]int
}

func genCase() caseE {
	n := rand.Intn(7) + 2 // 2..8
	edges := make([]edge, 0, n-1)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
		parent[i] = p
	}
	root := rand.Intn(n) + 1
	// compute distances via BFS from root
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{root}
	dist[root] = 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return caseE{n: n, edges: edges, root: root, dist: [][]int{dist}}
}

func runCase(bin string, c caseE) error {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	inW := bufio.NewWriter(stdin)
	outR := bufio.NewReader(stdout)

	fmt.Fprintln(inW, c.n)
	for _, e := range c.edges {
		fmt.Fprintf(inW, "%d %d\n", e.u, e.v)
	}
	inW.Flush()

	dist := c.dist[0]
	queries := 0
	for {
		line, err := outR.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queries++
			if queries > 2*c.n {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			parts := strings.Fields(line)
			if len(parts) != 3 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid query: %s", line)
			}
			u, _ := strconv.Atoi(parts[1])
			v, _ := strconv.Atoi(parts[2])
			if u < 1 || u > c.n || v < 1 || v > c.n {
				cmd.Process.Kill()
				return fmt.Errorf("query out of range")
			}
			du := dist[u]
			dv := dist[v]
			ans := u
			if dv < du {
				ans = v
			}
			fmt.Fprintln(inW, ans)
			inW.Flush()
		} else if strings.HasPrefix(line, "!") {
			parts := strings.Fields(line)
			if len(parts) != 2 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid answer: %s", line)
			}
			guess, _ := strconv.Atoi(parts[1])
			if guess != c.root {
				cmd.Process.Kill()
				return fmt.Errorf("wrong answer: expected %d got %d", c.root, guess)
			}
			stdin.Close()
			cmd.Wait()
			return nil
		} else {
			cmd.Process.Kill()
			return fmt.Errorf("invalid output: %s", line)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		c := genCase()
		if err := runCase(bin, c); err != nil {
			fmt.Printf("case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
