package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct {
	u, v  int
	color byte
}

type testCaseB struct {
	n     int
	edges []edge
}

func solveB(tc testCaseB) (int, []int) {
	n := tc.n
	G := make([][]edge, n+1)
	for _, e := range tc.edges {
		G[e.u] = append(G[e.u], edge{e.v, 0, e.color})
		G[e.v] = append(G[e.v], edge{e.u, 0, e.color})
	}
	const INF = int(1e9)
	resR, resB := 0, 0
	wrongR, wrongB := false, false
	Rcol := make([]int, n+1)
	Bcol := make([]int, n+1)
	visR := make([]bool, n+1)
	visB := make([]bool, n+1)
	var resultR, resultB []int
	// R-coloring
	for i := 1; i <= n; i++ {
		if !visR[i] {
			queue := []int{i}
			visR[i] = true
			Rcol[i] = 0
			comp := []int{i}
			for qi := 0; qi < len(queue); qi++ {
				v := queue[qi]
				for _, e := range G[v] {
					if !visR[e.u] {
						if e.color == 'R' {
							Rcol[e.u] = Rcol[v]
						} else {
							Rcol[e.u] = Rcol[v] ^ 1
						}
						visR[e.u] = true
						queue = append(queue, e.u)
						comp = append(comp, e.u)
					}
				}
			}
			var part0, part1 []int
			for _, v := range comp {
				if Rcol[v] == 0 {
					part0 = append(part0, v)
				} else {
					part1 = append(part1, v)
				}
			}
			if len(part0) < len(part1) {
				resultR = append(resultR, part0...)
				resR += len(part0)
			} else {
				resultR = append(resultR, part1...)
				resR += len(part1)
			}
		}
	}
	// B-coloring
	for i := 1; i <= n; i++ {
		if !visB[i] {
			queue := []int{i}
			visB[i] = true
			Bcol[i] = 0
			comp := []int{i}
			for qi := 0; qi < len(queue); qi++ {
				v := queue[qi]
				for _, e := range G[v] {
					if !visB[e.u] {
						if e.color == 'B' {
							Bcol[e.u] = Bcol[v]
						} else {
							Bcol[e.u] = Bcol[v] ^ 1
						}
						visB[e.u] = true
						queue = append(queue, e.u)
						comp = append(comp, e.u)
					}
				}
			}
			var part0, part1 []int
			for _, v := range comp {
				if Bcol[v] == 0 {
					part0 = append(part0, v)
				} else {
					part1 = append(part1, v)
				}
			}
			if len(part0) < len(part1) {
				resultB = append(resultB, part0...)
				resB += len(part0)
			} else {
				resultB = append(resultB, part1...)
				resB += len(part1)
			}
		}
	}
	for u := 1; u <= n; u++ {
		for _, e := range G[u] {
			v := e.u
			if e.color == 'R' && Rcol[u] != Rcol[v] {
				wrongR = true
			}
			if e.color == 'B' && Rcol[u] == Rcol[v] {
				wrongR = true
			}
			if e.color == 'B' && Bcol[u] != Bcol[v] {
				wrongB = true
			}
			if e.color == 'R' && Bcol[u] == Bcol[v] {
				wrongB = true
			}
		}
	}
	if wrongR && wrongB {
		return -1, nil
	}
	if wrongB {
		resB = INF
	}
	if wrongR {
		resR = INF
	}
	if resR < resB {
		return resR, resultR
	}
	return resB, resultB
}

func applyMoves(tc testCaseB, moves []int) []edge {
	// copy edges
	edges := make([]edge, len(tc.edges))
	copy(edges, tc.edges)
	for _, v := range moves {
		for i := range edges {
			if edges[i].u == v || edges[i].v == v {
				if edges[i].color == 'R' {
					edges[i].color = 'B'
				} else {
					edges[i].color = 'R'
				}
			}
		}
	}
	return edges
}

func allSameColor(edges []edge) bool {
	if len(edges) == 0 {
		return true
	}
	c := edges[0].color
	for _, e := range edges {
		if e.color != c {
			return false
		}
	}
	return true
}

func runCaseB(bin string, tc testCaseB, minimal int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %c\n", e.u, e.v, e.color))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	if fields[0] == "-1" {
		if minimal != -1 {
			return fmt.Errorf("expected solution of size %d", minimal)
		}
		return nil
	}
	var k int
	fmt.Sscan(fields[0], &k)
	if k != minimal {
		return fmt.Errorf("expected size %d got %d", minimal, k)
	}
	if len(fields)-1 != k {
		return fmt.Errorf("expected %d vertices got %d", k, len(fields)-1)
	}
	moves := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Sscan(fields[i+1], &moves[i])
	}
	edges := applyMoves(tc, moves)
	if !allSameColor(edges) {
		return fmt.Errorf("moves do not unify edge colors")
	}
	return nil
}

func genCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	seen := make(map[[2]int]bool)
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
		key := [2]int{u, v}
		if seen[key] {
			continue
		}
		seen[key] = true
		c := byte('R')
		if rng.Intn(2) == 0 {
			c = 'B'
		}
		edges = append(edges, edge{u, v, c})
	}
	return testCaseB{n, edges}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseB(rng)
		minK, _ := solveB(tc)
		if err := runCaseB(bin, tc, minK); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
