package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type edge struct{ x, y int }

func solveE(n int, edges []edge) (int, []int) {
	adj := make([][]bool, n+1)
	for i := range adj {
		adj[i] = make([]bool, n+1)
	}
	for _, e := range edges {
		adj[e.x][e.y] = true
		adj[e.y][e.x] = true
	}
	visited := make([]bool, n+1)
	sizes := []int{}
	for v := 1; v <= n; v++ {
		if visited[v] {
			continue
		}
		queue := []int{v}
		visited[v] = true
		size := 0
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			size++
			for u := 1; u <= n; u++ {
				if !visited[u] && !adj[cur][u] && u != cur {
					visited[u] = true
					queue = append(queue, u)
				}
			}
		}
		sizes = append(sizes, size)
	}
	sort.Ints(sizes)
	return len(sizes), sizes
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(5)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(8) + 2
		m := rand.Intn(n*(n-1)/2 + 1)
		exists := make(map[[2]int]bool)
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			for {
				x := rand.Intn(n) + 1
				y := rand.Intn(n) + 1
				if x == y || exists[[2]int{x, y}] || exists[[2]int{y, x}] {
					continue
				}
				exists[[2]int{x, y}] = true
				edges[i] = edge{x, y}
				break
			}
		}
		expectedK, expectedSizes := solveE(n, edges)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e.x, e.y)
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed to run binary:", err)
			os.Exit(1)
		}
		parts := strings.Fields(string(out))
		if len(parts) < 1 {
			fmt.Printf("Test %d produced no output\n", t+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(parts[0])
		if err != nil || k != expectedK {
			fmt.Printf("Test %d failed: expected %d components got %s\n", t+1, expectedK, parts[0])
			os.Exit(1)
		}
		if len(parts)-1 != k {
			fmt.Printf("Test %d failed: expected %d sizes got %d values\n", t+1, k, len(parts)-1)
			os.Exit(1)
		}
		for i := 0; i < k; i++ {
			val, err := strconv.Atoi(parts[i+1])
			if err != nil || val != expectedSizes[i] {
				fmt.Printf("Test %d size %d mismatch: expected %d got %s\n", t+1, i+1, expectedSizes[i], parts[i+1])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
