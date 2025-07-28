package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1706E.go")
	return runProg(ref, input)
}

type edge struct{ u, v int }

func genGraph(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(4) + 2 // 2..5
	// start with a tree for connectivity
	edges := make([][2]int, 0)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
		parent[i] = p
	}
	maxE := n * (n - 1) / 2
	extra := rng.Intn(maxE - (n - 1) + 1)
	used := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		used[e] = true
	}
	for len(edges) < (n-1)+extra {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		e := [2]int{u, v}
		if used[e] {
			continue
		}
		used[e] = true
		edges = append(edges, e)
	}
	return n, edges
}

func genCase(rng *rand.Rand) string {
	n, edges := genGraph(rng)
	m := len(edges)
	q := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const T = 100
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", T)
	for i := 0; i < T; i++ {
		input.WriteString(genCase(rng))
	}
	expect, err := runRef(input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "reference solver failed:", err)
		os.Exit(1)
	}
	got, err := runProg(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "candidate failed:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		fmt.Printf("output mismatch\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", input.String(), expect, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
