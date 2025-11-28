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

func run(cmdPath string, input string) (string, error) {
	cmd := exec.Command(cmdPath)
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{i, rand.Intn(i-1) + 1})
	}
	return edges
}

func genTests() []string {
	rand.Seed(time.Now().UnixNano())
	tests := make([]string, 0, 100)
	// Test cases logic
	for i := 0; i < 50; i++ {
		x := rand.Intn(4) + 1 // p up to 4 (n=16), small enough for quick check
		n := 1 << x
		edges := genTree(n)
		var sb strings.Builder
		sb.WriteString("1\n") // T=1
		sb.WriteString(strconv.Itoa(x) + "\n")
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func verify(input, output string) error {
	inSc := bufio.NewScanner(strings.NewReader(input))
	inSc.Split(bufio.ScanWords)

	// Parse input
	if !inSc.Scan() {
		return fmt.Errorf("empty input")
	}
	t, _ := strconv.Atoi(inSc.Text())
	if t != 1 {
		return fmt.Errorf("expected T=1, got %d", t)
	}

	if !inSc.Scan() {
		return fmt.Errorf("expected p")
	}
	p, _ := strconv.Atoi(inSc.Text())
	n := 1 << p

	// Parse edges to build adjacency
	// edges are 0-indexed in problem? No, 1-indexed u, v.
	// We need to map input edge index (0..n-2) to verify output edge weights.
	type Edge struct {
		to, id int
	}
	adj := make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		if !inSc.Scan() {
			return fmt.Errorf("expected u")
		}
		u, _ := strconv.Atoi(inSc.Text())
		if !inSc.Scan() {
			return fmt.Errorf("expected v")
		}
		v, _ := strconv.Atoi(inSc.Text())
		adj[u] = append(adj[u], Edge{v, i})
		adj[v] = append(adj[v], Edge{u, i})
	}

	// Parse output
	outSc := bufio.NewScanner(strings.NewReader(output))
	outSc.Split(bufio.ScanWords)

	if !outSc.Scan() {
		return fmt.Errorf("expected root")
	}
	root, _ := strconv.Atoi(outSc.Text())
	if root < 1 || root > n {
		return fmt.Errorf("invalid root %d", root)
	}

	// Node weights
	nodeVals := make([]int, n+1)
	seen := make(map[int]bool)
	maxVal := 2*n - 1

	for i := 1; i <= n; i++ {
		if !outSc.Scan() {
			return fmt.Errorf("expected node weight %d", i)
		}
		val, _ := strconv.Atoi(outSc.Text())
		if val < 1 || val > maxVal {
			return fmt.Errorf("node value %d out of range [1, %d]", val, maxVal)
		}
		if seen[val] {
			return fmt.Errorf("duplicate value %d", val)
		}
		seen[val] = true
		nodeVals[i] = val
	}

	// Edge weights
	edgeVals := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		if !outSc.Scan() {
			return fmt.Errorf("expected edge weight %d", i)
		}
		val, _ := strconv.Atoi(outSc.Text())
		if val < 1 || val > maxVal {
			return fmt.Errorf("edge value %d out of range [1, %d]", val, maxVal)
		}
		if seen[val] {
			return fmt.Errorf("duplicate value %d", val)
		}
		seen[val] = true
		edgeVals[i] = val
	}

	// Verify path XOR sums from root
	// DFS from root
	var dfs func(u, p, currentXor int) error
	dfs = func(u, p, currentXor int) error {
		currentXor ^= nodeVals[u]
		if currentXor > maxVal {
			return fmt.Errorf("path xor to %d is %d > %d", u, currentXor, maxVal)
		}
		for _, e := range adj[u] {
			if e.to != p {
				edgeVal := edgeVals[e.id]
				if err := dfs(e.to, u, currentXor^edgeVal); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := dfs(root, -1, 0); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
