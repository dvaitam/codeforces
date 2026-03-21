package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// validate checks whether the candidate output is a valid answer for 1379E.
// n people, k imbalanced. Each person has either 0 or 2 parents.
// s[i] = child of person i (0 for the root who has no child).
// A person with 2 parents is imbalanced if #ancestors(parent1) >= 2*#ancestors(parent2)
// or vice versa.
func validate(n, k int, output string) error {
	line := strings.TrimSpace(output)
	if strings.ToUpper(line) == "NO" || strings.HasPrefix(strings.ToUpper(line), "NO\n") || line == "NO" {
		// Candidate says NO. We need to check if a solution actually exists.
		// Use the feasibility check.
		if n == 1 && k == 0 {
			return fmt.Errorf("candidate said NO but answer exists for n=1 k=0")
		}
		if n%2 == 0 {
			return nil // correct: no solution for even n
		}
		L := (n + 1) / 2
		if !feasible(L, k) {
			return nil // correct
		}
		return fmt.Errorf("candidate said NO but a solution exists for n=%d k=%d", n, k)
	}

	lines := strings.Split(strings.TrimSpace(line), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("expected YES + array, got %q", line)
	}
	if strings.TrimSpace(strings.ToUpper(lines[0])) != "YES" {
		return fmt.Errorf("first line should be YES or NO, got %q", lines[0])
	}

	// Check feasibility: if no solution should exist, candidate shouldn't say YES
	if n%2 == 0 {
		return fmt.Errorf("candidate said YES but n=%d is even", n)
	}
	L := (n + 1) / 2
	if n > 1 && !feasible(L, k) {
		return fmt.Errorf("candidate said YES but no solution exists for n=%d k=%d", n, k)
	}

	fields := strings.Fields(lines[1])
	if len(fields) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		v, err := strconv.Atoi(fields[i-1])
		if err != nil {
			return fmt.Errorf("parse error at position %d: %v", i, err)
		}
		if v < 0 || v > n {
			return fmt.Errorf("s[%d]=%d out of range", i, v)
		}
		s[i] = v
	}

	// Find root (s[i]==0)
	rootCount := 0
	for i := 1; i <= n; i++ {
		if s[i] == 0 {
			rootCount++
		}
	}
	if rootCount != 1 {
		return fmt.Errorf("expected exactly 1 root, found %d", rootCount)
	}

	// Build children: for each node, collect parents (nodes whose child is this node)
	parents := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if s[i] != 0 {
			parents[s[i]] = append(parents[s[i]], i)
		}
	}

	// Each node must have 0 or 2 parents
	for i := 1; i <= n; i++ {
		if len(parents[i]) != 0 && len(parents[i]) != 2 {
			return fmt.Errorf("node %d has %d parents (must be 0 or 2)", i, len(parents[i]))
		}
	}

	// Check it forms a valid tree: compute ancestor counts via DFS from leaves
	// ancestors[i] = number of ancestors of i (including i itself)
	ancestors := make([]int, n+1)
	// Process in topological order: leaves first
	// indegree = number of parents
	computed := make([]bool, n+1)
	queue := make([]int, 0)
	for i := 1; i <= n; i++ {
		if len(parents[i]) == 0 {
			ancestors[i] = 1
			computed[i] = true
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		child := s[v]
		if child == 0 {
			continue
		}
		// Check if both parents of child are computed
		allDone := true
		for _, p := range parents[child] {
			if !computed[p] {
				allDone = false
				break
			}
		}
		if allDone && !computed[child] {
			total := 1
			for _, p := range parents[child] {
				total += ancestors[p]
			}
			ancestors[child] = total
			computed[child] = true
			queue = append(queue, child)
		}
	}

	// Check all nodes are computed (tree is connected)
	for i := 1; i <= n; i++ {
		if !computed[i] {
			return fmt.Errorf("node %d not reachable (cycle or disconnected)", i)
		}
	}

	// Count imbalanced
	imbalanced := 0
	for i := 1; i <= n; i++ {
		if len(parents[i]) == 2 {
			a := ancestors[parents[i][0]]
			b := ancestors[parents[i][1]]
			if a >= 2*b || b >= 2*a {
				imbalanced++
			}
		}
	}

	if imbalanced != k {
		return fmt.Errorf("expected %d imbalanced, got %d", k, imbalanced)
	}

	return nil
}

func isPow2(x int) bool {
	return x > 0 && (x&(x-1)) == 0
}

func feasible(L, k int) bool {
	if L < 1 || k < 0 {
		return false
	}
	if L == 1 {
		return k == 0
	}
	if k > L-2 {
		return false
	}
	if isPow2(L) {
		return k == 0 || k >= 2
	}
	if L == 5 {
		return k == 1 || k == 3
	}
	return k >= 1
}

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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]

	// Generate test cases: systematic small + random
	type tc struct{ n, k int }
	var cases []tc

	// All small cases
	for n := 1; n <= 15; n++ {
		for k := 0; k <= n; k++ {
			cases = append(cases, tc{n, k})
		}
	}

	// Random medium cases
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 50; i++ {
		n := rng.Intn(99) + 1
		k := rng.Intn(n)
		cases = append(cases, tc{n, k})
	}

	for idx, c := range cases {
		input := fmt.Sprintf("%d %d\n", c.n, c.k)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d (n=%d k=%d): %v\n", idx+1, c.n, c.k, err)
			os.Exit(1)
		}
		if err := validate(c.n, c.k, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (n=%d k=%d) failed: %v\noutput:\n%s\n", idx+1, c.n, c.k, err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
