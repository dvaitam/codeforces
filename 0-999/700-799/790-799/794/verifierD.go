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
	cmd := exec.Command("go", "build", "-o", ref, "794D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genGraph(n, m int) []Edge {
	edges := make([]Edge, 0, m)
	set := map[[2]int]struct{}{}
	// ensure connectivity with a tree
	for i := 2; i <= n; i++ {
		u := i - 1
		v := i
		edges = append(edges, Edge{u, v})
		set[[2]int{u, v}] = struct{}{}
		set[[2]int{v, u}] = struct{}{}
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
		set[[2]int{u, v}] = struct{}{}
		set[[2]int{v, u}] = struct{}{}
		edges = append(edges, Edge{u, v})
	}
	return edges
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2 // at least 2 nodes
		maxM := n * (n - 1) / 2
		extra := rand.Intn(maxM - (n - 1) + 1)
		m := (n - 1) + extra
		edges := genGraph(n, m)
		tests = append(tests, Test{n: n, m: m, edges: edges})
	}
	return tests
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
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
