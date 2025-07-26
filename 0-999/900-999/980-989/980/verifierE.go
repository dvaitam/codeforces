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

func solveE(n, k int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	parent := make([]int, n+1)
	queue := []int{n}
	parent[n] = 0
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			queue = append(queue, u)
		}
	}
	need := n - k
	good := make([]bool, n+1)
	good[n] = true
	count := 1
	for i := n - 1; i >= 1 && count < need; i-- {
		if good[i] {
			continue
		}
		path := []int{}
		x := i
		for !good[x] {
			path = append(path, x)
			x = parent[x]
		}
		if count+len(path) <= need {
			for _, v := range path {
				good[v] = true
			}
			count += len(path)
		}
	}
	res := make([]string, 0, k)
	for i := 1; i <= n; i++ {
		if !good[i] {
			res = append(res, strconv.Itoa(i))
		}
	}
	return strings.Join(res, " ") + "\n"
}

func randomTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		k := rand.Intn(n) + 1
		edges := randomTree(n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		output := solveE(n, k, edges)
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
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
