package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solve(n int, edges [][2]int) (bool, []int) {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	total := 1 << n
	dp := make([]bool, total*n)
	prev := make([]int8, total*n)
	for v := 0; v < n; v++ {
		mask := 1 << v
		dp[mask*n+v] = true
		prev[mask*n+v] = -1
	}
	full := total - 1
	found := false
	endV := -1
	for mask := 1; mask < total && !found; mask++ {
		for v := 0; v < n && !found; v++ {
			if !dp[mask*n+v] {
				continue
			}
			if mask == full {
				found = true
				endV = v
				break
			}
			for _, u := range adj[v] {
				if mask&(1<<u) == 0 {
					m2 := mask | (1 << u)
					idx := m2*n + u
					if !dp[idx] {
						dp[idx] = true
						prev[idx] = int8(v)
					}
				}
			}
		}
	}
	if !found {
		return false, nil
	}
	path := make([]int, n)
	mask := full
	v := endV
	for i := n - 1; i >= 0; i-- {
		path[i] = v + 1
		p := prev[mask*n+v]
		mask ^= 1 << v
		v = int(p)
	}
	return true, path
}

func genGraph() (int, [][2]int) {
	n := rand.Intn(5) + 2 // 2..6
	var edges [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(2) == 0 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	return n, edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n, edges := genGraph()
		m := len(edges)
		expOk, expPath := solve(n, edges)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, string(out))
			os.Exit(1)
		}
		lines := strings.Fields(strings.TrimSpace(string(out)))
		if !expOk {
			if len(lines) != 1 || lines[0] != "No" {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d expected No got %s\n", i+1, string(out))
				os.Exit(1)
			}
			continue
		}
		if len(lines) != n+1 || lines[0] != "Yes" {
			fmt.Fprintf(os.Stderr, "wrong format on test %d\n", i+1)
			os.Exit(1)
		}
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(lines[j+1])
			if err != nil || val != expPath[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	fmt.Println("Accepted")
}
