package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
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

func solve(n, m, k int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n+1)
	c := make([]int, n+2)
	cnt := 1
	cur := 1
	visited[1] = true
	c[1] = 1
	for {
		found := false
		for _, nb := range adj[cur] {
			if !visited[nb] {
				visited[nb] = true
				cnt++
				c[cnt] = nb
				cur = nb
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	end := c[cnt]
	for i := 1; i <= cnt; i++ {
		for _, nb := range adj[c[i]] {
			if nb == end {
				length := cnt - i + 1
				var out bytes.Buffer
				fmt.Fprintln(&out, length)
				for j := i; j <= cnt; j++ {
					if j > i {
						out.WriteByte(' ')
					}
					fmt.Fprintf(&out, "%d", c[j])
				}
				return strings.TrimSpace(out.String())
			}
		}
	}
	return ""
}

func generateCases() []testCase {
	rand.Seed(4)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		k := rand.Intn(3) + 2
		n := k + rand.Intn(3) + 1
		if n < k+1 {
			n = k + 1
		}
		edges := make([][2]int, 0, n*(n-1)/2)
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				edges = append(edges, [2]int{i, j})
			}
		}
		m := len(edges)
		rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d\n", n, m, k)
		for _, e := range edges {
			fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
		}
		expected := solve(n, m, k, edges)
		cases[t] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
