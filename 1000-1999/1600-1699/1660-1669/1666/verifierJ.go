package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	inputData := generateTestCase()

	n, matrix, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse generated input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, inputData)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refParents, err := parseParents(refOut, n)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}
	refCost, err := validateAndCost(n, matrix, refParents)
	if err != nil {
		fail("reference solution invalid: %v", err)
	}

	candOut, err := runProgram(candidate, inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candParents, err := parseParents(candOut, n)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}
	candCost, err := validateAndCost(n, matrix, candParents)
	if err != nil {
		fail("invalid tree: %v", err)
	}

	if candCost != refCost {
		fmt.Printf("Input:\n%s\n", string(inputData))
		fmt.Printf("Reference Output:\n%s\n", refOut)
		fmt.Printf("Candidate Output:\n%s\n", candOut)
		fail("cost mismatch: expected %d got %d", refCost, candCost)
	}

	fmt.Println("OK")
}

func generateTestCase() []byte {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(50) + 1 // Use smaller N for faster testing, or up to 200
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", n)

	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			val := rand.Intn(1000) // Random cost
			matrix[i][j] = val
			matrix[j][i] = val
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Fprintf(&buf, "%d", matrix[i][j])
			if j < n-1 {
				fmt.Fprint(&buf, " ")
			}
		}
		fmt.Fprintln(&buf)
	}
	return buf.Bytes()
}

func parseInput(data []byte) (int, [][]int64, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, err
	}
	if n <= 0 {
		return 0, nil, fmt.Errorf("n must be positive, got %d", n)
	}
	matrix := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		matrix[i] = make([]int64, n+1)
		for j := 1; j <= n; j++ {
			if _, err := fmt.Fscan(reader, &matrix[i][j]); err != nil {
				return 0, nil, err
			}
		}
	}
	return n, matrix, nil
}

func buildReference() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to determine current file path")
	}
	dir := filepath.Dir(currentFile)
	refSource := filepath.Join(dir, "1666J.go")

	tmp, err := os.CreateTemp("", "1666J-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseParents(out string, n int) ([]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if _, err := fmt.Fscan(reader, &parent[i]); err != nil {
			return nil, fmt.Errorf("expected %d integers, parsed %d: %v", n, i-1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("unexpected extra output: %s", extra)
	}
	return parent, nil
}

func validateAndCost(n int, c [][]int64, parent []int) (int64, error) {
	if len(parent) != n+1 {
		return 0, fmt.Errorf("parent array has invalid length")
	}
	left := make([]int, n+1)
	right := make([]int, n+1)
	adj := make([][]int, n+1)

	root := 0
	rootCount := 0
	for i := 1; i <= n; i++ {
		p := parent[i]
		if p < 0 || p > n {
			return 0, fmt.Errorf("parent of node %d is out of range", i)
		}
		if p == i {
			return 0, fmt.Errorf("node %d cannot be its own parent", i)
		}
		if p == 0 {
			root = i
			rootCount++
			if rootCount > 1 {
				return 0, fmt.Errorf("multiple roots detected")
			}
			continue
		}
		adj[i] = append(adj[i], p)
		adj[p] = append(adj[p], i)
		if i < p {
			if left[p] != 0 {
				return 0, fmt.Errorf("node %d has multiple left children", p)
			}
			left[p] = i
		} else if i > p {
			if right[p] != 0 {
				return 0, fmt.Errorf("node %d has multiple right children", p)
			}
			right[p] = i
		} else {
			return 0, fmt.Errorf("node %d has invalid parent %d", i, p)
		}
	}
	if rootCount != 1 {
		return 0, fmt.Errorf("expected exactly one root, found %d", rootCount)
	}

	if err := checkConnectivity(root, n, adj); err != nil {
		return 0, err
	}
	if err := validateBST(root, 0, n+1, left, right); err != nil {
		return 0, err
	}

	total := int64(0)
	for i := 1; i <= n; i++ {
		dist := bfsDistances(i, n, adj)
		for j := i + 1; j <= n; j++ {
			if dist[j] == -1 {
				return 0, fmt.Errorf("tree is not connected")
			}
			total += c[i][j] * int64(dist[j])
		}
	}
	return total, nil
}

func checkConnectivity(root, n int, adj [][]int) error {
	visited := make([]bool, n+1)
	queue := []int{root}
	visited[root] = true
	count := 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		count++
		for _, to := range adj[v] {
			if !visited[to] {
				visited[to] = true
				queue = append(queue, to)
			}
		}
	}
	if count != n {
		return fmt.Errorf("tree is not connected")
	}
	return nil
}

func validateBST(node, low, high int, left, right []int) error {
	if node <= low || node >= high {
		return fmt.Errorf("node %d violates BST constraints", node)
	}
	if left[node] != 0 {
		if err := validateBST(left[node], low, node, left, right); err != nil {
			return err
		}
	}
	if right[node] != 0 {
		if err := validateBST(right[node], node, high, left, right); err != nil {
			return err
		}
	}
	return nil
}

func bfsDistances(start, n int, adj [][]int) []int {
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	queue := []int{start}
	dist[start] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				queue = append(queue, to)
			}
		}
	}
	return dist
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
