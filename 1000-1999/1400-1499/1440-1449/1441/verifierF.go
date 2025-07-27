package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type edge struct{ u, v int }
type testCase struct {
	n      int
	m      int
	coords [][2]int
	edges  []edge
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, c := range tc.coords {
		sb.WriteString(fmt.Sprintf("%d %d\n", c[0], c[1]))
	}
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String()
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1441F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 2
	if n%2 == 1 {
		n++
	}
	coords := make([][2]int, n)
	for i := 0; i < n; i++ {
		coords[i] = [2]int{rng.Intn(11), rng.Intn(11)}
	}
	// build tree for connectivity
	edges := make([]edge, 0)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	m := len(edges)
	extra := rng.Intn(n)
	seen := make(map[[2]int]bool)
	for _, e := range edges {
		seen[[2]int{e.u, e.v}] = true
		seen[[2]int{e.v, e.u}] = true
	}
	for i := 0; i < extra; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if seen[[2]int{u, v}] {
			continue
		}
		seen[[2]int{u, v}] = true
		seen[[2]int{v, u}] = true
		edges = append(edges, edge{u, v})
	}
	m = len(edges)
	return testCase{n, m, coords, edges}
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{genCase(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		exp, err := runExe(oracle, tc.Input())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.Input())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
