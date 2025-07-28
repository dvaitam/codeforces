package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct{ u, v, w int }
type test struct {
	n     int
	edges []edge
	set   []int
	k     float64
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildOracle() (string, error) {
	oracle := "./oracleB"
	cmd := exec.Command("go", "build", "-o", oracle, "1723B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return oracle, nil
}

func genTests() []test {
	rand.Seed(2)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Intn(8) + 2
		maxM := n * (n - 1) / 2
		m := n - 1 + rand.Intn(maxM-(n-1)+1)
		edges := make([]edge, 0, m)
		used := make(map[[2]int]bool)
		for j := 0; j < n-1; j++ {
			w := rand.Intn(100) + 1
			edges = append(edges, edge{j, j + 1, w})
			used[[2]int{j, j + 1}] = true
		}
		for len(edges) < m {
			u := rand.Intn(n)
			v := rand.Intn(n)
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			if used[[2]int{u, v}] {
				continue
			}
			w := rand.Intn(100) + 1
			edges = append(edges, edge{u, v, w})
			used[[2]int{u, v}] = true
		}
		setSize := rand.Intn(n) + 1
		perm := rand.Perm(n)[:setSize]
		nodes := make([]int, setSize)
		copy(nodes, perm)
		k := rand.Float64() * 1000
		tests[i] = test{n, edges, nodes, k}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}
		input.WriteString(fmt.Sprintf("%d %.6f\n", len(tc.set), tc.k))
		for _, v := range tc.set {
			input.WriteString(fmt.Sprintf("%d\n", v))
		}
		expect, err := run(oracle, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
