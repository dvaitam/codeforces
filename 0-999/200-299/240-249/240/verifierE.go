package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edgeData struct {
	from, to   int
	needRepair bool
}

type testCase struct {
	input  string
	expect int
	n      int
	edges  []edgeData
}

type refEdge struct {
	from, to int
	cost     int
}

const inf = int(1e9)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkAnswer(tc, out); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\nexpected minimal repairs: %d\n", i+1, err, tc.input, out, tc.expect)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkAnswer(tc testCase, output string) error {
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		return fmt.Errorf("empty output")
	}
	tokens := strings.Fields(output)
	if tokens[0] == "-1" {
		if tc.expect != -1 {
			return fmt.Errorf("expected %d repairs but got -1", tc.expect)
		}
		if len(tokens) != 1 {
			return fmt.Errorf("unexpected extra tokens after -1")
		}
		return nil
	}
	if tc.expect == -1 {
		return fmt.Errorf("expected -1 but got %s", tokens[0])
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil || k < 0 {
		return fmt.Errorf("invalid number of repairs %q", tokens[0])
	}
	if len(tokens) != k+1 {
		return fmt.Errorf("reported %d roads but provided %d identifiers", k, len(tokens)-1)
	}
	if k != tc.expect {
		return fmt.Errorf("expected %d repairs but reported %d", tc.expect, k)
	}
	selected := make(map[int]struct{}, k)
	for i := 0; i < k; i++ {
		id, err := strconv.Atoi(tokens[i+1])
		if err != nil {
			return fmt.Errorf("invalid road index %q", tokens[i+1])
		}
		if id < 1 || id > len(tc.edges) {
			return fmt.Errorf("road index %d out of range", id)
		}
		if _, ok := selected[id]; ok {
			return fmt.Errorf("road index %d listed multiple times", id)
		}
		if !tc.edges[id-1].needRepair {
			return fmt.Errorf("road %d is already good and cannot be repaired", id)
		}
		selected[id] = struct{}{}
	}
	if !allReachable(tc.n, tc.edges, selected) {
		return fmt.Errorf("not all cities reachable from the capital with reported repairs")
	}
	return nil
}

func allReachable(n int, edges []edgeData, repaired map[int]struct{}) bool {
	adj := make([][]int, n+1)
	for idx, e := range edges {
		if !e.needRepair {
			adj[e.from] = append(adj[e.from], e.to)
			continue
		}
		if _, ok := repaired[idx+1]; ok {
			adj[e.from] = append(adj[e.from], e.to)
		}
	}
	vis := make([]bool, n+1)
	queue := []int{1}
	vis[1] = true
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range adj[u] {
			if !vis[v] {
				vis[v] = true
				queue = append(queue, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return false
		}
	}
	return true
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []testCase {
	rand.Seed(42)
	var tests []testCase
	tests = append(tests, newTestCase(1, nil))
	tests = append(tests, newTestCase(2, []edgeData{
		{from: 1, to: 2, needRepair: false},
	}))
	tests = append(tests, newTestCase(2, nil))
	tests = append(tests, newTestCase(3, []edgeData{
		{from: 1, to: 2, needRepair: true},
		{from: 2, to: 3, needRepair: false},
	}))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTestCase())
	}
	return tests
}

func randomTestCase() testCase {
	n := rand.Intn(10) + 1
	maxPossible := n * (n - 1)
	limit := n*3 + rand.Intn(n+1)
	if limit > maxPossible {
		limit = maxPossible
	}
	m := 0
	if limit > 0 {
		m = rand.Intn(limit + 1)
	}
	edges := make([]edgeData, 0, m)
	used := make(map[int]struct{})
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		key := (u-1)*n + (v - 1)
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		edges = append(edges, edgeData{
			from:       u,
			to:         v,
			needRepair: rand.Intn(2) == 1,
		})
	}
	return newTestCase(n, edges)
}

func newTestCase(n int, edges []edgeData) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		c := 0
		if e.needRepair {
			c = 1
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.from, e.to, c))
	}
	edgesCopy := make([]edgeData, len(edges))
	copy(edgesCopy, edges)
	expect := calcMinRepairs(n, edgesCopy)
	return testCase{
		input:  sb.String(),
		expect: expect,
		n:      n,
		edges:  edgesCopy,
	}
}

func calcMinRepairs(n int, edges []edgeData) int {
	if n == 0 {
		return 0
	}
	refEdges := make([]refEdge, 0, len(edges))
	for _, e := range edges {
		c := 0
		if e.needRepair {
			c = 1
		}
		refEdges = append(refEdges, refEdge{
			from: e.from - 1,
			to:   e.to - 1,
			cost: c,
		})
	}
	cost, ok := directedMST(n, 0, refEdges)
	if !ok {
		return -1
	}
	return cost
}

func directedMST(n, root int, edges []refEdge) (int, bool) {
	total := 0
	for {
		in := make([]int, n)
		pre := make([]int, n)
		id := make([]int, n)
		vis := make([]int, n)
		for i := range in {
			in[i] = inf
			id[i] = -1
			vis[i] = -1
		}
		for _, e := range edges {
			if e.from == e.to {
				continue
			}
			if e.cost < in[e.to] {
				in[e.to] = e.cost
				pre[e.to] = e.from
			}
		}
		in[root] = 0
		for v := 0; v < n; v++ {
			if v == root {
				continue
			}
			if in[v] == inf {
				return 0, false
			}
		}
		cnt := 0
		for v := 0; v < n; v++ {
			total += in[v]
			u := v
			for vis[u] != v && id[u] == -1 && u != root {
				vis[u] = v
				u = pre[u]
			}
			if u != root && id[u] == -1 {
				for x := pre[u]; x != u; x = pre[x] {
					id[x] = cnt
				}
				id[u] = cnt
				cnt++
			}
		}
		if cnt == 0 {
			return total, true
		}
		for v := 0; v < n; v++ {
			if id[v] == -1 {
				id[v] = cnt
				cnt++
			}
		}
		newEdges := make([]refEdge, 0, len(edges))
		for _, e := range edges {
			u := id[e.from]
			v := id[e.to]
			if u != v {
				newEdges = append(newEdges, refEdge{
					from: u,
					to:   v,
					cost: e.cost - in[e.to],
				})
			}
		}
		root = id[root]
		n = cnt
		edges = newEdges
	}
}
