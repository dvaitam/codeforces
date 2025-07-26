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
	n     int
	edges []Edge
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
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
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "955F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genTree(rng *rand.Rand, n int) []Edge {
	edges := make([]Edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, Edge{p, i})
	}
	return edges
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(5))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		edges := genTree(rng, n)
		tests = append(tests, Test{n: n, edges: edges})
	}
	tests = append(tests,
		Test{n: 1, edges: nil},
		Test{n: 2, edges: []Edge{{1, 2}}},
		Test{n: 3, edges: []Edge{{1, 2}, {1, 3}}},
		Test{n: 4, edges: []Edge{{1, 2}, {2, 3}, {3, 4}}},
		Test{n: 5, edges: []Edge{{1, 2}, {1, 3}, {3, 4}, {4, 5}}},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:%sexpected:%sgot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
