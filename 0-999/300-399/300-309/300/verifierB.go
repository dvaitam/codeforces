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

type edge struct{ a, b int }
type test struct {
	n     int
	edges []edge
}

func genTests() []test {
	rand.Seed(2)
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := (rand.Intn(5) + 1) * 3 // 3..15
		students := rand.Perm(n)
		groups := make([][]int, n/3)
		for i := 0; i < n/3; i++ {
			groups[i] = []int{students[3*i] + 1, students[3*i+1] + 1, students[3*i+2] + 1}
		}
		var edges []edge
		for _, g := range groups {
			for i := 0; i < 3; i++ {
				for j := i + 1; j < 3; j++ {
					if rand.Intn(2) == 0 {
						edges = append(edges, edge{g[i], g[j]})
					}
				}
			}
		}
		tests = append(tests, test{n, edges})
	}
	return tests
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, len(t.edges)))
	for _, e := range t.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.a, e.b))
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func verifyOutput(out string, t test) bool {
	out = strings.TrimSpace(out)
	if out == "-1" {
		return false // we always generate solvable cases
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	groups := [][]int{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return false
		}
		g := make([]int, 3)
		for i, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil || v < 1 || v > t.n {
				return false
			}
			g[i] = v
		}
		groups = append(groups, g)
	}
	if len(groups) != t.n/3 {
		return false
	}
	count := make([]int, t.n+1)
	for _, g := range groups {
		for _, v := range g {
			count[v]++
		}
	}
	for i := 1; i <= t.n; i++ {
		if count[i] != 1 {
			return false
		}
	}
	groupOf := make([]int, t.n+1)
	for i, g := range groups {
		for _, v := range g {
			groupOf[v] = i
		}
	}
	for _, e := range t.edges {
		if groupOf[e.a] != groupOf[e.b] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("test %d: run error %v\n", i+1, err)
			os.Exit(1)
		}
		if !verifyOutput(out, t) {
			fmt.Printf("test %d failed. input:\n%soutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
