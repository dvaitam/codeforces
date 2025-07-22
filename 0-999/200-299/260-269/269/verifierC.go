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

type testCase struct {
	input string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "269C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	// build tree to ensure connectivity
	edges := make([][3]int, 0, m)
	parent := rng.Perm(n - 1)
	for i := 2; i <= n; i++ {
		p := parent[i-2] + 1
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{p, i, w})
	}
	used := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] < e[1] {
			used[[2]int{e[0], e[1]}] = true
		} else {
			used[[2]int{e[1], e[0]}] = true
		}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if used[[2]int{u, v}] {
			continue
		}
		used[[2]int{u, v}] = true
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{u, v, w})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return testCase{input: sb.String()}
}

func runCase(bin, oracle string, tc testCase) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(tc.input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
	cases := make([]testCase, 0, 101)
	for i := 0; i < 101; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, oracle, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
