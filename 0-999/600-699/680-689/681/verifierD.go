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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

type edge struct{ u, v int }

type testCase struct {
	n     int
	edges []edge
	gifts []int
}

func hasSolution(tc testCase) bool {
	n := tc.n
	children := make([][]int, n+1)
	parent := make([]int, n+1)
	for _, e := range tc.edges {
		children[e.u] = append(children[e.u], e.v)
		parent[e.v] = e.u
	}
	for u := 1; u <= n; u++ {
		if tc.gifts[u] != u {
			p := parent[u]
			if p == 0 || tc.gifts[p] != tc.gifts[u] {
				return false
			}
		}
	}
	return true
}

func validateOutput(tc testCase, out string) error {
	toks := strings.Fields(strings.TrimSpace(out))
	if len(toks) == 0 {
		return fmt.Errorf("empty output")
	}
	if toks[0] == "-1" {
		if len(toks) != 1 {
			return fmt.Errorf("invalid -1 output format")
		}
		if hasSolution(tc) {
			return fmt.Errorf("reported -1 but a solution exists")
		}
		return nil
	}

	if !hasSolution(tc) {
		return fmt.Errorf("solution exists check failed: expected -1 for this case")
	}

	k, err := strconv.Atoi(toks[0])
	if err != nil || k < 1 || k > tc.n {
		return fmt.Errorf("invalid k")
	}
	if len(toks) != k+1 {
		return fmt.Errorf("expected %d listed vertices, got %d", k, len(toks)-1)
	}

	pos := make([]int, tc.n+1)
	for i := range pos {
		pos[i] = -1
	}
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(toks[i+1])
		if err != nil || v < 1 || v > tc.n {
			return fmt.Errorf("invalid listed vertex")
		}
		if pos[v] != -1 {
			return fmt.Errorf("duplicate vertex in list")
		}
		pos[v] = i
	}

	parent := make([]int, tc.n+1)
	for _, e := range tc.edges {
		parent[e.v] = e.u
	}

	for u := 1; u <= tc.n; u++ {
		bestPos := int(1e9)
		bestNode := -1
		x := u
		for x != 0 {
			if pos[x] != -1 && pos[x] < bestPos {
				bestPos = pos[x]
				bestNode = x
			}
			x = parent[x]
		}
		if bestNode == -1 {
			return fmt.Errorf("man %d has no ancestor in list", u)
		}
		if bestNode != tc.gifts[u] {
			return fmt.Errorf("wish mismatch for %d: expected %d got %d", u, tc.gifts[u], bestNode)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, testCase) {
	n := rng.Intn(6) + 1
	edges := make([]edge, 0, n-1)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
		parent[i] = p
	}
	gifts := make([]int, n+1)
	for i := 1; i <= n; i++ {
		// Generate valid gifts as required by statement:
		// each a[i] must be an ancestor of i (possibly itself).
		u := i
		anc := []int{u}
		for parent[u] != 0 {
			u = parent[u]
			anc = append(anc, u)
		}
		gifts[i] = anc[rng.Intn(len(anc))]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d", gifts[i]))
		if i < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String(), testCase{n: n, edges: edges, gifts: gifts}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ngot:\n%s\ninput:\n%s", i+1, err, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
