package main

import (
	"bufio"
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

func solveC(deg, xor []int) []edge {
	n := len(deg)
	q := make([]int, 0, n)
	for i, d := range deg {
		if d == 1 {
			q = append(q, i)
		}
	}
	edges := make([]edge, 0)
	head := 0
	for head < len(q) {
		u := q[head]
		head++
		if deg[u] != 1 {
			continue
		}
		v := xor[u]
		edges = append(edges, edge{u, v})
		deg[u]--
		xor[u] ^= v
		deg[v]--
		xor[v] ^= u
		if deg[v] == 1 {
			q = append(q, v)
		}
	}
	return edges
}

func genTree(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(9) + 2 // at least 2
	edges := make([]edge, 0, n-1)
	for i := 1; i < n; i++ {
		j := rng.Intn(i)
		edges = append(edges, edge{i, j})
	}
	return n, edges
}

func edgesToInput(n int, edges []edge) (string, []int, []int) {
	deg := make([]int, n)
	xor := make([]int, n)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
		xor[e.u] ^= e.v
		xor[e.v] ^= e.u
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", deg[i], xor[i]))
	}
	return sb.String(), deg, xor
}

func parseEdges(out string) ([]edge, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return nil, fmt.Errorf("no output")
	}
	var m int
	if _, err := fmt.Sscan(scanner.Text(), &m); err != nil {
		return nil, err
	}
	edges := make([]edge, 0, m)
	for scanner.Scan() {
		var a, b int
		if _, err := fmt.Sscan(scanner.Text(), &a, &b); err != nil {
			return nil, err
		}
		edges = append(edges, edge{a, b})
	}
	if len(edges) != m {
		return nil, fmt.Errorf("expected %d edges got %d", m, len(edges))
	}
	return edges, nil
}

func normalize(es []edge) []edge {
	res := make([]edge, len(es))
	for i, e := range es {
		if e.u > e.v {
			e.u, e.v = e.v, e.u
		}
		res[i] = e
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].u == res[j].u {
			return res[i].v < res[j].v
		}
		return res[i].u < res[j].u
	})
	return res
}

func runCase(bin string, input string, expected []edge) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got, err := parseEdges(out.String())
	if err != nil {
		return err
	}
	nexp := normalize(expected)
	ngot := normalize(got)
	if len(nexp) != len(ngot) {
		return fmt.Errorf("expected %d edges got %d", len(nexp), len(ngot))
	}
	for i := range nexp {
		if nexp[i] != ngot[i] {
			return fmt.Errorf("edge mismatch")
		}
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
		n, edges := genTree(rng)
		input, deg, xor := edgesToInput(n, edges)
		exp := solveC(append([]int(nil), deg...), append([]int(nil), xor...))
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
