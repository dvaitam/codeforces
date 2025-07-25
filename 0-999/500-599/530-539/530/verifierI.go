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

func runCandidate(bin string, input string) (string, error) {
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

func expected(n int, groups [][]int) []int {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for _, g := range groups {
		for i := 0; i < len(g); i++ {
			for j := i + 1; j < len(g); j++ {
				u, v := g[i], g[j]
				adj[u][v] = true
				adj[v][u] = true
			}
		}
	}
	xs := make([]int, n)
	var dfs func(int, int) bool
	dfs = func(i, C int) bool {
		if i == n {
			return true
		}
		for v := 1; v <= C; v++ {
			ok := true
			for j := 0; j < i; j++ {
				if xs[j] == v && adj[i][j] {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			xs[i] = v
			if dfs(i+1, C) {
				return true
			}
			xs[i] = 0
		}
		return false
	}
	for C := 1; C <= n; C++ {
		for i := range xs {
			xs[i] = 0
		}
		if dfs(0, C) {
			res := append([]int(nil), xs...)
			return res
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for caseIdx := 0; caseIdx < 100; caseIdx++ {
		n := rng.Intn(5) + 2
		k := rng.Intn(n) + 1
		groups := make([][]int, k)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i := 0; i < k; i++ {
			m := rng.Intn(n-1) + 2
			perm := rng.Perm(n)
			group := make([]int, m)
			for j := 0; j < m; j++ {
				group[j] = perm[j]
			}
			groups[i] = group
			fmt.Fprintf(&sb, "%d", m)
			for _, idx := range group {
				fmt.Fprintf(&sb, " %d", idx+1)
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := expected(n, groups)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d values got %d\n", caseIdx+1, n, len(fields))
			os.Exit(1)
		}
		got := make([]int, n)
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid int\n", caseIdx+1)
				os.Exit(1)
			}
			got[i] = v
		}
		for i := 0; i < n; i++ {
			if got[i] != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", caseIdx+1, want, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
