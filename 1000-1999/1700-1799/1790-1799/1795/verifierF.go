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

// Brute-force solver for 1795F (small n).
// Simulate the game: chips move in round-robin order 1..k,1..k,...
// Each chip greedily extends into unvisited territory via BFS/DFS maximizing moves.
// We try all possible move sequences via DFS to find the maximum number of moves.

func solveCase(n int, edges [][2]int, k int, chips []int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}

	// BFS/backtracking: try all possible moves
	colored := make([]bool, n+1)
	pos := make([]int, k+1)
	for i := 1; i <= k; i++ {
		pos[i] = chips[i-1]
		colored[chips[i-1]] = true
	}

	var best int
	var dfs func(move int)
	dfs = func(move int) {
		if move > best {
			best = move
		}
		chipIdx := (move % k) + 1 // 1-indexed chip for this move (0-indexed move)
		v := pos[chipIdx]
		for _, u := range adj[v] {
			if !colored[u] {
				colored[u] = true
				oldPos := pos[chipIdx]
				pos[chipIdx] = u
				dfs(move + 1)
				pos[chipIdx] = oldPos
				colored[u] = false
			}
		}
	}
	dfs(0)
	return best
}

func genTest(rng *rand.Rand) (string, int, [][2]int, int, []int) {
	n := rng.Intn(8) + 1
	// Generate random tree
	perm := rng.Perm(n)
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		u := perm[rng.Intn(i)] + 1
		v := perm[i] + 1
		edges[i-1] = [2]int{u, v}
	}
	k := rng.Intn(n) + 1
	chipPerm := rng.Perm(n)
	chips := make([]int, k)
	for i := 0; i < k; i++ {
		chips[i] = chipPerm[i] + 1
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for i, c := range chips {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c))
	}
	sb.WriteByte('\n')
	return sb.String(), n, edges, k, chips
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for tc := 0; tc < 200; tc++ {
		input, n, edges, k, chips := genTest(rng)
		expectedVal := solveCase(n, edges, k, chips)
		expectedOut := strconv.Itoa(expectedVal)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
