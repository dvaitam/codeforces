package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func compileRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		log.Fatal("REFERENCE_SOURCE_PATH environment variable is not set")
	}
	out := filepath.Join(os.TempDir(), "refF_"+fmt.Sprint(time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", out, src)
	if o, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build error: %v\n%s", err, o)
	}
	return out, nil
}

func runBin(bin, input string) (string, error) {
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

func generateCase(r *rand.Rand) string {
	k := r.Intn(3) + 1
	n := 2*k + r.Intn(5) // ensure n >= 2k
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, i))
	}
	nodes := r.Perm(n)
	specials := nodes[:2*k]
	for i, v := range specials {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

// parseInput extracts n, k, adjacency list, and special nodes from input
func parseInput(input string) (n, k int, adj [][]int, specials []int) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	fields := strings.Fields(lines[0])
	n, _ = strconv.Atoi(fields[0])
	k, _ = strconv.Atoi(fields[1])
	adj = make([][]int, n+1)
	for i := 1; i < n; i++ {
		f := strings.Fields(lines[i])
		a, _ := strconv.Atoi(f[0])
		b, _ := strconv.Atoi(f[1])
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	sFields := strings.Fields(lines[n])
	specials = make([]int, 2*k)
	for i := 0; i < 2*k; i++ {
		specials[i], _ = strconv.Atoi(sFields[i])
	}
	return
}

// isOnPath checks if node x is on the path from u to v in the tree
func isOnPath(adj [][]int, n, u, v, x int) bool {
	// BFS to find path from u to v
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -1
	}
	parent[u] = u
	queue := []int{u}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == v {
			break
		}
		for _, nb := range adj[cur] {
			if parent[nb] == -1 {
				parent[nb] = cur
				queue = append(queue, nb)
			}
		}
	}
	// Trace path from v back to u
	cur := v
	for cur != u {
		if cur == x {
			return true
		}
		cur = parent[cur]
	}
	return cur == x
}

// validateOutput checks if candidate output is a valid solution
func validateOutput(input, output string, refM int) error {
	n, k, adj, specials := parseInput(input)
	_ = n
	lines := strings.Split(output, "\n")
	if len(lines) < 2+k {
		return fmt.Errorf("output too short: expected at least %d lines, got %d", 2+k, len(lines))
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid m: %v", err)
	}
	if m != refM {
		return fmt.Errorf("m=%d but optimal is %d", m, refM)
	}
	// Parse settlement cities
	dFields := strings.Fields(lines[1])
	if len(dFields) != m {
		return fmt.Errorf("expected %d settlement cities, got %d", m, len(dFields))
	}
	dSet := make(map[int]bool)
	for _, f := range dFields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid city: %v", err)
		}
		dSet[v] = true
	}
	// Parse triplets
	specialSet := make(map[int]bool)
	for _, s := range specials {
		specialSet[s] = true
	}
	usedSpecials := make(map[int]bool)
	for i := 0; i < k; i++ {
		fields := strings.Fields(lines[2+i])
		if len(fields) != 3 {
			return fmt.Errorf("triplet %d: expected 3 integers, got %q", i+1, lines[2+i])
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		x, err3 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("triplet %d: parse error", i+1)
		}
		// u and v must be special nodes
		if !specialSet[u] || !specialSet[v] {
			return fmt.Errorf("triplet %d: %d or %d not special", i+1, u, v)
		}
		if usedSpecials[u] || usedSpecials[v] {
			return fmt.Errorf("triplet %d: duplicate special node", i+1)
		}
		usedSpecials[u] = true
		usedSpecials[v] = true
		// x must be a settlement city
		if !dSet[x] {
			return fmt.Errorf("triplet %d: city %d not in settlement cities", i+1, x)
		}
		// x must be on path from u to v
		if !isOnPath(adj, n, u, v, x) {
			return fmt.Errorf("triplet %d: city %d not on path from %d to %d", i+1, x, u, v)
		}
	}
	// All special nodes must be used
	if len(usedSpecials) != 2*k {
		return fmt.Errorf("not all special nodes used: %d of %d", len(usedSpecials), 2*k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(r)
		exp, err := runBin(ref, in)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBin(bin, in)
		if err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		// Parse reference m (first line)
		refLines := strings.Split(exp, "\n")
		refM, _ := strconv.Atoi(strings.TrimSpace(refLines[0]))
		if err := validateOutput(in, out, refM); err != nil {
			fmt.Printf("test %d failed: %v\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, err, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
