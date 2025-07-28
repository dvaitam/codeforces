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
}

func run(bin string, input string) (string, error) {
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
	oracle := "./oracleA"
	cmd := exec.Command("go", "build", "-o", oracle, "1723A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return oracle, nil
}

func genTests() []test {
	rand.Seed(1)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Intn(8) + 2 // 2..9
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
		tests[i] = test{n, edges}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
