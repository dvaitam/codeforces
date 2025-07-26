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

func disjoint(a, b []int) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			return false
		}
		if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}
	return true
}

func expected(days [][]int) string {
	m := len(days)
	adj := make([][]int, m)
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			if i == j {
				continue
			}
			if disjoint(days[i], days[j]) {
				adj[i] = append(adj[i], j)
			}
		}
	}
	vis := make([]int, m)
	var dfs func(int) bool
	dfs = func(u int) bool {
		vis[u] = 1
		for _, v := range adj[u] {
			if vis[v] == 1 {
				return true
			}
			if vis[v] == 0 {
				if dfs(v) {
					return true
				}
			}
		}
		vis[u] = 2
		return false
	}
	for i := 0; i < m; i++ {
		if vis[i] == 0 {
			if dfs(i) {
				return "impossible"
			}
		}
	}
	return "possible"
}

func generateCase(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 1
	n := rng.Intn(8) + m + 1 // ensure n > max s_i
	days := make([][]int, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", m, n))
	for i := 0; i < m; i++ {
		s := rng.Intn(n-1) + 1
		perm := rng.Perm(n)
		lst := perm[:s]
		sort.Ints(lst)
		days[i] = lst
		sb.WriteString(fmt.Sprintf("%d", s))
		for _, v := range lst {
			sb.WriteString(fmt.Sprintf(" %d", v+1))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expected(days)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
