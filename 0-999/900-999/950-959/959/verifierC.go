package main

import (
	"bufio"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Graph struct {
	n   int
	adj [][]int
}

func (g *Graph) isTree() bool {
	if g.n == 0 {
		return false
	}
	if len(g.adj) != g.n+1 {
		return false
	}
	visited := make([]bool, g.n+1)
	count := 0
	var dfs func(int, int)
	dfs = func(u, p int) {
		visited[u] = true
		count++
		for _, v := range g.adj[u] {
			if v == p {
				continue
			}
			if !visited[v] {
				dfs(v, u)
			}
		}
	}
	dfs(1, -1)
	return count == g.n
}

func (g *Graph) mahmoudAns() int {
	depth := make([]int, g.n+1)
	var dfs func(int, int, int)
	dfs = func(u, p, d int) {
		depth[u] = d
		for _, v := range g.adj[u] {
			if v != p {
				dfs(v, u, d+1)
			}
		}
	}
	dfs(1, -1, 0)
	even, odd := 0, 0
	for i := 1; i <= g.n; i++ {
		if depth[i]%2 == 0 {
			even++
		} else {
			odd++
		}
	}
	if even < odd {
		return even
	}
	return odd
}

func (g *Graph) trueMVC() int {
	// dp[u][0]: u is not in cover
	// dp[u][1]: u is in cover
	dp := make([][2]int, g.n+1)
	var dfs func(int, int)
	dfs = func(u, p int) {
		dp[u][0] = 0
		dp[u][1] = 1
		for _, v := range g.adj[u] {
			if v == p {
				continue
			}
			dfs(v, u)
			dp[u][0] += dp[v][1]
			val := dp[v][0]
			if dp[v][1] < val {
				val = dp[v][1]
			}
			dp[u][1] += val
		}
	}
	dfs(1, -1)
	if dp[1][0] < dp[1][1] {
		return dp[1][0]
	}
	return dp[1][1]
}

func check(n int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)

	// First tree
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	firstToken := scanner.Text()
	
	if n < 6 {
		if firstToken != "-1" {
			return fmt.Errorf("expected -1 for n < 6, got %s", firstToken)
		}
	} else {
		if firstToken == "-1" {
			return fmt.Errorf("unexpected -1 for n >= 6")
		}
		
		adj := make([][]int, n+1)
		var u, v int
		fmt.Sscan(firstToken, &u)
		if !scanner.Scan() { return fmt.Errorf("incomplete edge") }
		fmt.Sscan(scanner.Text(), &v)
		if u < 1 || u > n || v < 1 || v > n { return fmt.Errorf("node index out of bounds: %d %d", u, v) }
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		
		for i := 0; i < n-2; i++ {
			if !scanner.Scan() { return fmt.Errorf("incomplete tree 1") }
			fmt.Sscan(scanner.Text(), &u)
			if !scanner.Scan() { return fmt.Errorf("incomplete tree 1") }
			fmt.Sscan(scanner.Text(), &v)
			if u < 1 || u > n || v < 1 || v > n { return fmt.Errorf("node index out of bounds: %d %d", u, v) }
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		
		g1 := &Graph{n: n, adj: adj}
		if !g1.isTree() {
			return fmt.Errorf("first graph is not a tree")
		}
		
		mah := g1.mahmoudAns()
		trueAns := g1.trueMVC()
		if mah == trueAns {
			return fmt.Errorf("first tree: Mahmoud's algorithm (%d) gave correct answer (%d), expected incorrect", mah, trueAns)
		}
	}

	// Second tree
	if !scanner.Scan() { return fmt.Errorf("missing second tree") }
	t2Start := scanner.Text()
	if t2Start == "-1" {
	    return fmt.Errorf("second tree should always exist")
	}
	
	adj2 := make([][]int, n+1)
	var u, v int
	fmt.Sscan(t2Start, &u)
	if !scanner.Scan() { return fmt.Errorf("incomplete edge t2") }
	fmt.Sscan(scanner.Text(), &v)
	if u < 1 || u > n || v < 1 || v > n { return fmt.Errorf("node index out of bounds: %d %d", u, v) }
	adj2[u] = append(adj2[u], v)
	adj2[v] = append(adj2[v], u)
	
	for i := 0; i < n-2; i++ {
		if !scanner.Scan() { return fmt.Errorf("incomplete tree 2") }
		fmt.Sscan(scanner.Text(), &u)
		if !scanner.Scan() { return fmt.Errorf("incomplete tree 2") }
		fmt.Sscan(scanner.Text(), &v)
		if u < 1 || u > n || v < 1 || v > n { return fmt.Errorf("node index out of bounds: %d %d", u, v) }
		adj2[u] = append(adj2[u], v)
		adj2[v] = append(adj2[v], u)
	}
	
	g2 := &Graph{n: n, adj: adj2}
	if !g2.isTree() {
		return fmt.Errorf("second graph is not a tree")
	}
	
	mah2 := g2.mahmoudAns()
	trueAns2 := g2.trueMVC()
	if mah2 != trueAns2 {
		return fmt.Errorf("second tree: Mahmoud's algorithm (%d) gave incorrect answer (%d), expected correct", mah2, trueAns2)
	}

	return nil
}

func genCaseC(rng *rand.Rand) string {
	if rng.Float32() < 0.2 {
	    return fmt.Sprintf("%d\n", rng.Intn(4) + 2) // 2..5
	}
	return fmt.Sprintf("%d\n", rng.Intn(20) + 6) // 6..25
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 50; i++ {
		in := genCaseC(rng)
		n := 0
		fmt.Sscan(strings.TrimSpace(in), &n)
		
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed (n=%d): %v\n", i+1, n, err)
			os.Exit(1)
		}
		
		if err := check(n, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed (n=%d): %v\nOutput:\n%s\n", i+1, n, err, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}