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
	n     int
	edges [][2]int
}

func solve(n int, edges [][2]int) int {
	g := make([][]int, n)
	for _, e := range edges {
		x, y := e[0], e[1]
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	w := make([]int, n)
	for i := 0; i < n; i++ {
		w[i] = len(g[i]) - 1
	}
	down := make([]int, n)
	ans := 0
	var dfs func(v, p int)
	dfs = func(v, p int) {
		max1, max2 := 0, 0
		for _, u := range g[v] {
			if u == p {
				continue
			}
			dfs(u, v)
			val := down[u]
			if val > max1 {
				max2 = max1
				max1 = val
			} else if val > max2 {
				max2 = val
			}
		}
		if cur := w[v] + max1 + max2; cur > ans {
			ans = cur
		}
		down[v] = w[v] + max1
	}
	dfs(0, -1)
	return ans + 2
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			x := rng.Intn(j + 1)
			y := j + 1
			edges[j] = [2]int{x, y}
		}
		tests = append(tests, testCase{n: n, edges: edges})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d\n", t.n)
		for _, e := range t.edges {
			input += fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1)
		}
		want := fmt.Sprintf("%d", solve(t.n, t.edges))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
