package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refPath = "2000-2999/2100-2199/2130-2139/2137/2137G.go"

type edge struct{ u, v int }

type testCase struct {
	n, m, q int
	edges   []edge
	queries [][2]int
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 2 // 2..11
	order := rng.Perm(n)  // topo order
	pos := make([]int, n)
	for i, v := range order {
		pos[v] = i
	}
	var edges []edge
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float64() < 0.25 {
				u := order[i]
				v := order[j]
				edges = append(edges, edge{u, v})
			}
		}
	}
	if len(edges) == 0 {
		// ensure at least one edge
		u := order[0]
		v := order[len(order)-1]
		if u == v {
			u = 0
			v = 1
		}
		edges = append(edges, edge{u, v})
	}
	m := len(edges)

	q := rng.Intn(40) + 1
	red := make([]bool, n)
	queries := make([][2]int, 0, q)
	for len(queries) < q {
		if rng.Float64() < 0.35 {
			// type 1 if possible
			candidates := make([]int, 0, n)
			for i := 0; i < n; i++ {
				if !red[i] {
					candidates = append(candidates, i)
				}
			}
			if len(candidates) == 0 {
				continue
			}
			u := candidates[rng.Intn(len(candidates))]
			red[u] = true
			queries = append(queries, [2]int{1, u})
		} else {
			u := rng.Intn(n)
			queries = append(queries, [2]int{2, u})
		}
	}
	return testCase{n: n, m: m, q: len(queries), edges: edges, queries: queries}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.q)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u+1, e.v+1)
		}
		for _, qu := range tc.queries {
			fmt.Fprintf(&sb, "%d %d\n", qu[0], qu[1]+1)
		}
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout running %s", path)
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/2137G_binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input := buildInput(cases)

	refAbs, _ := filepath.Abs(refPath)
	candAbs := cand
	if !filepath.IsAbs(candAbs) {
		candAbs, _ = filepath.Abs(cand)
	}

	expected, err := runProgram(refAbs, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run reference: %v\n", err)
		os.Exit(1)
	}
	actual, err := runProgram(candAbs, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	expTokens := strings.Fields(expected)
	actTokens := strings.Fields(actual)
	if len(actTokens) < len(expTokens) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(actTokens), len(expTokens))
		os.Exit(1)
	}
	for i := range expTokens {
		if strings.ToUpper(expTokens[i]) != strings.ToUpper(actTokens[i]) {
			fmt.Fprintf(os.Stderr, "mismatch at answer %d: expected %s got %s\n", i+1, expTokens[i], actTokens[i])
			os.Exit(1)
		}
	}
	if len(actTokens) != len(expTokens) {
		fmt.Fprintf(os.Stderr, "extra output tokens detected\n")
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}
