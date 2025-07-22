package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ a, b int }

func runSolution(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(9) + 2 // 2..10
	center := rng.Intn(n) + 1
	var pairs []pair
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if i == center || j == center {
				continue
			}
			pairs = append(pairs, pair{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	m := rng.Intn(len(pairs) + 1)
	banned := pairs[:m]

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, p := range banned {
		fmt.Fprintf(&sb, "%d %d\n", p.a, p.b)
	}
	return sb.String()
}

func verifyB(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return fmt.Errorf("input parse: %v", err)
	}
	banned := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a > b {
			a, b = b, a
		}
		banned[[2]int{a, b}] = true
	}

	outR := bufio.NewReader(strings.NewReader(output))
	var s int
	if _, err := fmt.Fscan(outR, &s); err != nil {
		return fmt.Errorf("output parse s: %v", err)
	}
	if s != n-1 {
		return fmt.Errorf("expected %d edges, got %d", n-1, s)
	}
	deg := make([]int, n+1)
	edges := make(map[[2]int]bool)
	center := 0
	for i := 0; i < s; i++ {
		var a, b int
		if _, err := fmt.Fscan(outR, &a, &b); err != nil {
			return fmt.Errorf("parse edge %d: %v", i+1, err)
		}
		if a < 1 || a > n || b < 1 || b > n || a == b {
			return fmt.Errorf("invalid edge %d: %d %d", i+1, a, b)
		}
		x, y := a, b
		if x > y {
			x, y = y, x
		}
		if banned[[2]int{x, y}] {
			return fmt.Errorf("edge %d between %d %d is banned", i+1, a, b)
		}
		if edges[[2]int{x, y}] {
			return fmt.Errorf("duplicate edge %d %d", a, b)
		}
		edges[[2]int{x, y}] = true
		deg[a]++
		deg[b]++
	}
	for v := 1; v <= n; v++ {
		if deg[v] == n-1 {
			center = v
			break
		}
	}
	if center == 0 {
		return fmt.Errorf("no center vertex with degree n-1")
	}
	for e := range edges {
		if e[0] != center && e[1] != center {
			return fmt.Errorf("edge %d %d does not use center %d", e[0], e[1], center)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []string{
		"2 0\n",
		"3 1\n1 2\n",
	}
	for i := len(cases); i < 100; i++ {
		cases = append(cases, generateCaseB(rng))
	}

	for i, tc := range cases {
		out, err := runSolution(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyB(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
