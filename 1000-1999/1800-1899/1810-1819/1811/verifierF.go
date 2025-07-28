package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
)

type testF struct {
	n     int
	m     int
	edges [][2]int
}

func genTests() []testF {
	rand.Seed(181106)
	tests := make([]testF, 100)
	for i := range tests {
		k := rand.Intn(4) + 3 // k from 3..6
		n := k * k
		m := n + k
		edgeSet := make(map[[2]int]struct{})
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			u := rand.Intn(n)
			v := rand.Intn(n)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			pair := [2]int{u, v}
			if _, ok := edgeSet[pair]; ok {
				continue
			}
			edgeSet[pair] = struct{}{}
			edges = append(edges, pair)
		}
		tests[i] = testF{n: n, m: m, edges: edges}
	}
	return tests
}

func solve(tc testF) string {
	n := tc.n
	m := tc.m
	adj := make([][]int, n)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	k := int(math.Round(math.Sqrt(float64(n))))
	if k*k != n || k < 3 || m != n+k {
		return "NO"
	}

	deg := make([]int, n)
	for i := 0; i < n; i++ {
		deg[i] = len(adj[i])
	}
	central := []int{}
	isCentral := make([]bool, n)
	valid := true
	for i, d := range deg {
		if d == 4 {
			central = append(central, i)
			isCentral[i] = true
		} else if d != 2 {
			valid = false
			break
		}
	}
	if !valid || len(central) != k {
		return "NO"
	}

	centralEdges := map[[2]int]struct{}{}
	for _, u := range central {
		cnt := 0
		for _, v := range adj[u] {
			if isCentral[v] {
				cnt++
				if u < v {
					centralEdges[[2]int{u, v}] = struct{}{}
				}
			}
		}
		if cnt != 2 {
			valid = false
			break
		}
	}
	if !valid || len(centralEdges) != k {
		return "NO"
	}

	visited := make([]bool, n)
	queue := []int{central[0]}
	visited[central[0]] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if isCentral[v] && !visited[v] {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}
	cntVisited := 0
	for _, u := range central {
		if visited[u] {
			cntVisited++
		}
	}
	if cntVisited != k {
		return "NO"
	}

	visitedAll := make([]bool, n)
	compCount := 0
	for i := 0; i < n; i++ {
		if !visitedAll[i] {
			compCount++
			q := []int{i}
			visitedAll[i] = true
			size := 0
			centralCnt := 0
			for len(q) > 0 {
				u := q[0]
				q = q[1:]
				size++
				if isCentral[u] {
					centralCnt++
				}
				for _, v := range adj[u] {
					if isCentral[u] && isCentral[v] {
						continue
					}
					if !visitedAll[v] {
						visitedAll[v] = true
						q = append(q, v)
					}
				}
			}
			if centralCnt != 1 || size != k {
				valid = false
				break
			}
		}
	}
	if !valid || compCount != k {
		return "NO"
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0]+1, e[1]+1)
		}
	}

	expected := make([]string, len(tests))
	for i, tc := range tests {
		expected[i] = solve(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		got := scanner.Text()
		if got != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
