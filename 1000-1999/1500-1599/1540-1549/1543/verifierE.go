package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

func genEdges(n int) [][2]int {
	N := 1 << n
	edges := make([][2]int, 0, n*N/2)
	for i := 0; i < N; i++ {
		for b := 0; b < n; b++ {
			j := i ^ (1 << b)
			if i < j {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	return edges
}

func permEdges(edges [][2]int, perm []int) [][2]int {
	res := make([][2]int, len(edges))
	for i, e := range edges {
		res[i] = [2]int{perm[e[0]], perm[e[1]]}
	}
	return res
}

func validate(n int, edges [][2]int, scanner *bufio.Scanner) error {
	N := 1 << n
	// Read permutation P
	P := make([]int, N)
	usedP := make([]bool, N)
	for i := 0; i < N; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected permutation value, got EOF")
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid permutation value %q: %v", scanner.Text(), err)
		}
		if val < 0 || val >= N {
			return fmt.Errorf("permutation value %d out of range [0, %d)", val, N)
		}
		if usedP[val] {
			return fmt.Errorf("duplicate permutation value %d", val)
		}
		usedP[val] = true
		P[i] = val
	}

	// Build adjacency for input graph
	adj := make([][]int, N)
	isEdge := make([]map[int]bool, N)
	for i := 0; i < N; i++ {
		isEdge[i] = make(map[int]bool)
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		isEdge[u][v] = true
		isEdge[v][u] = true
	}

	// Verify P is a valid hypercube mapping
	// If (i, j) is an edge in simple hypercube, (P[i], P[j]) must be an edge in input graph
	for i := 0; i < N; i++ {
		for b := 0; b < n; b++ {
			j := i ^ (1 << b)
			if i < j {
				if !isEdge[P[i]][P[j]] {
					return fmt.Errorf("edge (%d, %d) in simple hypercube maps to non-existent edge (%d, %d) in input", i, j, P[i], P[j])
				}
			}
		}
	}

	// Coloring
	isPower := (n > 0) && (n&(n-1) == 0)
	if !scanner.Scan() {
		return fmt.Errorf("expected coloring start, got EOF")
	}
	firstColorStr := scanner.Text()
	if !isPower {
		if firstColorStr != "-1" {
			return fmt.Errorf("expected -1 coloring for n=%d, got %q", n, firstColorStr)
		}
		return nil
	}

	if firstColorStr == "-1" {
		return fmt.Errorf("expected valid coloring for n=%d, got -1", n)
	}

	colors := make([]int, N)
	c, err := strconv.Atoi(firstColorStr)
	if err != nil {
		return fmt.Errorf("invalid color value %q: %v", firstColorStr, err)
	}
	colors[0] = c
	for i := 1; i < N; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected color for vertex %d, got EOF", i)
		}
		colors[i], err = strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid color value %q for vertex %d: %v", scanner.Text(), i, err)
		}
	}

	// Check coloring properties
	for i := 0; i < N; i++ {
		if colors[i] < 0 || colors[i] >= n {
			return fmt.Errorf("color %d of vertex %d out of range [0, %d)", colors[i], i, n)
		}
	}

	for i := 0; i < N; i++ {
		seen := make([]bool, n)
		for _, v := range adj[i] {
			col := colors[v]
			if seen[col] {
				return fmt.Errorf("vertex %d has multiple neighbors with color %d", i, col)
			}
			seen[col] = true
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(42)

	testNs := []int{1, 2, 3, 4, 8}
	var input bytes.Buffer
	fmt.Fprintln(&input, len(testNs))
	type tc struct {
		n     int
		edges [][2]int
	}
	tests := make([]tc, len(testNs))

	for i, n := range testNs {
		base := genEdges(n)
		perm := rand.Perm(1 << n)
		pe := permEdges(base, perm)
		tests[i] = tc{n: n, edges: pe}
		fmt.Fprintln(&input, n)
		for _, e := range pe {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)

	for i, tc := range tests {
		if err := validate(tc.n, tc.edges, scanner); err != nil {
			fmt.Fprintf(os.Stderr, "test case %d (n=%d) failed: %v\n", i+1, tc.n, err)
			os.Exit(1)
		}
	}

	if scanner.Scan() {
		fmt.Fprintf(os.Stderr, "trailing output after last test case: %q\n", scanner.Text())
		os.Exit(1)
	}

	fmt.Println("all tests passed")
}
