package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const inf = 1000000000

func bfs(start int, g [][]int) []int {
	n := len(g)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = inf
	}
	dist[start] = 0
	q := list.New()
	q.PushBack(start)
	for q.Len() > 0 {
		u := q.Remove(q.Front()).(int)
		for _, v := range g[u] {
			if dist[v] > dist[u]+1 {
				dist[v] = dist[u] + 1
				q.PushBack(v)
			}
		}
	}
	return dist
}

func solveCase(n, m, a, b int, edges [][2]int, buses [][2]int) string {
	g := make([][]int, n)
	for _, e := range edges {
		g[e[0]-1] = append(g[e[0]-1], e[1]-1)
	}
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(i, g)
	}
	busAdj := make([][]int, n)
	for _, pr := range buses {
		s := pr[0] - 1
		t := pr[1] - 1
		d := dist[s][t]
		if d >= inf {
			continue
		}
		for v := 0; v < n; v++ {
			if dist[s][v]+dist[v][t] == d {
				busAdj[v] = append(busAdj[v], t)
			}
		}
	}
	da := make([]int, n)
	for i := range da {
		da[i] = inf
	}
	da[a-1] = 0
	q := list.New()
	q.PushBack(a - 1)
	for q.Len() > 0 {
		u := q.Remove(q.Front()).(int)
		for _, v := range busAdj[u] {
			if da[v] > da[u]+1 {
				da[v] = da[u] + 1
				q.PushBack(v)
			}
		}
	}
	if da[b-1] >= inf {
		return "-1"
	}
	return fmt.Sprintf("%d", da[b-1])
}

type test struct{ input, expected string }

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	fixed := []struct {
		n, m, a, b int
		edges      [][2]int
		buses      [][2]int
	}{
		{2, 1, 1, 2, [][2]int{{1, 2}}, [][2]int{{1, 2}}},
		{3, 3, 1, 3, [][2]int{{1, 2}, {2, 3}, {1, 3}}, [][2]int{{1, 3}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", f.n, f.m, f.a, f.b))
		for _, e := range f.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", len(f.buses)))
		for _, pr := range f.buses {
			sb.WriteString(fmt.Sprintf("%d %d\n", pr[0], pr[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(f.n, f.m, f.a, f.b, f.edges, f.buses)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 2
		m := rng.Intn(n*(n-1)) + 1
		edgeSet := make(map[[2]int]bool)
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			e := [2]int{u, v}
			if edgeSet[e] {
				continue
			}
			edgeSet[e] = true
			edges = append(edges, e)
		}
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		k := rng.Intn(3) + 1
		buses := make([][2]int, k)
		for i := 0; i < k; i++ {
			s := rng.Intn(n) + 1
			t := rng.Intn(n) + 1
			buses[i] = [2]int{s, t}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, len(edges), a, b))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", k))
		for _, pr := range buses {
			sb.WriteString(fmt.Sprintf("%d %d\n", pr[0], pr[1]))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solveCase(n, len(edges), a, b, edges, buses)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
