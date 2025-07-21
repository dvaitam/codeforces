package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type inputB struct {
	n int
	p []int
	d []int
}

func parseInputB(s string) (inputB, error) {
	rdr := bufio.NewReader(strings.NewReader(s))
	var n int
	if _, err := fmt.Fscan(rdr, &n); err != nil {
		return inputB{}, err
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(rdr, &p[i]); err != nil {
			return inputB{}, err
		}
	}
	d := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(rdr, &d[i]); err != nil {
			return inputB{}, err
		}
	}
	return inputB{n: n, p: p, d: d}, nil
}

func canReach(inp inputB) bool {
	n := inp.n
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		di := inp.d[i]
		for _, j := range []int{i - di, i + di} {
			if j >= 0 && j < n {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}
	visited := make([]bool, n)
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		queue := []int{i}
		visited[i] = true
		comp := []int{i}
		for q := 0; q < len(queue); q++ {
			u := queue[q]
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					queue = append(queue, v)
					comp = append(comp, v)
				}
			}
		}
		initVals := make([]int, len(comp))
		target := make([]int, len(comp))
		for k, idx := range comp {
			initVals[k] = idx + 1
			target[k] = inp.p[idx]
		}
		sort.Ints(initVals)
		sort.Ints(target)
		for k := range initVals {
			if initVals[k] != target[k] {
				return false
			}
		}
	}
	return true
}

func verifyB(input, output string) error {
	inp, err := parseInputB(input)
	if err != nil {
		return fmt.Errorf("invalid input: %v", err)
	}
	expect := canReach(inp)
	ans := strings.ToUpper(strings.TrimSpace(output))
	if ans != "YES" && ans != "NO" {
		return fmt.Errorf("output must be YES or NO")
	}
	if expect && ans != "YES" {
		return fmt.Errorf("expected YES got %s", ans)
	}
	if !expect && ans != "NO" {
		return fmt.Errorf("expected NO got %s", ans)
	}
	return nil
}

func runCase(bin, tc string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verifyB(tc, out.String())
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	p := make([]int, n)
	for i, v := range perm {
		p[i] = v + 1
	}
	d := make([]int, n)
	for i := range d {
		d[i] = rng.Intn(n) + 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", p[i])
	}
	b.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", d[i])
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
