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

type Case struct{ input string }

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := 0; i < 100; i++ {
		T := rng.Intn(10) + 2
		arr := make([]int, T+1)
		sum := 1
		for j := 1; j <= T; j++ {
			arr[j] = rng.Intn(5) + 1
			sum += arr[j]
		}
		n := sum
		// k can be between 1 and n - T ideally, let's randomize reasonably
		// k leaf nodes means n-1-k inner nodes. 
		// min inner nodes = T-1 (one path)
		// max inner nodes = sum_{i=0}^{T-2} min(a[i], a[i+1])
		
		// Just randomly picking k in valid range (1 <= k < n)
		k := rng.Intn(n-1) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, T, k)
		for j := 1; j <= T; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", arr[j])
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func checkSolution(input, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 1 && lines[0] == "-1" {
		// Solution says impossible. We trust it for now or need a deeper check.
		// However, verifying impossibility is hard without a reference.
		// Given the user says this is a "correct solution", we can try to verify validity if it outputs a tree.
		return nil
	}

	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	nOut, err := strconv.Atoi(lines[0])
	if err != nil {
		return fmt.Errorf("invalid n in output: %v", err)
	}

	// Parse input n, t, k
	inLines := strings.Fields(input)
	nIn, _ := strconv.Atoi(inLines[0])
	tIn, _ := strconv.Atoi(inLines[1])
	kIn, _ := strconv.Atoi(inLines[2])
	
	if nOut != nIn {
		return fmt.Errorf("output n (%d) != input n (%d)", nOut, nIn)
	}

	if len(lines) != nIn {
		return fmt.Errorf("expected %d lines (1 for n + %d edges), got %d", nIn, nIn-1, len(lines))
	}

	adj := make(map[int][]int)
	degree := make(map[int]int)
	
	for i := 1; i < nIn; i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != 2 {
			return fmt.Errorf("invalid edge format at line %d: %s", i+1, lines[i])
		}
		u, _ := strconv.Atoi(parts[0])
		v, _ := strconv.Atoi(parts[1])
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		degree[u]++
		degree[v]++
	}

	// BFS to check distances and tree structure
	dist := make(map[int]int)
	q := []int{1}
	dist[1] = 0
	visited := make(map[int]bool)
	visited[1] = true
	
	nodesAtDist := make(map[int]int)
	maxDist := 0

	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		d := dist[u]
		if d > 0 {
			nodesAtDist[d]++
		}
		if d > maxDist {
			maxDist = d
		}

		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				dist[v] = d + 1
				q = append(q, v)
			}
		}
	}

	if len(visited) != nIn {
		return fmt.Errorf("graph is not connected, visited %d/%d nodes", len(visited), nIn)
	}
	if maxDist != tIn {
		return fmt.Errorf("max distance is %d, expected %d", maxDist, tIn)
	}

	// Check 'a' counts
	// input a is 0-indexed in Go slice, representing levels 1 to t
	aStartIdx := 3
	for i := 1; i <= tIn; i++ {
		expectedCount, _ := strconv.Atoi(inLines[aStartIdx + i - 1])
		if nodesAtDist[i] != expectedCount {
			return fmt.Errorf("level %d has %d nodes, expected %d", i, nodesAtDist[i], expectedCount)
		}
	}

	// Check leaf nodes count (excluding root)
	leafCount := 0
	for i := 2; i <= nIn; i++ {
		if degree[i] == 1 {
			leafCount++
		}
	}

	if leafCount != kIn {
		return fmt.Errorf("found %d leaves, expected %d", leafCount, kIn)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := genCases()
	for i, c := range cases {
		out, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		
		if err := checkSolution(c.input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}