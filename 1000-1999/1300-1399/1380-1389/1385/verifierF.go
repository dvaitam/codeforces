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

type testCase struct {
	n, k     int
	edges    [][2]int
	input    string
	expected string
}

func solveCase(n, k int, edges [][2]int) int {
	if k == 1 {
		return n - 1
	}
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	leafVec := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			p := adj[i][0]
			leafVec[p] = append(leafVec[p], i)
		}
	}
	leafCnt := make([]int, n+1)
	for i := 1; i <= n; i++ {
		leafCnt[i] = len(leafVec[i])
	}
	queue := make([]int, 0)
	for i := 1; i <= n; i++ {
		if leafCnt[i] >= k {
			queue = append(queue, i)
		}
	}
	head := 0
	ans := 0
	for head < len(queue) {
		v := queue[head]
		head++
		if leafCnt[v] < k {
			continue
		}
		times := leafCnt[v] / k
		ans += times
		for i := 0; i < times*k; i++ {
			u := leafVec[v][len(leafVec[v])-1]
			leafVec[v] = leafVec[v][:len(leafVec[v])-1]
			deg[u] = 0
			leafCnt[v]--
		}
		deg[v] -= times * k
		if deg[v] == 1 {
			var parent int
			for _, to := range adj[v] {
				if deg[to] > 0 {
					parent = to
					break
				}
			}
			if parent > 0 {
				leafVec[parent] = append(leafVec[parent], v)
				leafCnt[parent]++
				if leafCnt[parent] >= k {
					queue = append(queue, parent)
				}
			}
		}
		if leafCnt[v] >= k {
			queue = append(queue, v)
		}
	}
	return ans
}

func buildCase(n, k int, edges [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	res := solveCase(n, k, edges)
	return testCase{n: n, k: k, edges: edges, input: sb.String(), expected: fmt.Sprintf("%d\n", res)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(15) + 2
	k := rng.Intn(n-1) + 1
	edges := make([][2]int, 0, n-1)
	// generate random tree
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return buildCase(n, k, edges)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.expected)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
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
	cases := make([]testCase, 0)
	cases = append(cases, buildCase(2, 1, [][2]int{{1, 2}}))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
