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

type Edge struct{ to, d int }

type Case struct {
	n   int
	M   int
	adj [][]Edge
}

func generateCase(rng *rand.Rand) Case {
	n := rng.Intn(5) + 2 //2..6
	Ms := []int{3, 7, 9, 11, 13, 17, 19}
	M := Ms[rng.Intn(len(Ms))]
	adj := make([][]Edge, n)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		d := rng.Intn(9) + 1
		adj[i] = append(adj[i], Edge{p, d})
		adj[p] = append(adj[p], Edge{i, d})
	}
	return Case{n, M, adj}
}

func pathDigits(c Case, start, goal int) []int {
	n := c.n
	parent := make([]int, n)
	digit := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	queue := []int{start}
	parent[start] = start
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		if u == goal {
			break
		}
		for _, e := range c.adj[u] {
			if parent[e.to] == -1 {
				parent[e.to] = u
				digit[e.to] = e.d
				queue = append(queue, e.to)
			}
		}
	}
	if parent[goal] == -1 {
		return nil
	}
	var digits []int
	for v := goal; v != start; v = parent[v] {
		digits = append([]int{digit[v]}, digits...)
	}
	return digits
}

func countPairs(c Case) int {
	cnt := 0
	for u := 0; u < c.n; u++ {
		for v := 0; v < c.n; v++ {
			if u == v {
				continue
			}
			digits := pathDigits(c, u, v)
			val := 0
			for _, d := range digits {
				val = (val*10 + d) % c.M
			}
			if val == 0 {
				cnt++
			}
		}
	}
	return cnt
}

func runCase(bin string, c Case) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", c.n, c.M)
	for u := 0; u < c.n; u++ {
		for _, e := range c.adj[u] {
			if e.to > u { // avoid duplicates
				fmt.Fprintf(&sb, "%d %d %d\n", u, e.to, e.d)
			}
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expect := countPairs(c)
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
