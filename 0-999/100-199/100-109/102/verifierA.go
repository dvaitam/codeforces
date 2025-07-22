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

type testCaseA struct {
	n, m   int
	prices []int
	edges  [][2]int
}

func expectedA(tc testCaseA) string {
	n := tc.n
	con := make([][]bool, n)
	for i := range con {
		con[i] = make([]bool, n)
	}
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		con[u][v] = true
		con[v][u] = true
	}
	best := int(1 << 60)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if !con[i][j] {
				continue
			}
			for k := j + 1; k < n; k++ {
				if con[i][k] && con[j][k] {
					sum := tc.prices[i] + tc.prices[j] + tc.prices[k]
					if sum < best {
						best = sum
					}
				}
			}
		}
	}
	if best == int(1<<60) {
		return "-1"
	}
	return fmt.Sprintf("%d", best)
}

func generateCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 3 // 3..10
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	prices := make([]int, n)
	for i := range prices {
		prices[i] = rng.Intn(20) + 1
	}
	edgeSet := make(map[[2]int]struct{})
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := edgeSet[key]; ok {
			continue
		}
		edgeSet[key] = struct{}{}
		edges = append(edges, key)
	}
	tc := testCaseA{n: n, m: m, prices: prices, edges: edges}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, p := range prices {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	input := sb.String()
	expected := expectedA(tc)
	return input, expected
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseA(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
