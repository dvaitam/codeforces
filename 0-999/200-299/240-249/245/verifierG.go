package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func check(id int, adj [][]int) int {
	N := len(adj)
	a := make([]int, N)
	for _, nei := range adj[id] {
		a[nei] = 1
	}
	a[id] = 1
	m := 0
	for _, nei := range adj[id] {
		for _, nn := range adj[nei] {
			if a[nn]%2 == 0 {
				a[nn] += 2
				if a[nn] > m {
					m = a[nn]
				}
			}
		}
	}
	r := 0
	for i := 0; i < N; i++ {
		if a[i] == m {
			r++
		}
	}
	return r
}

func generateCase(rng *rand.Rand) (string, string) {
	N := rng.Intn(5) + 2
	names := make([]string, N)
	for i := 0; i < N; i++ {
		names[i] = fmt.Sprintf("u%d", i)
	}
	edges := make(map[[2]int]bool)
	adj := make([][]int, N)
	// ensure each node has at least one friend
	for i := 0; i < N-1; i++ {
		edges[[2]int{i, i + 1}] = true
	}
	extra := rng.Intn(N)
	for k := 0; k < extra; k++ {
		a := rng.Intn(N)
		b := rng.Intn(N)
		if a == b {
			b = (b + 1) % N
		}
		if a > b {
			a, b = b, a
		}
		edges[[2]int{a, b}] = true
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(edges)))
	for e := range edges {
		a, b := e[0], e[1]
		sb.WriteString(fmt.Sprintf("%s %s\n", names[a], names[b]))
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	// compute expected output
	var keys []string
	for _, n := range names {
		keys = append(keys, n)
	}
	sort.Strings(keys)
	var exp strings.Builder
	exp.WriteString(fmt.Sprintf("%d\n", N))
	for _, name := range keys {
		id := -1
		for i, nm := range names {
			if nm == name {
				id = i
				break
			}
		}
		val := check(id, adj)
		exp.WriteString(fmt.Sprintf("%s %d\n", name, val))
	}
	return sb.String(), strings.TrimSpace(exp.String())
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
