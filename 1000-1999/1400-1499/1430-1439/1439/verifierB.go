package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct{ u, v int }

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func brute(n int, edges []Edge, k int) (hasClique bool, cliques [][]int, subsets [][]int) {
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}
	for _, e := range edges {
		adj[e.u][e.v] = true
		adj[e.v][e.u] = true
	}
	for mask := 1; mask < (1 << n); mask++ {
		verts := []int{}
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				verts = append(verts, i)
			}
		}
		if len(verts) == k {
			ok := true
			for i := 0; i < k && ok; i++ {
				for j := i + 1; j < k; j++ {
					if !adj[verts[i]][verts[j]] {
						ok = false
						break
					}
				}
			}
			if ok {
				hasClique = true
				cliques = append(cliques, append([]int(nil), verts...))
			}
		}
		good := len(verts) > 0
		if good {
			for _, v := range verts {
				deg := 0
				for _, u := range verts {
					if u != v && adj[v][u] {
						deg++
					}
				}
				if deg < k {
					good = false
					break
				}
			}
		}
		if good {
			subsets = append(subsets, append([]int(nil), verts...))
		}
	}
	return
}

func parseOutput(out string, n int, edges []Edge, k int) error {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return fmt.Errorf("no output")
	}
	t := sc.Text()
	if t == "-1" {
		// verify none exists
		hasC, _, sub := brute(n, edges, k)
		if hasC || len(sub) > 0 {
			return fmt.Errorf("answer exists but got -1")
		}
		return nil
	}
	tp, err := strconv.Atoi(t)
	if err != nil {
		return fmt.Errorf("invalid first token")
	}
	switch tp {
	case 1:
		if !sc.Scan() {
			return fmt.Errorf("missing size")
		}
		sz, _ := strconv.Atoi(sc.Text())
		verts := make([]int, sz)
		for i := 0; i < sz; i++ {
			if !sc.Scan() {
				return fmt.Errorf("missing vertex")
			}
			v, _ := strconv.Atoi(sc.Text())
			if v < 1 || v > n {
				return fmt.Errorf("vertex out of range")
			}
			verts[i] = v - 1
		}
		// check subset property
		for _, v := range verts {
			deg := 0
			for _, u := range verts {
				if u != v && isEdge(edges, v, u) {
					deg++
				}
			}
			if deg < k {
				return fmt.Errorf("subset property fail")
			}
		}
	case 2:
		verts := make([]int, k)
		for i := 0; i < k; i++ {
			if !sc.Scan() {
				return fmt.Errorf("missing vertex")
			}
			v, _ := strconv.Atoi(sc.Text())
			if v < 1 || v > n {
				return fmt.Errorf("vertex out of range")
			}
			verts[i] = v - 1
		}
		// check clique
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				if !isEdge(edges, verts[i], verts[j]) {
					return fmt.Errorf("not clique")
				}
			}
		}
	default:
		return fmt.Errorf("invalid type %v", tp)
	}
	return nil
}

func isEdge(edges []Edge, a, b int) bool {
	if a > b {
		a, b = b, a
	}
	for _, e := range edges {
		if e.u == a && e.v == b {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/bin")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(4) + 2 //2..5
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges)
		k := rand.Intn(n) + 1
		edgeMap := make(map[[2]int]bool)
		edges := make([]Edge, 0, m)
		for len(edges) < m {
			a := rand.Intn(n)
			b := rand.Intn(n)
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			if edgeMap[[2]int{a, b}] {
				continue
			}
			edgeMap[[2]int{a, b}] = true
			edges = append(edges, Edge{a, b})
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), k))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u+1, e.v+1))
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d exec error %v\n", tcase, err)
			return
		}
		if err := parseOutput(out, n, edges, k); err != nil {
			fmt.Printf("test %d failed: %v\n", tcase, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
