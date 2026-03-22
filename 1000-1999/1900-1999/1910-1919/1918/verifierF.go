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

// ── embedded reference solver ──────────────────────────────────────────────

func referenceSolve(n, k int, parent []int) int {
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		children[parent[i-2]] = append(children[parent[i-2]], i)
	}

	paths := make([]int, 0, n)

	var dfs func(int) int
	dfs = func(u int) int {
		maxH := 0
		for _, v := range children[u] {
			h := dfs(v) + 1
			if h > maxH {
				if maxH > 0 {
					paths = append(paths, maxH)
				}
				maxH = h
			} else {
				paths = append(paths, h)
			}
		}
		return maxH
	}

	maxH := dfs(1)
	if maxH > 0 {
		paths = append(paths, maxH)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(paths)))

	saved := 0
	limit := k + 1
	if limit > len(paths) {
		limit = len(paths)
	}
	for i := 0; i < limit; i++ {
		saved += paths[i]
	}

	return 2*(n-1) - saved
}

// ── verifier harness ───────────────────────────────────────────────────────

func generateCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(10) + 2
	k := rng.Intn(10)
	parent := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parent[i-2] = rng.Intn(i-1) + 1
	}
	return n, k, parent
}

func runCase(bin string, n, k int, parent []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, p := range parent {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", p))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", referenceSolve(n, k, parent))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		n, k, parent := generateCase(rng)
		if err := runCase(bin, n, k, parent); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
