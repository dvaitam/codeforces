package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(6) + 5 // 5 to 10
		perm := make([]int, n)
		perm[0] = 1
		p := rand.Perm(n - 1)
		for i, v := range p {
			perm[i+1] = v + 2
		}
		edgeMap := map[[2]int]bool{}
		for i := 0; i < n; i++ {
			a := perm[i]
			b := perm[(i+1)%n]
			if a > b {
				a, b = b, a
			}
			edgeMap[[2]int{a, b}] = true
			a2 := perm[i]
			b2 := perm[(i+2)%n]
			if a2 > b2 {
				a2, b2 = b2, a2
			}
			edgeMap[[2]int{a2, b2}] = true
		}
		edges := make([][2]int, 0, len(edgeMap))
		for e := range edgeMap {
			edges = append(edges, e)
		}
		rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
		var buf bytes.Buffer
		fmt.Fprintln(&buf, n)
		for _, e := range edges {
			fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
		}
		cases[t] = testCase{input: buf.String()}
	}
	return cases
}

func verify(input, output string) error {
	// Parse input
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scanInt := func() (int, bool) {
		if scanner.Scan() {
			v, _ := strconv.Atoi(scanner.Text())
			return v, true
		}
		return 0, false
	}

	n, ok := scanInt()
	if !ok {
		return fmt.Errorf("failed to parse n")
	}

	adj := make(map[int]map[int]bool)
	for {
		u, ok := scanInt()
		if !ok {
			break
		}
		v, ok := scanInt()
		if !ok {
			break
		}
		if adj[u] == nil {
			adj[u] = make(map[int]bool)
		}
		if adj[v] == nil {
			adj[v] = make(map[int]bool)
		}
		adj[u][v] = true
		adj[v][u] = true
	}

	// Parse output
	outFields := strings.Fields(output)
	if len(outFields) == 0 {
		return fmt.Errorf("empty output")
	}
	if len(outFields) == 1 && outFields[0] == "-1" {
		return fmt.Errorf("output is -1, but a solution is guaranteed to exist")
	}

	if len(outFields) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(outFields))
	}

	p := make([]int, n)
	seen := make(map[int]bool)
	for i, s := range outFields {
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid number: %s", s)
		}
		if v < 1 || v > n {
			return fmt.Errorf("number out of range: %d", v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate number: %d", v)
		}
		seen[v] = true
		p[i] = v
	}

	// Check edges
	hasEdge := func(u, v int) bool {
		return adj[u][v]
	}

	for i := 0; i < n; i++ {
		u := p[i]
		v1 := p[(i+1)%n]
		v2 := p[(i+2)%n]

		if !hasEdge(u, v1) {
			return fmt.Errorf("edge %d-%d (dist 1) missing in input", u, v1)
		}
		if !hasEdge(u, v2) {
			return fmt.Errorf("edge %d-%d (dist 2) missing in input", u, v2)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc.input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\nactual:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}