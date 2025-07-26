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

func buildRef() (string, error) {
	ref := "./refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1187G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	m := rng.Intn(n*(n-1)/2) + n - 1
	k := rng.Intn(3) + 1
	c := rng.Intn(5) + 1
	d := rng.Intn(5) + 1
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	// first build tree
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		e := [2]int{p, i}
		used[e] = true
		used[[2]int{i, p}] = true
		edges = append(edges, e)
	}
	// add extra edges
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if used[[2]int{u, v}] {
			continue
		}
		used[[2]int{u, v}] = true
		used[[2]int{v, u}] = true
		edges = append(edges, [2]int{u, v})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", n, m, k, c, d))
	for i := 0; i < k; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(n-1)+2))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		out, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
