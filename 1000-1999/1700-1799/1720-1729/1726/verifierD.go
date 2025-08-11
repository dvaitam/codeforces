package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct{ u, v int }

type Test struct {
	n, m  int
	edges []Edge
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1726D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genGraph(n, m int) []Edge {
	edges := make([]Edge, 0, m)
	// tree edges
	for i := 2; i <= n; i++ {
		edges = append(edges, Edge{i - 1, i})
	}
	set := map[[2]int]struct{}{}
	for _, e := range edges {
		set[[2]int{e.u, e.v}] = struct{}{}
		set[[2]int{e.v, e.u}] = struct{}{}
	}
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		if _, ok := set[[2]int{u, v}]; ok {
			continue
		}
		edges = append(edges, Edge{u, v})
		set[[2]int{u, v}] = struct{}{}
		set[[2]int{v, u}] = struct{}{}
	}
	return edges
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		maxM := n * (n - 1) / 2
		m := n - 1 + rand.Intn(min(3, maxM-(n-1))+1)
		if m > maxM {
			m = maxM
		}
		edges := genGraph(n, m)
		tests = append(tests, Test{n: n, m: m, edges: edges})
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func components(n int, edges []Edge, s string, color byte) int {
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		fa, fb := find(a), find(b)
		if fa != fb {
			parent[fb] = fa
		}
	}
	for i, e := range edges {
		if s[i] == color {
			union(e.u, e.v)
		}
	}
	cnt := 0
	for i := 1; i <= n; i++ {
		if find(i) == i {
			cnt++
		}
	}
	return cnt
}

func sumComponents(n int, edges []Edge, s string) int {
	return components(n, edges, s, '0') + components(n, edges, s, '1')
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp = strings.TrimSpace(exp)
		got = strings.TrimSpace(got)
		if len(got) != tc.m {
			fmt.Printf("Test %d failed\nInput:%sExpected length %d Got length %d\n", i+1, input, tc.m, len(got))
			os.Exit(1)
		}
		if strings.IndexFunc(got, func(r rune) bool { return r != '0' && r != '1' }) != -1 {
			fmt.Printf("Test %d failed\nInput:%sOutput contains invalid characters: %s\n", i+1, input, got)
			os.Exit(1)
		}
		expVal := sumComponents(tc.n, tc.edges, exp)
		gotVal := sumComponents(tc.n, tc.edges, got)
		if expVal != gotVal {
			fmt.Printf("Test %d failed\nInput:%sExpected value %d Got value %d\nExpected:%s\nGot:%s\n", i+1, input, expVal, gotVal, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
