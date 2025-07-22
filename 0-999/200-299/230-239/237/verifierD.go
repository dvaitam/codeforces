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

type testCase struct {
	input    string
	expected string
}

func solve(n int, edges [][2]int) string {
	c := make([]int, n+1)
	for i, e := range edges {
		x, y := e[0], e[1]
		if c[x] == 0 {
			c[x] = i + 1
		}
		if c[y] == 0 {
			c[y] = i + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("2 %d %d\n", e[0], e[1]))
	}
	for i, e := range edges {
		u, v := e[0], e[1]
		idx := i + 1
		if idx != c[u] {
			sb.WriteString(fmt.Sprintf("%d %d\n", idx, c[u]))
		}
		if idx != c[v] {
			sb.WriteString(fmt.Sprintf("%d %d\n", idx, c[v]))
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	exp := solve(n, edges)
	return testCase{input: sb.String(), expected: exp}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(tc.expected)
	if got != expect {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, generateCase(rand.New(rand.NewSource(1))))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
