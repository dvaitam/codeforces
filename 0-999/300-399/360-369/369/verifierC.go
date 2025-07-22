package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type edge struct {
	to, t int
}

type caseC struct {
	n      int
	edges  [][3]int
	input  string
	expect []int
}

func computeSet(n int, edges [][3]int) []int {
	adj := make([][]edge, n+1)
	for _, e := range edges {
		a, b, c := e[0], e[1], e[2]
		adj[a] = append(adj[a], edge{b, c})
		adj[b] = append(adj[b], edge{a, c})
	}
	ans := make([]int, 0)
	var dfs func(u, p int) bool
	dfs = func(u, p int) bool {
		selected := false
		for _, e := range adj[u] {
			if e.to == p {
				continue
			}
			child := dfs(e.to, u)
			if e.t == 2 && !child {
				ans = append(ans, e.to)
				child = true
			}
			if child {
				selected = true
			}
		}
		return selected
	}
	dfs(1, 0)
	sort.Ints(ans)
	return ans
}

func generateCase(rng *rand.Rand) caseC {
	n := rng.Intn(9) + 2 // at least 2
	edges := make([][3]int, n-1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		t := rng.Intn(2) + 1
		edges[i-2] = [3]int{p, i, t}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", p, i, t))
	}
	expect := computeSet(n, edges)
	return caseC{n, edges, sb.String(), expect}
}

func runCase(bin string, c caseC) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(c.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := strings.NewReader(out.String())
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("failed to read k: %v\n%s", err, out.String())
	}
	got := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &got[i]); err != nil {
			return fmt.Errorf("failed to read node %d: %v\n%s", i+1, err, out.String())
		}
	}
	sort.Ints(got)
	if k != len(c.expect) {
		return fmt.Errorf("expected %d nodes got %d", len(c.expect), k)
	}
	for i := range got {
		if got[i] != c.expect[i] {
			return fmt.Errorf("mismatch at %d: expected %d got %d", i, c.expect[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
