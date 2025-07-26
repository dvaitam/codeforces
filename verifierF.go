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

type Test struct {
	input  string
	output string
}

func solveF(n int, edges [][2]int) string {
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	dist := make([]int, n)
	q := make([]int, n)
	ans := make([]int, n)
	for s := 0; s < n; s++ {
		for i := 0; i < n; i++ {
			dist[i] = -1
		}
		head, tail := 0, 0
		dist[s] = 0
		q[tail] = s
		tail++
		for head < tail {
			u := q[head]
			head++
			for _, v := range g[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q[tail] = v
					tail++
				}
			}
		}
		maxd := 0
		for i := 0; i < n; i++ {
			if dist[i] > maxd {
				maxd = dist[i]
			}
		}
		ans[s] = maxd
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = strconv.Itoa(ans[i])
	}
	return strings.Join(res, " ") + "\n"
}

func randomGraph(n int) [][2]int {
	edges := make([][2]int, 0)
	// start with tree edges
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		edges = append(edges, [2]int{p, i})
	}
	// optionally add extra edges
	extra := rand.Intn(n)
	for i := 0; i < extra; i++ {
		a := rand.Intn(n)
		b := rand.Intn(n)
		if a != b {
			edges = append(edges, [2]int{a, b})
		}
	}
	return edges
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(8) + 1
		edges := randomGraph(n)
		m := len(edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		input := sb.String()
		output := solveF(n, edges)
		tests = append(tests, Test{input: input, output: output})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Test %d failed. Input: %q\nExpected: %qGot: %q\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
