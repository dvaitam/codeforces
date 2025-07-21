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

func run(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func bfs(n, m int, starts [][2]int) (int, int) {
	total := n * m
	visited := make([]bool, total)
	q := make([]int, 0, total)
	var last int
	for _, st := range starts {
		x, y := st[0]-1, st[1]-1
		idx := x*m + y
		if !visited[idx] {
			visited[idx] = true
			q = append(q, idx)
		}
	}
	for head := 0; head < len(q); head++ {
		idx := q[head]
		last = idx
		r := idx / m
		c := idx % m
		if r > 0 {
			ni := (r-1)*m + c
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
		if r+1 < n {
			ni := (r+1)*m + c
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
		if c > 0 {
			ni := r*m + c - 1
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
		if c+1 < m {
			ni := r*m + c + 1
			if !visited[ni] {
				visited[ni] = true
				q = append(q, ni)
			}
		}
	}
	return last/m + 1, last%m + 1
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	k := rng.Intn(min(3, n*m)) + 1
	starts := make([][2]int, k)
	used := make(map[[2]int]bool)
	for i := 0; i < k; i++ {
		for {
			x := rng.Intn(n) + 1
			y := rng.Intn(m) + 1
			if !used[[2]int{x, y}] {
				used[[2]int{x, y}] = true
				starts[i] = [2]int{x, y}
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d\n", k)
	for i, st := range starts {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d %d", st[0], st[1])
	}
	sb.WriteByte('\n')
	input := sb.String()
	x, y := bfs(n, m, starts)
	expected := fmt.Sprintf("%d %d", x, y)
	return input, expected
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	cases = append(cases, [2]string{"1 1\n1\n1 1\n", "1 1"})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for i, tc := range cases {
		out, err := run(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
