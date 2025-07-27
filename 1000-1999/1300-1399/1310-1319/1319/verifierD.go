package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type edge struct{ u, v int }

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "1319D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2 // 2..6
	// build edges map
	edgeMap := make(map[[2]int]bool)
	// ensure cycle for strong connectivity
	for i := 1; i <= n; i++ {
		u := i
		v := i%n + 1
		edgeMap[[2]int{u, v}] = true
	}
	// random extra edges
	for u := 1; u <= n; u++ {
		for v := 1; v <= n; v++ {
			if u == v {
				continue
			}
			if rng.Float64() < 0.3 {
				edgeMap[[2]int{u, v}] = true
			}
		}
	}
	edges := make([]edge, 0, len(edgeMap))
	for k := range edgeMap {
		edges = append(edges, edge{k[0], k[1]})
	}
	// path
	klen := rng.Intn(n-1) + 2
	perm := rng.Perm(n)
	path := make([]int, klen)
	for i := 0; i < klen; i++ {
		path[i] = perm[i] + 1
	}
	for i := 0; i < klen-1; i++ {
		edgeMap[[2]int{path[i], path[i+1]}] = true
	}
	// rebuild edges
	edges = edges[:0]
	for k := range edgeMap {
		edges = append(edges, edge{k[0], k[1]})
	}
	m := len(edges)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	fmt.Fprintln(&sb, klen)
	for i, v := range path {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
