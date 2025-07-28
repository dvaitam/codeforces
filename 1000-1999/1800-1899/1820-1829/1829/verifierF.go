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

func solveF(n int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = len(adj[i])
	}
	leafNeighbors := make([]int, n+1)
	for i := 1; i <= n; i++ {
		cnt := 0
		for _, v := range adj[i] {
			if deg[v] == 1 {
				cnt++
			}
		}
		leafNeighbors[i] = cnt
	}
	x := 0
	y := 0
	for i := 1; i <= n; i++ {
		if leafNeighbors[i] > 0 {
			x++
			y = leafNeighbors[i]
		}
	}
	return fmt.Sprintf("%d %d", x, y)
}

func genSnowflake() (int, [][2]int, int, int) {
	x := rand.Intn(5) + 2
	y := rand.Intn(5) + 2
	n := 1 + x + x*y
	edges := make([][2]int, 0, x+x*y)
	center := 1
	curr := 2
	for i := 0; i < x; i++ {
		v := curr
		curr++
		edges = append(edges, [2]int{center, v})
		for j := 0; j < y; j++ {
			w := curr
			curr++
			edges = append(edges, [2]int{v, w})
		}
	}
	return n, edges, x, y
}

func genTestsF() ([]string, string) {
	const t = 100
	rand.Seed(1)
	var input strings.Builder
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n, edges, x, y := genSnowflake()
		fmt.Fprintf(&input, "%d %d\n", n, len(edges))
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		expected[i] = fmt.Sprintf("%d %d", x, y)
	}
	return expected, input.String()
}

func runBinary(path, in string) ([]string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(&out)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	expected, input := genTestsF()
	lines, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if len(lines) != len(expected) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(expected), len(lines))
		os.Exit(1)
	}
	for i, exp := range expected {
		if lines[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
