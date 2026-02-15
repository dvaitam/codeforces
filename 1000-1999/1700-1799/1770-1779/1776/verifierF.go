package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Edge struct {
	u, v int
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) (string, int, int, []Edge) {
	n := rng.Intn(10) + 3 // n >= 3, keep it small for speed (3 to 12)
	// Ensure connectivity: generate a spanning tree first
	perm := rng.Perm(n)
	edges := make([]Edge, 0)
	for i := 1; i < n; i++ {
		u := perm[i] + 1
		v := perm[rng.Intn(i)] + 1
		edges = append(edges, Edge{u, v})
	}
	
	// Add extra random edges
	maxEdges := n * (n - 1) / 2
	extra := 0
	if maxEdges > len(edges) {
		extra = rng.Intn(maxEdges - len(edges) + 1)
	}
	
	existing := make(map[[2]int]bool)
	for _, e := range edges {
		u, v := e.u, e.v
		if u > v { u, v = v, u }
		existing[[2]int{u, v}] = true
	}
	
	for i := 0; i < extra; i++ {
		for {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v { continue }
			if u > v { u, v = v, u }
			if !existing[[2]int{u, v}] {
				existing[[2]int{u, v}] = true
				edges = append(edges, Edge{u, v})
				break
			}
			if len(existing) == maxEdges { break }
		}
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String(), n, len(edges), edges
}

func checkConnectivity(n int, edges []Edge) bool {
	if n == 0 { return true }
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	
	visited := make(map[int]bool)
	q := []int{1}
	visited[1] = true
	count := 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		count++
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	return count == n
}

func validate(n int, edges []Edge, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	
	if !scanner.Scan() {
		return fmt.Errorf("missing k")
	}
	kStr := scanner.Text()
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	
	if k < 2 || k > len(edges) {
		return fmt.Errorf("k out of range: %d (m=%d)", k, len(edges))
	}
	
	assignments := make([]int, len(edges))
	for i := 0; i < len(edges); i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing assignment for edge %d", i)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid assignment for edge %d: %v", i, err)
		}
		if val < 1 || val > k {
			return fmt.Errorf("assignment %d out of range [1, %d]", val, k)
		}
		assignments[i] = val
	}
	
	// Check 1: Each company is NOT connected
	for c := 1; c <= k; c++ {
		companyEdges := make([]Edge, 0)
		for i, e := range edges {
			if assignments[i] == c {
				companyEdges = append(companyEdges, e)
			}
		}
		if checkConnectivity(n, companyEdges) {
			return fmt.Errorf("company %d graph is fully connected", c)
		}
	}
	
	// Check 2: Any pair of companies IS connected
	for c1 := 1; c1 <= k; c1++ {
		for c2 := c1 + 1; c2 <= k; c2++ {
			pairEdges := make([]Edge, 0)
			for i, e := range edges {
				if assignments[i] == c1 || assignments[i] == c2 {
					pairEdges = append(pairEdges, e)
				}
			}
			if !checkConnectivity(n, pairEdges) {
				return fmt.Errorf("union of companies %d and %d is NOT connected", c1, c2)
			}
		}
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		input, n, _, edges := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		
		if err := validate(n, edges, got); err != nil {
			fmt.Printf("Test %d failed\nInput:\n%s\nOutput:\n%s\nError: %v\n", i+1, input, got, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}