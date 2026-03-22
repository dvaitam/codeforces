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

func runBin(path string, input []byte) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return out.Bytes(), nil
}

func compileBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifierC-bin-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, nil, nil
}

// solveC solves Codeforces 1098C (tree with min branching coefficient and given subtree-size sum).
// Returns "No" or "Yes\np2 p3 ... pn".
func solveC(n int, s int64) string {
	minSum := int64(2*n - 1)
	maxSum := int64(n) * int64(n+1) / 2
	if s < minSum || s > maxSum {
		return "No"
	}

	a := make([]int, n+2)
	b := make([]int, n+2)

	chk := func(x int) bool {
		for i := 0; i <= n+1; i++ {
			b[i] = 0
		}
		a[1] = 1
		idx := 2
		depth := 1
		t := int64(1)
		rem := s - 1
		for idx <= n {
			t *= int64(x)
			if t > int64(n) {
				t = int64(n)
			}
			depth++
			b[depth] = 0
			limit := int(t)
			for j := 0; j < limit && idx <= n; j++ {
				a[idx] = depth
				rem -= int64(depth)
				b[depth]++
				idx++
				if rem < 0 {
					return false
				}
			}
		}
		if rem < 0 {
			return false
		}
		j := n
		for rem > 0 {
			depth++
			for j > 1 && b[a[j]] == 1 {
				j--
			}
			if j <= 1 {
				return false
			}
			inc := depth - a[j]
			if int64(inc) > rem {
				inc = int(rem)
			}
			rem -= int64(inc)
			b[a[j]]--
			a[j] += inc
			j--
		}
		return rem == 0
	}

	lo, hi := 1, n
	ans := n
	for lo <= hi {
		mid := (lo + hi) / 2
		if chk(mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	chk(ans)

	// Sort depths for assignment
	a[1] = 1
	depths := make([]int, n-1)
	for i := 2; i <= n; i++ {
		depths[i-2] = a[i]
	}
	sortInts(depths)
	for i := 2; i <= n; i++ {
		a[i] = depths[i-2]
	}

	children := make([]int, n+2)
	parents := make([]int, n+1)
	ptr := 1
	for i := 2; i <= n; i++ {
		for a[ptr] != a[i]-1 || children[ptr] == ans {
			ptr++
		}
		parents[i] = ptr
		children[ptr]++
	}

	var sb strings.Builder
	sb.WriteString("Yes\n")
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(parents[i]))
	}
	return sb.String()
}

func sortInts(a []int) {
	// Simple insertion sort (n <= 10^5 in real usage, but test cases are small)
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

// verifyTree checks that the candidate output represents a valid tree
// with n vertices, correct branching coefficient, and sum of subtree sizes = s.
func verifyTree(n int, s int64, out string) error {
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	first := strings.TrimSpace(lines[0])

	// Compute expected answer to know if "Yes" or "No" is correct
	refOut := solveC(n, s)
	refFirst := strings.TrimSpace(strings.Split(refOut, "\n")[0])

	if strings.EqualFold(first, "No") {
		if refFirst == "Yes" {
			return fmt.Errorf("candidate says No but solution exists")
		}
		return nil
	}

	if !strings.EqualFold(first, "Yes") {
		return fmt.Errorf("first line should be Yes or No, got %q", first)
	}

	if refFirst == "No" {
		return fmt.Errorf("candidate says Yes but no solution exists")
	}

	if len(lines) < 2 {
		return fmt.Errorf("missing parent array")
	}

	parts := strings.Fields(lines[1])
	if len(parts) != n-1 {
		return fmt.Errorf("expected %d parents, got %d", n-1, len(parts))
	}

	parent := make([]int, n+1)
	childCount := make([]int, n+1)
	children := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		v, err := strconv.Atoi(parts[i])
		if err != nil {
			return fmt.Errorf("parse error for parent[%d]: %v", i+2, err)
		}
		if v < 1 || v > n || v == i+2 {
			return fmt.Errorf("parent[%d]=%d invalid", i+2, v)
		}
		parent[i+2] = v
		childCount[v]++
		children[v] = append(children[v], i+2)
	}

	// Verify it's a valid tree rooted at 1 using BFS
	visited := make([]bool, n+1)
	visited[1] = true
	queue := []int{1}
	visitCount := 1
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, c := range children[u] {
			if visited[c] {
				return fmt.Errorf("cycle detected at vertex %d", c)
			}
			visited[c] = true
			visitCount++
			queue = append(queue, c)
		}
	}
	if visitCount != n {
		return fmt.Errorf("tree not connected: visited %d of %d vertices", visitCount, n)
	}

	// Compute subtree sizes using topological order (reverse BFS)
	order := make([]int, 0, n)
	visited2 := make([]bool, n+1)
	visited2[1] = true
	queue2 := []int{1}
	for len(queue2) > 0 {
		u := queue2[0]
		queue2 = queue2[1:]
		order = append(order, u)
		for _, c := range children[u] {
			if !visited2[c] {
				visited2[c] = true
				queue2 = append(queue2, c)
			}
		}
	}
	subtreeSize := make([]int, n+1)
	for i := 1; i <= n; i++ {
		subtreeSize[i] = 1
	}
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		if parent[u] != 0 {
			subtreeSize[parent[u]] += subtreeSize[u]
		}
	}

	var totalSum int64
	for i := 1; i <= n; i++ {
		totalSum += int64(subtreeSize[i])
	}

	if totalSum != s {
		return fmt.Errorf("sum of subtree sizes = %d, expected %d", totalSum, s)
	}

	// Compute expected branching coefficient from reference
	refLines := strings.Split(refOut, "\n")
	refParts := strings.Fields(refLines[1])
	refChildCount := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		v, _ := strconv.Atoi(refParts[i])
		refChildCount[v]++
	}
	refMaxBranch := 0
	for i := 1; i <= n; i++ {
		if refChildCount[i] > refMaxBranch {
			refMaxBranch = refChildCount[i]
		}
	}

	maxBranch := 0
	for i := 1; i <= n; i++ {
		if childCount[i] > maxBranch {
			maxBranch = childCount[i]
		}
	}

	if maxBranch != refMaxBranch {
		return fmt.Errorf("branching coefficient = %d, expected %d", maxBranch, refMaxBranch)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := compileBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "compile error: %v\n", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	rand.Seed(3)

	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		maxSum := n * (n + 1) / 2
		s := rand.Intn(maxSum-n+1) + n
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, s)
		input := buf.Bytes()

		out, err := runBin(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		if err := verifyTree(n, int64(s), string(out)); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput: %d %d\nGot:\n%s", t, err, n, s, string(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
