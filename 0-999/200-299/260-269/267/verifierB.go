package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

func solveDomino(x, y []int) string {
	n := len(x)
	var a [7][7]int
	deg := [7]int{}
	for i := 0; i < n; i++ {
		a[x[i]][y[i]]++
		a[y[i]][x[i]]++
		deg[x[i]]++
		deg[y[i]]++
	}
	start := x[0]
	odd := 0
	for i := 0; i <= 6; i++ {
		if deg[i]%2 != 0 {
			odd++
			start = i
		}
	}
	if odd > 2 {
		return "No solution"
	}
	ans := make([]int, 0, n+1)
	var dfs func(int)
	dfs = func(u int) {
		for v := 0; v <= 6; v++ {
			for a[u][v] > 0 {
				a[u][v]--
				a[v][u]--
				dfs(v)
			}
		}
		ans = append(ans, u)
	}
	dfs(start)
	if len(ans) != n+1 {
		return "No solution"
	}
	for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
		ans[i], ans[j] = ans[j], ans[i]
	}
	used := make([]bool, n)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		u, v := ans[i], ans[i+1]
		found := false
		for j := 0; j < n; j++ {
			if used[j] {
				continue
			}
			if u == x[j] && v == y[j] {
				fmt.Fprintf(&sb, "%d +\n", j+1)
				used[j] = true
				found = true
				break
			} else if u == y[j] && v == x[j] {
				fmt.Fprintf(&sb, "%d -\n", j+1)
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return "No solution"
		}
	}
	return strings.TrimSpace(sb.String())
}

func buildCase(dominoes [][2]int) testCase {
	n := len(dominoes)
	var in strings.Builder
	x := make([]int, n)
	y := make([]int, n)
	fmt.Fprintf(&in, "%d\n", n)
	for i, d := range dominoes {
		x[i] = d[0]
		y[i] = d[1]
		fmt.Fprintf(&in, "%d %d\n", d[0], d[1])
	}
	exp := solveDomino(x, y)
	if exp != "No solution" {
		exp += "\n"
	}
	return testCase{in.String(), exp}
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(2))
	cases := make([]testCase, 0, 120)
	cases = append(cases, buildCase([][2]int{{1, 1}}))
	cases = append(cases, buildCase([][2]int{{1, 2}, {2, 1}}))
	for i := 0; i < 118; i++ {
		n := rng.Intn(8) + 1
		dom := make([][2]int, n)
		for j := 0; j < n; j++ {
			dom[j] = [2]int{rng.Intn(7), rng.Intn(7)}
		}
		cases = append(cases, buildCase(dom))
	}
	return cases
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected:\n%s\n got:\n%s", tc.expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
