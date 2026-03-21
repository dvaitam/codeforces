package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Correct solver adapted from the accepted solution.
// Returns the expected count and set of nodes to flip.
func solve(n int, edges [][2]int, init, goal []int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	var ans []int
	var dfs func(u, p, oddFlips, evenFlips, depth int)
	dfs = func(u, p, oddFlips, evenFlips, depth int) {
		cur := init[u]
		if depth%2 == 1 {
			cur ^= oddFlips
		} else {
			cur ^= evenFlips
		}

		if cur != goal[u] {
			ans = append(ans, u)
			if depth%2 == 1 {
				oddFlips ^= 1
			} else {
				evenFlips ^= 1
			}
		}

		for _, v := range adj[u] {
			if v != p {
				dfs(v, u, oddFlips, evenFlips, depth+1)
			}
		}
	}

	dfs(1, 0, 0, 0, 0)
	return len(ans)
}

// Verify that flipping the given set of nodes transforms init into goal.
func verify(n int, edges [][2]int, init, goal []int, flips []int) error {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	depth := make([]int, n+1)
	var dfs func(u, p int)
	dfs = func(u, p int) {
		for _, v := range adj[u] {
			if v != p {
				depth[v] = depth[u] + 1
				dfs(v, u)
			}
		}
	}
	dfs(1, 0)

	cur := make([]int, n+1)
	copy(cur, init)

	flipSet := make(map[int]bool)
	for _, f := range flips {
		if f < 1 || f > n {
			return fmt.Errorf("flip node %d out of range", f)
		}
		if flipSet[f] {
			return fmt.Errorf("duplicate flip node %d", f)
		}
		flipSet[f] = true
	}

	// For each flip node, it flips itself and all nodes at even distance from it in the subtree.
	// Actually, flipping node u flips all nodes v such that depth[v] and depth[u] have the same parity.
	// Wait, need to understand the problem properly.
	// Problem 429A: Each operation selects a node u, and flips all nodes at even distance from u.
	// "even distance" in the tree means same parity of depth.
	// So flipping u flips all nodes with depth%2 == depth[u]%2.

	// Count flips affecting even-depth and odd-depth nodes
	evenFlips := 0
	oddFlips := 0
	for _, f := range flips {
		if depth[f]%2 == 0 {
			evenFlips++
		} else {
			oddFlips++
		}
	}

	// Each even-depth flip toggles all even-depth nodes
	// Each odd-depth flip toggles all odd-depth nodes
	for i := 1; i <= n; i++ {
		val := init[i]
		if depth[i]%2 == 0 {
			val ^= evenFlips % 2
		} else {
			val ^= oddFlips % 2
		}
		if val != goal[i] {
			return fmt.Errorf("node %d: expected %d got %d after flips", i, goal[i], val)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", i, p))
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('0' + byte(rng.Intn(2)))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('0' + byte(rng.Intn(2)))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) (int, [][2]int, []int, []int) {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(r, &edges[i][0], &edges[i][1])
	}
	init := make([]int, n+1)
	goal := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &init[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &goal[i])
	}
	return n, edges, init, goal
}

func runCase(exe, input string) error {
	n, edges, init, goal := parseInput(input)
	expectedCount := solve(n, edges, init, goal)

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	// Parse candidate output
	scanner := strings.NewReader(strings.TrimSpace(out.String()))
	var cnt int
	if _, err := fmt.Fscan(scanner, &cnt); err != nil {
		return fmt.Errorf("cannot parse count: %v", err)
	}
	if cnt != expectedCount {
		return fmt.Errorf("expected count %d got %d", expectedCount, cnt)
	}

	flips := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		tok := ""
		if _, err := fmt.Fscan(scanner, &tok); err != nil {
			return fmt.Errorf("cannot parse flip %d: %v", i, err)
		}
		v, err := strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("cannot parse flip value: %v", err)
		}
		flips[i] = v
	}

	if err := verify(n, edges, init, goal, flips); err != nil {
		return fmt.Errorf("verification failed: %v", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
