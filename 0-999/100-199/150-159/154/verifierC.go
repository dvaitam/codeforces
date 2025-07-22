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

type edge struct{ u, v int }

func solveC(n int, edges []edge) int64 {
	z := make([]uint64, n+1)
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= n; i++ {
		z[i] = rng.Uint64()
	}
	h := make([]uint64, n+1)
	for _, e := range edges {
		h[e.u] += z[e.v]
		h[e.v] += z[e.u]
	}
	hs := make([]uint64, n)
	for i := 1; i <= n; i++ {
		hs[i-1] = h[i]
	}
	sort.Slice(hs, func(i, j int) bool { return hs[i] < hs[j] })
	var ans int64
	for i := 0; i < n; {
		j := i + 1
		for j < n && hs[j] == hs[i] {
			j++
		}
		c := int64(j - i)
		ans += c * (c - 1) / 2
		i = j
	}
	for _, e := range edges {
		if h[e.u]-z[e.v] == h[e.v]-z[e.u] {
			ans++
		}
	}
	return ans
}

func genCaseC(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(10) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	used := make(map[[2]int]bool)
	edges := make([]edge, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if used[[2]int{u, v}] {
			continue
		}
		used[[2]int{u, v}] = true
		edges = append(edges, edge{u, v})
	}
	return n, edges
}

func runCaseC(bin string, n int, edges []edge) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&input, "%d %d\n", e.u, e.v)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	expected := fmt.Sprintf("%d", solveC(n, edges))
	if strings.TrimSpace(string(out)) != expected {
		return fmt.Errorf("expected %s got %s", expected, strings.TrimSpace(string(out)))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, e := genCaseC(rng)
		if err := runCaseC(bin, n, e); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
